package configs

import (
	"log"
	"os"
	"path/filepath"
)

func ConfigPath(app string, cfgType ...string) string {
	cfgSuffix := "default"
	for _, ct := range cfgType {
		cfgSuffix = ct
	}
	sep := string(filepath.Separator)
	switch cfgSuffix {
	default:
		log.Fatalf("unknown config type: %v", cfgSuffix)
	case "dev":
		root := os.Getenv("GOPATH")
		return root + sep + "src" + sep + "github.com" + sep + "Galdoba" + sep + "ffproc" + sep + "configs" + sep + "procman-dev.json"
	case "default":
		root, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		return root + sep + ".config" + sep + app + sep + "config.json"
	}
	return "???"
}
