package mfscan

import (
	"fmt"
	"testing"
)

func inputs() []string {
	return []string{
		`C:\Users\Admin\Downloads\1.mp4`,
		`c:\Users\Admin\TorrentDownloads\out.mp4`,
	}
}

func TestSize(t *testing.T) {
	for _, inp := range inputs() {
		scaned := NewProfile()
		scaned.ConsumeFile(inp)
		size := scaned.Size()
		fmt.Println("size:", size)
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
}
