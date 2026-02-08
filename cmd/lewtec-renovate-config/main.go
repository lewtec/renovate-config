package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	preset       string
	noPr         bool
	includeForks bool
)

const (
	defaultPreset = "github>lewtec/renovate-config:base"
)

func main() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(handler))

	var rootCmd = &cobra.Command{
		Use:   "renovate-config [owner] [repo]",
		Short: "Add renovate-config preset to a GitHub repository",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			owner := args[0]

			ghAvailable, ghError := checkGhCli()
			if !ghAvailable {
				slog.Warn("GitHub CLI warning", "warning", ghError)
				if !noPr {
					slog.Warn("PR creation will not be available. Use --no-pr to suppress this warning")
				}
			}

			if len(args) == 2 {
				repo := args[1]
				os.Exit(processRepository(owner, repo, preset, noPr, ghAvailable, ghError))
			} else {
				if !ghAvailable {
					ReportError(fmt.Errorf(ghError), "GitHub CLI not available")
					os.Exit(1)
				}
				repos, err := getRepositories(owner, includeForks)
				if err != nil {
					ReportError(err, "Error fetching repos")
					os.Exit(1)
				}
				slog.Info("Found repositories", "count", len(repos), "repos", strings.Join(repos, ", "))

				bar := progressbar.Default(int64(len(repos)))

				failures := 0
				for _, repo := range repos {
					if res := processRepository(owner, repo, preset, noPr, ghAvailable, ghError); res != 0 {
						failures++
					}
					bar.Add(1)
				}
				if failures > 0 {
					slog.Error("Finished with failures", "failures", failures, "total", len(repos))
					os.Exit(1)
				}
				slog.Info("Successfully processed all repositories", "count", len(repos))
			}
		},
	}

	rootCmd.Flags().StringVar(&preset, "preset", defaultPreset, "Preset reference")
	rootCmd.Flags().BoolVar(&noPr, "no-pr", false, "Do not create a PR, just push the branch")
	rootCmd.Flags().BoolVar(&includeForks, "include-forks", false, "Include forked repositories when processing all repos")

	if err := rootCmd.Execute(); err != nil {
		ReportError(err, "Execution error")
		os.Exit(1)
	}
}
