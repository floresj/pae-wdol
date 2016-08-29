package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/extrame/xls"
	"github.com/floresj/pae-wdol/model"
)

var (
	wageRegexp                = regexp.MustCompile("(\\d{5})\\s-\\s(.*)\\b(\\d{1,}.\\d{2})")
	holidayRegexp             = regexp.MustCompile("HOLIDAYS: A minimum of (\\w.+) paid holidays")
	revisionRegexp            = regexp.MustCompile("Revision\\sNo.:\\s(\\d+)")
	revisionDateRegexp        = regexp.MustCompile("Date Of Revision:\\s(\\d+/\\d+/\\d+)")
	determinationNumberRegexp = regexp.MustCompile("Wage Determination No.:\\s(\\d+-\\d+)")
	stateRegexp               = regexp.MustCompile("State[s]?:\\s(.+)")
	vacationRegexp            = regexp.MustCompile("(?m)^VACATION: ((.+)\\n){1,}(?:HOLIDAYS:)")
)

func ReadSource(path string) ([]model.WageDeterminationRequest, error) {
	wb, err := xls.Open(path, "utf-8")
	if err != nil {
		return nil, err
	}
	sheet := wb.GetSheet(0)
	if sheet != nil {

	}
	requests := []model.WageDeterminationRequest{}
	total := (int(sheet.MaxRow))
	for i := 0; i <= total; i++ {
		row := sheet.Rows[uint16(i)]
		city := row.Cols[0]
		county := row.Cols[1]
		state := row.Cols[2]

		link := row.Cols[6].String(wb)
		url := strings.Split(link[0], "(")[0]
		fmt.Printf("\rPulling Wage Determinations %v of %v from [%v]",
			i, total, url)
		if url == "" || !strings.Contains(url, "http") {
			continue
		}

		req := model.WageDeterminationRequest{Url: url}
		req.State = state.String(wb)[0]
		req.County = county.String(wb)[0]
		req.Location = city.String(wb)[0]
		rc, err := PullWageDetermination(url)
		if err != nil {
			req.Error = err
			requests = append(requests, req)
			fmt.Printf("\n\tUnable to Pull [%v] Error [%v]\n", url, req.Error)
			continue
		}
		w, err := ExtractWageDetermination(rc)
		if err != nil {
			fmt.Printf("Error Extracting Wage Determination")
			req.Error = err
			requests = append(requests, req)
			continue
		}
		w.State = state.String(wb)[0]
		req.WageDetermination = *w
		w.Url = url
		requests = append(requests, req)
	}
	return requests, err
}

func PullWageDetermination(path string) (io.ReadCloser, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unable to access url [%v]. Received response: [%v]",
			path, resp.Status)
	}
	return resp.Body, err
}

func ExtractWageDetermination(r io.ReadCloser) (*model.WageDetermination, error) {
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	data := string(b)

	wd := model.WageDetermination{}
	wd.Holiday = ExtractHoliday(data)
	wd.RevisionNumber = ExtractRevisionNumber(data)
	wd.RevisionDate = ExtractRevisionDate(data)
	wd.DeterminationNumber = ExtractDeterminationNumber(data)
	wd.State = ExtractState(data)
	wd.Vacation = ExtractVacation(data)

	matches := wageRegexp.FindAllStringSubmatch(data, -1)
	for _, match := range matches {
		if len(match) == 0 {
			continue
		}
		wage := ExtractWage(match)
		wd.Wages = append(wd.Wages, wage)
	}
	return &wd, err
}

func ExtractHoliday(data string) string {
	matches := holidayRegexp.FindStringSubmatch(data)
	if len(matches) == 0 {
		return ""
	}
	return strings.Trim(matches[1], " \r\n")
}

func ExtractRevisionNumber(data string) string {
	matches := revisionRegexp.FindStringSubmatch(data)
	if len(matches) == 0 {
		return ""
	}
	return strings.Trim(matches[1], " \r\n")
}

func ExtractRevisionDate(data string) string {
	matches := revisionDateRegexp.FindStringSubmatch(data)
	if len(matches) == 0 {
		return ""
	}
	return strings.Trim(matches[1], " \r\n")
}

func ExtractState(data string) string {
	matches := stateRegexp.FindStringSubmatch(data)
	if len(matches) == 0 {
		return ""
	}
	return strings.Trim(matches[1], " \r\n")
}

func ExtractDeterminationNumber(data string) string {
	matches := determinationNumberRegexp.FindStringSubmatch(data)
	if len(matches) == 0 {
		return ""
	}
	return strings.Trim(matches[1], " \r\n")
}

func ExtractVacation(data string) string {
	match := vacationRegexp.FindString(data)
	match = strings.Replace(match, "HOLIDAYS:", "", 1)
	match = strings.Replace(match, "VACATION:", "", 1)
	match = strings.Replace(match, "\n", "", -1)
	match = strings.Replace(match, "\r", "", -1)
	return strings.Trim(match, " \r\n")
}

func ExtractWage(m []string) model.Wage {
	w := model.Wage{}
	w.OccupationCode = strings.Trim(m[1], " \r\n")
	w.Title = strings.Trim(m[2], " \r\n")
	w.Rate = strings.Trim(m[3], " \r\n")
	w.Raw = strings.Trim(m[0], " \r\n")
	return w
}

func ReadCodes(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	rawCodes := strings.Split(string(data), "\n")
	var codes []string
	for _, code := range rawCodes {
		codes = append(codes, strings.Trim(code, " "))
	}
	return codes, nil
}

func toInt(v string) interface{} {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return v
	}
	return i
}

func toFloat(v string) interface{} {
	f, err := strconv.ParseFloat(v, 2)
	if err != nil {
		return v
	}
	return f
}
