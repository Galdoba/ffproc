package process

import (
	"fmt"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/ffproc/internal/bridge"
	"github.com/Galdoba/ffproc/internal/pkg/sourcefiles"
	"github.com/Galdoba/ffproc/internal/pkg/table"
	"github.com/Galdoba/ffproc/internal/pkg/ticket"
)

type Project struct {
	Key    string
	Ticket *ticket.Ticket
	Data   bridge.Bridge
	//Path   string
}

type Process struct {
	Projects []*Project
}

func New() *Process {
	pr := Process{}
	return &pr
}

func (pr *Process) AddProject(tk *ticket.Ticket, data *bridge.Bridge, path string) error {
	keys := append([]string{}, data.Key)
	// fmt.Println(data.Key)
	for _, v := range data.SourceTableMap {
		// fmt.Println(v.ProcessType + "--" + table.Key(v))
		keys = appendUnique(keys, v.ProcessType+"--"+table.Key(v))
	}
	for _, sources := range data.TableSourceMap {
		for _, src := range sources {
			// fmt.Println(src.Key())
			if strings.Contains(keys[0], src.Key()) {
				continue
			}
			keys = appendUnique(keys, src.Key())

		}
	}
	if len(keys) != 1 {
		return fmt.Errorf("non unique keys in project: %v", keys)
	}

	for _, project := range pr.Projects {
		if project.Key == keys[0] {
			return fmt.Errorf("non unique keys in process")
		}
	}

	pr.Projects = append(pr.Projects, &Project{
		Key:    data.Key,
		Ticket: tk,
		Data:   *data,
		//Path:   path,
	})
	return nil
}

func appendUnique(slice []string, element string) []string {
	for _, s := range slice {
		if s == element {
			return slice
		}
	}
	slice = append(slice, element)

	return slice
}

// ////////////////////////////////////////////////////////////////////////
func scanBasicData(src *sourcefiles.Sourcefile) string {
	cmnd, _ := command.New(
		command.CommandLineArguments("mfline", fmt.Sprintf("show -l %v", src.Dir+src.Name)),
		command.AddBuffer("buf"),
	)
	cmnd.Run()
	out := cmnd.Buffer("buf")
	result := strings.TrimSuffix(out.String(), "\n")
	return result
}

/////////////////////////////////////////////////////////////////////////

/*
Scan Sources if needed
Prompt shell
Save shell
*/

func (pr *Process) Run() error {
	//SCAN
	for _, project := range pr.Projects {
		tk := project.Ticket
		for _, sources := range project.Data.TableSourceMap {
			for _, src := range sources {

				if tk.Info_Tags[src.Name][ticket.SOURCE_DATA] == "" {
					tk.Info_Tags[src.Name][ticket.SOURCE_DATA] = scanBasicData(src)
				}

			}
		}
	}

	return fmt.Errorf("code me")
}
