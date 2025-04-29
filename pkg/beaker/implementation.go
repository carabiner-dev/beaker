// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc

package beaker

import (
	"context"
	"fmt"
	"path/filepath"

	v0 "github.com/in-toto/attestation/go/predicates/test_result/v0"
	v1 "github.com/in-toto/attestation/go/v1"
	"sigs.k8s.io/release-utils/util"

	"github.com/carabiner-dev/beaker/internal/git"
)

type launcherImplementation interface {
	InitAttestation(context.Context, *Options) (*v0.TestResult, error)
	RunLaunchPack(context.Context, *Options, *LaunchPack) ([]byte, error)
}

type defaultLauncherImplementation struct{}

func (dli *defaultLauncherImplementation) RunLaunchPack(ctx context.Context, opts *Options, pack *LaunchPack) ([]byte, error) {
	output, _, err := pack.Runner.Run(ctx)
	if err != nil {
		return nil, fmt.Errorf("runner error: %w", err)
	}
	return output, nil
}

func (dli *defaultLauncherImplementation) InitAttestation(_ context.Context, opts *Options) (*v0.TestResult, error) {
	if !util.Exists(filepath.Join(opts.WorkDir, ".git")) {
		return nil, nil
	}

	locator, err := git.RepoVCSLocator(opts.WorkDir)
	if err != nil {
		return nil, fmt.Errorf("reading VCS locator: %w", err)
	}

	// Get the repo version
	tagPlus, commit, err := git.RepoVersion(opts.WorkDir)
	if err != nil {
		return nil, fmt.Errorf("computing git commit: %w", err)
	}

	// Build the attestation
	att := &v0.TestResult{
		Configuration: []*v1.ResourceDescriptor{
			{
				Name: tagPlus,
				Uri:  locator,
				Digest: map[string]string{
					"sha1":      commit,
					"gitCommit": commit,
				},
			},
		},
		// Url: "",
	}
	return att, nil
}
