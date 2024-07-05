package ticket

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Galdoba/ffproc/internal/pkg/define"
)

func TestTicket(t *testing.T) {
	source := `Garri_uayld_s01e01--SER--Garri_uayld_s01e01_PRT240627183618_SER_00915_18.mp4`
	mfline := `1v0{#HD#25#[SAR=1:1_DAR=16:9]#8007#ns};2a0{#5.1#48#341}a1{#stereo#48#130};1;0;w0`
	srt := `Garri_uayld_s01e01--SER--Garri_uayld_s01e01_PRT240627183618_SER_00915_18.RUS.srt`
	tk := New(`Гарри Уайлд. 01 сезон. 01 серия (Замена)`, define.PROJ_TYPE_SER)
	tk.AddSource(source, mfline)
	tk.AddSource(srt, `0;0;0;1;w1`)
	tk.AddTag(source, "{v:0}:size", "HD")
	tk.AddTag(source, "{a:0}:channels", "6")
	tk.AddTag(source, "{a:0}:layout", "51")
	tk.AddTag(source, "{a:0}:lang", "rus")
	tk.AddTag(source, "{a:1}:channels", "2")
	tk.AddTag(source, "{a:1}:layout", "stereo")
	tk.AddTag(source, "{a:1}:lang", "eng")
	tk.AddTag(source, "{a:0}:atempo", "25/(25)")
	tk.AddTag(source, "{a:1}:atempo", "25/(25)")
	tk.AddRequest("Langs", "RUS ENG")
	tk.AddRequest("Revision", "1")
	tk.AddRequest("Video_Scale", "HD")
	tk.AddTag(srt, "{s:0}:type", "hardsub")
	tk.AddTag(PROCESS_DATA, PROCESS_DEST, `//nas/ROOT/EDIT/_amedia/Garri_uayld_s01/`)
	tk.AddTag(PROCESS_DATA, ARCHIVE_DEST, `//192.168.31.4/root/IN/_AMEDIA/_DONE/Garri_uayld_s01/`)
	tk.AddTag(PROCESS_DATA, EPISODES_EXPECTED, `01`)
	tk.AddTag(PROCESS_DATA, SEASON_EXPECTED, `01`)
	tk.AddTag(PROCESS_REQUEST, "Langs", `RUS ENG`)
	tk.AddTag(source, SEASON_NUM, `01`)
	tk.AddTag(source, EPISODE_NUM, `01`)
	tk.AddTag(srt, SEASON_NUM, `01`)
	tk.AddTag(srt, EPISODE_NUM, `01`)
	tk.Processable = true
	bt, _ := json.MarshalIndent(tk, "", "  ")
	fmt.Println(string(bt))
	if err := tk.ValidateWith(&pseudoValidator{}); err != nil {
		fmt.Println(err.Error())
	}

}
