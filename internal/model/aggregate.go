package model

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// Entry is one resolved budget cell: the amount to write for a category in a
// table, and a comment explaining the amount when it needs one.
type Entry struct {
	Table    string
	Category string
	Amount   int64
	Comment  string
}

// Result is a month of transactions mapped onto the budget.
type Result struct {
	// Entries covers every category in the mapper, in budget order.
	Entries []Entry
	// Ignored transactions matched the mapper's ignore list.
	Ignored []Transaction
	// Unmapped transactions matched nothing and had no fallback to absorb them.
	Unmapped []Transaction
}

// target is a single cell row in the budget: a table plus a category in it.
type target struct {
	table    string
	category string
}

// rule is where one bank category lands in the budget.
type rule struct {
	target target
	negate bool
}

// Aggregate maps transactions onto budget cells.
//
// The mapper is inverted into a bank category -> budget cell lookup once, so
// each transaction is placed with a single map hit instead of a scan over every
// section. Every category in the mapper gets an entry - a zero one if nothing
// matched - so importing the same month twice overwrites values left behind by
// the previous run instead of leaving stale numbers in the sheet.
func (m *Mapper) Aggregate(transactions []Transaction) (Result, error) {
	rules, err := m.rules()
	if err != nil {
		return Result{}, err
	}

	incomeFallback, err := m.fallback(m.Income.Fallback)
	if err != nil {
		return Result{}, fmt.Errorf("income fallback: %w", err)
	}

	expenseFallback, err := m.fallback(m.Expenses.Fallback)
	if err != nil {
		return Result{}, fmt.Errorf("expenses fallback: %w", err)
	}

	totals := make(map[target]int64)
	for _, section := range m.Sections() {
		for category := range section.Categories {
			totals[target{table: section.Table, category: category}] = 0
		}
	}

	var result Result
	absorbed := make(map[target][]Transaction)

	for _, transaction := range transactions {
		bankCategory := normalize(transaction.Category)

		if m.Ignored(bankCategory) {
			result.Ignored = append(result.Ignored, transaction)
			continue
		}

		if matched, ok := rules[bankCategory]; ok {
			totals[matched.target] += matched.sign(transaction.Amount)
			continue
		}

		catchAll := incomeFallback
		if transaction.Amount < 0 {
			catchAll = expenseFallback
		}
		if catchAll == nil {
			result.Unmapped = append(result.Unmapped, transaction)
			continue
		}

		totals[catchAll.target] += catchAll.sign(transaction.Amount)
		absorbed[catchAll.target] = append(absorbed[catchAll.target], transaction)
	}

	for _, section := range m.Sections() {
		for _, category := range slices.Sorted(maps.Keys(section.Categories)) {
			cell := target{table: section.Table, category: category}
			result.Entries = append(result.Entries, Entry{
				Table:    cell.table,
				Category: cell.category,
				Amount:   totals[cell],
				Comment:  summarize(absorbed[cell]),
			})
		}
	}

	return result, nil
}

// rules inverts the mapper into a bank category -> budget cell lookup, and
// rejects a bank category that two budget categories both claim - that would
// silently count the same transactions twice.
func (m *Mapper) rules() (map[string]rule, error) {
	rules := make(map[string]rule)

	for _, section := range m.Sections() {
		for budgetCategory, bankCategories := range section.Categories {
			for _, bankCategory := range bankCategories {
				key := normalize(bankCategory)
				if existing, taken := rules[key]; taken {
					return nil, fmt.Errorf(
						"bank category %q is mapped to both %q and %q",
						bankCategory, existing.target.category, budgetCategory,
					)
				}
				rules[key] = rule{
					target: target{table: section.Table, category: budgetCategory},
					negate: section.Negate,
				}
			}
		}
	}

	return rules, nil
}

// fallback resolves a fallback category name to the cell it lives in. An empty
// name means the section has no fallback; a name the budget does not know is an
// error, because silently dropping transactions is worse than refusing to run.
func (m *Mapper) fallback(category string) (*rule, error) {
	if strings.TrimSpace(category) == "" {
		return nil, nil
	}

	for _, section := range m.Sections() {
		for budgetCategory := range section.Categories {
			if normalize(budgetCategory) == normalize(category) {
				return &rule{
					target: target{table: section.Table, category: budgetCategory},
					negate: section.Negate,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("category %q is not in the budget", category)
}

func (r rule) sign(amount int64) int64 {
	if r.negate {
		return -amount
	}
	return amount
}

// summarize renders one line per bank category with its total, for the comment
// on a fallback cell.
func summarize(transactions []Transaction) string {
	if len(transactions) == 0 {
		return ""
	}

	totals := make(map[string]int64)
	for _, transaction := range transactions {
		category := strings.TrimSpace(transaction.Category)
		if category == "" {
			category = "(no category)"
		}
		totals[category] += transaction.Amount
	}

	lines := make([]string, 0, len(totals))
	for _, category := range slices.Sorted(maps.Keys(totals)) {
		lines = append(lines, fmt.Sprintf("• %s: %.2f kr.", category, float64(totals[category])/100))
	}

	return strings.Join(lines, "\n")
}
