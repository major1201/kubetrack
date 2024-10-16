package cache

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/major1201/kubetrack/utils/slicex"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type resourceBuilder[T runtime.Object] struct {
	gi  *globalInformer
	err error

	gvr           schema.GroupVersionResource
	gvk           schema.GroupVersionKind
	model         runtime.Object
	isNamespaced  bool
	objType       reflect.Type
	objRegistered bool

	clusterSet map[ClusterID]clusterSetState
}

type clusterSetState struct{}

var emptyStruct clusterSetState

func NewResourceBuilder[T runtime.Object](gi GlobalInformer) ResourceBuilder[T] {
	giInst := gi.(*globalInformer)

	res := &resourceBuilder[T]{
		gi:         giInst,
		clusterSet: make(map[ClusterID]clusterSetState),
	}

	typ := reflect.TypeOf(*new(T)).Elem() // non-ptr type
	res.ForModel(reflect.New(typ).Interface().(T))
	return res
}

func (r *resourceBuilder[T]) ForResource(gvr schema.GroupVersionResource) ResourceBuilder[T] {
	r.gvr = gvr
	r.gvk = schema.GroupVersionKind{}
	r.model = nil
	return r
}

func (r *resourceBuilder[T]) ForKind(gvk schema.GroupVersionKind) ResourceBuilder[T] {
	r.gvk = gvk
	r.gvr = schema.GroupVersionResource{}
	r.model = nil
	return r
}

func (r *resourceBuilder[T]) ForModel(model runtime.Object) ResourceBuilder[T] {
	r.model = model
	r.gvk = schema.GroupVersionKind{}
	r.gvr = schema.GroupVersionResource{}
	return r
}

func (r *resourceBuilder[T]) Clusters(clusters ...ClusterID) ResourceBuilder[T] {
	for _, cluster := range clusters {
		r.clusterSet[cluster] = emptyStruct
	}
	return r
}

func (r *resourceBuilder[T]) AllClusters() ResourceBuilder[T] {
	for id := range r.gi.clusterMap {
		r.clusterSet[id] = emptyStruct
	}
	return r
}

func (r *resourceBuilder[T]) List(opts ...ListOption[T]) (list []T, err error) {
	r.ensureKindResource()
	if r.err != nil {
		err = r.err
		return
	}

	listConfig := &ListConfig[T]{}
	for _, opt := range opts {
		opt(listConfig)
	}

	// raw list
	var unstList []*unstructured.Unstructured
	for clusterID := range r.clusterSet {
		for unit, entity := range r.gi.clusterMap[clusterID].informerMap {
			if unit.Resource != r.gvr {
				continue
			}
			if unit.Namespace != "" && listConfig.namespaces.Len() > 0 && !listConfig.namespaces.Contains(unit.Namespace) {
				continue
			}
			informer := entity.informer.ForResource(r.gvr).Informer()
			if !informer.HasSynced() {
				err = errors.Errorf("clusterID: %s has not synced yet", clusterID)
				return
			}
			for _, unstIf := range informer.GetStore().List() {
				unstList = append(unstList, unstIf.(*unstructured.Unstructured))
			}
		}
	}

	// set total count
	size := len(unstList)
	if listConfig.totalCount != nil {
		*(listConfig.totalCount) = size
	}

	// prefilter
	for _, fn := range listConfig.preFiltersFns {
		unstList = slicex.FilterInplace(unstList, fn)
	}
	size = len(unstList)

	// pagination
	// 1. sorting
	if listConfig.preSortingLessFunc != nil {
		sort.Slice(unstList, func(i, j int) bool {
			return listConfig.preSortingLessFunc(unstList[i], unstList[j])
		})
	}
	// 2. limiting
	if listConfig.limit > 0 {
		if listConfig.offset >= size {
			unstList = nil
		} else {
			ubound := min(listConfig.offset+listConfig.limit, size)
			unstList = unstList[listConfig.offset:ubound:ubound]
		}
	}
	size = len(unstList)

	list = make([]T, len(unstList))
	if reflect.TypeOf((*unstructured.Unstructured)(nil)) == reflect.TypeOf(*new(T)) {
		list = make([]T, len(unstList))
		for i, unst := range unstList {
			var a any = unst
			list[i] = a.(T)
		}
	} else {
		for i, unst := range unstList {
			var out runtime.Object

			if r.objRegistered {
				out, err = r.gi.scheme.ConvertToVersion(unst, r.gi.groupVersioner)
				if err != nil {
					err = errors.Wrap(err, "convert to version failed")
					return
				}
			} else {
				out = unst
			}

			list[i] = out.(T)
		}
	}

	// filter with function
	for _, fn := range listConfig.filtersFns {
		list = slicex.FilterInplace(list, fn)
	}
	return
}

func (r *resourceBuilder[T]) Get(namespace, name string) (obj T, err error) {
	r.ensureKindResource()
	if r.err != nil {
		err = r.err
		return
	}

	// raw get
	clusterID := r.getFirstClusterID()
	unit := ResourceUnit{Namespace: namespace, Resource: r.gvr}
	entity, ok := r.gi.clusterMap[clusterID].informerMap[unit]
	if !ok {
		unit.Namespace = "" // try all namespaces
		entity, ok = r.gi.clusterMap[clusterID].informerMap[unit]
		if !ok {
			err = errors.Errorf("clusterID: %s for resource %s in namespace %s not watched", clusterID, r.gvr.String(), namespace)
			return
		}
	}
	informer := entity.informer.ForResource(r.gvr).Informer()
	if !informer.HasSynced() {
		err = errors.Errorf("clusterID: %s has not synced yet", clusterID)
		return
	}

	item, exist, e := informer.GetStore().GetByKey(r.getKeyName(namespace, name))
	if e != nil {
		err = errors.Wrap(e, "get store object failed: ")
		return
	}
	if !exist {
		return
	}

	var out runtime.Object
	if r.objRegistered {
		out, err = r.gi.scheme.ConvertToVersion(item.(runtime.Object), r.gi.groupVersioner)
		if err != nil {
			err = errors.Wrap(err, "convert to version failed")
			return
		}
	} else {
		out = item.(runtime.Object)
	}

	obj = out.(T)
	return
}

func (r *resourceBuilder[T]) getFirstClusterID() (clusterID ClusterID) {
	for item := range r.clusterSet {
		return item
	}
	return
}

func (r *resourceBuilder[T]) ensureKindResource() {
	if r.err != nil {
		return
	}

	if len(r.clusterSet) == 0 {
		r.err = errors.New("no cluster selected")
		return
	}

	r.ensureKindResourceByClusterID(r.getFirstClusterID())
}

func (r *resourceBuilder[T]) ensureKindResourceByClusterID(clusterID ClusterID) {
	if r.err != nil {
		return
	}

	cluster := r.gi.clusterMap[clusterID]
	if cluster == nil {
		r.err = errors.Errorf("cluster not found: id=%s", clusterID)
		return
	}
	kubeClient := cluster.client
	var mapping *meta.RESTMapping

	switch {
	case !r.gvk.Empty() && !r.gvr.Empty():
		break
	case !r.gvk.Empty() && r.gvr.Empty():
		if mapping, r.err = kubeClient.KindToMapping(r.gvk); r.err == nil {
			r.gvr = mapping.Resource
		}
	case r.gvk.Empty() && !r.gvr.Empty():
		if mapping, r.err = kubeClient.ResourceToMapping(r.gvr); r.err == nil {
			r.gvk = mapping.GroupVersionKind
		}
	case r.model != nil:
		if mapping, r.err = kubeClient.ModelToMapping(r.model); r.err == nil {
			r.gvk = mapping.GroupVersionKind
			r.gvr = mapping.Resource
		}
	default:
		r.err = errors.New("resource type not set")
	}

	if r.err != nil {
		return
	}

	if mapping != nil {
		r.isNamespaced = mapping.Scope.Name() == meta.RESTScopeNameNamespace
	}

	// get reflect type
	r.objType, r.objRegistered = r.gi.scheme.AllKnownTypes()[r.gvk]
	if !r.objRegistered {
		r.objType = reflect.TypeOf(unstructured.Unstructured{})
	}
	r.objType = reflect.PtrTo(r.objType)
}

func (r *resourceBuilder[T]) getKeyName(namespace, name string) string {
	if r.isNamespaced {
		return fmt.Sprintf("%s/%s", namespace, name)
	}
	return name
}
