package xlsx

import (
	"budgetpipe/tables"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Budget struct {
	workbook *excelize.File
	path     string
	sheet    string
}

type CellInput struct {
	Category string
	Month    string
	Value    int64
	Comment  string
}

type TableBounds struct {
	startCol, startRow int
	endCol, endRow     int
}

type TableCategories struct {
	Income   []string
	Fixed    []string
	Variable []string
}

func NewBudget(path string, sheet string) (*Budget, error) {
	wb, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	if err := wb.Close(); err != nil {
		return nil, err
	}

	current := wb.GetSheetName(0)
	if current == "" {
		return nil, fmt.Errorf("template has no sheets")
	}
	if current != sheet {
		if err := wb.SetSheetName(current, sheet); err != nil {
			return nil, err
		}
	}

	return &Budget{
		workbook: wb,
		sheet:    sheet,
	}, nil
}

func (b *Budget) WriteToCell(input CellInput) error {
	address, err := b.dateCellAddressByCategory(input.Category, input.Month)
	if err != nil {
		return err
	}

	err = b.workbook.SetCellFloat(b.sheet, address, formatDanishAmount(input.Value), 2, 64)
	if err != nil {
		return err
	}

	if input.Comment != "" {
		comments, err := b.workbook.GetComments(b.sheet)
		if err != nil {
			return err
		}

		for _, comment := range comments {
			if comment.Cell == address {
				err = b.workbook.DeleteComment(b.sheet, comment.Cell)
				if err != nil {
					return err
				}
			}
		}

		err = b.workbook.AddComment(b.sheet, excelize.Comment{
			Author: "Budget",
			Cell:   address,
			Text:   input.Comment,
			Height: 100,
			Width:  400,
		})

		if err != nil {
			return err
		}
	}

	slog.Info("writing transaction", "category", input.Category, "amount", formatDanishAmount(input.Value), "month", input.Month, "cell", address)

	return nil
}

func (b *Budget) Save() error  { return b.workbook.Save() }
func (b *Budget) Close() error { return b.workbook.Close() }

func (b *Budget) TableCategories() (TableCategories, error) {
	incomeCategories, err := b.getTableCategories(tables.Income)
	if err != nil {
		return TableCategories{}, fmt.Errorf("getting table categories: %w", err)
	}

	variableCategories, err := b.getTableCategories(tables.Variable)
	if err != nil {
		return TableCategories{}, fmt.Errorf("getting table categories: %w", err)
	}

	fixedCategories, err := b.getTableCategories(tables.Fixed)
	if err != nil {
		return TableCategories{}, fmt.Errorf("getting table categories: %w", err)
	}

	return TableCategories{
		Income:   incomeCategories,
		Variable: variableCategories,
		Fixed:    fixedCategories,
	}, nil
}

func (b *Budget) WriteCellByCategoryAndMonth(tableName string, category string, month string, value int32) error {
	monthCol, err := b.getMonthCol(tableName, month)
	if err != nil {
		return err
	}

	categoryRow, err := b.getCategoryRow(tableName, category)
	if err != nil {
		return err
	}

	cell, err := excelize.JoinCellName(monthCol, categoryRow)
	if err != nil {
		return err
	}

	if err := b.workbook.SetCellValue(b.sheet, cell, value); err != nil {
		return err
	}

	return nil
}

// TableCategories returns the category labels in a named Excel table,
// in sheet order. Header row and the trailing Total row are skipped.
func (b *Budget) getTableCategories(tableName string) ([]string, error) {
	table, err := b.getTableByName(tableName)
	if err != nil {
		return nil, err
	}

	bounds, err := b.getTableBounds(table)
	if err != nil {
		return nil, err
	}

	var categories []string
	for row := bounds.startRow + 1; row <= bounds.endRow; row++ { // +1 skips the header row
		cell, err := excelize.CoordinatesToCellName(bounds.startCol, row)
		if err != nil {
			return nil, err
		}
		val, err := b.workbook.GetCellValue(b.sheet, cell)
		if err != nil {
			return nil, err
		}
		val = strings.TrimSpace(val)
		if val == "" || strings.EqualFold(val, "Total") {
			continue
		}
		categories = append(categories, val)
	}
	return categories, nil
}

func (b *Budget) getTableByName(tableName string) (*excelize.Table, error) {
	excelTables, err := b.workbook.GetTables(b.sheet)
	if err != nil {
		return nil, err
	}

	for i := range excelTables {
		if strings.EqualFold(excelTables[i].Name, tableName) {
			return &excelTables[i], nil
		}
	}

	fmt.Println("Found tables:")
	for _, table := range excelTables {
		fmt.Printf("  %s (%s)\n", table.Name, table.Range)
	}

	return nil, fmt.Errorf(
		"table %q not found on sheet %q",
		tableName,
		b.sheet,
	)
}

func (b *Budget) getTableBounds(table *excelize.Table) (*TableBounds, error) {
	parts := strings.Split(table.Range, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected table range %q", table.Range)
	}
	startCol, startRow, err := excelize.CellNameToCoordinates(parts[0])
	if err != nil {
		return nil, err
	}
	endCol, endRow, err := excelize.CellNameToCoordinates(parts[1])
	if err != nil {
		return nil, err
	}

	return &TableBounds{startCol, startRow, endCol, endRow}, nil
}

func (b *Budget) getMonthCol(tableName string, month string) (col string, err error) {
	month = strings.TrimSpace(strings.ToUpper(month))

	table, err := b.getTableByName(tableName)
	if err != nil {
		return "", err
	}

	bounds, err := b.getTableBounds(table)
	if err != nil {
		return "", err
	}

	for col := bounds.startCol; col <= bounds.endCol; col++ {
		cell, err := excelize.CoordinatesToCellName(col, bounds.startRow)
		if err != nil {
			return "", err
		}

		val, err := b.workbook.GetCellValue(b.sheet, cell)
		if err != nil {
			return "", err
		}

		val = strings.TrimSpace(strings.ToUpper(val))
		if strings.EqualFold(val, month) {
			return excelize.ColumnNumberToName(col)
		}
	}

	return "", fmt.Errorf("no cell found on sheet %q for %q", b.sheet, month)
}

func (b *Budget) getCategoryRow(tableName, category string) (int, error) {
	category = strings.TrimSpace(strings.ToLower(category))

	table, err := b.getTableByName(tableName)
	if err != nil {
		return -1, err
	}
	bounds, err := b.getTableBounds(table)
	if err != nil {
		return -1, err
	}

	for row := bounds.startRow + 1; row <= bounds.endRow; row++ {
		cell, err := excelize.CoordinatesToCellName(bounds.startCol, row)
		if err != nil {
			return -1, err
		}
		val, err := b.workbook.GetCellValue(b.sheet, cell)
		if err != nil {
			return -1, err
		}
		if strings.TrimSpace(strings.ToLower(val)) == category {
			return row, nil
		}
	}

	return -1, fmt.Errorf("category %q not found in table %q", category, tableName)
}

func (b *Budget) dateCellAddressByCategory(category string, date string) (string, error) {
	categoryCell, err := b.findCell(category)
	if err != nil {
		return "", err
	}

	dateCell, err := b.findCell(date)
	if err != nil {
		return "", err
	}

	dateCol, _, err := excelize.SplitCellName(dateCell)
	if err != nil {
		return "", err
	}

	_, catRow, err := excelize.SplitCellName(categoryCell)
	if err != nil {
		return "", err
	}
	return excelize.JoinCellName(dateCol, catRow)
}

func (b *Budget) findCell(value string) (string, error) {
	value = strings.TrimSpace(value)
	pattern := `(?i)^\s*` + regexp.QuoteMeta(value) + `\s*$`
	cells, err := b.workbook.SearchSheet(b.sheet, pattern, true)
	if err != nil {
		return "", err
	}

	if len(cells) == 0 {
		return "", fmt.Errorf("could not find %q", value)
	}

	return cells[0], nil
}

func formatDanishAmount(danishCent int64) float64 {
	return float64(danishCent) / 100
}
