package types

type HostInfo struct {
	Hostname string `json:"hostname"`
	Version  string `json:"version"`
	Build    string `json:"build"`
	OS       string `json:"os"`
}

type Host struct {
	ID string `json:"id"`
	// embed
	HostInfo
	Status string `json:"status"`
}
