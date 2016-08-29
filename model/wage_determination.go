package model

import (
	"fmt"
)

type WageDetermination struct {
	DeterminationNumber string
	RevisionNumber      string
	RevisionDate        string
	State               string
	Holiday             string
	Vacation            string
	Wages               Wages
	Url                 string
}

func (wd WageDetermination) String() string {
	return fmt.Sprintf("Wage Determination\n\tNumber: [%v]\n\tRevision Number: [%v]\n\tRevision Date: [%v]\n\tState: [%v]\n\tHoliday: [%v]", wd.DeterminationNumber, wd.RevisionNumber, wd.RevisionDate, wd.State, wd.Holiday)
}

type Wages []Wage

// Filter returns a new slice containing wages matching the occupation codes
func (w Wages) Filter(codes []string) []Wage {
	var filtered []Wage
	for _, wage := range w {
		for _, code := range codes {
			if wage.Title == code {
				filtered = append(filtered, wage)
			}
		}
	}
	return filtered
}
