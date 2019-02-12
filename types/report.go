package types

type Report struct {
	ID        string `json:"id"  gorm:"primary_key;type:uuid;"`
	HostID    string `json:"host_id"`
	Host      Host
	Detection struct {
		Description      string `json:"detection"`
		PID              int    `json:"pid"`
		TID              int    `json:"tid"`
		ProcessName      string `json:"process_name"`
		ProcessImagePath string `json:"process_image_path"`
		ProcessCmdLine   string `json:"process_cmd_line"`
	} `json:"detection"`
	Error struct {
		Description string `json:"description"`
	} `json:"error,omitempty"`
}
