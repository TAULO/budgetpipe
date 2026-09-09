# budgetpipe

A small CLI that pours your bank's monthly CSV export into an Excel budget.

Everything lives in one working directory: the bank exports under `data/`, a
`mapper.json` that says which bank category feeds which budget category, and
`budget.xlsx` itself. `budgetpipe` reads the first two and writes the numbers
into the third — one cell per budget category, per month.

```
init
├── create directory
├── copy template
├── create CSV files
└── create empty mapper

sync
├── read Excel categories
├── read mapper
├── reconcile categories
└── save mapper

add
├── read CSV transactions
├── read mapper
├── categorize transactions
└── update budget.xlsx
```

## Install

```sh
go build -o budget .
```

The budget template is embedded in the binary, so the single file is all you
need.

## Quick start

```sh
budget init -C ~/budget          # scaffold a working directory
# open ~/budget/budget.xlsx and type your category names in the first column
budget sync -C ~/budget          # pull those categories into mapper.json
# open ~/budget/mapper.json and list the bank categories for each of them
cp ~/Downloads/export.csv ~/budget/data/aug.csv
budget add -C ~/budget -m august # write August into the budget
```

`add` prints what it did and warns about every transaction it could not place:

```
2026/09/09 14:19:37 INFO writing transaction category=Løn amount=28952.24 month=AUG cell=J6
...
2026/09/09 14:19:37 WARN unmapped transaction category="Hår- og hudpleje" amount=-240
imported aug: 23 cells written, 0 ignored, 4 unmapped
```

## Commands

| Command | What it does |
| --- | --- |
| `init` | Creates the working directory: `budget.xlsx` from the template, an empty `data/<month>.csv` for all twelve months, and a `mapper.json` holding every category found in the workbook. Refuses to run if the directory is not empty. |
| `sync` | Reads the category labels out of `budget.xlsx` and adds the ones `mapper.json` does not know yet, with an empty bank category list. Existing mappings are left alone. Run it whenever you add a row to the budget. |
| `add` | Reads `data/<month>.csv`, maps every transaction through `mapper.json`, and writes one value per budget category into the month's column. |

### Flags

| Flag | Commands | Default | Meaning |
| --- | --- | --- | --- |
| `-C`, `--dir` | all | current directory | The working directory holding `budget.xlsx`, `mapper.json` and `data/`. |
| `-S`, `--sheet` | all | current year | The sheet to read and write. `init` names the new sheet this. |
| `-m`, `--month` | `add` | required | The month to import. |
| `-r`, `--reset` | `sync` | off | Throw the current `mapper.json` away and rebuild it from the workbook. |

Months are accepted in Danish or English, full or abbreviated — `august`,
`aug`, `oktober`, `okt`, `october` and `oct` all work. The canonical short form
is what names the data file (`data/aug.csv`) and what is matched against the
budget's column headers (`AUG`).

## The working directory

```
~/budget
├── budget.xlsx     the budget itself; the only file you edit by hand in Excel
├── mapper.json     bank category -> budget category
└── data
    ├── jan.csv     one bank export per month
    ├── feb.csv
    └── ...
```

## mapper.json

Each section maps **budget** categories (the keys, matching the row labels in
the workbook) to the **bank** categories that feed them (the values, matching
the category column in the CSV).

```json
{
  "expenses": {
    "fixed": {
      "Husleje": ["Husleje"],
      "Internet og Telefon": ["Telefon, internet, streaming og TV"]
    },
    "variable": {
      "Andet (Diverse)": [],
      "Dagligvarer": ["Dagligvarer"]
    },
    "investment": {
      "Indskud til Nordnet": ["Opsparing og investering (Andet)"]
    },
    "fallback": "Andet (Diverse)"
  },
  "income": {
    "categories": {
      "Løn": ["Løn, dagpenge og pension"],
      "Overførelser": ["Anden indtægt", "Anden udgift"]
    },
    "fallback": ""
  },
  "ignore": []
}
```

- **The four sections** — `income.categories`, `expenses.fixed`,
  `expenses.variable` and `expenses.investment` — correspond to the four tables
  in the workbook (`tblIncome`, `tblFixed`, `tblVariable`, `tblInvestment`).
- **`fallback`** names the budget category that absorbs transactions no mapping
  claims: negative amounts go to the expenses fallback, positive ones to the
  income fallback. The cell gets a comment listing the bank categories that
  landed there and their totals, so you can see what to map next. Leave a
  fallback empty and those transactions are reported as unmapped instead —
  nothing is written for them.
- **`ignore`** lists bank categories to drop entirely — internal transfers
  between your own accounts, typically. Ignore wins over any mapping.
- Category names are matched **case- and whitespace-insensitively**, so a stray
  space in a bank export does not break a mapping.
- A bank category may feed **only one** budget category. Listing it twice is an
  error rather than a silent double count:
  `bank category "Dagligvarer" is mapped to both "Dagligvarer" and "Andet (Diverse)"`.

## What `add` writes

- **Amounts as the bank exported them.** Expenses stay negative; nothing is
  flipped on the way in. Values are written as kroner with two decimals.
- **Every category in the mapper, including the ones that saw no
  transactions.** They get a zero, which means re-importing a month overwrites
  what the previous run left behind instead of stranding stale numbers in the
  sheet. Importing the same month twice is safe.
- **Nothing outside those cells.** Totals, averages and sparklines are formulas
  in the workbook and are left to Excel.

Close `budget.xlsx` in Excel before running `add` — Excel holds a lock on the
file while it is open.

## Input formats

### Bank CSV

Semicolon-separated, Danish number format (`-1.234,56`), one file per month:

```
Dato;Tekst;Beløb;Saldo;Afstemt;Kontonummer;Kontonavn;Hovedkategori;Kategori;Kommentar
31.08.2026;Lønoverførsel;28.566,24;86.566,24;;7670 2746936;Lønkonto;Indtægter;Løn, dagpenge og pension;
```

Read from each row: the date (column 1, `dd.mm.yyyy`), the amount (column 3),
the category (column 9) and the comment (column 10). Rows whose date or amount
do not parse are skipped, which takes care of the header row. All rows in a
file must fall in the same month.

### Budget workbook

`add` and `sync` find their way around by Excel table names, not cell
addresses, so you can move the tables and insert rows freely. The workbook
needs:

- a sheet named after the year (or whatever `--sheet` says),
- the tables `tblIncome`, `tblFixed`, `tblVariable` and `tblInvestment`,
- category labels in each table's first column, and
- month headers `JAN`…`DEC` in each table's header row.

A row labelled `Total` is skipped, so a table can carry its own total row.

## Project layout

```
cmd/              cobra commands: init, sync, add
internal/model/   the mapper, and mapping transactions onto budget cells
internal/csv/     bank export parsing
internal/xlsx/    reading and writing the workbook
months/           month parsing and the file name / column header forms
tables/           the Excel table names
assets/           the embedded budget template
```

The mapping itself lives in `internal/model/aggregate.go`: the mapper is
inverted once into a bank category → budget cell lookup, then each transaction
is placed with a single map hit.
