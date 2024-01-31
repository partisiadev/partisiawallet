package evm

type Explorer struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Standard string `json:"standard"`
	Icon     string `json:"icon,omitempty"`
}
