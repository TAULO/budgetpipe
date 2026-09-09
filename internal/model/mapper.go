package model

import (
	"budgetpipe/tables"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Expenses struct {
	Fixed      map[string][]string `json:"fixed"`
	Variable   map[string][]string `json:"variable"`
	Investment map[string][]string `json:"investment"`
	Fallback   string              `json:"fallback"`
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

// Section is one budget table's worth of mapping: the table it feeds, the
// budget category -> bank category mapping, and the category that absorbs
// transactions nothing else claims.
//
// Categories is the live map from the Mapper, so writing to it updates the
// mapper itself.
type Section struct {
	Table      string
	Categories map[string][]string
	Fallback   string
}

func NewMapper() Mapper {
	m := Mapper{}
	m.Sections() // initializes the category maps
	return m
}

// Sections lists every mapping section in the order it appears in the budget.
// Adding a table to the budget means adding it here and nowhere else.
func (m *Mapper) Sections() []Section {
	if m.Income.Categories == nil {
		m.Income.Categories = map[string][]string{}
	}
	if m.Expenses.Fixed == nil {
		m.Expenses.Fixed = map[string][]string{}
	}
	if m.Expenses.Variable == nil {
		m.Expenses.Variable = map[string][]string{}
	}
	if m.Expenses.Investment == nil {
		m.Expenses.Investment = map[string][]string{}
	}
	if m.Ignore == nil {
		m.Ignore = []string{}
	}

	return []Section{
		{Table: tables.Income, Categories: m.Income.Categories, Fallback: m.Income.Fallback},
		{Table: tables.Fixed, Categories: m.Expenses.Fixed, Fallback: m.Expenses.Fallback},
		{Table: tables.Variable, Categories: m.Expenses.Variable, Fallback: m.Expenses.Fallback},
		{Table: tables.Investment, Categories: m.Expenses.Investment, Fallback: m.Expenses.Fallback},
	}
}

// SyncCategories adds an empty mapping for every budget category that the
// mapper does not know yet. Existing mappings are left untouched.
func (m *Mapper) SyncCategories(categoriesByTable map[string][]string) {
	for _, section := range m.Sections() {
		for _, category := range categoriesByTable[section.Table] {
			if _, exists := section.Categories[category]; !exists {
				section.Categories[category] = []string{}
			}
		}
	}
}

// Ignored reports whether a bank category should be skipped entirely.
func (m *Mapper) Ignored(bankCategory string) bool {
	for _, ignored := range m.Ignore {
		if normalize(ignored) == normalize(bankCategory) {
			return true
		}
	}
	return false
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

// normalize makes bank and budget category names comparable regardless of the
// whitespace and casing the bank exports them with.
func normalize(category string) string {
	return strings.ToLower(strings.TrimSpace(category))
}
