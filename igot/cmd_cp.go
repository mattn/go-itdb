package main

import (
	"github.com/gonuts/commander"
	"github.com/mattn/go-itdb"
)

func make_cmd_cp(iPod *itdb.IPod) *commander.Command {
	cmd_cp := func(cmd *commander.Command, args []string) error {
		for _, f := range args {
			err := iPod.CopyTrack(f)
			if err != nil {
				return err
			}
		}
		return iPod.Write()
	}

	return &commander.Command{
		Run:       cmd_cp,
		UsageLine: "cp [file...]",
		Short:     "copy track into iPod",
	}
}
