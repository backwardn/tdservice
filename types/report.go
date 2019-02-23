package types

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type Detection struct {
	Description      string   `json:"detection"`
	PID              int      `json:"pid"`
	TID              int      `json:"tid"`
	ProcessName      string   `json:"process_name"`
	ProcessImagePath string   `json:"process_image_path"`
	ProcessCmdLine   string   `json:"process_cmd_line"`
	Timestamp        int64    `json:"timestamp"` // seconds since epoch
	Severity         int      `json:"severity"`
	ProfileName      string   `json:"profile_name"`
	CVEIDs           []string `json:"cve_ids"`
	ThreatClass      string   `json:"threat_class"`
}

type Report struct {
	ID            string         `json:"id"  gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt     time.Time      `json:"-"`
	DeletedAt     *time.Time     `json:"-"`
	HostID        string         `json:"host_id" gorm:"type:uuid;"`
	Host          Host           `json:"-" gorm:"association_autoupdate:false;association_autocreate:false"`
	Detection     Detection      `json:"detection"`
	DetectionJSON postgres.Jsonb `json:"-"`
	Error         struct {
		Description string `json:"description"`
	} `json:"error,omitempty"`
}
