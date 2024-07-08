package mfscan

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
)

const (
	ScanBasic     = "basic"
	ScanInterlace = "interlace"
	ScanSilence   = "silence"
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
	err = prof.validate()
	if err != nil {
		return fmt.Errorf("validation error: %v", err.Error())
	}
	//fmt.Println(prof.Short())

	return nil
}

type mfscanInfo struct {
	Format     *Format   `json:"format"`
	Streams    []*Stream `json:"streams,omitempty"`
	warnings   []string
	streamInfo map[string]string
	short      string
	long       string
	chanLayout string
	Lang       map[string]string
	ChanNum    map[string]string
	VideoSize  map[string]string
	VideoFPS   map[string]string
}

// type mfscanInfo interface {
// 	Warnings() []string
// 	Short() string
// 	Long() string
// 	AudioLayout() string
// 	ConfirmScan(string) error
// 	ConsumeFile(string) error
// 	ScanBasic(string) error
// }

// func (mp *mfscanInfo) ScanBasic(sourceFile string) error {
// 	if strings.Contains(sourceFile, " ") {
// 		return fmt.Errorf("can't perform scan: filepath contains space")
// 	}

// 	if mp.scanCompleted(ScanBasic) {
// 		return fmt.Errorf("can't perform scan: scan was already completed")
// 	}
// 	err := mp.ConsumeFile(sourceFile)
// 	if err != nil {
// 		return err
// 	}
// 	if mp.confirmScan(ScanBasic) != nil {
// 		return err
// 	}
// 	return nil
// }

// func filterSilence(report string) (string, []SilenceSegment) {
// 	silenceData := []SilenceSegment{}
// 	lines := strings.Split(report, "\n")
// 	total := 0.0
// 	timestamp := ""
// 	for _, line := range lines {
// 		if strings.Contains(line, "silence_end") {
// 			sil := SilenceSegment{}
// 			flds := strings.Fields(line)
// 			if flds[3] == "silence_end:" {
// 				end, err := strconv.ParseFloat(flds[4], 64)
// 				if err != nil {
// 					fmt.Println(line)
// 					panic(err.Error() + line)
// 				}
// 				sil.SilenceEnd = end
// 			}
// 			if flds[6] == "silence_duration:" {
// 				dur, err := strconv.ParseFloat(flds[7], 64)
// 				if err != nil {
// 					fmt.Println(line)
// 					panic(err.Error() + line)
// 				}
// 				sil.SilenceDuration = dur
// 			}

// 			sil.SilenceStart = sil.SilenceEnd - sil.SilenceDuration
// 			if sil.SilenceStart < 0 {
// 				sil.SilenceStart = 0
// 			}
// 			total += sil.SilenceDuration
// 			silenceData = append(silenceData, sil)
// 		}
// 		if strings.Contains(line, "time=") {
// 			data := strings.Fields(line)
// 			for _, time := range data {
// 				if strings.HasPrefix(line, "time=") {
// 					time = strings.TrimPrefix(time, "time=")
// 					timestamp = time
// 				}
// 			}
// 		}
// 	}
// 	total = float64(int(total*1000)) / 1000
// 	return fmt.Sprintf("%v: %v seconds of silence in %v segments detected", timestamp, total, len(silenceData)), silenceData
// }

// func filterIdet(report string) map[string]int {
// 	rMap := make(map[string]int)
// 	lines := strings.Split(report, "\n")
// 	for _, ln := range lines {
// 		if !strings.Contains(ln, "Parsed_idet") {
// 			continue
// 		}
// 		if strings.Contains(ln, "Repeated ") {
// 			words := strings.Split(ln, " ")
// 			fn := 0
// 			for _, w := range words {
// 				n, err := strconv.Atoi(w)
// 				if err != nil {
// 					continue
// 				}
// 				fn++
// 				rMap["T"] += n
// 				switch fn {
// 				case 2, 3:
// 					rMap["I"] += n
// 				}
// 			}
// 		}
// 		if strings.Contains(ln, "Single ") {
// 			words := strings.Split(ln, " ")
// 			fn := 0
// 			for _, w := range words {
// 				n, err := strconv.Atoi(w)
// 				if err != nil {
// 					continue
// 				}
// 				fn++
// 				rMap["T"] += n
// 				switch fn {
// 				case 1, 2:
// 					rMap["I"] += n
// 				case 3:
// 					rMap["P"] += n
// 				}
// 			}
// 		}
// 		if strings.Contains(ln, "Multi ") {
// 			words := strings.Split(ln, " ")
// 			fn := 0
// 			for _, w := range words {
// 				n, err := strconv.Atoi(w)
// 				if err != nil {
// 					continue
// 				}
// 				fn++
// 				rMap["T"] += n
// 				switch fn {
// 				case 1, 2:
// 					rMap["I"] += n
// 				case 3:
// 					rMap["P"] += n
// 				}
// 			}
// 		}
// 	}
// 	return rMap
// }

// func (mp *mfscanInfo) scanCompleted(scanType string) bool {
// 	for _, scanned := range mp.ScansCompleted {
// 		if scanType == scanned {
// 			return true
// 		}
// 	}
// 	return false
// }

/*
TZ:
Утилита должна оценить файл и составить его рабочий Профиль

//
Input #0, mov,mp4,m4a,3gp,3g2,mj2, from 'Barri_4s_treyler_a_teka.mp4':

	Duration: 00:01:26.08, start: 0.000000, bitrate: 15338 kb/s
	  Stream #0:0(eng): Video: h264 (Main) (avc1 / 0x31637661), yuv420p, 1920x1080 [SAR 1:1 DAR 16:9], 14831 kb/s, 25 fps, 25 tbr, 25k tbn, 50 tbc (default)
	  Stream #0:1(eng): Audio: aac (LC) (mp4a / 0x6134706D), 48000 Hz, stereo, fltp, 317 kb/s (default)

//
[1100]
[1{#1920x1080#25#5822#[]#[[SAR 1:1 DAR 16:9]]#5.3%};2{#5.1#48.0#341}{#stereo#48.0#129};0;1]

ffprobe -v quiet -of json -show_format -show_streams Barri_4s_treyler_a_teka.mp4




Структура Профиля:
краткий:
[ABCDE-F], где
a = количество видеопотоков
b = количество аудиопотоков
c = количество потоков данных
d = количество srt потоков
e = количество хз чего
f = количество замечаний от библиотеки парса

Сепараторы:
- = логический раздел (Абзац)
; = потоки (Строка)
# = данные внутри потоков (Слово)

развернутый:
[A#a#b#c#d#e#f;B#g#h#i;C          ], где
A = количество видеопотоков     формат: eHex                (eHex = целое число с базой 32)          int
a = размер видео                формат: IxI                 (I = целое число)                        int
b = fps видео                   формат: F                   (F = десятичное число)                   float64
c = битрейт                     формат: I                   (I = целое число)                        int
d = SARDAR внешний              формат: [SAR I:I DAR I:I]   (I = целое число)                        []int
e = SARDAR внутренний           формат: [d]                 (d = формат внутреннего SARDAR)          []int
f = интерлейс                   формат: F                   (F = средний % не прогрессивных кадров)  float64
B = количество аудиопотоков     формат: eHex                (eHex = целое число с базой 32)          int
g = раскладка каналов           формат: S                   (S = текст ключ-значение по таблице)     string
h = герцовка                    формат: F                   (F = десятичное число kHz)               float64
i = битрейт                     формат: I                   (I = целое число значение kbit/s)        int
C = количество потоков данных   формат: eHex                (eHex = целое число с базой 32)          int??????????

e = битрейт      формат: [I]   (I = целое число)
d = количество srt потоков
e = количество хз чего
f = количество замечаний от библиотеки парса

Определения:
1.Профиль - форматированная информация о внутренней медиа структуре файла.
2.Аргументы = файлы к которым нужно составить профайл
*/

// validate - проверяет данные на целостность и формирует отчетные строки
func (p *mfscanInfo) validate() error {
	vSNum := 0
	aSNum := 0
	dSNum := 0
	sSNum := 0
	p.streamInfo = make(map[string]string)
	for _, stream := range p.Streams {
		switch stream.Codec_type {
		default:
			return fmt.Errorf("unknown codec type: %v", stream.Codec_type)
		case "video":
			p.validateVideo(stream, vSNum)
			vSNum++
		case "audio":
			p.validateAudio(stream, aSNum)
			aSNum++
		case "data":
			dSNum++
		case "subtitle":
			sSNum++
		}
	}

	switch vSNum {
	case 1:
	case 0:
		if aSNum+sSNum == 0 {
			p.warnings = append(p.warnings, fmt.Sprintf("file ==> no video, audio and subtitle streams detected"))
		}
	default:
		p.warnings = append(p.warnings, fmt.Sprintf("file ==> %v video streams detected", vSNum))
	}
	switch aSNum {
	case 0:
	case 1, 2:
	default:
		p.warnings = append(p.warnings, fmt.Sprintf("file ==> %v audio streams detected", aSNum))
	}
	switch sSNum {
	case 0:
	case 1:
		p.warnings = append(p.warnings, fmt.Sprintf("file ==> %v subtitle stream detected", sSNum))
	default:
		p.warnings = append(p.warnings, fmt.Sprintf("file ==> %v subtitle streams detected", sSNum))
	}
	p.short = fmt.Sprintf("%v%v%v%v-%v", ehex(vSNum), ehex(aSNum), ehex(dSNum), ehex(sSNum), ehex(len(p.warnings)))
	p.combineLong(vSNum, aSNum, dSNum, sSNum)
	p.chanLayout = strings.TrimSuffix(p.chanLayout, " ")
	return nil
}

func (p *mfscanInfo) validateVideo(stream *Stream, vSNum int) {
	currentBlock := fmt.Sprintf("0:v:%v", vSNum)
	vSize := fmt.Sprintf("%vx%v", stream.Width, stream.Height)
	switch vSize {
	case "720x576":
		vSize = "SD"
	case "1920x1080":
		vSize = "HD"
	case "3840x2160":
		vSize = "4K"
	default:
		p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> bad video size: [%v]", vSNum, vSize))
		vSize = "[" + vSize + "]"
	}
	p.streamInfo[currentBlock] = "#" + vSize

	fps := stream.R_frame_rate
	fps_block := fpsToFloat(fps)
	switch fps {
	case "24/1":
	case "25/1":
	case "24000/1001":
	case "2997/125": //about valid
	default:
		p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> bad fps: [%v]", vSNum, fps_block))
	}
	p.streamInfo[currentBlock] += "#" + fmt.Sprintf("%v", fps_block)

	sardar := "SAR=" + stream.Sample_aspect_ratio + " DAR=" + stream.Display_aspect_ratio
	sardar = strings.ReplaceAll(sardar, " ", "_")
	switch sardar {
	case "SAR=1:1_DAR=16:9", "":
	case "SAR=_DAR=":
		p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> blank SAR DAR data present", vSNum))
		sardar = "???"
	case "SAR=1:1_DAR=1024:429", "SAR=1:1_DAR=37:20", "SAR=1:1_DAR=160:67":
		//p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> bad (but possible) SAR DAR: [%v]", vSNum, sardar))
	default:
		p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> bad SAR DAR: [%v]", vSNum, sardar))
	}

	p.streamInfo[currentBlock] += "#" + fmt.Sprintf("[%v]", sardar)

	btrate := stream.Bit_rate
	switch btrate {
	case "":
		p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> blank bitRate data present", vSNum))
		btrate = "???"
	default:
		btr, err := strconv.Atoi(btrate)
		if err != nil {
			panic(err.Error())
		}
		if btr == 0 {
			panic(fmt.Sprintf("stream [0:v:%v] ==> bitRate: [%v]", vSNum, btr))
		}
		//p.warnings = append(p.warnings, fmt.Sprintf("stream [0:v:%v] ==> bitRate: [%v]", vSNum, btr))
		btrate = fmt.Sprintf("%v", btr/1000)
	}
	p.streamInfo[currentBlock] += "#" + fmt.Sprintf("%v", btrate)
	prog := stream.Progressive_frames_pct
	switch prog {
	case 0.0:
		p.streamInfo[currentBlock] += `#ns`
	default:
		p.streamInfo[currentBlock] += `#` + fmt.Sprintf("%v", prog) + "%"
		if prog < 95 && prog > 0 {
			text := fmt.Sprintf("stream [0:v:%v] ==> interlace suspected (%v", vSNum, 100-prog) + "%)"
			p.warnings = append(p.warnings, text)
		}
	}
}

func (p *mfscanInfo) validateAudio(stream *Stream, aSNum int) {
	currentBlock := fmt.Sprintf("0:a:%v", aSNum)
	chan_lay := stream.Channel_layout
	channel_num := stream.Channels
	switch chan_lay {
	default:
		p.warnings = append(p.warnings, fmt.Sprintf("stream [0:a:%v] ==> unknown channel layout provided (%v)", aSNum, chan_lay))
		chan_lay = chan_lay + ":" + fmt.Sprintf("%vch", channel_num)
	case "":
		switch stream.Channels {
		default:
			p.warnings = append(p.warnings, fmt.Sprintf("stream [0:a:%v] ==> no channel layout provided (can not guess: %v channels)", aSNum, channel_num))
			chan_lay = "???"
		case 1:
			chan_lay = "*mono"
		case 2:
			chan_lay = "*stereo"
		case 6:
			chan_lay = "*5.1"
		}

	case "1 channels (FL)", "1 channels (LFE)", "1 channels (BL)", "1 channels (FR)", "1 channels (BR)", "5 channels (FL+FR+LFE+SL+SR)", "downmix":
		p.warnings = append(p.warnings, fmt.Sprintf("stream [0:a:%v] ==> unusual channel layout provided (%v)", aSNum, chan_lay))
		chan_lay = "*:" + fmt.Sprintf("%vch", channel_num)
	case "5.1":
	case "5.1(side)":
		chan_lay = "5.1"
	case "mono":
	case "stereo":
	}

	p.streamInfo[currentBlock] += "#" + fmt.Sprintf("%v", chan_lay)
	p.ChanNum[currentBlock] = fmt.Sprintf("%v", channel_num)
	p.chanLayout += fmt.Sprintf("%v ", ehex(channel_num))

	switch channel_num {
	case 1, 2, 6:
	default:
		p.warnings = append(p.warnings, fmt.Sprintf("stream [0:a:%v] ==> unusual number of channels (%v)", aSNum, channel_num))
		//p.streamInfo[currentBlock] += ":" + fmt.Sprintf("%vch", ehex(channel_num))
	}

	hz := stream.Sample_rate
	switch hz {
	default:
		p.warnings = append(p.warnings, fmt.Sprintf("stream [0:a:%v] ==> sample rate [%v Hz]", aSNum, hz))
	case "48000":
	}
	p.streamInfo[currentBlock] += "#" + hzFormat(hz)

	bitRt := stream.Bit_rate
	bits, err := strconv.Atoi(bitRt)
	if err != nil {
		switch bitRt {
		case "":
			p.warnings = append(p.warnings, fmt.Sprintf("stream [0:a:%v] ==> no bitrate provided", aSNum))
			p.streamInfo[currentBlock] += "#???"
			return
		}
		if bitRt != "" {
			//panic(fmt.Sprintf("stream [0:a:%v] ==> bad bitrate  (%v): %v", aSNum, bitRt, err.Error()))
			p.warnings = append(p.warnings, fmt.Sprintf("stream [0:a:%v] ==> bad bitrate  (%v): %v", aSNum, bitRt, err.Error()))
		}
	}
	switch bits {
	case 0:
	default:
		if bits < 80000 {
			p.warnings = append(p.warnings, fmt.Sprintf("stream [0:a:%v] ==> silence suspected: bitrate is extreamly low [%v b/s]", aSNum, bitRt))
		}
	}
	p.streamInfo[currentBlock] += "#" + fmt.Sprintf("%v", bits/1000)

}

func (p *mfscanInfo) combineLong(vSNum, aSNum, dSNum, sSNum int) {
	p.long = ""
	for _, stTp := range []string{"v", "a", "d", "s"} {
		switch stTp {
		case "v":
			p.long += fmt.Sprintf("%v", vSNum)
		case "a":
			p.long += fmt.Sprintf("%v", aSNum)
		case "d":
			p.long += fmt.Sprintf("%v", dSNum)
		case "s":
			p.long += fmt.Sprintf("%v", sSNum)
		}
		//[1{#1920x1080#25#5822#[]#[[SAR 1:1 DAR 16:9]]#5.3%};2{#5.1#48.0#341}{#stereo#48.0#129};0;1]
		for i := 0; i < 50; i++ {
			if val, ok := p.streamInfo[fmt.Sprintf("0:%v:%v", stTp, i)]; ok {
				p.long += fmt.Sprintf("%v%v{%v}", stTp, i, val)
			}
		}
		p.long += ";"
	}
	//p.long = strings.TrimSuffix(p.long, ";")
	p.long += fmt.Sprintf("w%v", len(p.warnings))
	p.long = strings.ReplaceAll(p.long, " ", "_")
}

func ehex(i int) string {
	switch i {
	default:
		return "?"
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
		return fmt.Sprintf("%v", i)
	case 10:
		return "A"
	case 11:
		return "B"
	case 12:
		return "C"
	case 13:
		return "D"
	case 14:
		return "E"
	case 15:
		return "F"
	case 16:
		return "G"
	case 17:
		return "H"
	case 18:
		return "J"
	case 19:
		return "K"
	case 20:
		return "L"
	case 21:
		return "M"
	case 22:
		return "N"
	case 23:
		return "P"
	case 24:
		return "Q"
	case 25:
		return "R"
	case 26:
		return "S"
	case 27:
		return "T"
	case 28:
		return "U"
	case 29:
		return "V"
	case 30:
		return "W"
	case 31:
		return "X"
	case 32:
		return "Y"
	case 33:
		return "Z"
	}
}

func (p *mfscanInfo) Warnings() []string {
	return p.warnings
}

func (p *mfscanInfo) Short() string {
	return p.short
}

func (p *mfscanInfo) Long() string {
	return p.long
}

func (p *mfscanInfo) ByStream() map[string]string {
	return p.streamInfo
}

func (p *mfscanInfo) AudioLayout() string {
	return p.chanLayout
}

func fpsToFloat(fps string) float64 {
	data := strings.Split(fps, "/")
	i1, _ := strconv.Atoi(data[0])
	i2, _ := strconv.Atoi(data[1])
	fl := float64(i1) / float64(i2)
	fli := float64(int(fl*1000)) / 1000
	return fli
}

func hzFormat(hz string) string {
	h, err := strconv.Atoi(hz)
	if err != nil {
		return "???"
	}
	return fmt.Sprintf("%v", float64(h)/1000)
}

// func (mp *mfscanInfo) confirmScan(scan string) error {
// 	if mp.scanCompleted(scan) {
// 		return fmt.Errorf("confirmation failed: was already confirmed")
// 	}
// 	switch scan {
// 	default:
// 		return fmt.Errorf("can't confirm %v scan: unknown or unimplemented scan type")
// 	case ScanBasic:
// 		if len(mp.ScansCompleted) != 0 {
// 			return fmt.Errorf("can't confirm %v scan: other scan data must not exist")
// 		}
// 	case ScanInterlace, ScanSilence:
// 		if !mp.scanCompleted(ScanBasic) {
// 			return fmt.Errorf("can't confirm %v scan: basic scan not completed")
// 		}

// 	}
// 	if err := mp.validate(); err != nil {
// 		return fmt.Errorf("can't validate profile: %v", err.Error())
// 	}
// 	mp.ScansCompleted = appendUnique(mp.ScansCompleted, scan)
// 	return nil
// }

func appendUnique(slice []string, newElem string) []string {
	for i := range slice {
		if slice[i] == newElem {
			return slice
		}
	}
	return append(slice, newElem)
}

// func (sr *mfscanInfo) String() string {
// 	s := "file: " + sr.Format.Filename + "\n"
// 	s += fmt.Sprintf("streams total: %v", len(sr.Streams))
// 	return s
// }

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
