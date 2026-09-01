package model

type expenses struct {
	Fixed    map[string][]string `json:"fixed"`
	Variable map[string][]string `json:"variable"`
}
type Mapper struct {
	Expenses map[string][]expenses `json:"expenses"`
	Income   map[string][]string   `json:"income"`
	Ignore   []string              `json:"ignore"`
	Fallback string                `json:"fallback"`
}
