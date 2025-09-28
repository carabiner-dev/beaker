// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc

package beaker

import (
	"errors"
	"fmt"
	"io"

	"sigs.k8s.io/release-utils/helpers"
)

type OptFn func(*Options) error

type Options struct {
	Writer  io.Writer
	WorkDir string
	Attest  bool
}

func WithWriter(w io.Writer) OptFn {
	return func(o *Options) error {
		if w == nil {
			return errors.New("passed writer is nil")
		}
		o.Writer = w
		return nil
	}
}

func WithWorkDir(path string) OptFn {
	return func(o *Options) error {
		if !helpers.IsDir(path) {
			return fmt.Errorf("workind dir does not exist: %q", path)
		}
		o.WorkDir = path
		return nil
	}
}

func WithAttest(doattest bool) OptFn {
	return func(o *Options) error {
		o.Attest = doattest
		return nil
	}
}
