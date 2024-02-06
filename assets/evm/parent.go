package evm

type Parent struct {
	Type    string   `json:"type"`
	Chain   string   `json:"chain"`
	Bridges []Bridge `json:"bridges"`
}
