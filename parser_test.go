// Copyright 2025 Sergey Vinogradov
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flenv

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserPrintHelp(t *testing.T) {
	var (
		b bool
		i int
		s string
	)

	p := New(
		WithAppName("test-app"),
		WithAppVersion("1.2.3"),
	)
	p.Bool(&b, "test-bool-flag", "Test bool flag")
	p.Int(&i, "test-int-flag", "Test int flag").Required()
	p.String(&s, "test-string-flag", "Test string flag")

	buf := bytes.NewBuffer(nil)
	p.printHelp(buf)

	const helpMessage = "Usage: test-app --test-int-flag=INT [--help] [--test-bool-flag] [--test-string-flag=STRING] [--version]\n\n" +
		"Flags:\n" +
		"  --help                     Show help message\n" +
		"  --test-bool-flag           Test bool flag [$TEST_BOOL_FLAG]\n" +
		"  --test-int-flag=INT        Test int flag (required) [$TEST_INT_FLAG]\n" +
		"  --test-string-flag=STRING  Test string flag [$TEST_STRING_FLAG]\n" +
		"  --version                  Show application version\n"

	assert.Equal(t, helpMessage, buf.String())
}

func TestParserPrintError(t *testing.T) {
	p := New()

	buf := bytes.NewBuffer(nil)
	p.printErrs(buf, []error{errors.New("test-error")})

	assert.Equal(t, "test-error\n\nUse '--help' flag for more info.\n", buf.String())
}

func TestParserPrintVersion(t *testing.T) {
	p := New(
		WithAppVersion("1.2.3"),
	)

	buf := bytes.NewBuffer(nil)
	p.printVersion(buf)

	assert.Equal(t, "1.2.3\n", buf.String())
}

func TestParserRegisterExistingFlag(t *testing.T) {
	var v string

	p := New()
	p.String(&v, "test-flag", "Test flag")
	assert.Panics(t, func() {
		p.String(&v, "test-flag", "Test flag")
	})
}

func TestParserParse(t *testing.T) {
	t.Run("ValueFromEnvError", func(t *testing.T) {
		t.Setenv("TEST_FLAG", "abc")

		var i int
		p := New()
		p.Int(&i, "test-flag", "Test flag")
		errs := p.parse(nil)
		assert.Len(t, errs, 1)
	})

	t.Run("NonexistentFlag", func(t *testing.T) {
		p := New()
		errs := p.parse([]string{"--nonexistent-flag", "abc"})
		assert.Len(t, errs, 1)
	})

	t.Run("UnexpectedArgument", func(t *testing.T) {
		var i int
		p := New()
		p.Int(&i, "test-flag", "Test flag")

		errs := p.parse([]string{"--test-flag", "10", "abc"})
		assert.Len(t, errs, 1)
	})

	t.Run("UnexpectedValue", func(t *testing.T) {
		var i int
		p := New()
		p.Int(&i, "test-flag", "Test flag")

		errs := p.parse([]string{"--test-flag", "abc"})
		assert.Len(t, errs, 1)
	})

	t.Run("MalformedFlag", func(t *testing.T) {
		var i int
		p := New()
		p.Int(&i, "test-flag", "Test flag")

		errs := p.parse([]string{"--", "test-flag", "10"})
		assert.Len(t, errs, 1)
	})

	t.Run("Toggle", func(t *testing.T) {
		var b bool
		p := New()
		p.Bool(&b, "test-flag", "Test flag")

		errs := p.parse([]string{"--test-flag"})
		assert.Empty(t, errs)
		assert.True(t, b)
	})

	t.Run("EqualsSignFormat", func(t *testing.T) {
		var i int
		p := New()
		p.Int(&i, "test-flag", "Test flag")

		errs := p.parse([]string{"--test-flag=10"})
		assert.Empty(t, errs)
		assert.Equal(t, 10, i)
	})

	t.Run("TwoArgsFormat", func(t *testing.T) {
		var i int
		p := New()
		p.Int(&i, "test-flag", "Test flag")

		errs := p.parse([]string{"--test-flag", "10"})
		assert.Empty(t, errs)
		assert.Equal(t, 10, i)
	})
}

func TestParserCheckRequiredFlags(t *testing.T) {
	t.Run("NoRequiredFlags", func(t *testing.T) {
		var i int
		p := New()
		p.Int(&i, "test-flag", "Test flag")

		parseErrs := p.parse(nil)
		require.Empty(t, parseErrs)

		checkErrs := p.checkRequiredFlags()
		assert.Empty(t, checkErrs)
	})

	t.Run("RequiredFlagNotSet", func(t *testing.T) {
		var i int
		p := New()
		p.Int(&i, "test-flag", "Test flag").Required()

		parseErrs := p.parse(nil)
		require.Empty(t, parseErrs)

		checkErrs := p.checkRequiredFlags()
		assert.Len(t, checkErrs, 1)
	})

	t.Run("RequiredFlagSet", func(t *testing.T) {
		var i int
		p := New()
		p.Int(&i, "test-flag", "Test flag").Required()

		parseErrs := p.parse([]string{"--test-flag=10"})
		require.Empty(t, parseErrs)

		checkErrs := p.checkRequiredFlags()
		assert.Empty(t, checkErrs)
	})
}
