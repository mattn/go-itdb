package main

import (
	"github.com/gonuts/commander"
	"github.com/mattn/go-itdb"
	"strconv"
)

func make_cmd_rm(iPod *itdb.IPod) *commander.Command {
	cmd_rm := func(cmd *commander.Command, args []string) error {
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				return err
			}
			err = iPod.RemoveTrack(id)
			if err != nil {
				return err
			}
		}
		return iPod.Write()
	}

	return &commander.Command{
		Run:       cmd_rm,
		UsageLine: "rm [id...]",
		Short:     "remove track from iPod",
	}
}
