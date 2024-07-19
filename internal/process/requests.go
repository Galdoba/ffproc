package process

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Galdoba/ffproc/internal/pkg/sourcefiles"
	"github.com/Galdoba/ffproc/internal/pkg/ticket"
)

func (pr *Project) MakeRequests() error {
	//entrySourcesMap := pr.Data.TableSourceMap
	tk := pr.Ticket
	totalVideoNum := 0
	totalAudioNum := 0
	totalAudioChans := 0
	totalSubtitleNum := 0
	totals := make(map[string][4]int)
	ep := ""
	for _, stComb := range streamTypeCombinations(pr) {
		if ep != stComb.episode {
			totals[ep] = [4]int{totalVideoNum, totalAudioNum, totalAudioChans, totalSubtitleNum}
			ep = stComb.episode
			totalVideoNum = 0
			totalAudioNum = 0
			totalAudioChans = 0
			totalSubtitleNum = 0
		}
		totalVideoNum += stComb.numbers[0]
		totalAudioNum += stComb.numbers[1]
		totalSubtitleNum += stComb.numbers[2]
		fmt.Println(stComb)
		for _, chans := range stComb.channels {
			chn, err := strconv.Atoi(chans)
			if err != nil {
				fmt.Println("bad data", chans)
			}
			totalAudioChans += chn
		}

	}
	totals[ep] = [4]int{totalVideoNum, totalAudioNum, totalAudioChans, totalSubtitleNum}

	for k, total := range totals {
		if k == "" && total[0]+total[1]+total[2]+total[3] == 0 {
			continue
		}

		switch k {
		case "":
		default:
			k = k + ": "
		}
		videoRequest(tk, k, total[0])
		audioRequest(tk, k, total[1], total[2])
		subtitleRequest(tk, k, total[3])
		if k != "" {
			tk.AddRequest(fmt.Sprintf("define sources for %v", k), "AUTO")
		}

	}
	panic("с реквестами не закончил")
	return fmt.Errorf("???")

}

func videoRequest(tk *ticket.Ticket, ep string, vNum int) error {
	switch vNum {
	case 0:
		tk.AddRequest(fmt.Sprintf("%vomit video", ep), "AUTO")
	case 1:
		tk.AddRequest(fmt.Sprintf("%vconfirm video", ep), "AUTO")
	default:
		tk.AddRequest(fmt.Sprintf("%vselect source video", ep), "USER")
	}
	return nil
}

func audioRequest(tk *ticket.Ticket, ep string, aNum, cNum int) error {
	switch aNum {
	case 0:
		tk.AddRequest(fmt.Sprintf("%vomit audio", ep), "AUTO")
	case 1:
		switch cNum {
		default:
			tk.AddRequest(fmt.Sprintf("%vinspect audio", ep), "USER")
		case 2, 6:
			tk.AddRequest(fmt.Sprintf("%vconfirm audio", ep), "AUTO")
		}
	case 2:
		tk.AddRequest(fmt.Sprintf("%vconfirm audio", ep), "USER")
		switch cNum {
		case 2:
			tk.AddRequest(fmt.Sprintf("%vaudio merge mono", ep), "USER")
		}
	default:
		tk.AddRequest(fmt.Sprintf("%vselect audio", ep), "USER")
		if aNum > 5 {
			tk.AddRequest(fmt.Sprintf("%vmap audio channels", ep), "USER")
		}
	}

	// switch aNum {
	// case 2, 4, 6, 8, 12, 0:
	// default:
	// 	tk.AddRequest(fmt.Sprintf("%vmap audio channels", ep), "USER")
	// }

	if cNum == 2 && aNum == 2 {
		tk.AddRequest(fmt.Sprintf("%vmerge 2 mono to stereo", ep), "AUTO")
	}
	if cNum == 6 && aNum == 6 {
		tk.AddRequest(fmt.Sprintf("%vmerge 6 mono to 5.1", ep), "AUTO")
	}
	return nil
}

func subtitleRequest(tk *ticket.Ticket, ep string, sNum int) error {
	switch sNum {
	case 0:
	case 1:
		tk.AddRequest(fmt.Sprintf("%vselect subtitle type", ep), "USER")
	case 2:
		tk.AddRequest(fmt.Sprintf("%vinspect subtitles", ep), "USER")
		tk.AddRequest(fmt.Sprintf("%vselect subtitle type", ep), "USER")
	}
	return nil
}

func streamTypeCombinations(pr *Project) []streamTypeCombination {
	stc := []streamTypeCombination{}
	entrySourcesMap := pr.Data.TableSourceMap
	tk := pr.Ticket

	for _, sources := range entrySourcesMap {
		chans := []string{}

		for _, src := range sources {
			vn, an, sn := 0, 0, 0
			ep := src.EpisodeTag()
			data := tk.Info_Tags[src.Name]
			for k, v := range data {
				if strings.Contains(k, ":v:") {
					vn++
				}
				if strings.Contains(k, ":a:") {
					an++
					chans = append(chans, v)
				}
				if strings.Contains(k, ":s:") {
					sn++
				}
			}
			stc = append(stc, streamTypeCombination{key: src.Name, episode: ep, numbers: [3]int{vn, an, sn}, channels: chans})
		}

	}
	return stc
}

type streamTypeCombination struct {
	key      string
	episode  string
	numbers  [3]int
	channels []string
}

func (pr *Project) sources() []*sourcefiles.Sourcefile {
	for _, sources := range pr.Data.TableSourceMap {
		return sources
	}
	return nil
}
