// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2026 Carabiner Systems, Inc

package npm

import (
	"bufio"
	"bytes"
	"context"
	"regexp"
	"strings"

	testresult "github.com/in-toto/attestation/go/predicates/test_result/v0"
	intoto "github.com/in-toto/attestation/go/v1"
)

// tapLine matches a TAP test point line, e.g.:
//
//	ok 1 - description
//	not ok 2 - description # TODO reason
//
// Leading whitespace is allowed so that subtests are picked up too.
var tapLine = regexp.MustCompile(`^\s*(not )?ok\s+\d+(?:\s*-?\s*(.*))?$`)

// ParseResults extracts test names and pass/fail status from TAP output
// emitted by the underlying npm test framework.
func (r *Runner) ParseResults(_ context.Context, att *testresult.TestResult, res []byte) (*testresult.TestResult, error) {
	if att == nil {
		att = &testresult.TestResult{
			Result:        "pass",
			Configuration: []*intoto.ResourceDescriptor{},
			PassedTests:   []string{},
			FailedTests:   []string{},
		}
	} else {
		att.Result = "pass"
		att.PassedTests = []string{}
		att.FailedTests = []string{}
	}

	scanner := bufio.NewScanner(bytes.NewReader(res))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		m := tapLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		name := strings.TrimSpace(m[2])
		// Strip trailing TAP directives like "# SKIP" or "# TODO".
		if i := strings.Index(name, "#"); i >= 0 {
			name = strings.TrimSpace(name[:i])
		}
		if name == "" {
			continue
		}

		if m[1] == "" {
			att.PassedTests = append(att.PassedTests, name)
		} else {
			att.FailedTests = append(att.FailedTests, name)
		}
	}

	if len(att.GetFailedTests()) > 0 {
		att.Result = "fail"
	}

	return att, nil
}
