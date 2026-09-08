package xlsx

import (
	"bytes"
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

func NewBudget(template []byte, sheet string) (*Budget, error) {
	wb, err := excelize.OpenReader(bytes.NewReader(template))
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

func (b *Budget) Write(cell string, val string) error {
	return b.workbook.SetCellValue(b.sheet, cell, val)
}

func (b *Budget) Save() error                     { return b.workbook.Save() }
func (b *Budget) SaveAs(destination string) error { return b.workbook.SaveAs(destination) }
func (b *Budget) Close() error                    { return b.workbook.Close() }

// TableCategories returns the category labels in a named Excel table,
// in sheet order. Header row and the trailing Total row are skipped.
func (b *Budget) TableCategories(tableName string) ([]string, error) {
	tables, err := b.workbook.GetTables(b.sheet)
	if err != nil {
		return nil, err
	}

	var target *excelize.Table
	for i := range tables {
		if strings.EqualFold(tables[i].Name, tableName) {
			target = &tables[i]
			break
		}
	}
	if target == nil {
		allTables, _ := b.workbook.GetTables(b.sheet)
		for _, table := range allTables {
			fmt.Printf("Found table: %s (%s)\n", table.Name, table.Range)
		}
		return nil, fmt.Errorf(
			"table %q not found on sheet %q",
			tableName,
			b.sheet,
		)
	}

	parts := strings.Split(target.Range, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected table range %q", target.Range)
	}
	startCol, startRow, err := excelize.CellNameToCoordinates(parts[0])
	if err != nil {
		return nil, err
	}
	_, endRow, err := excelize.CellNameToCoordinates(parts[1])
	if err != nil {
		return nil, err
	}

	var categories []string
	for row := startRow + 1; row <= endRow; row++ { // +1 skips the header row
		cell, err := excelize.CoordinatesToCellName(startCol, row)
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
