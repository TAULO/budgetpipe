package model

import "budgetpipe/internal/xlsx"

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

func (m *Mapper) AddCategories(categories xlsx.TableCategories) {
	for _, category := range categories.Income {
		m.Income.Categories[category] = []string{}
	}

	for _, category := range categories.Fixed {
		m.Expenses.Fixed[category] = []string{}
	}

	for _, category := range categories.Variable {
		m.Expenses.Variable[category] = []string{}
	}
}
