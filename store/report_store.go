package store

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/floresj/pae-wdol/model"
)

type ReportStore struct {
	WageStore *WageStore
}

func (s ReportStore) GetReports(path string, codes []string) []model.Report {
	source, err := os.Open(path)
	if err != nil {
		os.Exit(1)
	}

	r := csv.NewReader(source)
	var reports []model.Report
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		report := model.Report{}
		report.Load(record)

		if err != nil {
			log.Fatal("Unable to query")
		}

		wages, err := s.WageStore.Pull(report.Link)
		if err != nil {

		}

		wagesToFind := []model.Wage{}
		for _, wage := range wages {
			for _, code := range codes {
				if wage.Title == code {
					wagesToFind = append(wagesToFind, wage)
				}
			}
		}
		report.Wages = wagesToFind
		reports = append(reports, report)
	}

	return reports
}

func (s ReportStore) Query(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", nil
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s ReportStore) Output(reports []model.Report, codes []string, path string) error {
	output, err := os.Create(path)
	if err != nil {
		return err
	}
	w := csv.NewWriter(output)
	header := []string{
		"Location",
		"Revision Date",
		"Revision Number",
		"Wage Determination",
		"Holiday",
	}

	for _, code := range codes {
		header = append(header, code)
	}
	w.Write(header)
	for _, report := range reports {

		row := []string{
			report.Location,
			report.RevisionDate,
			report.RevisionNumber,
			report.WageDetermination,
		}
		x := 0
		for _, wage := range report.Wages {
			if x == 0 {
				row = append(row, wage.Holiday)
				x = x + 1
			}
			row = append(row, wage.Rate)
		}

		w.Write(row)
	}
	w.Flush()
	return nil
}

func NewReportStore() *ReportStore {
	wageStore := &WageStore{}
	reportStore := &ReportStore{wageStore}
	return reportStore
}
