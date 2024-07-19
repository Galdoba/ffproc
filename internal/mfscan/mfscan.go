package mfscan

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Galdoba/devtools/cli/command"
)

const (
	ScanBasic     = "basic"
	ScanInterlace = "interlace"
	ScanSilence   = "silence"
	stVideo       = "video"
	stAudio       = "audio"
	stSubs        = "subtitle"
	tagLanguage   = "language"
)

func NewProfile() *mfscanInfo {
	sr := &mfscanInfo{}
	//	sr.Format = &Format{}
	return sr
}

func (prof *mfscanInfo) ConsumeFile(path string) error {
	stdout, stderr, err := command.Execute("ffprobe "+fmt.Sprintf("-v quiet -print_format json -show_format -show_streams -show_programs %v", path), command.Set(command.BUFFER_ON))
	if err != nil {
		if err.Error() != "exit status 1" {
			return fmt.Errorf("execution error: %v", err.Error())
		}
	}
	if stderr != "" {
		fmt.Println("stderr:")
		fmt.Println(stderr)
		panic("неожиданный выхлоп")
		//
	}
	data := []byte(stdout)
	if len(data) == 0 {
		flbts, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("file reading error: %v", err.Error())
		}
		if len(flbts) == 0 {
			return fmt.Errorf("file empty: %v", path)
		}
		check, _ := command.New(
			command.CommandLineArguments("ffprobe", fmt.Sprintf("-hide_banner "+fmt.Sprintf("-i %v", path))),
			//command.Set(command.TERMINAL_ON),
			command.Set(command.BUFFER_ON),
		)
		check.Run()
		checkOut := check.StdOut() + check.StdErr()
		if checkOut != "" {
			return fmt.Errorf("can't read: %v", checkOut)
		}
	}
	err = json.Unmarshal(data, &prof)
	if err != nil {
		return fmt.Errorf("can't unmarshal data from file: %v (%v)\n%v", err.Error(), path, string(data))
	}

	return nil
}

type mfscanInfo struct {
	Format  *Format   `json:"format"`
	Streams []*Stream `json:"streams,omitempty"`
}

type Format struct {
	Bit_rate         string            `json:"bit_rate,omitempty"`
	Duration         string            `json:"duration,omitempty"`
	Filename         string            `json:"filename,omitempty"`
	Format_long_name string            `json:"format_long_name,omitempty"`
	Format_name      string            `json:"format_name,omitempty"`
	Nb_programs      int               `json:"nb_programs,omitempty"`
	Nb_streams       int               `json:"nb_streams,omitempty"`
	Probe_score      int               `json:"probe_score,omitempty"`
	Size             string            `json:"size,omitempty"`
	Start_time       string            `json:"start_time,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
}

type Stream struct {
	Avg_frame_rate         string                  `json:"avg_frame_rate,omitempty"`
	Bit_rate               string                  `json:"bit_rate,omitempty"`
	Bits_per_raw_sample    string                  `json:"bits_per_raw_sample,omitempty"`
	Bits_per_sample        int                     `json:"bits_per_sample,omitempty"`
	Channel_layout         string                  `json:"channel_layout,omitempty"`
	Channels               int                     `json:"channels,omitempty"`
	Chroma_location        string                  `json:"chroma_location,omitempty"`
	Closed_captions        int                     `json:"closed_captions,omitempty"`
	Codec_long_name        string                  `json:"codec_long_name,omitempty"`
	Codec_name             string                  `json:"codec_name,omitempty"`
	Codec_tag              string                  `json:"codec_tag,omitempty"`
	Codec_tag_string       string                  `json:"codec_tag_string,omitempty"`
	Codec_time_base        string                  `json:"codec_time_base,omitempty"`
	Codec_type             string                  `json:"codec_type,omitempty"`
	Coded_height           int                     `json:"coded_height,omitempty"`
	Coded_width            int                     `json:"coded_width,omitempty"`
	Color_primaries        string                  `json:"color_primaries,omitempty"`
	Color_range            string                  `json:"color_range,omitempty"`
	Color_space            string                  `json:"color_space,omitempty"`
	Color_transfer         string                  `json:"color_transfer,omitempty"`
	Display_aspect_ratio   string                  `json:"display_aspect_ratio,omitempty"`
	Divx_packed            string                  `json:"divx_packed,omitempty"`
	Dmix_mode              string                  `json:"dmix_mode,omitempty"`
	Duration               string                  `json:"duration,omitempty"`
	Duration_ts            int                     `json:"duration_ts,omitempty"`
	Field_order            string                  `json:"field_order,omitempty"`
	Has_b_frames           int                     `json:"has_b_frames,omitempty"`
	Height                 int                     `json:"height,omitempty"`
	Id                     string                  `json:"id,omitempty"`
	Index                  int                     `json:"index,omitempty"`
	Is_avc                 string                  `json:"is_avc,omitempty"`
	Level                  int                     `json:"level,omitempty"`
	Loro_cmixlev           string                  `json:"loro_cmixlev,omitempty"`
	Loro_surmixlev         string                  `json:"loro_surmixlev,omitempty"`
	Ltrt_cmixlev           string                  `json:"ltrt_cmixlev,omitempty"`
	Ltrt_surmixlev         string                  `json:"ltrt_surmixlev,omitempty"`
	Max_bit_rate           string                  `json:"max_bit_rate,omitempty"`
	Nal_length_size        string                  `json:"nal_length_size,omitempty"`
	Nb_frames              string                  `json:"nb_frames,omitempty"`
	Pix_fmt                string                  `json:"pix_fmt,omitempty"`
	Profile                string                  `json:"profile,omitempty"`
	Progressive_frames_pct float64                 `json:"progressive_frames_pct,omitempty"`
	Quarter_sample         string                  `json:"quarter_sample,omitempty"`
	R_frame_rate           string                  `json:"r_frame_rate,omitempty"`
	Refs                   int                     `json:"refs,omitempty"`
	Sample_aspect_ratio    string                  `json:"sample_aspect_ratio,omitempty"`
	Sample_fmt             string                  `json:"sample_fmt,omitempty"`
	Sample_rate            string                  `json:"sample_rate,omitempty"`
	Start_pts              int                     `json:"start_pts,omitempty"`
	Start_time             string                  `json:"start_time,omitempty"`
	Time_base              string                  `json:"time_base,omitempty"`
	Width                  int                     `json:"width,omitempty"`
	Side_data_list         []Side_data_list_struct `json:"side_data_list,omitempty"`
	Tags                   map[string]string       `json:"tags,omitempty"`
	Disposition            map[string]int          `json:"disposition,omitempty"`
	// SilenceData            []SilenceSegment        `json:"silence_segments,omitempty"`
}

type Side_data_list_struct struct {
	Side_data map[string]string
}

func (si *mfscanInfo) Size() []string {
	sizes := []string{}
	for _, stream := range si.Streams {
		if stream.Codec_type != stVideo {
			continue
		}
		sizes = append(sizes, fmt.Sprintf("%vx%v", stream.Width, stream.Height))
	}
	return sizes
}

func (si *mfscanInfo) Fps() []string {
	fps := []string{}
	for _, stream := range si.Streams {
		if stream.Codec_type != stVideo {
			continue
		}
		fps = append(fps, stream.Avg_frame_rate)
	}
	return fps
}

func (si *mfscanInfo) Channels() []int {
	channels := []int{}
	for _, stream := range si.Streams {
		if stream.Codec_type != stAudio {
			continue
		}
		channels = append(channels, stream.Channels)
	}
	return channels
}

func (si *mfscanInfo) Layout() []string {
	layout := []string{}
	for _, stream := range si.Streams {
		if stream.Codec_type != stAudio {
			continue
		}
		layout = append(layout, stream.Channel_layout)
	}
	return layout
}

func (si *mfscanInfo) Langs() []string {
	layout := []string{}
	for _, stream := range si.Streams {
		if stream.Codec_type != stAudio {
			continue
		}
		layout = append(layout, stream.Tags[tagLanguage])
	}
	return layout
}

func (si *mfscanInfo) StreamTypes() []string {
	strms := []string{}
	for _, stream := range si.Streams {
		switch stream.Codec_type {
		case stVideo, stAudio, stSubs:
			strms = append(strms, stream.Codec_type)
		}
	}
	return strms
}
