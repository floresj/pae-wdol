package store

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/floresj/pae-wdol/model"
)

var (
	//wageRegexp    = regexp.MustCompile("(\\d{5})\\s-\\s(.*)(\\d\\d.\\d\\d)")
	wageRegexp    = regexp.MustCompile("(\\d{5})\\s-\\s(.*)\\b(\\d{1,}.\\d{2})")
	holidayRegexp = regexp.MustCompile("HOLIDAYS: A minimum of (\\w.+) paid holidays")
)

type WageStore struct {
}

func (s WageStore) Pull(url string) ([]model.Wage, error) {
	fmt.Printf("Pull Wage Determination from [%v]\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	content := string(b)
	holidayMatch := holidayRegexp.FindStringSubmatch(content)

	var holiday string
	if len(holidayMatch) > 0 {
		holiday = strings.Trim(holidayMatch[1], " ")
		fmt.Printf("Found Holiday [%v]\n", holiday)
	} else {
		fmt.Printf("Could not found holiday information")
	}

	wages := []model.Wage{}
	matches := wageRegexp.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		wage := s.toWage(match)
		wage.Holiday = holiday
		wages = append(wages, wage)
	}
	return wages, nil
}

func (s WageStore) toWage(w []string) model.Wage {
	wage := model.Wage{}
	occupationCode := strings.Trim(w[1], " ")
	title := strings.Trim(w[2], " ")
	rate := strings.Trim(w[3], " ")

	wage.OccupationCode = occupationCode
	wage.Title = title
	wage.Rate = rate

	return wage
}
