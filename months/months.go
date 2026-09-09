package months

import (
	"fmt"
	"slices"
	"strings"
)

// Month is a calendar month in the two forms this tool needs: the key used for
// data file names (data/aug.csv) and the header used in the budget tables (AUG).
type Month struct {
	key string
}

// keys are in calendar order and double as the data file names.
var keys = []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}

// names holds the extra spellings accepted on the command line, Danish first.
var names = map[string][]string{
	"jan": {"januar", "january"},
	"feb": {"februar", "february"},
	"mar": {"marts", "march"},
	"apr": {"april"},
	"may": {"maj", "may"},
	"jun": {"juni", "june"},
	"jul": {"juli", "july"},
	"aug": {"august"},
	"sep": {"september"},
	"oct": {"oktober", "october", "okt"},
	"nov": {"november"},
	"dec": {"december"},
}

func All() []Month {
	all := make([]Month, 0, len(keys))
	for _, key := range keys {
		all = append(all, Month{key: key})
	}
	return all
}

// Parse resolves a user-supplied month - Danish or English, full or abbreviated -
// to a canonical month.
func Parse(value string) (Month, error) {
	value = strings.ToLower(strings.TrimSpace(value))

	if slices.Contains(keys, value) {
		return Month{key: value}, nil
	}

	for key, aliases := range names {
		if slices.Contains(aliases, value) {
			return Month{key: key}, nil
		}
	}

	return Month{}, fmt.Errorf("unknown month %q, expected one of %s", value, strings.Join(keys, ", "))
}

// Key is the month as used in data file names, e.g. "aug".
func (m Month) Key() string { return m.key }

// Header is the month as it appears in the budget table headers, e.g. "AUG".
func (m Month) Header() string { return strings.ToUpper(m.key) }

func (m Month) String() string { return m.key }
