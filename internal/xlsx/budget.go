package xlsx

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Budget struct {
	workbook *excelize.File
	path     string
	sheet    string
}

func NewBudget(path string, sheet string) (*Budget, error) {
	workbook, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	return &Budget{
		workbook: workbook,
		path:     path,
		sheet:    sheet,
	}, nil
}

func (b *Budget) WriteCellFloat(category string, date string, value int64) error {
	address, err := b.dateCellAddressByCategory(category, date)
	if err != nil {
		return err
	}

	err = b.workbook.SetCellFloat(b.sheet, address, formatDanishAmount(value), 2, 64)
	if err != nil {
		return err
	}

	return nil
}

func (b *Budget) WriteCellFloatWithComment(category string, date string, value int64, comment string) error {
	address, err := b.dateCellAddressByCategory(category, date)
	if err != nil {
		return err
	}

	err = b.workbook.SetCellFloat(b.sheet, address, formatDanishAmount(value), 2, 64)
	if err != nil {
		return err
	}

	err = b.workbook.AddComment(b.sheet, excelize.Comment{
		Author: "Budget",
		Cell:   address,
		Text:   comment,
		Height: 100,
		Width:  400,
	})

	if err != nil {
		return err
	}

	return nil
}

func (b *Budget) Save() error  { return b.workbook.Save() }
func (b *Budget) Close() error { return b.workbook.Close() }

func (b *Budget) GetCategories() ([]string, error) {
	rows, err := b.workbook.GetRows(b.sheet)
	if err != nil {
		return nil, err
	}

	var categories []string
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) == 0 {
			continue
		}
		category := strings.TrimSpace(rows[i][0])
		if category == "" {
			continue
		}

		cell := fmt.Sprintf("A%d", i+1)
		styleID, err := b.workbook.GetCellStyle(b.sheet, cell)
		if err != nil {
			return nil, err
		}

		style, err := b.workbook.GetStyle(styleID)
		if err != nil {
			return nil, err
		}

		// RULES:
		// Skip section headers and totals — they're bold or italic.
		// Real categories are plain text (no bold, no italic).
		if style.Font != nil && (style.Font.Bold || style.Font.Italic) {
			continue
		}

		categories = append(categories, category)
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
