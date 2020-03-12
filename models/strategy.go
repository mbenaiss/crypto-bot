package models

type Strategy struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Pair        string `json:"pair"`
	Steps       []Step `json:"step"`
}

type Step struct {
	Repartition string `json:"repartition"`
	Limit       string `json:"limit"`
	Type        string `json:"type"`
}
