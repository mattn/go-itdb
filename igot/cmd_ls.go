package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/gonuts/commander"
	"github.com/mattn/go-itdb"
)

func make_cmd_ls(iPod *itdb.IPod) *commander.Command {
	cmd_ls := func(cmd *commander.Command, args []string) error {
		tracks, err := iPod.Tracks()
		if err != nil {
			return err
		}
		for _, t := range tracks {
			ct.ChangeColor(ct.Yellow, false, ct.None, false)
			fmt.Print(t.Id)
			ct.ResetColor()

			fmt.Print(" ")

			ct.ChangeColor(ct.Green, false, ct.None, false)
			fmt.Print(t.Title)
			ct.ResetColor()

			fmt.Print(" ")

			ct.ChangeColor(ct.Cyan, true, ct.None, false)
			fmt.Println(t.Artist)
			ct.ResetColor()
		}
		return nil
	}

	return &commander.Command{
		Run:       cmd_ls,
		UsageLine: "ls [options]",
		Short:     "list tracks",
	}
}
