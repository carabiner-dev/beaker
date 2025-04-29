// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/carabiner-dev/beaker/pkg/beaker"
	"github.com/carabiner-dev/beaker/pkg/runners/golang"
	"github.com/spf13/cobra"
)

type runOptions struct {
	configFile string
}

// Validates the options in context with arguments
func (ro *runOptions) Validate() error {
	errs := []error{}
	// if to.SpecPath == "" {
	// 	errs = append(errs, errors.New("spec path not defined"))
	// }

	// for _, val := range to.VarSubstitutions {
	// 	if !strings.Contains(val, "=") {
	// 		errs = append(errs, fmt.Errorf("variable substitution not well formed: %q", val))
	// 	}
	// }
	return errors.Join(errs...)
}

// AddFlags adds the subcommands flags
func (ro *runOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&ro.configFile, "runner", "r", "", "test runner to configure",
	)
	cmd.PersistentFlags().StringVarP(
		&ro.configFile, "config", "c", ".beaker.yaml", "path to configuration file",
	)
}

func addRun(parentCmd *cobra.Command) {
	opts := &runOptions{}
	attCmd := &cobra.Command{
		Short: "executes a test runner and captures the results",
		Use:   "run",
		// Example:           fmt.Sprintf(`%s snap --var REPO=example spec.yaml`, appname),
		SilenceUsage:      false,
		SilenceErrors:     true,
		PersistentPreRunE: initLogging,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Validate the options
			if err := opts.Validate(); err != nil {
				return err
			}
			cmd.SilenceUsage = true

			launcher, err := beaker.New()
			if err != nil {
				return fmt.Errorf("creating launcher")
			}

			gorunner, err := golang.New()
			if err != nil {
				return fmt.Errorf("creating runner: %w", err)
			}

			return launcher.Test(
				context.Background(),
				&beaker.LaunchPack{
					Runner: gorunner,
					Parser: gorunner,
				},
			)
		},
	}
	opts.AddFlags(attCmd)
	parentCmd.AddCommand(attCmd)
}
