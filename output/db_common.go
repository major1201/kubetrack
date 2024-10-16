package output

import (
	"time"

	"github.com/major1201/kubetrack/gormutils"
	"gorm.io/datatypes"
)

type Events struct {
	gormutils.Model

	EventTime time.Time `json:"event_time"`
	Source    string    `json:"source"`
	EventType string    `json:"event_type"`

	APIVersion string `json:"api_version" gorm:"type:varchar(64);index"`
	Kind       string `json:"kind" gorm:"type:varchar(64);index"`
	Namespace  string `json:"namespace" gorm:"type:varchar(64);index"`
	Name       string `json:"name" gorm:"type:varchar(64);index"`
	UID        string `json:"uid" gorm:"type:varchar(64);index"`

	Fields  datatypes.JSON `json:"fields"`
	Message string         `json:"message" gorm:"type:text"`

	Object    datatypes.JSON `json:"object"`
	Diff      string         `json:"diff" gorm:"type:text"`
	JsonPatch datatypes.JSON `json:"json_patch"`
}
