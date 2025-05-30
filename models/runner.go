// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc

package models

import (
	"context"

	testresult "github.com/in-toto/attestation/go/predicates/test_result/v0"
)

type TestRunner interface {
	Run(context.Context) ([]byte, bool, error)
}

type ResultsParser interface {
	ParseResults(context.Context, *testresult.TestResult, []byte) (*testresult.TestResult, error)
}
