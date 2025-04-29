// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc

package golang

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	testresult "github.com/in-toto/attestation/go/predicates/test_result/v0"
	intoto "github.com/in-toto/attestation/go/v1"
)

type testLine struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Output  string    `json:"Output"`
	Test    string    `json:"Test"`
	Elapsed float32   `json:"Elapsed"`
}

// ParseResults parses the structures output of the go tests
func (r *Runner) ParseResults(ctx context.Context, att *testresult.TestResult, res []byte) (*testresult.TestResult, error) {
	dec := json.NewDecoder(bytes.NewReader(res))
	if att == nil {
		att = &testresult.TestResult{
			Result:        "pass", // will change below if tests fail
			Configuration: []*intoto.ResourceDescriptor{},
			PassedTests:   []string{},
			FailedTests:   []string{},
		}
	} else {
		att.Result = "pass"
		att.PassedTests = []string{}
		att.FailedTests = []string{}
	}

	for {
		var result testLine

		err := dec.Decode(&result)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Parsing line: %v", err)
			continue
		}

		if result.Test == "" {
			continue
		}

		switch result.Action {
		case "fail":
			att.FailedTests = append(att.FailedTests, result.Test)
		case "pass":
			att.PassedTests = append(att.PassedTests, result.Test)
		}
	}

	if len(att.GetFailedTests()) > 0 {
		att.Result = "fail"
	}

	return att, nil
}
