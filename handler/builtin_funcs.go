package handler

import (
	"fmt"
	"strings"

	"github.com/major1201/kubetrack/kube"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type BuiltInFunc func(obj *unstructured.Unstructured, toVersioned func() (runtime.Object, error)) string

var builtInFuncMap = map[string]BuiltInFunc{
	"PodStatus":                 PodStatus,
	"PodStatusWithRestartCount": PodStatusWithRestartCount,
	"NodeStatus":                NodeStatus,
	"FindNodeRoles":             FindNodeRoles,
}

func PodStatus(obj *unstructured.Unstructured, toVersioned func() (runtime.Object, error)) string {
	podIf, err := toVersioned()
	if err != nil {
		return fmt.Sprintf("Error: PodStatus() toVersioned failed, err=%s", err.Error())
	}

	pod, ok := podIf.(*corev1.Pod)
	if !ok {
		return fmt.Sprintf("Error: PodStatus(), not a pod object, gvk=%s", obj.GroupVersionKind().String())
	}
	return kube.PodStatus(pod)
}

func PodStatusWithRestartCount(obj *unstructured.Unstructured, toVersioned func() (runtime.Object, error)) string {
	podIf, err := toVersioned()
	if err != nil {
		return fmt.Sprintf("Error: PodStatus() toVersioned failed, err=%s", err.Error())
	}

	pod, ok := podIf.(*corev1.Pod)
	if !ok {
		return fmt.Sprintf("Error: PodStatus(), not a pod object, gvk=%s", obj.GroupVersionKind().String())
	}

	totalContainers := len(pod.Spec.Containers)
	readyContainers := 0

	var restartCount int32
	for _, cs := range pod.Status.ContainerStatuses {
		restartCount += cs.RestartCount
		if cs.Ready && cs.State.Running != nil {
			readyContainers++
		}
	}
	return kube.PodStatus(pod) + Ternary(readyContainers == totalContainers, "", fmt.Sprintf(" (%d/%d)", readyContainers, totalContainers))
}

func NodeStatus(obj *unstructured.Unstructured, toVersioned func() (runtime.Object, error)) string {
	nodeIf, err := toVersioned()
	if err != nil {
		return fmt.Sprintf("Error: NodeStatus() toVersioned failed, err=%s", err.Error())
	}

	node, ok := nodeIf.(*corev1.Node)
	if !ok {
		return fmt.Sprintf("Error: NodeStatus(), not a node object, gvk=%s", obj.GroupVersionKind().String())
	}
	return strings.Join(kube.NodeStatus(node), ",")
}

func FindNodeRoles(obj *unstructured.Unstructured, toVersioned func() (runtime.Object, error)) string {
	nodeIf, err := toVersioned()
	if err != nil {
		return fmt.Sprintf("Error: NodeStatus() toVersioned failed, err=%s", err.Error())
	}

	node, ok := nodeIf.(*corev1.Node)
	if !ok {
		return fmt.Sprintf("Error: NodeStatus(), not a node object, gvk=%s", obj.GroupVersionKind().String())
	}
	return strings.Join(kube.FindNodeRoles(node), ",")
}
