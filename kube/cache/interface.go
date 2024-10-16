package cache

import (
	"time"

	"github.com/major1201/kubetrack/kube"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

type ClusterID string

func (cid ClusterID) String() string {
	return string(cid)
}

type GlobalInformer interface {
	// cluster cares
	GetCluster(id ClusterID) Cluster
	ListClusters() (clusters []Cluster)
	AddCluster(clusterID ClusterID, client kube.Client, defaultResync time.Duration, tweakListOptions dynamicinformer.TweakListOptionsFunc, watchUnits []ResourceUnitWithHandlers)
	RemoveCluster(clusterID ClusterID)
	HasSynced(clusterID ClusterID, resource schema.GroupVersionResource) bool
	ClusterHasSynced(clusterID ClusterID) bool
	ClusterSyncMap(clusterID ClusterID) map[ResourceUnit]bool
	AllSynced() bool

	// resource
	//Resources() ResourceBuilder // TODO: put it back until golang support type parameters in methods
}

type Cluster interface {
	ID() ClusterID
	KubeClient() kube.Client
}

type ResourceUnit struct {
	Namespace string // leave empty to watch all namespaces
	Resource  schema.GroupVersionResource
}

type ResourceUnitWithHandlers struct {
	ResourceUnit

	ResourceEventHandlers []ClusterResourceEventHandler
	WatchErrorHandlers    []ClusterWatchErrorHandler
}

func BuildResourceUnitWithHandlersSlice(units []ResourceUnit) (res []ResourceUnitWithHandlers) {
	res = make([]ResourceUnitWithHandlers, len(units))
	for i, unit := range units {
		res[i].ResourceUnit = unit
	}
	return
}

type ClusterResourceEventHandler interface {
	OnAdd(cluster Cluster, obj any)
	OnUpdate(cluster Cluster, oldObj, newObj any)
	OnDelete(cluster Cluster, obj any)
}

type ClusterWatchErrorHandler func(cluster Cluster, r *cache.Reflector, err error)

type ResourceBuilder[T runtime.Object] interface {
	// resource type
	ForResource(gvr schema.GroupVersionResource) ResourceBuilder[T]
	ForKind(gvk schema.GroupVersionKind) ResourceBuilder[T]

	// include clusters
	Clusters(clusters ...ClusterID) ResourceBuilder[T]
	AllClusters() ResourceBuilder[T]

	// lister
	List(opts ...ListOption[T]) (list []T, err error)
	Get(namespace, name string) (obj T, err error)
}
