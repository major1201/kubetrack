package cache

import (
	"github.com/major1201/kubetrack/utils"
	"github.com/major1201/kubetrack/utils/setx"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type (
	PreFilterFunction                func(obj *unstructured.Unstructured) bool
	FilterFunction[T runtime.Object] func(obj T) bool
)

type ListConfig[T runtime.Object] struct {
	// 预过滤
	preFiltersFns []PreFilterFunction

	// 过滤结果
	filtersFns []FilterFunction[T]

	namespaces setx.Set[string]

	// pagination
	preSortingLessFunc func(a, b *unstructured.Unstructured) bool
	offset             int
	limit              int
	totalCount         *int
}

type ListOption[T runtime.Object] func(config *ListConfig[T])

func InNamespaces[T runtime.Object](namespaces ...string) ListOption[T] {
	return func(config *ListConfig[T]) {
		config.namespaces = setx.NewHashSetFromSlice(namespaces)
		config.preFiltersFns = append(config.preFiltersFns, func(obj *unstructured.Unstructured) bool {
			return utils.ContainsString(namespaces, obj.GetNamespace())
		})
	}
}

func WithLabelSelector[T runtime.Object](selector labels.Selector) ListOption[T] {
	return func(config *ListConfig[T]) {
		config.preFiltersFns = append(config.preFiltersFns, func(obj *unstructured.Unstructured) bool {
			return selector.Matches(labels.Set(obj.GetLabels()))
		})
	}
}

func WithLabelsMap[T runtime.Object](m map[string]string) ListOption[T] {
	selector := labels.Set(m).AsSelector()
	return WithLabelSelector[T](selector)
}

func WithLabelSelectorStr[T runtime.Object](selector string) (opt ListOption[T], err error) {
	m, err := labels.ConvertSelectorToLabelsMap(selector)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	opt = WithLabelsMap[T](m)
	return
}

func WithMetaLabelSelector[T runtime.Object](selector *metav1.LabelSelector) (opt ListOption[T], err error) {
	s, err := metav1.LabelSelectorAsSelector(selector)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	opt = WithLabelSelector[T](s)
	return
}

func WithOwnerReference[T runtime.Object](gvk schema.GroupVersionKind, name string) ListOption[T] {
	return func(config *ListConfig[T]) {
		config.preFiltersFns = append(config.preFiltersFns, func(obj *unstructured.Unstructured) bool {
			for _, ref := range obj.GetOwnerReferences() {
				if ref.APIVersion == gvk.GroupVersion().String() && ref.Kind == gvk.Kind && ref.Name == name {
					return true
				}
			}
			return false
		})
	}
}

func FromMetaListOptions[T runtime.Object](opts metav1.ListOptions) (res []ListOption[T]) {
	if opts.LabelSelector != "" {
		m, err := labels.ConvertSelectorToLabelsMap(opts.LabelSelector)
		if err != nil {
			// should filter all resources
			m = map[string]string{"must-not-exist": "must-not-exist-filter-all"}
		}
		res = append(res, WithLabelsMap[T](m))
	}
	return
}

func WithCustomFilter[T runtime.Object](fn FilterFunction[T]) ListOption[T] {
	return func(config *ListConfig[T]) {
		config.filtersFns = append(config.filtersFns, fn)
	}
}

func WithCustomPreFilter[T runtime.Object](fn PreFilterFunction) ListOption[T] {
	return func(config *ListConfig[T]) {
		config.preFiltersFns = append(config.preFiltersFns, fn)
	}
}

func WithNonTerminatedPods[T runtime.Object](config *ListConfig[T]) {
	config.preFiltersFns = append(config.preFiltersFns, func(obj *unstructured.Unstructured) bool {
		return obj.GetDeletionTimestamp() == nil
	})
}

func WithNames[T runtime.Object](names ...string) ListOption[T] {
	return func(config *ListConfig[T]) {
		config.preFiltersFns = append(config.preFiltersFns, func(obj *unstructured.Unstructured) bool {
			return utils.ContainsString(names, obj.GetName())
		})
	}
}

func WithPagination[T runtime.Object](less func(a, b *unstructured.Unstructured) bool, offset, limit int, totalCount *int) ListOption[T] {
	return func(config *ListConfig[T]) {
		config.preSortingLessFunc = less
		config.offset = offset
		config.limit = limit
		config.totalCount = totalCount
	}
}
