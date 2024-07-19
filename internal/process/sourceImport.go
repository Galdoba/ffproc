package process

import (
	"fmt"

	"github.com/Galdoba/ffproc/internal/mfscan"
	"github.com/Galdoba/ffproc/internal/pkg/sourcefiles"
	"github.com/Galdoba/ffproc/internal/pkg/ticket"
)

func (pr *Project) ImportStreamData(src *sourcefiles.Sourcefile) error {
	mprof := mfscan.NewProfile()
	err := mprof.ConsumeFile(src.Dir + src.Name)
	if err != nil {
		fmt.Println(err)
	}
	pr.Ticket.Info_Tags[src.Name] = make(map[string]string)
	if pr.Ticket.Info_Tags[src.Name][ticket.SOURCE_DATA] == "" {
		pr.Ticket.Info_Tags[src.Name][ticket.SOURCE_DATA] = scanBasicData(src)
	}
	vn, an, sn := 0, 0, 0
	for _, stream := range mprof.Streams {
		data := StreamData{}
		switch stream.Codec_type {
		case "video":
			data = videoData(fmt.Sprintf(":v:%v", vn), stream.Avg_frame_rate, stream.Width, stream.Height)
			vn++
		case "audio":
			data = audioData(fmt.Sprintf(":a:%v", an), stream.Channel_layout, stream.Tags["language"], stream.Channels)
			an++
		case "subtitle":
			data = subsData(fmt.Sprintf(":s:%v", sn))
		default:
			continue
		}
		pr.Ticket.AddTag(src.Name, data.Key, stringFromStreamData(data))
		switch pr.Ticket.TicketType {
		case "SER":
			pr.Ticket.AddTag(src.Name, "episode", src.EpisodeTag())

		}
	}
	return nil
}

type StreamData struct {
	Key      string
	Type     string
	Width    int
	Height   int
	Fps      string
	Layout   string
	Lang     string
	Channels int
}

func videoData(key, fps string, width, height int) StreamData {
	return StreamData{
		Key:    key,
		Type:   "video",
		Fps:    fps,
		Width:  width,
		Height: height,
	}
}

func audioData(key, layout, lang string, channels int) StreamData {
	return StreamData{
		Key:      key,
		Type:     "audio",
		Layout:   layout,
		Lang:     lang,
		Channels: channels,
	}
}

func subsData(key string) StreamData {
	return StreamData{
		Key:  key,
		Type: "subtitle",
	}
}

func stringFromStreamData(sd StreamData) string {
	switch sd.Type {
	case "video":
		return fmt.Sprintf("%v|%vx%v", sd.Fps, sd.Width, sd.Height)
	case "audio":
		return fmt.Sprintf("%v", sd.Channels)
	case "subtitle":
		return fmt.Sprintf("subtitles")
	default:

	}
	return "????"
}
