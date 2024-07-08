package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Galdoba/ffproc/configs"
	"github.com/Galdoba/ffproc/db/spreadsheet"
	"github.com/Galdoba/ffproc/internal/bridge"
	"github.com/Galdoba/ffproc/internal/pkg/process"
	"github.com/Galdoba/ffproc/internal/pkg/sourcefiles"
	"github.com/Galdoba/ffproc/internal/pkg/table"
	"github.com/Galdoba/ffproc/internal/pkg/ticket"
	"github.com/Galdoba/ffproc/pkg/survey"
	"github.com/urfave/cli/v2"
)

var cfg *configs.Procman

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
			cfg = &configs.Procman{}
			cfgPath := configs.ConfigPath(c.App.Name, "dev")
			// cfg.Link = "https://docs.google.com/spreadsheets/d/1Waa58usrgEal2Da6tyayaowiWujpm0rzd06P5ASYlsg/edit?gid=250314867#gid=250314867"
			// cfg.Path = "c:/Users/pemaltynov/.galdoba/ffproc/db/db.csv"
			// cfg.TicketStorage = "c:/Users/pemaltynov/.galdoba/ffproc/tickets/"
			// bt, _ := json.MarshalIndent(cfg, "", "  ")
			// fmt.Println(string(bt))
			// fmt.Println(cfgPath)
			bt, err := os.ReadFile(cfgPath)
			// fmt.Println(string(bt))
			if err != nil {
				return fmt.Errorf("read config failed: %v", err)
			}

			err = json.Unmarshal(bt, cfg)
			if err != nil {
				return fmt.Errorf("unmarshal config failed: %v", err)
			}
			fmt.Println(cfg)
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
			////////////
			fmt.Println("update tabledata...")
			// fmt.Println(cfg)
			db, err := spreadsheet.New(cfg.Path)
			if err != nil {
				return fmt.Errorf("db.New: %v", err)
			}
			if err := db.CurlUpdate(cfg.Link); err != nil {
				return fmt.Errorf("db.Update: %v", err)
			}
			tableCompiled, err := table.CompileTableData(db)
			if err != nil {
				return err
			}
			///////////////////
			fmt.Println("define source projects...")
			projects := sourcefiles.SplitByKeys(sources)
			fmt.Println("connect source projects with table...")
			processList := process.New()
			for projectKey, projData := range projects {
				for _, tableData := range tableCompiled.Entries {

					if commonKey(projectKey, table.Key(tableData)) == "" {
						continue
					}
					if sourcefiles.ProcessType(projData) != tableData.ProcessType {
						continue
					}

					br := bridge.New(tableData, projData)

					fmt.Println("create ticket")
					tkName := tableData.ProcessType + "--" + projectKey
					fmt.Println("create ticket", tkName)
					tkPath := cfg.TicketStorage + tkName + ".json"
					tk, err := ticket.Load(tkPath)
					switch err {
					case ticket.NoTicket:
						tk = ticket.New(projectKey, tableData.ProcessType)
					case nil:
						fmt.Println("load existing:", tkPath)
					default:
						return fmt.Errorf("ticket.Load: %v", err)
					}
					ticket.Save(tk, cfg.TicketStorage)
					// if err := ticket.Save(tk, cfg.TicketStorage); err != nil {
					// 	fmt.Printf("error Save Ticket: %v")
					// }

					// fmt.Println("----")
					// fmt.Println("PROCESS TICKET:")
					// fmt.Println(tk, br)
					if err := processList.AddProject(tk, br, tkPath); err != nil {
						return fmt.Errorf("AddProject: %v", err)
					}
				}
			}
			fmt.Println("---------------")
			fmt.Println(processList)
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

func commonKey(projectKey, tableKey string) string {
	if strings.HasPrefix(tableKey, projectKey) {
		return projectKey
	}
	return ""
}
