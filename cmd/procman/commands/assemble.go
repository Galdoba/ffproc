package commands

import (
	"fmt"
	"os"
	"sort"

	"github.com/Galdoba/ffproc/configs"
	"github.com/Galdoba/ffproc/internal/pkg/sourcefiles"
	"github.com/Galdoba/ffproc/pkg/survey"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var cfg configs.Procman

func Assemble() *cli.Command {
	cm := cli.Command{
		Name:        "assemble",
		Aliases:     []string{},
		Usage:       "creates command for ffmpeg for processing file(s)",
		UsageText:   "",
		Description: "",
		Args:        false,
		ArgsUsage:   "",
		Category:    "",
		BashComplete: func(*cli.Context) {
		},
		Before: func(c *cli.Context) error {

			cfgPath := configs.ConfigPath(c.App.Name, "dev")
			bt, err := os.ReadFile(cfgPath)
			if err != nil {
				return fmt.Errorf("read config failed: %v", err)
			}
			err = yaml.Unmarshal(bt, cfg)
			if err != nil {
				return fmt.Errorf("unmarshal config failed: %v", err)
			}
			return nil
		},
		After: func(*cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Start action: assemble")
			fmt.Println("check arguments: file")
			args := c.Args().Slice()
			sort.Strings(args)
			if err := survey.Files(args...); err != nil {
				return fmt.Errorf("survey failed: %v", err)
			}
			sources, err := sourcefiles.Inspect(args)
			if err != nil {
				return err
			}
			projects := sourcefiles.SplitByKeys(sources)
			for k, v := range projects {
				fmt.Println(k, ":")
				for _, src := range v {
					fmt.Println(" ", src)
				}
			}

			// db, err := spreadsheet.New("path")
			// if err != nil {
			// 	return err
			// }
			// db.CurlUpdate("url")
			// table, err := table.CompileTableData(db)
			// if err != nil {
			// 	return err
			// }

			/*
				inspectSourceFiles
				collectRequests
				Process
			*/

			fmt.Println("TODO: check arguments: scan sources")
			fmt.Println("TODO: check arguments: user prompt if needed")
			fmt.Println("TODO: check arguments: assemble")
			fmt.Println("TODO: check arguments: output")
			return nil
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			return nil
		},
		Subcommands:            []*cli.Command{},
		Flags:                  []cli.Flag{},
		SkipFlagParsing:        false,
		HideHelp:               false,
		HideHelpCommand:        false,
		Hidden:                 false,
		UseShortOptionHandling: false,
		HelpName:               "",
		CustomHelpTemplate:     "",
	}
	return &cm
}
