// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc

package shell

import (
	"bytes"
	"context"
	"fmt"

	"sigs.k8s.io/release-utils/command"
)

type Options struct {
	WorkDir string
	Command string
	Args    []string
	Env     map[string]string
}

type OptFn func(*Options) error

func WithWorkDir(path string) OptFn {
	return func(o *Options) error {
		o.WorkDir = path
		// TODO: CHeck path
		return nil
	}
}

func WithCommand(cmdline string) OptFn {
	return func(o *Options) error {
		o.Command = cmdline
		return nil
	}
}

func WithArguments(args []string) OptFn {
	return func(o *Options) error {
		o.Args = args
		// TODO: Check args
		return nil
	}
}

func WithEnv(env map[string]string) OptFn {
	return func(o *Options) error {
		o.Env = env
		return nil
	}
}

// New returns a new shell runner configured with the passed options
func New(funcs ...OptFn) (*Runner, error) {
	opts := Options{
		Args: []string{},
		Env:  map[string]string{},
	}
	for _, f := range funcs {
		if err := f(&opts); err != nil {
			return nil, err
		}
	}

	return &Runner{
		Options: opts,
	}, nil
}

type Runner struct {
	Options Options
}

// Run runs the tests
func (r *Runner) Run(context.Context) ([]byte, bool, error) {
	var cmd *command.Command
	if r.Options.WorkDir == "" {
		cmd = command.New(
			r.Options.Command,
			r.Options.Args...,
		)
	} else {
		cmd = command.NewWithWorkDir(
			r.Options.WorkDir,
			r.Options.Command,
			r.Options.Args...,
		)
	}

	if len(r.Options.Env) > 0 {
		var envs = []string{}
		for k, val := range r.Options.Env {
			envs = append(envs, fmt.Sprintf("%s=%s", k, val))
		}
		cmd = cmd.Env(envs...)
	}

	var b bytes.Buffer
	cmd = cmd.AddOutputWriter(&b)

	status, err := cmd.Run()
	if err != nil {
		return nil, false, fmt.Errorf("shelling out to command: %w", err)
	}

	return b.Bytes(), status.Success(), nil
}
