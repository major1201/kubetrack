package cache

import (
	"testing"
	"time"

	"github.com/major1201/kubetrack/kube"
	"github.com/major1201/kubetrack/utils/goutils"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	podGVR  = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	podGVK  = schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}
	nodeGVR = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "nodes"}
	//chrGVR  = schema.GroupVersionResource{Group: "dp.byted.org", Version: "v1", Resource: "clickhousereplicas"}
	crdGVR = schema.GroupVersionResource{Group: "crd.projectcalico.org", Version: "v1", Resource: "ippools"}
)

func TestResourceBuilder_List(t *testing.T) {
	ta := assert.New(t)

	kubeClient, err := kube.NewClient()
	goutils.Must(err)

	gi := NewGlobalInformer(kube.GetScheme())
	gi.AddCluster("1", kubeClient, 30*time.Second, nil, []ResourceUnitWithHandlers{
		{ResourceUnit: ResourceUnit{Resource: podGVR}},
		{ResourceUnit: ResourceUnit{Resource: nodeGVR}},
	})

	waitUntilAllSynced(gi)

	// case 1: pod gvr
	pods, err := NewResourceBuilder[*corev1.Pod](gi).ForResource(podGVR).AllClusters().List()
	ta.NoError(err)
	goutils.Must(err)
	println("case 1:", pods)

	// case 2: crd gvr
	//list, err := NewResourceBuilder[*unstructured.Unstructured](gi).ForResource(chrGVR).AllClusters().List()
	//ta.NoError(err)
	//goutils.Must(err)
	//chrs := list.([]*dpv1.ClickHouseReplica)
	//println("case 2:", chrs)

	// case 3: pod gvk
	pods, err = NewResourceBuilder[*corev1.Pod](gi).ForKind(podGVK).AllClusters().List()
	ta.NoError(err)
	goutils.Must(err)
	println("case 3:", pods)

	// case 4: pod model
	pods, err = NewResourceBuilder[*corev1.Pod](gi).AllClusters().List()
	ta.NoError(err)
	goutils.Must(err)
	println("case 4:", pods)

	// case 5: not namespaced gvr
	nodes, err := NewResourceBuilder[*corev1.Node](gi).AllClusters().List()
	ta.NoError(err)
	goutils.Must(err)
	println("case 5:", nodes)

	//gi.RemoveCluster("1")
}

func TestResourceBuilder_Get(t *testing.T) {
	ta := assert.New(t)

	kubeClient, err := kube.NewClient()
	goutils.Must(err)

	gi := NewGlobalInformer(kube.GetScheme())
	gi.AddCluster("1", kubeClient, 30*time.Second, nil, []ResourceUnitWithHandlers{
		{ResourceUnit: ResourceUnit{Resource: podGVR}},
		{ResourceUnit: ResourceUnit{Resource: crdGVR}},
		{ResourceUnit: ResourceUnit{Resource: nodeGVR}},
	})

	waitUntilAllSynced(gi)

	var obj any
	obj, err = NewResourceBuilder[*corev1.Pod](gi).AllClusters().Get("chtest", "testpod")
	ta.NoError(err)
	goutils.Must(err)
	println(obj.(*corev1.Pod))
}
