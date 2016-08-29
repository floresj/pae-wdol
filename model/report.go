package model

import "strings"

type Report struct {
	Location          string
	County            string
	State             string
	WageDetermination string
	RevisionNumber    string
	RevisionDate      string
	Wages             []Wage
	Link              string
}

func (r *Report) Load(record []string) {
	r.Location = strings.Trim(record[0], " ")
	r.County = strings.Trim(record[1], " ")
	r.State = strings.Trim(record[2], " ")
	r.WageDetermination = strings.Trim(record[3], " ")
	r.RevisionNumber = strings.Trim(record[4], " ")
	r.RevisionDate = strings.Trim(record[5], " ")
	r.Link = strings.Trim(record[6], " ")
}
