// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc

package golang

import (
	"context"
	"fmt"

	"github.com/carabiner-dev/beaker/pkg/runners/shell"
	"sigs.k8s.io/release-utils/util"
)

type Options struct {
	WorkDir string
}

func WithWorkDir(path string) OptFn {
	return func(o *Options) error {
		if !util.IsDir(path) {
			return fmt.Errorf("workind dir does not exist: %q", path)
		}
		o.WorkDir = path
		return nil
	}
}

type OptFn func(*Options) error

// New returns a new go runner
func New(funcs ...OptFn) (*Runner, error) {
	opts := Options{
		WorkDir: ".",
	}

	for _, f := range funcs {
		if err := f(&opts); err != nil {
			return nil, err
		}
	}
	shellrunner, err := shell.New(
		shell.WithWorkDir(opts.WorkDir),
		shell.WithCommand("go"),
		shell.WithArguments([]string{"test", "-json", "./..."}),
	)
	if err != nil {
		return nil, err
	}
	return &Runner{
		Options: opts,
		runner:  shellrunner,
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
