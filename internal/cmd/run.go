// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/util"

	"github.com/carabiner-dev/beaker/pkg/beaker"
)

type runOptions struct {
	configFile string
	workDir    string
	attest     bool
	outputPath string
}

// Validates the options in context with arguments
func (ro *runOptions) Validate() error {
	errs := []error{}
	if !util.IsDir(ro.workDir) {
		errs = append(errs, errors.New("working directory does not exist"))
	}

	if ro.outputPath == "" {
		errs = append(errs, errors.New("output path is required"))
	}
	return errors.Join(errs...)
}

// AddFlags adds the subcommands flags
func (ro *runOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&ro.configFile, "runner", "r", "", "test runner to configure",
	)
	cmd.PersistentFlags().StringVarP(
		&ro.workDir, "dir", "d", ".", "path to codebase",
	)
	cmd.PersistentFlags().StringVarP(
		&ro.configFile, "config", "c", ".beaker.yaml", "path to configuration file",
	)
	cmd.PersistentFlags().BoolVarP(
		&ro.attest, "attest", "a", true, "output the entire in-toto statement (instead of predicate)",
	)
	cmd.PersistentFlags().StringVarP(
		&ro.outputPath, "output", "o", "tests.intoto.json", "path to file to write the predicate or attestation",
	)
}

func addRun(parentCmd *cobra.Command) {
	opts := &runOptions{}
	attCmd := &cobra.Command{
		Short:             "executes a test runner and captures the results",
		Use:               "run",
		SilenceUsage:      false,
		SilenceErrors:     true,
		PersistentPreRunE: initLogging,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if opts.workDir == "." || opts.workDir == "" {
					opts.workDir = args[0]
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Validate the options
			if err := opts.Validate(); err != nil {
				return err
			}
			cmd.SilenceUsage = true

			f, err := os.Create(opts.outputPath)
			if err != nil {
				return fmt.Errorf("opening file: %w", err)
			}

			defer func() {
				f.Close() //nolint:errcheck,gosec
			}()

			launcher, err := beaker.New(
				beaker.WithAttest(opts.attest),
				beaker.WithWorkDir(opts.workDir),
			)
			if err != nil {
				return fmt.Errorf("creating launcher")
			}

			pack, err := beaker.LaunchPackFromRepo(opts.workDir)
			if err != nil {
				return fmt.Errorf("automatically building launchpack: %w", err)
			}

			return launcher.Test(context.Background(), pack)
		},
	}
	opts.AddFlags(attCmd)
	parentCmd.AddCommand(attCmd)
}
