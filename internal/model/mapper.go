package model

type Expenses struct {
	Fixed    map[string][]string `json:"fixed"`
	Variable map[string][]string `json:"variable"`
	Fallback string              `json:"fallback"`
}

type Income struct {
	Categories map[string][]string `json:"categories"`
	Fallback   string              `json:"fallback"`
}

type Mapper struct {
	Expenses Expenses `json:"expenses"`
	Income   Income   `json:"income"`
	Ignore   []string `json:"ignore"`
}

func NewMapper() Mapper {
	return Mapper{
		Expenses: Expenses{Fixed: map[string][]string{}, Variable: map[string][]string{}},
		Income:   Income{Categories: map[string][]string{}},
		Ignore:   []string{},
	}
}
