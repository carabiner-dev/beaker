// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc

package golang

import (
	"context"

	"github.com/carabiner-dev/beaker/pkg/runners/shell"
)

type Options struct {
}

type OptFn func(*Options) error

// New returns a new go runner
func New(funcs ...OptFn) (*Runner, error) {
	shellrunner, err := shell.New(
		shell.WithCommand("go"),
		shell.WithArguments([]string{"test", "./...", "-json"}),
	)
	if err != nil {
		return nil, err
	}
	return &Runner{
		runner: shellrunner,
	}, nil
}

// Runner implements a Test runner to execute go tests
type Runner struct {
	Options Options
	runner  *shell.Runner
}

// Run runs the tests
func (r *Runner) Run(ctx context.Context) ([]byte, bool, error) {
	return r.runner.Run(ctx)
}
