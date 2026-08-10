package xlsx

import (
	"fmt"
	"regexp"

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

func (b *Budget) WriteCellFloat(address string, value float64) {
	b.workbook.SetCellFloat(b.sheet, address, value, 2, 64)
}

func (b *Budget) Save() error  { return b.workbook.Save() }
func (b *Budget) Close() error { return b.workbook.Close() }

func (b *Budget) DateCellAddressByCategory(category string, date string) (string, error) {
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
	pattern := "(?i)" + regexp.QuoteMeta(value)
	cells, err := b.workbook.SearchSheet(b.sheet, pattern, true)
	if err != nil {
		return "", err
	}

	if len(cells) == 0 {
		return "", fmt.Errorf("could not find %q", value)
	}

	return cells[0], nil
}
