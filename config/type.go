package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubeTrackConfiguration struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// +optional
	Rules []Rule `json:"rules,omitempty" protobuf:"bytes,2,opt,name=rules"`

	// +optional
	Output []Output `json:"output"`
}

type Rule struct {
	ObjectSelector

	RecordEvents bool `json:"recordEvents,omitempty"`

	// the fields you cares about
	CareFields []Field `json:"careFields,omitempty"`

	OnCreate EventAction `json:"onCreate,omitempty"`

	OnDelete EventAction `json:"onDelete,omitempty"`

	OnUpdate EventAction `json:"onUpdate,omitempty"`
}

type ObjectSelector struct {
	metav1.TypeMeta

	Namespaces []string `json:"namespaces,omitempty"`

	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

type EventAction struct {
	// save the full resource object when the event happens
	SaveFullObject bool `json:"saveFullObject,omitempty"`

	// save the google cmp string or not
	SaveCmp bool `json:"saveCmp,omitempty"`

	// save the json patch result of the diff or not
	SaveJsonPatch bool `json:"saveJsonPatch,omitempty"`
}

type FieldType string

const (
	FieldTypeJsonPath   FieldType = "jsonpath"
	FieldTypeGoTemplate FieldType = "gotemplate"
	FieldTypeBuiltIn    FieldType = "builtin"
)

type Field struct {
	Name string    `json:"name"`
	Type FieldType `json:"type,omitempty"`
	Expr string    `json:"expr,omitempty"`
}

type Output struct {
	Log      *OutputLog
	Mysql    *OutputMysql
	Postgres *OutputPostgres
}

type OutputLog struct{}

type OutputMysql struct {
	DSN string `json:"dsn"`
}

type OutputPostgres struct {
	DSN string `json:"dsn"`
}
