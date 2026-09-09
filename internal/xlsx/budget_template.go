package xlsx

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
)

func WriteTemplate(template []byte, sheet, destination string) error {
	wb, err := excelize.OpenReader(bytes.NewReader(template))
	if err != nil {
		return err
	}
	defer wb.Close()

	current := wb.GetSheetName(0)
	if current == "" {
		return fmt.Errorf("template has no sheets")
	}
	if current != sheet {
		if err := wb.SetSheetName(current, sheet); err != nil {
			return err
		}
	}

	return wb.SaveAs(destination)
}
