package spreadsheet

import (
	"fmt"
	"testing"
)

func TestSpreadsheet(t *testing.T) {

	path3 := `c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\pkg\spreadsheet\v2\isCsv.csv`
	path4 := `c:\Users\pemaltynov\go\src\github.com\Galdoba\ffstuff\pkg\spreadsheet\v2\isNotExist.csv`

	paths := []string{path3, path4}
	for _, path := range paths {
		sp, err := New(path)
		if err != nil {
			fmt.Println("err:", err.Error())
			continue
		}
		fmt.Println(sp)
		fmt.Println("========")
		errUpd := sp.CurlUpdate(`https://docs.google.com/spreadsheets/d/1Waa58usrgEal2Da6tyayaowiWujpm0rzd06P5ASYlsg/edit?gid=250314867#gid=250314867`)
		if errUpd != nil {
			fmt.Println("errUpd:", errUpd.Error())
		}
		fmt.Println(sp)
	}

}
