// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc

package beaker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/carabiner-dev/beaker/pkg/runners/golang"
	"google.golang.org/protobuf/encoding/protojson"
	"sigs.k8s.io/release-utils/util"
)

func New(funcs ...OptFn) (*Launcher, error) {
	opts := Options{
		Writer:  os.Stdout,
		WorkDir: ".",
	}
	for _, f := range funcs {
		if err := f(&opts); err != nil {
			return nil, err
		}
	}
	return &Launcher{
		impl:    &defaultLauncherImplementation{},
		Options: opts,
	}, nil
}

type Launcher struct {
	impl    launcherImplementation
	Options Options
}

// Test launches the test suite defined in the launch pack
func (l *Launcher) Test(ctx context.Context, pack *LaunchPack) error {
	att, err := l.impl.InitAttestation(ctx, &l.Options)
	if err != nil {
		return fmt.Errorf("initializing attestation: %w", err)
	}

	output, _, err := pack.Runner.Run(ctx)
	if err != nil {
		return fmt.Errorf("runner error: %w", err)
	}

	att, err = pack.Parser.ParseResults(ctx, att, output)
	if err != nil {
		return fmt.Errorf("parsing results: %w", err)
	}

	if l.Options.Writer == nil {
		return fmt.Errorf("tests ran successfully but no writer was configured")
	}

	encoder := protojson.MarshalOptions{
		Multiline:         true,
		Indent:            "  ",
		EmitDefaultValues: true,
	}

	jdata, err := encoder.Marshal(att)
	if err != nil {
		return fmt.Errorf("marshalling attestation: %w", err)
	}

	if _, err := l.Options.Writer.Write(jdata); err != nil {
		return fmt.Errorf("wiriting attestation data: %w", err)
	}

	return nil
}

// LaunchPackFromRepo reads a codebase and returns a launchpack
func LaunchPackFromRepo(path string) (*LaunchPack, error) {
	switch {
	case util.Exists(filepath.Join(path, "go.mod")):
		gorunner, err := golang.New(golang.WithWorkDir(path))
		if err != nil {
			return nil, fmt.Errorf("initializing go launchpack: %w", err)
		}
		return &LaunchPack{
			Runner: gorunner,
			Parser: gorunner,
		}, nil
	default:
		return nil, errors.New("unable to detect the language ecosystem")
	}
}
