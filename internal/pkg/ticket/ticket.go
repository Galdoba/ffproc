package ticket

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Galdoba/ffproc/internal/pkg/define"
)

const (
	DEFAULT = iota
	// TYPE_TRL          = "--TRL--"
	// TYPE_FILM         = "--FILM--"
	// TYPE_SER          = "--SER--"
	SOURCE_DATA       = "*source_profile"
	PROCESS_DATA      = "*Process data"
	PROCESS_REQUEST   = "*Process request"
	ARCHIVE_DEST      = "path_archive"
	PROCESS_DEST      = "path_process"
	SEASON_NUM        = "season "
	EPISODE_NUM       = "episode"
	EPISODES_EXPECTED = "expected_episodes"
	SEASON_EXPECTED   = "season"
)

type Ticket struct {
	Name        string                       `json:"Table Name"`
	TicketType  string                       `json:"Type"`
	BaseName    string                       `json:"Output Base,omitempty"`
	Agent       string                       `json:"Agent,omitempty"`
	Date        string                       `json:"Deliver Date,omitempty"`
	StartTime   time.Time                    `json:"Started at"`
	CloseTime   string                       `json:"Closed at,omitempty"`
	Info_Tags   map[string]map[string]string `json:"Info,omitempty"`                 //карта [имя_файла]набор_тэгов
	Processable bool                         `json:"Ready for processing,omitempty"` //карта [имя_файла]набор_тэгов
}

func New(name, tType string) *Ticket {
	tk := Ticket{}
	tk.Name = name
	tk.TicketType = tType
	tk.StartTime = time.Now()
	tk.Info_Tags = make(map[string]map[string]string)
	tk.Info_Tags[PROCESS_DATA] = make(map[string]string)
	return &tk
}

var SourceExists error = errors.New("source exists")

func (tk *Ticket) AddSource(source, mfdata string) error {
	if _, ok := tk.Info_Tags[source]; ok {
		return SourceExists
	}
	tk.Info_Tags[source] = make(map[string]string)
	tk.Info_Tags[source][SOURCE_DATA] = mfdata
	return nil
}

//AddTag - добавляет тэг с информацией
//менять и удалять тэг нельзя
func (tk *Ticket) AddTag(source, key, val string) error {
	if _, ok := tk.Info_Tags[source]; !ok {
		return fmt.Errorf("source '%v' was not added", source)
	}
	sourcemap := tk.Info_Tags[source]
	if val, ok := sourcemap[key]; ok {
		if val != "" {
			return fmt.Errorf("tag '%v' for source '%v' exists", key, source)
		}
	}
	tk.Info_Tags[source][key] = val
	return nil
}

//AddRequest - задел на будущее
//менять и удалять нельзя
func (tk *Ticket) AddRequest(key, val string) error {
	if _, ok := tk.Info_Tags[PROCESS_REQUEST]; !ok {
		tk.Info_Tags[PROCESS_REQUEST] = make(map[string]string)
	}
	requestmap := tk.Info_Tags[PROCESS_REQUEST]

	if val, ok := requestmap[key]; ok {
		if val != "" {
			return fmt.Errorf("request '%v' exists", key)
		}
	}
	tk.Info_Tags[PROCESS_REQUEST][key] = val
	return nil
}

func (tk *Ticket) ValidateWith(validator Validator) error {
	tk.Processable = false
	if _, ok := tk.Info_Tags[PROCESS_DATA]; !ok {
		return fmt.Errorf("process data absent")
	}
	procdata := tk.Info_Tags[PROCESS_DATA]
	if _, ok := procdata[PROCESS_DEST]; !ok {
		return fmt.Errorf("process destination absent")
	}
	if _, ok := procdata[ARCHIVE_DEST]; !ok {
		return fmt.Errorf("archive destination absent")
	}
	switch tk.TicketType {
	default:
		return fmt.Errorf("ticket type invalid")
	case define.PROJ_TYPE_TRL:
		if _, ok := procdata[SEASON_NUM]; !ok {
			return fmt.Errorf("season number absent")
		}
		if _, ok := procdata[EPISODE_NUM]; !ok {
			return fmt.Errorf("episode number absent")
		}
		if !validator.ValidateAgent(procdata[EPISODES_EXPECTED]) {
			return fmt.Errorf("invalid expected episodes: [%v]", procdata[EPISODES_EXPECTED])
		}
	case define.PROJ_TYPE_FLM:
		if _, ok := procdata[SEASON_NUM]; ok {
			return fmt.Errorf("season number not expected")
		}
		if _, ok := procdata[EPISODE_NUM]; !ok {
			return fmt.Errorf("episode number not expected")
		}
	}
	if !validator.ValidateAgent(tk.Agent) {
		return fmt.Errorf("invalid agent: %v", tk.Agent)
	}
	if !validator.ValidateProcessPath(procdata[PROCESS_DEST]) {
		return fmt.Errorf("invalid processing destination: %v", procdata[PROCESS_DEST])
	}
	if !validator.ValidateArchivePath(procdata[ARCHIVE_DEST]) {
		return fmt.Errorf("invalid archive destination: %v", procdata[ARCHIVE_DEST])
	}
	tk.Processable = true
	return nil
}

//Validator - ожидаем реализацию трех валидаторов ФИЛЬМ, СЕРИАЛ, ТРЕЙЛЕР
type Validator interface {
	ValidateAgent(string) bool
	ValidateProcessPath(string) bool
	ValidateArchivePath(string) bool
	ValidateEpisodes(string) bool
}

type pseudoValidator struct{}

func (pv *pseudoValidator) ValidateAgent(s string) bool {
	return true
}
func (pv *pseudoValidator) ValidateProcessPath(s string) bool {
	return true
}
func (pv *pseudoValidator) ValidateArchivePath(s string) bool {
	return true
}
func (pv *pseudoValidator) ValidateEpisodes(s string) bool {
	return true
}

var NoTicket error = errors.New("ticket not found")

func Load(path string) (*Ticket, error) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, NoTicket
	}
	bt, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("can't read file: %v", err)
	}
	tk := &Ticket{}
	if err := json.Unmarshal(bt, tk); err != nil {
		return nil, fmt.Errorf("can't unmarshal: %v", err)
	}
	return tk, nil
}

func Save(tk *Ticket, dir string) error {
	fi, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("dir not found: %v", dir)
	}
	if !fi.IsDir() {
		return fmt.Errorf("not a directory: %v", dir)
	}
	bt, err := json.MarshalIndent(tk, "", "  ")
	if err != nil {
		return fmt.Errorf("can't marshal %v: %v", tk.Name, err)
	}
	filepath := dir + tk.Name + "--" + tk.TicketType + ".json"
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return fmt.Errorf("can't open %v: %v", filepath, err)
	}
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("can't truncate %v: %v", filepath, err)
	}
	if _, err := f.Write(bt); err != nil {
		return fmt.Errorf("can't write %v: %v", filepath, err)
	}
	return nil
}
