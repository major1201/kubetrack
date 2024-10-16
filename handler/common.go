package handler

import (
	"bytes"
	"fmt"

	"github.com/major1201/kubetrack/config"
	"github.com/major1201/kubetrack/kube"
	"github.com/major1201/kubetrack/log"
	"github.com/major1201/kubetrack/tmpl"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/jsonpath"
)

func ObjectInDifferentTypes(obj any) (runtime.Object, meta.Type, metav1.Object, error) {
	objRuntime, ok := obj.(runtime.Object)
	if !ok {
		return nil, nil, nil, errors.New("convert runtime object failed")
	}
	objType, err := meta.TypeAccessor(obj)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "convert to meta.Type failed")
	}
	objMeta, err := meta.Accessor(obj)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "convert to metav1.Object failed")
	}
	return objRuntime, objType, objMeta, nil
}

func ParseJsonPath(obj any, tmpl string) (string, error) {
	j := jsonpath.New("")
	j.AllowMissingKeys(true)
	if err := j.Parse(tmpl); err != nil {
		return "", errors.Wrap(err, "parse jsonpath failed")
	}

	buf := new(bytes.Buffer)
	if err := j.Execute(buf, obj); err != nil {
		return "", errors.Wrap(err, "execute jsonpath failed")
	}
	return buf.String(), nil
}

func BuildFieldsMap(unstr *unstructured.Unstructured, fields []config.Field) map[string]any {
	if unstr == nil {
		return nil
	}

	var runtimeObj runtime.Object
	toVersioned := func() (runtime.Object, error) {
		if runtimeObj != nil {
			return runtimeObj, nil
		}
		var err error
		runtimeObj, err = kube.GetScheme().ConvertToVersion(unstr, runtime.GroupVersioner(schema.GroupVersions(kube.GetScheme().PrioritizedVersionsAllGroups())))
		return runtimeObj, err
	}

	obj := unstr.Object
	mp := make(map[string]any, len(fields))
	for _, field := range fields {
		switch field.Type {
		case config.FieldTypeJsonPath:
			res, err := ParseJsonPath(obj, fmt.Sprintf("{ %s }", field.Expr))
			if err != nil {
				log.L.Error(err, "parse jsonpath failed", "key", field.Expr)
				continue
			}
			mp[field.Name] = res

		case config.FieldTypeGoTemplate:
			if result, err := tmpl.ExecuteTextTemplate(field.Expr, unstr); err != nil {
				mp[field.Name] = fmt.Sprintf("Error rendering go-template: %s", err.Error())
			} else {
				mp[field.Name] = result
			}

		case config.FieldTypeBuiltIn:
			fn, ok := builtInFuncMap[field.Expr]
			if !ok {
				mp[field.Name] = fmt.Sprintf("Unknown func: %s", field.Expr)
			} else {
				mp[field.Name] = fn(unstr, toVersioned)
			}
		}
	}
	return mp
}

func DisplayObjectReference(objectRef corev1.ObjectReference) string {
	return fmt.Sprintf("%s/%s,%s/%s,uid=%s", objectRef.APIVersion, objectRef.Kind, objectRef.Namespace, objectRef.Name, objectRef.UID)
}

func Ternary[T any](flag bool, trueVal T, falseVal T) T {
	if flag {
		return trueVal
	}
	return falseVal
}

func BuildObjectReference(unstr *unstructured.Unstructured) corev1.ObjectReference {
	return corev1.ObjectReference{
		APIVersion:      unstr.GetAPIVersion(),
		Kind:            unstr.GetKind(),
		UID:             unstr.GetUID(),
		Namespace:       unstr.GetNamespace(),
		Name:            unstr.GetName(),
		ResourceVersion: unstr.GetResourceVersion(),
	}
}
