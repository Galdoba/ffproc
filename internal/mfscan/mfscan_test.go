package mfscan

import (
	"fmt"
	"testing"
)

func inputs() []string {
	return []string{
		// `C:\Users\Admin\Downloads\1.mp4`,
		// `c:\Users\Admin\TorrentDownloads\out.mp4`,
		`\\192.168.31.4\buffer\IN\test2_s01e02--SER--some_name2.mp4`,
		`\\192.168.31.4\buffer\IN\test_s01e02--SER--some_name.mp4`,
	}
}

func TestSize(t *testing.T) {
	for _, inp := range inputs() {
		fmt.Println(inp)
		scaned := NewProfile()
		scaned.ConsumeFile(inp)
		size := scaned.Size()
		fmt.Println("size:", size)
		fps := scaned.Fps()
		fmt.Println("fps:", fps)
		layout := scaned.Layout()
		fmt.Println("layout:", layout)
		channels := scaned.Channels()
		fmt.Println("channels:", channels)
		langs := scaned.Langs()
		fmt.Println("langs:", langs)
		streamTypes := scaned.StreamTypes()
		fmt.Println("streamTypes:", streamTypes)
		if len(size) < 1 {
			fmt.Printf("input %v: no video stream detected\n", inp)
		}
		for _, sz := range size {
			switch sz {
			case "640x360", "1920x1080":
			default:
				t.Errorf("unexpected size %v", sz)
			}

		}
	}
	fmt.Println("=============")
}
