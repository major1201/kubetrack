package cache

import (
	"time"

	"github.com/major1201/kubetrack/kube"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

type globalInformer struct {
	clusterMap     map[ClusterID]*cluster
	scheme         *runtime.Scheme
	groupVersioner runtime.GroupVersioner
}

type informerEntity struct {
	informer dynamicinformer.DynamicSharedInformerFactory
	stopCh   chan struct{}

	resourceEventHandlers []ClusterResourceEventHandler
	watchErrorHandlers    []ClusterWatchErrorHandler
}

type cluster struct {
	id               ClusterID
	client           kube.Client
	defaultResync    time.Duration
	tweakListOptions dynamicinformer.TweakListOptionsFunc
	informerMap      map[ResourceUnit]informerEntity
}

func (c *cluster) ID() ClusterID {
	if c == nil {
		return ""
	}
	return c.id
}

func (c *cluster) KubeClient() kube.Client {
	if c == nil {
		return nil
	}
	return c.client
}

func (c *cluster) startInformers(gi *globalInformer) {
	for unit, entity := range c.informerMap {
		if entity.informer != nil {
			continue
		}

		informer := dynamicinformer.NewFilteredDynamicSharedInformerFactory(c.client.GetDynamicClient(), c.defaultResync, unit.Namespace, c.tweakListOptions)
		ch := make(chan struct{})
		informer.Start(ch)

		// add resources
		handlerWrapper := newClusterHandlerWrapper(c, unit, entity.resourceEventHandlers, entity.watchErrorHandlers)
		siInformer := informer.ForResource(unit.Resource).Informer()
		siInformer.AddEventHandler(handlerWrapper)
		_ = siInformer.SetWatchErrorHandler(handlerWrapper.WatchErrorHandler)
		go siInformer.Run(ch)

		entity.informer = informer
		entity.stopCh = ch
		c.informerMap[unit] = entity
	}
}

type clusterHandlerWrapper struct {
	cluster Cluster
	unit    ResourceUnit

	resourceEventHandlers []ClusterResourceEventHandler
	watchErrorHandlers    []ClusterWatchErrorHandler
}

func newClusterHandlerWrapper(cluster Cluster, unit ResourceUnit, resourceEventHandlers []ClusterResourceEventHandler, watchErrorHandlers []ClusterWatchErrorHandler) *clusterHandlerWrapper {
	return &clusterHandlerWrapper{
		cluster:               cluster,
		unit:                  unit,
		resourceEventHandlers: resourceEventHandlers,
		watchErrorHandlers:    watchErrorHandlers,
	}
}

func (w *clusterHandlerWrapper) OnAdd(obj interface{}) {
	for _, handler := range w.resourceEventHandlers {
		go handler.OnAdd(w.cluster, obj)
	}
}

func (w *clusterHandlerWrapper) OnUpdate(oldObj, newObj interface{}) {
	for _, handler := range w.resourceEventHandlers {
		go handler.OnUpdate(w.cluster, oldObj, newObj)
	}
}

func (w *clusterHandlerWrapper) OnDelete(obj interface{}) {
	for _, handler := range w.resourceEventHandlers {
		go handler.OnDelete(w.cluster, obj)
	}
}

func (w *clusterHandlerWrapper) WatchErrorHandler(r *cache.Reflector, err error) {
	for _, handler := range w.watchErrorHandlers {
		go handler(w.cluster, r, err)
	}
}

func NewGlobalInformer(scheme *runtime.Scheme) GlobalInformer {
	return &globalInformer{
		clusterMap:     make(map[ClusterID]*cluster),
		scheme:         scheme,
		groupVersioner: runtime.GroupVersioner(schema.GroupVersions(scheme.PrioritizedVersionsAllGroups())),
	}
}

func (gi *globalInformer) GetCluster(id ClusterID) Cluster {
	return gi.clusterMap[id]
}

func (gi *globalInformer) ListClusters() (clusters []Cluster) {
	for _, c := range gi.clusterMap {
		clusters = append(clusters, c)
	}
	return
}

func (gi *globalInformer) AddCluster(clusterID ClusterID, client kube.Client, defaultResync time.Duration, tweakListOptions dynamicinformer.TweakListOptionsFunc, watchUnits []ResourceUnitWithHandlers) {
	c := &cluster{
		id:               clusterID,
		client:           client,
		defaultResync:    defaultResync,
		tweakListOptions: tweakListOptions,
		informerMap:      make(map[ResourceUnit]informerEntity),
	}
	gi.clusterMap[clusterID] = c

	// add resources
	for _, unit := range watchUnits {
		c.informerMap[unit.ResourceUnit] = informerEntity{
			resourceEventHandlers: unit.ResourceEventHandlers,
			watchErrorHandlers:    unit.WatchErrorHandlers,
		}
	}
	c.startInformers(gi)
}

func (gi *globalInformer) RemoveCluster(clusterID ClusterID) {
	if ci, ok := gi.clusterMap[clusterID]; ok {
		for _, entity := range ci.informerMap {
			entity.stopCh <- struct{}{}
		}
	}
	delete(gi.clusterMap, clusterID)
}

func (gi *globalInformer) HasSynced(clusterID ClusterID, resource schema.GroupVersionResource) bool {
	c := gi.clusterMap[clusterID]
	if c == nil {
		return false
	}

	for _, entity := range c.informerMap {
		if !entity.informer.ForResource(resource).Informer().HasSynced() {
			return false
		}
	}
	return true
}

func (gi *globalInformer) ClusterHasSynced(clusterID ClusterID) bool {
	c := gi.clusterMap[clusterID]
	if c == nil {
		return false
	}

	for unit, entity := range c.informerMap {
		if !entity.informer.ForResource(unit.Resource).Informer().HasSynced() {
			return false
		}
	}
	return true
}

func (gi *globalInformer) ClusterSyncMap(clusterID ClusterID) map[ResourceUnit]bool {
	c := gi.clusterMap[clusterID]
	if c == nil {
		return nil
	}

	res := make(map[ResourceUnit]bool, len(c.informerMap))
	for unit, entity := range c.informerMap {
		res[unit] = entity.informer.ForResource(unit.Resource).Informer().HasSynced()
	}
	return res
}

func (gi *globalInformer) AllSynced() bool {
	for _, c := range gi.clusterMap {
		if !gi.ClusterHasSynced(c.id) {
			return false
		}
	}
	return true
}
