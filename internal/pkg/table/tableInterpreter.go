package table

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Galdoba/ffproc/pkg/translit"
)

type TableData struct {
	Entries map[string]Entry
}

type Entry struct {
	ID                 string
	RowNum             int
	Comment            string
	EditPath           string
	TrailerExpected    bool
	TrailerInProgress  bool
	TrailerClosed      bool
	TrailerReady       bool
	TrailerMaker       string
	SourceDownloading  bool //сиреневое
	IsReadyToProcess   bool
	EditInWork         bool
	EditDone           bool
	TicketClosed       bool //green
	Agent              string
	PublicationDate    string
	ProcessDestination string
	ProcessType        string
}

func NewEntry(dest string, data []string) Entry {
	if data[8] == "" {
		return Entry{}
	}
	ent := Entry{}
	ent.ProcessDestination = dest
	ent.Comment = data[0]
	ent.EditPath = data[1]
	switch data[2] {
	case "r", "к":
		ent.TrailerExpected = true
	case "y", "н":
		// ent.TrailerClosed = true
	}
	switch data[3] {
	case "r", "к":
		ent.TrailerInProgress = true
	case "y", "н":
		ent.TrailerReady = true
	case "g", "п":
		ent.TrailerClosed = true
	}
	ent.TrailerMaker = data[4]
	ent.ID = data[8]
	switch data[9] {
	case "v", "м":
		ent.SourceDownloading = true
	case "y", "н":
		ent.IsReadyToProcess = true
	case "g", "п":
		ent.TicketClosed = true
	}
	switch data[10] {
	case "r", "к":
		ent.EditInWork = true
	case "y", "н":
		ent.EditDone = true
	case "g", "п":
		ent.TicketClosed = true
	}
	switch data[11] {

	case "v", "м":
		ent.TicketClosed = true
	}

	ent.Agent = data[13]
	ent.PublicationDate = data[14]
	return ent

}

func Key(e Entry) string {
	kparts := strings.Split(e.ID, ":")
	key := translit.Process(kparts[1], translit.Title)
	return key
}

type Data interface {
	Data() [][]string
}

func CompileTableData(dt Data) (*TableData, error) {
	td := TableData{}
	td.Entries = make(map[string]Entry)
	data := dt.Data()
	dest := ""
	for i, row := range data {
		if len(row) != 15 {
			return nil, fmt.Errorf("unexpected data array size: have %v: expect 15", len(row))
		}
		if i != 0 && row[14] == "" {
			dest = row[8]
		}
		entry := NewEntry(dest, row)
		entry.RowNum = i
		if entry.ID == "" {
			continue
		}
		if !entry.TrailerInProgress && !entry.TrailerClosed && !entry.TrailerReady && !entry.SourceDownloading && !entry.IsReadyToProcess && !entry.EditInWork && !entry.EditDone && !entry.TicketClosed {
			continue
		}
		if entry.TrailerExpected {
			trlEntry := entry
			trlEntry.ProcessType = "TRL"
			trlEntry.ID = trlEntry.ProcessType + ":" + trlEntry.ID
			trlEntry.ProcessDestination = `//192.168.31.4/root/EDIT/@trailers_temp/`
			td.Entries[trlEntry.ID] = trlEntry
		}
		entry.ProcessType = processTypeOf(entry)
		switch entry.ProcessType {
		case "SER":
			entry.ProcessDestination = "[EDIT]/" + destinationSerial(entry)
		case "FLM":
			entry.ProcessDestination = "[EDIT]/" + destinationFilm(entry)
		}
		entry.ID = entry.ProcessType + ":" + entry.ID
		td.Entries[entry.ID] = entry
	}
	return &td, nil
}

func processTypeOf(entry Entry) string {
	keytr := translit.Process(entry.ID)
	if strings.Contains(keytr, "_sezon") {
		return "SER"
	}
	return "FLM"
}

func seasonOf(entry Entry) string {
	keytr := translit.Process(entry.ID)
	parts := strings.Split(keytr, "_sezon")
	words := strings.Split(parts[0], "_")
	return words[len(words)-1]
}

func destinationSerial(entry Entry) string {
	keytr := translit.Process(entry.ID, translit.Title)
	parts := strings.Split(keytr, "_sezon")
	words := strings.Split(parts[0], "_")
	words[len(words)-1] = "s" + words[len(words)-1]
	base := strings.Join(words, "_")
	baseSpl := strings.Split(base, "")
	baseSpl[0] = strings.ToUpper(baseSpl[0])
	base = strings.Join(baseSpl, "")
	agent := strings.TrimSuffix("_"+translit.Process(entry.Agent), "_")
	agent = strings.TrimSuffix("_"+translit.Process(entry.Agent), "_")
	dest := agent + "/" + base + "/"
	return dest
}

func destinationFilm(entry Entry) string {
	switch isDate(entry.ProcessDestination) {
	case true:
		date := strings.Split(entry.ProcessDestination, ".")
		dest := strings.TrimPrefix(date[2], "20") + "_" + date[1] + "_" + date[0] + "/"
		return dest
	default:
		return strings.TrimSuffix("_"+translit.Process(entry.Agent), "_")
	}
}

func isDate(str string) bool {
	parts := strings.Split(str, ".")
	if len(parts) != 3 {
		return false
	}
	for _, p := range parts {
		if _, err := strconv.Atoi(p); err != nil {
			return false
		}

	}
	return true
}
