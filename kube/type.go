package kube

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ResourceToMapping gvr to mapping
func (c *ClientImpl) ResourceToMapping(gvr schema.GroupVersionResource) (mapping *meta.RESTMapping, err error) {
	gvk, err := c.restMapper.KindFor(gvr)
	if err != nil {
		err = errors.Wrapf(err, "gvk not found for: %s", gvr.String())
		return
	}

	return c.KindToMapping(gvk)
}

// KindToMapping gvk to mapping
func (c *ClientImpl) KindToMapping(gvk schema.GroupVersionKind) (mapping *meta.RESTMapping, err error) {
	mapping, err = c.restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		err = errors.Wrapf(err, "gvk mapping failed: %s", gvk.String())
		return
	}
	return
}

// ModelToMapping model to mapping
func (c *ClientImpl) ModelToMapping(model runtime.Object) (mapping *meta.RESTMapping, err error) {
	var gvks []schema.GroupVersionKind
	gvks, _, err = GetScheme().ObjectKinds(model)
	if err != nil {
		err = errors.Wrap(err, "get model kinds error")
		return
	}
	if len(gvks) == 0 {
		err = errors.New("no gvk found from model")
		return
	}

	return c.KindToMapping(gvks[0])
}
