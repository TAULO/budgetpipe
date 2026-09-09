package model

import (
	"encoding/json"
	"fmt"
	"os"
)

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

func (m *Mapper) SyncCategories(income, fixed, variable []string) {
	for _, category := range income {
		if _, exists := m.Income.Categories[category]; !exists {
			m.Income.Categories[category] = []string{}
		}
	}

	for _, category := range fixed {
		if _, exists := m.Expenses.Fixed[category]; !exists {
			m.Expenses.Fixed[category] = []string{}
		}
	}

	for _, category := range variable {
		if _, exists := m.Expenses.Variable[category]; !exists {
			m.Expenses.Variable[category] = []string{}
		}
	}
}

func (m *Mapper) ToOriginalToOriginal(path string) error {
	starterMapper, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("building starter mapper: %w", err)
	}

	if err := os.WriteFile(path, starterMapper, 0o644); err != nil {
		return fmt.Errorf("creating mapper: %w", err)
	}

	return nil
}
