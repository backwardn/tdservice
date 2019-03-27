package types

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"errors"
)

type CVE struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type Detection struct {
	PID                int     `json:"pid"`
	TID                int     `json:"tid"`
	ProcessName        string  `json:"process_name"`
	ProcessPath        string  `json:"process_path"`
	Message            string  `json:"message"`
	Timestamp          int64   `json:"timestamp"` // seconds since epoch
	Severity           float32 `json:"severity"`
	ProfileDescription string  `json:"profile_description"`
	ProfileName        string  `json:"profile_name"`
	ProfileAuthor      string  `json:"profile_author"`
	ProfileDate        string  `json:"profile_date"`
	PluginOrigin       string  `json:"plugin_origin"`
	LastNDetections    int     `json:"last_n_detections"`
	AverageSeverity    float32 `json:"avg_severity_of_last_n_detections"`
	CVEIDs             []CVE   `json:"cve_ids"`
}

type Report struct {
	ID            string         `json:"id" gorm:"primary_key;type:uuid"`
	CreatedAt     time.Time      `sql:"type:timestamp"`
	DeletedAt     time.Time      `sql:"type:timestamp"`
	HostID        string         `json:"host_id" gorm:"type:uuid;"`
	Host          Host           `json:"-" gorm:"association_autoupdate:false;association_autocreate:false"`
	Detection     Detection      `json:"detection"`
	DetectionJSON postgres.Jsonb `json:"-"`
}

func (r *Report) BeforeCreate(scope *gorm.Scope) error {

	id, err := uuid.NewV4()
	if err != nil {
		return errors.New("unable to create uuid")
	}
	if err = scope.SetColumn("id", id.String()); err != nil {
		return err
	}
	return nil
}

func (r *Report) BeforeSave() (err error) {
	r.DetectionJSON.RawMessage, err = json.Marshal(&r.Detection)
	return
}

func (r *Report) AfterFind() (err error) {
	err = json.Unmarshal(r.DetectionJSON.RawMessage, &r.Detection)
	return
}
