package config

import (
	"os"

	"github.com/major1201/kubetrack/log"
	"github.com/major1201/kubetrack/utils/goutils"
	"github.com/major1201/kubetrack/utils/slicex"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

func (osel ObjectSelector) Match(obj runtime.Object) bool {
	tobj, err := meta.TypeAccessor(obj)
	if err != nil {
		log.L.Error(err, "convert to type meta failed")
		return false
	}

	oobj, err := meta.Accessor(obj)
	if err != nil {
		log.L.Error(err, "convert to object meta failed")
		return false
	}

	// check type meta
	if osel.APIVersion != tobj.GetAPIVersion() || osel.Kind != tobj.GetKind() {
		return false
	}

	// check namespaces
	namespace := oobj.GetNamespace()
	if len(osel.Namespaces) > 0 {
		if !slicex.Contains(osel.Namespaces, namespace) {
			return false
		}
	}

	// check excluded namespaces
	for _, nsWildcard := range osel.ExcludedNamespaces {
		if goutils.WildcardMatchSimple(nsWildcard, namespace) {
			return false
		}
	}

	// check label selector
	if osel.Selector == nil {
		return true
	}
	selector, err := metav1.LabelSelectorAsSelector(osel.Selector)
	if err != nil {
		log.L.Error(err, "LabelSelectorAsSelector failed")
		return false
	}
	if selector.Empty() {
		return true
	}
	return selector.Matches(labels.Set(oobj.GetLabels()))
}

func (er EventRule) Match(unstr *unstructured.Unstructured) bool {
	if unstr == nil {
		return false
	}

	// check namespaces
	namespace := unstr.GetNamespace()
	if len(er.Namespaces) > 0 {
		if !slicex.Contains(er.Namespaces, namespace) {
			return false
		}
	}

	// check excluded namespaces
	for _, nsWildcard := range er.ExcludedNamespaces {
		if goutils.WildcardMatchSimple(nsWildcard, namespace) {
			return false
		}
	}
	return true
}

// todo reloader
func LoadFromFile(path string) (KubeTrackConfiguration, error) {
	config := KubeTrackConfiguration{}

	// read all
	yamlByte, err := os.ReadFile(path)
	if err != nil {
		return config, errors.Wrap(err, "read config file error")
	}

	if err := yaml.Unmarshal(yamlByte, &config); err != nil {
		return config, errors.Wrap(err, "unmarshal config file error")
	}
	return config, nil
}
