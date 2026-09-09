package tables

const (
	Income     = "tblIncome"
	Fixed      = "tblFixed"
	Variable   = "tblVariable"
	Investment = "tblInvestment"
	Totals     = "tblTotals"
)

// All lists the tables that hold budget categories, in sheet order. Totals is
// left out: it is derived from the others by formulas in the sheet.
var All = []string{Income, Fixed, Variable, Investment}
