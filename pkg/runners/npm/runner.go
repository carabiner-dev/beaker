// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2026 Carabiner Systems, Inc

package npm

import (
	"context"
	"fmt"

	"sigs.k8s.io/release-utils/helpers"

	"github.com/carabiner-dev/beaker/pkg/runners/shell"
)

type Options struct {
	WorkDir string
}

func WithWorkDir(path string) OptFn {
	return func(o *Options) error {
		if !helpers.IsDir(path) {
			return fmt.Errorf("working dir does not exist: %q", path)
		}
		o.WorkDir = path
		return nil
	}
}

type OptFn func(*Options) error

// New returns a new npm runner
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
		shell.WithCommand("npm"),
		shell.WithArguments([]string{"test", "--", "--reporter=tap"}),
	)
	if err != nil {
		return nil, err
	}
	return &Runner{
		Options: opts,
		runner:  shellrunner,
	}, nil
}

// Runner implements a TestRunner that executes `npm test` and expects
// the underlying test framework to emit TAP on stdout.
type Runner struct {
	Options Options
	runner  *shell.Runner
}

// Run runs the tests
func (r *Runner) Run(ctx context.Context) (attestation []byte, pass bool, err error) {
	return r.runner.Run(ctx)
}
