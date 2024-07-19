package process

import (
	"fmt"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/ffproc/configs"
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

func (pr *Process) Run(cfg *configs.Procman) error {
	//SCAN
	totalProjects := len(pr.Projects)
	for pNum, project := range pr.Projects {
		tk := project.Ticket
		fmt.Printf("start project %v/%v (%v)\n", pNum+1, totalProjects, tk.Name)

		for _, sources := range project.Data.TableSourceMap {
			srcNum := len(sources)
			for sNum, src := range sources {
				fmt.Printf("import source %v/%v (%v)\n", sNum+1, srcNum, src.Name)
				if err := project.ImportStreamData(src); err != nil {
					fmt.Println("source importing failed:", err.Error())
				}
			}
		}
		fmt.Printf("TODO: form requests\n")
		project.MakeRequests()
		fmt.Printf("TODO: user confirm\n")
		fmt.Printf("update ticket\n")
		ticket.Save(tk, cfg.TicketStorage)
		fmt.Printf("TODO: assemble shell\n")
	}

	return fmt.Errorf("code me")
}

/*
шаблон продукта:
Priority=X
Source1="path1"
...
SourceN="pathN" (Data.TableSourceMap)
OutBase=""  (tk.Basename)
Editroot="" (cfg.EDIT)
ArchiveRoot="" (cfg.Archive)
InputDir="" (cfg.InputDir)
ProrgessDir="" (cfg.ProrgessDir)
DoneDir="" (cfg.DoneDir)
FC_VID_1="" (hsub+scale+setsar+unsharp+pad requests) scale=1920:-2,setsar=1/1,unsharp=3:3:0.3:3:3:0,pad=1920:1080:-1:-1
...
FC_VID_N="" (hsub+scale+setsar+unsharp+pad requests) scale=1920:-2,setsar=1/1,unsharp=3:3:0.3:3:3:0,pad=1920:1080:-1:-1
atempo=fps
FC_AUD_1="" (resample+atempo requests) aresample=48000,atempo=25/(atempo)
...
FC_AUD_N="" (resample+atempo requests) aresample=48000,atempo=25/(atempo)
AUDIOTAG_N="[LANG_N]+[LAYOUT_N]" (resample+atempo requests) aresample=48000,atempo=25/(atempo)


Revision_Num="" (user_defined)
[setup_paths]
[run ffmpeg]
[make_ready]
[notify]
[move_to_done]
[append_archivator]



PRIORITY=5
FILE="__________"
OUTBASE="__________"
ROOT="/mnt/pemaltynov"
EDIT="/ROOT/EDIT/_______"
EDIT_PATH="${ROOT}${EDIT}"
ARCHIVE_PATH="/mnt/pemaltynov/ROOT/IN/_______"
FC_AUD="aresample=48000,atempo=25/(__)"
AUDIO_OUT1="AUDIORUS__"
AUDIO_OUT2="AUDIOENG__"
REV=""
mkdir -p ${ARCHIVE_PATH}/_DONE/${OUTBASE}
mkdir -p ${EDIT_PATH}
clear && mv /home/pemaltynov/IN/${FILE} /home/pemaltynov/IN/_IN_PROGRESS/  && fflite -r 25 -i /home/pemaltynov/IN/_IN_PROGRESS/${FILE} \
-filter_complex "[0:v:0]split=2[vidHD][inProxy]; [inProxy]scale=iw/2:ih, setsar=(1/1)*2[vidHD_pr]; [0:a:0]${FC_AUD}[arus_in];  [arus_in]asplit=2[arus][arus_pr]; [0:a:1]${FC_AUD}[aeng_in];  [aeng_in]asplit=2[aeng][aeng_pr]" \
  -map "[vidHD]" -c:v libx264 -preset medium -crf 16 -pix_fmt yuv420p -g 0 -map_metadata -1 -map_chapters -1 ${EDIT_PATH}/${OUTBASE}${REV}_HD.mp4  \
  -map "[vidHD_pr]" -c:v libx264 -x264opts interlaced=1 -preset superfast -pix_fmt yuv420p  -b:v 2000k -maxrate 2000k -map_metadata -1 -map_chapters -1 ${EDIT_PATH}/${OUTBASE}${REV}_HD_proxy.mp4  \
  -map "[arus]" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 ${EDIT_PATH}/${OUTBASE}${REV}_${AUDIO_OUT1}.m4a \
  -map "[arus_pr]" -c:a ac3 -b:a 128k ${EDIT_PATH}/${OUTBASE}${REV}_${AUDIO_OUT1}_proxy.ac3  \
  -map "[aeng]" -c:a alac -compression_level 0 -map_metadata -1 -map_chapters -1 ${EDIT_PATH}/${OUTBASE}${REV}_${AUDIO_OUT2}.m4a \
  -map "[aeng_pr]" -c:a ac3 -b:a 128k ${EDIT_PATH}/${OUTBASE}${REV}_${AUDIO_OUT2}_proxy.ac3  \
  && touch ${EDIT_PATH}/${OUTBASE}${REV}.ready \
&& echo "${EDIT}/${OUTBASE}${REV}.ready" > /home/pemaltynov/IN/notifications/${OUTBASE}.done \
&& mv /home/pemaltynov/IN/_IN_PROGRESS/${FILE} /home/pemaltynov/IN/_DONE/  && at now + 10 hours <<< "mv /home/pemaltynov/IN/_DONE/${FILE} ${ARCHIVE_PATH}/_DONE/${OUTBASE}" \
&& clear && mv "$0" /home/pemaltynov/IN/_DONE/bash/

*/
