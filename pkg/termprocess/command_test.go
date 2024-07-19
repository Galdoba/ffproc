package termprocess

import (
	"fmt"
	"os"
	"testing"
)

func Test_Command(t *testing.T) {
	input := `\\192.168.31.4\buffer\IN\_DONE\Stalnaya_hvatka--FILM--IronClaw_2K_prores422hq.mov`
	tcmd, err := New(
		CommandLineArguments(fmt.Sprintf("ffprobe -v verbose -f lavfi -i amovie=%v,asetnsamples=48000,astats=metadata=1:reset=1 -show_entries frame=pkt_pts_time:frame_tags=lavfi.astats.Overall.RMS_level,lavfi.astats.1.RMS_level,lavfi.astats.2.RMS_level,lavfi.astats.3.RMS_level,lavfi.astats.4.RMS_level,lavfi.astats.5.RMS_level,lavfi.astats.6.RMS_level -of csv=p=0", input)),
		Set(TERMINAL_ON),
		WriteToFile(`\\192.168.31.4\buffer\IN\_DONE\log3.txt`),
		AddBuffer("aaa"),
		//\\192.168.31.4\buffer\IN\_DONE\
	)
	bf := tcmd.Buffer("aaa")
	if err != nil {
		t.Errorf("func New() error: %v", err.Error())
	}
	if err := tcmd.Run(); err != nil {
		t.Errorf("func Run() error: %v", err.Error())
	}

	fmt.Println("///")
	fmt.Println("OUT")
	fmt.Println(tcmd.StdOut())
	fmt.Println("///")
	fmt.Println("ERR")
	fmt.Println(tcmd.stErr)
	ex, err := os.Executable()
	fmt.Println("launch position:", ex, err)

	hn, err := os.Hostname()
	fmt.Println("Host Name:", hn, err)
	fmt.Println("///")
	fmt.Println("BUF")
	fmt.Println(bf.String())

}
