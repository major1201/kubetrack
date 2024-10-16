package kube

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/ghodss/yaml"
	"github.com/major1201/kubetrack/utils"
	"github.com/major1201/kubetrack/utils/slicex"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kubernetes/pkg/util/node"
)

func GetUnstructuredFromYAML(s, namespace string) (objs []*unstructured.Unstructured, err error) {
	re := regexp.MustCompilePOSIX("^---\n")
	for _, part := range re.Split(s, -1) {
		if utils.IsBlank(part) {
			continue
		}

		jsonBytes, err2 := yaml.YAMLToJSON([]byte(part))
		if err2 != nil {
			err = errors.WithStack(err2)
			return
		}

		if string(jsonBytes) == "null" {
			continue
		}

		obj := &unstructured.Unstructured{}
		if err = obj.UnmarshalJSON(jsonBytes); err != nil {
			err = errors.WithStack(err)
			return
		}

		// set namespace
		if namespace != "" && obj.GetNamespace() == "" {
			obj.SetNamespace(namespace)
		}

		objs = append(objs, obj)
	}
	return
}

// GetObjectsFromYAML get actual object from yaml text
func GetObjectsFromYAML(s, namespace string) (objs []runtime.Object, err error) {
	unstList, err := GetUnstructuredFromYAML(s, namespace)
	if err != nil {
		return
	}

	var retObjs []runtime.Object
	converter := runtime.ObjectConvertor(scheme)
	groupVersioner := runtime.GroupVersioner(schema.GroupVersions(scheme.PrioritizedVersionsAllGroups()))

	for _, unst := range unstList {
		outObj, err := converter.ConvertToVersion(unst, groupVersioner)
		if err != nil {
			outObj = unst
		}

		retObjs = append(retObjs, outObj)
	}

	objs = retObjs
	return
}

// PodStatus returns the pod's current status
// source code: https://github.com/kubernetes/kubernetes/blob/master/pkg/printers/internalversion/printers.go#L740
func PodStatus(pod *corev1.Pod) (reason string) {
	if pod == nil {
		return ""
	}

	reason = string(pod.Status.Phase)
	if pod.Status.Reason != "" {
		reason = pod.Status.Reason
	}

	initializing := false
	for i := range pod.Status.InitContainerStatuses {
		container := pod.Status.InitContainerStatuses[i]
		switch {
		case container.State.Terminated != nil && container.State.Terminated.ExitCode == 0:
			continue
		case container.State.Terminated != nil:
			// initialization is failed
			if len(container.State.Terminated.Reason) == 0 {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Init:Signal:%d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("Init:ExitCode:%d", container.State.Terminated.ExitCode)
				}
			} else {
				reason = "Init:" + container.State.Terminated.Reason
			}
			initializing = true
		case container.State.Waiting != nil && len(container.State.Waiting.Reason) > 0 && container.State.Waiting.Reason != "PodInitializing":
			reason = "Init:" + container.State.Waiting.Reason
			initializing = true
		default:
			reason = fmt.Sprintf("Init:%d/%d", i, len(pod.Spec.InitContainers))
			initializing = true
		}
		break
	}
	if !initializing {
		hasRunning := false
		for i := len(pod.Status.ContainerStatuses) - 1; i >= 0; i-- {
			container := pod.Status.ContainerStatuses[i]

			if container.State.Waiting != nil && container.State.Waiting.Reason != "" {
				reason = container.State.Waiting.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason != "" {
				reason = container.State.Terminated.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason == "" {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Signal:%d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("ExitCode:%d", container.State.Terminated.ExitCode)
				}
			} else if container.Ready && container.State.Running != nil {
				hasRunning = true
			}
		}

		// change pod status back to "Running" if there is at least one container still reporting as "Running" status
		if reason == "Completed" && hasRunning {
			if hasPodReadyCondition(pod.Status.Conditions) {
				reason = "Running"
			} else {
				reason = "NotReady"
			}
		}
	}

	if pod.DeletionTimestamp != nil && pod.Status.Reason == node.NodeUnreachablePodReason {
		reason = "Unknown"
	} else if pod.DeletionTimestamp != nil {
		reason = "Terminating"
	}

	return
}

func hasPodReadyCondition(conditions []corev1.PodCondition) bool {
	for _, condition := range conditions {
		if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func NodeStatus(node *corev1.Node) (status []string) {
	conditionMap := make(map[corev1.NodeConditionType]*corev1.NodeCondition)
	NodeAllConditions := []corev1.NodeConditionType{corev1.NodeReady}
	for i := range node.Status.Conditions {
		cond := node.Status.Conditions[i]
		conditionMap[cond.Type] = &cond
	}
	for _, validCondition := range NodeAllConditions {
		if condition, ok := conditionMap[validCondition]; ok {
			if condition.Status == corev1.ConditionTrue {
				status = append(status, string(condition.Type))
			} else {
				status = append(status, "Not"+string(condition.Type))
			}
		}
	}
	if len(status) == 0 {
		status = append(status, "Unknown")
	}
	if node.Spec.Unschedulable {
		status = append(status, "SchedulingDisabled")
	}
	return
}

const (
	// labelNodeRolePrefix is a label prefix for node roles
	// It's copied over to here until it's merged in core: https://github.com/kubernetes/kubernetes/pull/39112
	labelNodeRolePrefix = "node-role.kubernetes.io/"

	// nodeLabelRole specifies the role of a node
	nodeLabelRole = "kubernetes.io/role"
)

// FindNodeRoles returns the roles of a given node.
// The roles are determined by looking for:
// * a node-role.kubernetes.io/<role>="" label
// * a kubernetes.io/role="<role>" label
func FindNodeRoles(node *corev1.Node) []string {
	roles := sets.NewString()
	for k, v := range node.Labels {
		switch {
		case strings.HasPrefix(k, labelNodeRolePrefix):
			if role := strings.TrimPrefix(k, labelNodeRolePrefix); len(role) > 0 {
				roles.Insert(role)
			}

		case k == nodeLabelRole && v != "":
			roles.Insert(v)
		}
	}
	return roles.List()
}

func (c *ClientImpl) isReady(obj runtime.Object, withUpdate bool) (ready bool, retObj runtime.Object, err error) {
	switch obj := obj.(type) {
	case *corev1.ReplicationController:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.CoreV1().ReplicationControllers(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = replicationControllerReady(realObj)
	case *corev1.Pod:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.CoreV1().Pods(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = podReady(realObj)
	case *appsv1.Deployment:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.AppsV1().Deployments(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = deploymentReady(realObj)
	case *appsv1beta1.Deployment:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.AppsV1beta1().Deployments(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = deploymentBeta1Ready(realObj)
	case *appsv1beta2.Deployment:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.AppsV1beta2().Deployments(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = deploymentBeta2Ready(realObj)
	case *extensions.Deployment:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.ExtensionsV1beta1().Deployments(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = deploymentExtensionsReady(realObj)
	case *extensions.DaemonSet:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.ExtensionsV1beta1().DaemonSets(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = daemonSetExtensionsReady(realObj)
	case *appsv1.DaemonSet:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.AppsV1().DaemonSets(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = daemonSetReady(realObj)
	case *appsv1beta2.DaemonSet:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.AppsV1beta2().DaemonSets(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = daemonSetBeta2Ready(realObj)
	case *appsv1.StatefulSet:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.AppsV1().StatefulSets(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = statefulSetReady(realObj)
	case *appsv1beta1.StatefulSet:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.AppsV1beta1().StatefulSets(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = statefulSetBeta1Ready(realObj)
	case *appsv1beta2.StatefulSet:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.AppsV1beta2().StatefulSets(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = statefulSetBeta2Ready(realObj)
	case *extensions.ReplicaSet:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.ExtensionsV1beta1().ReplicaSets(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = replicaSetExtensionsReady(realObj)
	case *appsv1beta2.ReplicaSet:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.AppsV1beta2().ReplicaSets(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = replicaSetBeta2Ready(realObj)
	case *appsv1.ReplicaSet:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.AppsV1().ReplicaSets(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = replicaSetReady(realObj)
	case *corev1.PersistentVolumeClaim:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.CoreV1().PersistentVolumeClaims(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = persistentVolumeClaimReady(realObj)
	case *corev1.Service:
		realObj := obj.DeepCopy()
		if withUpdate {
			realObj, err = c.kubeClient.CoreV1().Services(obj.GetNamespace()).Get(c.ctx, obj.GetName(), metav1.GetOptions{})
			if err != nil {
				err = errors.WithStack(err)
				return
			}
			realObj.TypeMeta = obj.TypeMeta
			retObj = realObj
		}
		ready = serviceReady(realObj)
	case *unstructured.Unstructured:
		ready, retObj, err = c.isReadyUnstructured(obj)
		if err != nil {
			return
		}
	default:
		ready = true
	}

	if retObj == nil {
		retObj = obj
	}

	return
}

func (c *ClientImpl) isReadyUnstructured(unst *unstructured.Unstructured) (ready bool, realObj *unstructured.Unstructured, err error) {
	gvk := fmt.Sprintf("%s/%s", unst.GetAPIVersion(), unst.GetKind())

	type action struct {
		resource string
		f        func(unst *unstructured.Unstructured, result gjson.Result) bool
	}

	gvkActMap := map[string]action{
		"kubeflow.org/v1/Notebook": {
			resource: "notebooks",
			f: func(unst *unstructured.Unstructured, result gjson.Result) bool {
				return result.Get("status.readyReplicas").Int() > 0
			},
		},
	}

	if act, ok := gvkActMap[gvk]; ok {
		// get real object
		gv, e := schema.ParseGroupVersion(unst.GetAPIVersion())
		if e != nil {
			err = errors.WithStack(e)
			return
		}
		gvr := gv.WithResource(act.resource)
		realObj, e = c.dynamicClient.Resource(gvr).Namespace(unst.GetNamespace()).Get(c.ctx, unst.GetName(), metav1.GetOptions{})
		if e != nil {
			err = errors.WithStack(e)
			return
		}

		// parse gjson result
		jsonbytes, e := realObj.MarshalJSON()
		if e != nil {
			err = errors.WithStack(e)
			return
		}
		result := gjson.ParseBytes(jsonbytes)

		// return
		ready = act.f(realObj, result)
		return
	}

	ready = true
	realObj = unst
	return
}

// IsReady returns whether the runtime object is ready or not
func (c *ClientImpl) IsReady(obj runtime.Object) (ready bool, err error) {
	ready, _, err = c.isReady(obj, false)
	return
}

// IsReadyWithUpdate returns whether the runtime object is ready or not with update
func (c *ClientImpl) IsReadyWithUpdate(obj runtime.Object) (ready bool, retObj runtime.Object, err error) {
	return c.isReady(obj, true)
}

func replicationControllerReady(rc *corev1.ReplicationController) bool {
	return rc.Status.AvailableReplicas == rc.Status.Replicas
}

func podReady(pod *corev1.Pod) bool {
	if pod == nil {
		return false
	}
	if pod.GetDeletionTimestamp() != nil {
		return false
	}
	if len(pod.Status.Conditions) > 0 {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
				return true
			}
		}
	}
	return false
}

func deploymentReady(deploy *appsv1.Deployment) bool {
	return deploy.Status.AvailableReplicas == deploy.Status.Replicas
}

func deploymentBeta1Ready(deploy *appsv1beta1.Deployment) bool {
	return deploy.Status.AvailableReplicas == deploy.Status.Replicas
}

func deploymentBeta2Ready(deploy *appsv1beta2.Deployment) bool {
	return deploy.Status.AvailableReplicas == deploy.Status.Replicas
}

func deploymentExtensionsReady(deploy *extensions.Deployment) bool {
	return deploy.Status.AvailableReplicas == deploy.Status.Replicas
}

func daemonSetExtensionsReady(ds *extensions.DaemonSet) bool {
	return ds.Status.DesiredNumberScheduled == ds.Status.NumberAvailable
}

func daemonSetReady(ds *appsv1.DaemonSet) bool {
	return ds.Status.DesiredNumberScheduled == ds.Status.NumberAvailable
}

func daemonSetBeta2Ready(ds *appsv1beta2.DaemonSet) bool {
	return ds.Status.DesiredNumberScheduled == ds.Status.NumberAvailable
}

func statefulSetReady(ss *appsv1.StatefulSet) bool {
	return ss.Status.ReadyReplicas == ss.Status.Replicas
}

func statefulSetBeta1Ready(ss *appsv1beta1.StatefulSet) bool {
	return ss.Status.ReadyReplicas == ss.Status.Replicas
}

func statefulSetBeta2Ready(ss *appsv1beta2.StatefulSet) bool {
	return ss.Status.ReadyReplicas == ss.Status.Replicas
}

func replicaSetReady(rs *appsv1.ReplicaSet) bool {
	return rs.Status.AvailableReplicas == rs.Status.Replicas
}

func replicaSetBeta2Ready(rs *appsv1beta2.ReplicaSet) bool {
	return rs.Status.AvailableReplicas == rs.Status.Replicas
}

func replicaSetExtensionsReady(rs *extensions.ReplicaSet) bool {
	return rs.Status.AvailableReplicas == rs.Status.Replicas
}

func persistentVolumeClaimReady(pvc *corev1.PersistentVolumeClaim) bool {
	return pvc.Status.Phase == corev1.ClaimBound
}

func serviceReady(svc *corev1.Service) bool {
	if svc.Spec.Type == corev1.ServiceTypeExternalName {
		return true
	}

	// Make sure the service is not explicitly set to "None" before checking the IP
	if svc.Spec.ClusterIP == "" {
		return false
	}
	// This checks if the service has a LoadBalancer and that balancer has an Ingress defined
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer && svc.Status.LoadBalancer.Ingress == nil {
		return false
	}

	return true
}

// ListAllResources list all resources within namespace, if namespace is empty, the scale would be all namespaces
func (c *ClientImpl) ListAllResources(namespace string) (objs []unstructured.Unstructured, err error) {
	resList, err := c.GetDiscoveryClient().ServerPreferredResources()
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	// fetch all resource types
	var gvrList []schema.GroupVersionResource
	for _, group := range resList {
		for _, res := range group.APIResources {
			if !slicex.Contains(res.Verbs, "list") {
				// filter resources which can't be listed
				continue
			}
			if namespace != "" && !res.Namespaced {
				// filter not namespaced resources when specified namespace
				continue
			}

			gv, e := schema.ParseGroupVersion(group.GroupVersion)
			if e != nil {
				err = errors.WithStack(e)
				return
			}
			gvr := gv.WithResource(res.Name)
			gvrList = append(gvrList, gvr)
		}
	}

	// fetch all resources in parallel
	var wg sync.WaitGroup
	wg.Add(len(gvrList))
	for _, gvr := range gvrList {
		go func(gvr schema.GroupVersionResource) {
			defer wg.Done()

			rif := c.GetDynamicClient().Resource(gvr).Namespace(namespace)

			var next string
			for {
				// fetch all pages
				resp, e := rif.List(c.ctx, metav1.ListOptions{
					Limit:    250,
					Continue: next,
				})
				if e != nil {
					err = errors.Wrapf(e, "list resource failed: %s", gvr.String())
					return
				}
				objs = append(objs, resp.Items...)

				next = resp.GetContinue()
				if next == "" {
					break
				}
			}
		}(gvr)
	}
	wg.Wait()
	return
}
