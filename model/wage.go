package model

import "fmt"

type Wage struct {
	OccupationCode string
	Title          string
	Rate           string
	Holiday        string
	Raw            string
}

func (w Wage) String() string {
	return fmt.Sprintf("Wage:\n\tOccupationCode: [%v]\n\tTitle: [%v]\n\tRate: [%v]\n",
		w.OccupationCode, w.Title, w.Rate)
}
