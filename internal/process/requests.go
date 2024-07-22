package process

import (
	"fmt"
	"strings"

	"github.com/Galdoba/ffproc/internal/mfscan"
	"github.com/Galdoba/ffproc/internal/pkg/sourcefiles"
	"github.com/Galdoba/ffproc/internal/pkg/ticket"
)

// func (pr *Project) MakeRequests() error {
// 	//entrySourcesMap := pr.Data.TableSourceMap
// 	tk := pr.Ticket
// 	totalVideoNum := 0
// 	totalAudioNum := 0
// 	totalAudioChans := 0
// 	totalSubtitleNum := 0
// 	totals := make(map[string][4]int)
// 	ep := ""
// 	for _, stComb := range streamTypeCombinations(pr) {
// 		if ep != stComb.episode {
// 			totals[ep] = [4]int{totalVideoNum, totalAudioNum, totalAudioChans, totalSubtitleNum}
// 			ep = stComb.episode
// 			totalVideoNum = 0
// 			totalAudioNum = 0
// 			totalAudioChans = 0
// 			totalSubtitleNum = 0
// 		}
// 		totalVideoNum += stComb.numbers[0]
// 		totalAudioNum += stComb.numbers[1]
// 		totalSubtitleNum += stComb.numbers[2]
// 		fmt.Println(stComb)
// 		for _, chans := range stComb.channels {
// 			chn, err := strconv.Atoi(chans)
// 			if err != nil {
// 				fmt.Println("bad data", chans)
// 			}
// 			totalAudioChans += chn
// 		}

// 	}
// 	totals[ep] = [4]int{totalVideoNum, totalAudioNum, totalAudioChans, totalSubtitleNum}

// 	for k, total := range totals {
// 		if k == "" && total[0]+total[1]+total[2]+total[3] == 0 {
// 			continue
// 		}

// 		switch k {
// 		case "":
// 		default:
// 			k = k + ": "
// 		}
// 		videoRequest(tk, k, total[0])
// 		audioRequest(tk, k, total[1], total[2])
// 		subtitleRequest(tk, k, total[3])
// 		if k != "" {
// 			tk.AddRequest(fmt.Sprintf("define sources for %v", k), "AUTO")
// 		}

// 	}

// 	// panic("с реквестами не закончил")
// 	return fmt.Errorf("???")

// }

func (pr *Project) MakeRequests() error {
	// streamsVideo := pr.sourceStreams("v")
	pr.processData()
	if err := pr.videoRequest(); err != nil {
		return fmt.Errorf("video requests failed: %v", err)
	}

	return nil
}

func (pr *Project) processData() error {
	sourcesByEpisode := pr.sourcesByEpisode()
	for _, episode := range sourcesByEpisode {
		for i, src := range episode {
			ep := src.EpisodeTag()
			if ep != "" {
				ep += "_"
			}
			pr.Ticket.Info_Tags[ticket.PROCESS_DATA][fmt.Sprintf("%vsource_%v", ep, i+1)] = src.Name
		}
	}
	return nil
}

func (pr *Project) videoRequest() error {
	sourcesByEpisode := pr.sourcesByEpisode()
	for _, episode := range sourcesByEpisode {

		ep, streamMap := pr.collectStreams(episode, "video")
		if ep != "" {
			ep += " "
		}

		for key, stream := range streamMap {
			if stream.Width == 1920 && stream.Height == 1080 {
				pr.Ticket.AddRequest(fmt.Sprintf("%vvideo scale", ep), "NONE")
				switch len(streamMap) {
				case 0:
				default:
					pr.Ticket.AddRequest(fmt.Sprintf("%vvideo mapping", ep), key)
				}
			}

		}

		fmt.Println("REQUEST FOR EPISODE HERE", ep)
	}

	// for _, episodes := range pr.sourceStreamsByEpisodes("v") {
	// 	streams := []*mfscan.Stream{}
	// 	keys := []string{}
	// 	for k, _ := range episodes {
	// 		keys = appendUnique(keys, k)
	// 	}
	// 	for _, k := range keys {
	// 		streams = append(streams, episodes[k])
	// 		kSplit := strings.Split(k, "|")
	// 		ep := kSplit[0]
	// 		strKey := kSplit[1]
	// 		if ep != "" {
	// 			ep += " "
	// 		}
	// 		switch len(streams) {
	// 		case 0:
	// 			pr.Ticket.AddRequest(fmt.Sprintf("%vvideo source", ep), "NONE")
	// 		case 1:
	// 			pr.Ticket.AddRequest(fmt.Sprintf("%vvideo source", ep), strKey)
	// 			if streams[0].Width == 1920 && streams[0].Height == 1080 {
	// 				pr.Ticket.AddRequest(fmt.Sprintf("%vvideo scale", ep), "CONFIRMED")
	// 			} else {
	// 				pr.Ticket.AddRequest(fmt.Sprintf("%vvideo scale", ep), "DECIDION REQUIRED")
	// 			}
	// 			pr.Ticket.AddRequest(fmt.Sprintf("%vvideo interlace", ep), "NOT SCANNED")
	// 		default:
	// 			pr.Ticket.AddRequest(fmt.Sprintf("%vvideo source", ep), "MAPPING REQUIRED")
	// 		}
	// 	}

	// }
	return nil
}

func (pr *Project) collectStreams(sources []*sourcefiles.Sourcefile, stype string) (string, map[string]*mfscan.Stream) {
	streams := make(map[string]*mfscan.Stream)
	ep := ""
	for _, src := range sources {
		for key, stream := range pr.Data.SourceStreamMap[src] {
			switch stype {
			case "video", "audio", "subtitle":
				streams[key] = stream
			}
			// if stream.Codec_type == "audio" {
			// 	streams[key] = stream
			// 	fmt.Println(len(streams), "audio for", src.EpisodeTag(), src.Name)
			// }
		}
		ep = src.EpisodeTag()
	}
	return ep, streams
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

func (pr *Project) sourcesByEpisode() [][]*sourcefiles.Sourcefile {
	return sourcefiles.SortedByEpisode(pr.sources())
}

// func (pr *Project) sourceStreamsByEpisodes(mask string) []map[string]*mfscan.Stream {
// 	streamsSlice := []map[string]*mfscan.Stream{}

// 	vid, aud, sub := false, false, false
// 	if strings.Contains(mask, "v") {
// 		vid = true
// 	}
// 	if strings.Contains(mask, "a") {
// 		aud = true
// 	}
// 	if strings.Contains(mask, "s") {
// 		sub = true
// 	}
// 	episodes := pr.sourcesByEpisode()
// 	for _, sources := range episodes {
// 		for i, src := range sources {
// 			episodeStreamMap := make(map[string]*mfscan.Stream)
// 			srcStreams := pr.Data.SourceStreamMap[src]
// 			for k, v := range srcStreams {
// 				if v.Codec_type == "video" && vid {
// 					episodeStreamMap[fmt.Sprintf("%v|%v%v", src.EpisodeTag(), i, k)] = v
// 				}
// 				if v.Codec_type == "audio" && aud {
// 					episodeStreamMap[fmt.Sprintf("%v|%v%v", src.EpisodeTag(), i, k)] = v
// 				}
// 				if v.Codec_type == "subtitle" && sub {
// 					episodeStreamMap[fmt.Sprintf("%v|%v%v", src.EpisodeTag(), i, k)] = v
// 				}
// 			}
// 			streamsSlice = append(streamsSlice, episodeStreamMap)
// 		}
// 	}

// 	return streamsSlice
// }

// func mapSourceToStreams(sources []sourcefiles.Sourcefile) map[string]map[string]*mfscan.Stream {
// 	stsMap := make(map[string]map[string]*mfscan.Stream)
// 	return nil
// }

/*
source=>
	file1=>
		stream [0:v:0]
		stream [0:a:0]
		stream [0:a:1]
	file2=>
		stream [1:s:0]

map[string]map[string]Stream
*/
