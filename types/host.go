package types

type Host struct {
	ID       string `json:"id"`
	OS       string `json:"os"`
	Hostname string `json:"hostname"`
	Version  string `json:"version"`
	Build    string `json:"build"`
}
