package service

import (
	"fmt"

	"github.com/floresj/pae-wdol/model"
	"github.com/floresj/pae-wdol/styles"
	"github.com/tealeg/xlsx"
)

var (
	LocationCol            = 0
	CountyCol              = 1
	StateCol               = 2
	DeterminationNumberCol = 3
	RevisionNumberCol      = 4
	RevisionDateCol        = 5
	StatesCol              = 6
	CodeStartCol           = 7
	HolidayColOffset       = 1
	VacationColOffset      = 2
	UrlColOffset           = 3
	RawColOffset           = 4
	ErrorColOffset         = 5
)

func writeWageDeterminationSheet(f *xlsx.File, wd []model.WageDeterminationRequest, codes []string) error {

	sheet, err := f.AddSheet("Wage Determinations")
	if err != nil {
		return err
	}
	row := sheet.AddRow()
	row.SetHeightCM(1.9)
	styles.Header(row.AddCell()).SetValue("Location")
	styles.Header(row.AddCell()).SetValue("County")
	styles.Header(row.AddCell()).SetValue("State")
	styles.Header(row.AddCell()).SetValue("Determination Number")
	styles.Header(row.AddCell()).SetValue("Revision Number")
	styles.Header(row.AddCell()).SetValue("Revision Date")
	styles.Header(row.AddCell()).SetValue("State(s)")

	// Add headers for codes being searched
	for _, code := range codes {
		styles.Header(row.AddCell()).SetValue(code)
	}

	styles.Header(row.AddCell()).SetValue("Holidays")
	styles.Header(row.AddCell()).SetValue("Vacation")
	styles.Header(row.AddCell()).SetValue("Source URL")
	styles.Header(row.AddCell()).SetValue("Raw Rates")
	styles.Header(row.AddCell()).SetValue("Errors")

	for _, req := range wd {
		w := req.WageDetermination
		row := sheet.AddRow()
		row.SetHeightCM(1.8)
		styles.Center(row.AddCell()).SetValue(req.Location)
		styles.Center(row.AddCell()).SetValue(req.County)
		styles.Center(row.AddCell()).SetValue(req.State)
		styles.Center(row.AddCell()).SetValue(w.DeterminationNumber)
		styles.Center(row.AddCell()).SetValue(toInt(w.RevisionNumber))
		styles.Center(row.AddCell()).SetValue(w.RevisionDate)
		styles.Center(row.AddCell()).SetValue(w.State)

		var raw string
		for _, code := range codes {
			c := row.AddCell()
			filtered := w.Wages.Filter([]string{code})
			if len(filtered) > 1 {
				fmt.Printf("Wage with code [%v] contains more than one find. This is no bueno\n",
					code)
			}
			for _, f := range filtered {
				c.SetValue(toFloat(f.Rate))
				c.NumFmt = "$0.00"
				styles.Rate(c)
			}
			for _, f := range filtered {
				raw = raw + fmt.Sprintf("%v\n", f.Raw)
			}
		}
		styles.Center(row.AddCell()).SetValue(w.Holiday)
		styles.Wrap(row.AddCell()).SetValue(w.Vacation)
		styles.Center(row.AddCell()).SetValue(req.Url)
		styles.Wrap(styles.Center(row.AddCell())).SetValue(raw)
		styles.Wrap(row.AddCell()).SetValue(req.Error)
	}

	sheet.Col(LocationCol).Width = 14
	sheet.Col(CountyCol).Width = 14
	sheet.Col(StateCol).Width = 14
	sheet.Col(DeterminationNumberCol).Width = 18
	sheet.Col(RevisionNumberCol).Width = 14
	sheet.Col(RevisionDateCol).Width = 14
	sheet.Col(StatesCol).Width = 14

	codeColOffset := CodeStartCol + len(codes) - 1
	for i := CodeStartCol; i <= codeColOffset; i++ {
		sheet.Col(i).Width = 15
	}
	sheet.Col(codeColOffset + HolidayColOffset).Width = 14
	sheet.Col(codeColOffset + VacationColOffset).Width = 100
	sheet.Col(codeColOffset + UrlColOffset).Width = 40
	sheet.Col(codeColOffset + ErrorColOffset).Width = 60
	sheet.Col(codeColOffset + RawColOffset).Width = 60
	for x, row := range sheet.Rows {
		// Skip header
		if x == 0 {
			continue
		}
		c := row.Cells[codeColOffset+ErrorColOffset]
		if c.Value != "" {
			for _, cell := range row.Cells {
				cell = styles.Error(cell)
			}
		}
	}
	return nil
}

func ExportWageDeterminations(savePath string, wd []model.WageDeterminationRequest, codes []string) error {
	f := xlsx.NewFile()
	writeWageDeterminationSheet(f, wd, codes)
	f.Save(savePath)
	return nil
}
