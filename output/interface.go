package output

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

type Output interface {
	Name() string
	Write(out OutputStruct) error
}

type SourceType string

const (
	SourceTypeGeneral   SourceType = "general"
	SourceTypeEvent     SourceType = "event"
	SourceTypeKubetrack SourceType = "kubetrack"
)

type EventType string

const (
	EventTypeAdd    EventType = "add"
	EventTypeUpdate EventType = "update"
	EventTypeDelete EventType = "delete"
)

type OutputStruct struct {
	EventTime time.Time `json:"event_time"`

	ObjectRef corev1.ObjectReference

	EventType EventType      `json:"event_type"`
	Source    SourceType     `json:"source"`
	Object    map[string]any `json:"object"`
	Diff      string         `json:"diff"`
	JsonPatch string         `json:"json_patch"`
	Fields    map[string]any `json:"fields"`
	Message   string         `json:"message"` // event message
}
