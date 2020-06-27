package main

import (
	"fmt"
	"os"

	"github.com/porgull/go-search/pkg/algorithms"
	"github.com/porgull/go-search/pkg/environments"
	"github.com/porgull/go-search/pkg/search"
	"github.com/spf13/cobra"
)

type runFlagsCfg struct {
	on        string
	load      string
	with      string
	visualize bool
}

var (
	runFlags = &runFlagsCfg{}
)

var (
	runCmd = &cobra.Command{
		Use:   "run (--on <environment>|--load <env.json>) --with <algorithm>",
		Short: "run allows you to run and print diagnostics about the perfomance of a search algorithm.",
		Long:  "run allows you to run and print diagnostics about the perfomance of a search algorithm.",
		Run: func(cmd *cobra.Command, args []string) {
			var env environments.Environment
			var algo algorithms.Algorithm
			var err error

			if runFlags.on == "" && runFlags.load == "" {
				cmd.Help()
				os.Exit(1)
			} else if runFlags.on != "" && runFlags.load != "" {
				cmd.Help()
				os.Exit(1)
			} else if runFlags.on != "" {
				env, err = environments.GetEnvironment(runFlags.on)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not get pre-made environment %s: %s\n", runFlags.on, err.Error())
					os.Exit(1)
				}
			} else if runFlags.load != "" {
				f, err := os.Open(runFlags.load)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not open env file at %s: %s", runFlags.load, err.Error())
					os.Exit(1)
				}

				env, err = environments.LoadEnvironmentFrom(f)
				if err != nil {
					f.Close()
					fmt.Fprintf(os.Stderr, "Could not load environment from %s: %s\n", runFlags.load, err.Error())
					os.Exit(1)
				}
				f.Close()
			}

			if runFlags.with == "" {
				cmd.Help()
				os.Exit(1)
			} else {
				algo, err = algorithms.GetAlgorithm(runFlags.with)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Could not get algorithm %s: %s\n", runFlags.with, err.Error())
					os.Exit(1)
				}
			}

			sctx := search.Context{
				Visualize: runFlags.visualize,
			}

			result, err := algo.Run(sctx, env)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while running algorithm %s on %s: %s\n", runFlags.with, env.Name(), err.Error())
				os.Exit(1)
			}

			result.Print()
		},
	}
)

func init() {
	runCmd.PersistentFlags().StringVar(&runFlags.on, "on", "", "Use this pre-created environment to run the search algorithm")
	runCmd.PersistentFlags().StringVar(&runFlags.load, "load", "", "Load your own environment into memory")
	runCmd.PersistentFlags().StringVar(&runFlags.with, "with", "", "Algorithm to use to search")
	runCmd.PersistentFlags().BoolVar(&runFlags.visualize, "visualize", false, "Visualize the search process")
}

func init() {
	rootCmd.AddCommand(runCmd)
}
