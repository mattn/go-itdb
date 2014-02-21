package main

import (
	"fmt"
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
			fmt.Println(t.Title, t.Artist)
		}
		return nil
	}

	return &commander.Command{
		Run:       cmd_ls,
		UsageLine: "ls [options]",
		Short:     "list tracks",
	}
}
