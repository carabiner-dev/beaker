// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc

package beaker

import (
	"errors"
	"fmt"

	"github.com/carabiner-dev/beaker/pkg/run"
)

type LaunchPack struct {
	Runner run.TestRunner
	Parser run.ResultsParser
}

func (pack *LaunchPack) Verify() error {
	errs := []error{}
	if pack.Parser == nil {
		errs = append(errs, fmt.Errorf("LaunchPack has no parser set"))
	}

	if pack.Runner == nil {
		errs = append(errs, fmt.Errorf("LaunchPack has no runner set"))
	}
	return errors.Join(errs...)
}
