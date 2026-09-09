package xlsx

import (
	"bytes"

	"github.com/xuri/excelize/v2"
)

type BudgetTemplate struct{}

func NewBudgetTemplate(template []byte, destination string) error {
	wb, err := excelize.OpenReader(bytes.NewReader(template))
	if err != nil {
		return err
	}
	if err := wb.Close(); err != nil {
		return err
	}

	return wb.SaveAs(destination)
}
