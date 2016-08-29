package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/floresj/pae-wdol/service"
)

func main() {
	saveFilename := "./WageDeterminationReport.xlsx"
	configFilename := "./code.config.txt"
	sourceFilename := "./WageDeterminations.xls"
	codes, err := service.ReadCodes(configFilename)
	if err != nil {
		fmt.Printf("Unable to read configuration file [%v]. Error [%v]\n", configFilename, err)
		os.Exit(1)
	}

	fmt.Printf("Searching for the following determination codes: [%v]\n", codes)
	if len(codes) == 0 {
		fmt.Println("No codes specified. You must specify at least one determination code")
		os.Exit(1)
	}

	fmt.Printf("Reading file [%v] as source and pulling wage determinations\n", sourceFilename)

	wageDeterminations, err := service.ReadSource(sourceFilename)
	if err != nil {
		fmt.Printf("Error Opening `WageDeterminations.xls` [%v]\n", err)
	}

	fmt.Printf("Saving results to [%v]\n", sourceFilename)
	err = service.ExportWageDeterminations(saveFilename, wageDeterminations, codes)
	if err != nil {
		fmt.Printf("Unable to export results to xlsx. Error [%v]\n", err)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Done! Your file saved. Press any key to close")
	reader.ReadString('\n')
}
