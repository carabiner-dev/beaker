// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc

package run

import (
	"context"

	testresult "github.com/in-toto/attestation/go/predicates/test_result/v0"
)

type TestRunner interface {
	Run(context.Context) ([]byte, bool, error)
}

type ResultsParser interface {
	ParseResults(context.Context, []byte) (*testresult.TestResult, error)
}
