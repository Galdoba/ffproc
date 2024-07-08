package spreadsheet

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/devtools/csvp"
)

type sheet struct {
	path  string
	cells [][]string
}

/*
адресс должен быть доступен для чтения/записи
path не может быть папкой
*/

func New(path string) (*sheet, error) {
	if err := pathError(path); err != nil {
		return nil, fmt.Errorf("error checking path: %v", err)
	}
	sp := sheet{}
	sp.path = path
	bt, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &sp, nil
	}
	container, _ := csvp.FromString(string(bt))
	for _, entries := range container.Entries() {
		sp.cells = append(sp.cells, entries.Fields())
	}
	return &sp, nil
}

const (
	ERR_INVALID_BOUNDS_ROW = "#ERROR: row out of bounds"
	ERR_INVALID_BOUNDS_COL = "#ERROR: col out of bounds"
)

type Sheet interface {
	CurlUpdate(string) error
	Path() string
	Data() [][]string
}

func (sp *sheet) CurlUpdate(url string) error {
	request := fmt.Sprintf(`-s --use-ascii --proxy http://proxy.local:3128 %v -k --output`, curlCsvUrl(url))
	_, _, err := command.Execute("curl "+request+" "+sp.path+".tmp", command.Set(command.TERMINAL_ON))
	if err != nil {
		return fmt.Errorf("error updating: %v", err)
	}
	newPath := sp.path + ".tmp"
	oldPath := sp.path
	if err := os.Rename(newPath, oldPath); err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}
	bt, err := os.ReadFile(sp.path)
	if err != nil {
		return fmt.Errorf("can't read %v", sp.path)
	}
	container, err := csvp.FromString(string(bt))
	if err != nil {
		return fmt.Errorf("can't convert csv %v", sp.path)
	}
	for _, entries := range container.Entries() {
		sp.cells = append(sp.cells, entries.Fields())
	}
	return nil
}

func curlCsvUrl(url string) string {
	uparts := strings.Split(url, "/edit?gid")
	return uparts[0] + "/gviz/tq?tqx=out:csv"
}

func (sp *sheet) Cell(row, col int) string {
	if row < 0 {
		return ERR_INVALID_BOUNDS_ROW
	}
	if col < 0 {
		return ERR_INVALID_BOUNDS_COL
	}
	if row >= len(sp.cells) {
		return ERR_INVALID_BOUNDS_ROW
	}
	rowData := sp.cells[row]
	if col >= len(rowData) {
		return ERR_INVALID_BOUNDS_COL
	}
	return rowData[col]
}

func (sp *sheet) Size() (int, int) {
	w := len(sp.cells)
	if w == 0 {
		return 0, 0
	}
	return w, len(sp.cells[0])
}

func (sp *sheet) Path() string {
	return sp.path
}

func (sp *sheet) Data() [][]string {
	return sp.cells
}

////helpers
func pathError(path string) error {
	if !strings.HasSuffix(path, ".csv") {
		return fmt.Errorf("path is not csv")
	}
	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("os.Stat: %v", err)
	}
	if f.IsDir() {
		return fmt.Errorf("path is dir")
	}
	bt, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %v", err)
	}
	_, err = csvp.FromString(string(bt))
	if err != nil {
		return fmt.Errorf("read csv: %v", err)
	}
	return nil
}
