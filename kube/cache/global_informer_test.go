package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/major1201/kubetrack/kube"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

func TestNewGlobalInformer(t *testing.T) {
	ta := assert.New(t)

	kc, err := kube.NewClient()
	if err != nil {
		ta.Error(err)
		return
	}

	gi := NewGlobalInformer(kube.GetScheme())
	h := &testClusterResourceEventHandler{}
	gi.AddCluster("1", kc, time.Minute, nil, []ResourceUnitWithHandlers{
		{
			ResourceUnit:          ResourceUnit{Namespace: "chtest", Resource: podGVR},
			ResourceEventHandlers: []ClusterResourceEventHandler{h},
			WatchErrorHandlers:    []ClusterWatchErrorHandler{h.WatchErrorHandler},
		},
	})

	waitUntilAllSynced(gi)
	const (
		podNamespace = "chtest"
		podName      = "mypod1"
	)
	pod := &corev1.Pod{}
	pod.SetNamespace(podNamespace)
	pod.SetName(podName)
	pod.Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:  "main",
				Image: "nginx",
			},
		},
	}
}

type testClusterResourceEventHandler struct{}

func (h *testClusterResourceEventHandler) OnAdd(cluster Cluster, obj any) {
	fmt.Println(cluster.ID(), "OnAdd", obj)
}

func (h *testClusterResourceEventHandler) OnUpdate(cluster Cluster, oldObj, newObj any) {
	fmt.Println(cluster.ID(), "OnUpdate", oldObj, newObj)
}

func (h *testClusterResourceEventHandler) OnDelete(cluster Cluster, obj any) {
	fmt.Println(cluster.ID(), "OnDelete", obj)
}

func (h *testClusterResourceEventHandler) WatchErrorHandler(cluster Cluster, r *cache.Reflector, err error) {
	fmt.Println("WatchError", cluster.ID(), r, err)
}

func waitUntilAllSynced(gi GlobalInformer) {
	for i := 0; i < 30; i++ {
		if gi.AllSynced() {
			return
		}
		time.Sleep(time.Second)
	}
	panic("informer still not synced")
}
