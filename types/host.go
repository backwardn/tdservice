package types

import "time"

type HostInfo struct {
	Hostname string `json:"hostname"`
	Version  string `json:"version"`
	Build    string `json:"build"`
	OS       string `json:"os"`
}

type Host struct {
	ID        string     `json:"id" gorm:"primary_key;type:uuid;"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	// embed
	HostInfo
	Status string `json:"status"`
}
