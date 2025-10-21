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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFlag(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		var v bool
		f := NewBoolFlag(&v, "test-bool-flag", "Test bool flag")
		assert.Equal(t, "test-bool-flag", f.getName())
		assert.Equal(t, "--test-bool-flag", f.getShortDescription())
		assert.Equal(t, true, f.isBool)
	})

	t.Run("duration", func(t *testing.T) {
		var v time.Duration
		f := NewDurationFlag(&v, "test-duration-flag", "Test duration flag")
		assert.Equal(t, "test-duration-flag", f.getName())
		assert.Equal(t, "--test-duration-flag=DURATION", f.getShortDescription())
		assert.Equal(t, false, f.isBool)
	})

	t.Run("int", func(t *testing.T) {
		var v int
		f := NewIntFlag(&v, "test-int-flag", "Test int flag")
		assert.Equal(t, "test-int-flag", f.getName())
		assert.Equal(t, "--test-int-flag=INT", f.getShortDescription())
		assert.Equal(t, false, f.isBool)
	})

	t.Run("string", func(t *testing.T) {
		var v string
		f := NewStringFlag(&v, "test-string-flag", "Test string flag")
		assert.Equal(t, "test-string-flag", f.getName())
		assert.Equal(t, "--test-string-flag=STRING", f.getShortDescription())
		assert.Equal(t, false, f.isBool)
	})
}

func TestFlagLongDescription(t *testing.T) {
	t.Run("required", func(t *testing.T) {
		var s string
		f := NewStringFlag(&s, "test-flag", "Test flag").Placeholder("<test_placeholder>").Env("TEST_FLAG").Required()
		assert.Equal(t, "  --test-flag=<test_placeholder>\tTest flag (required) [$TEST_FLAG]", f.getLongDescription())
	})

	t.Run("default", func(t *testing.T) {
		var s string
		f := NewStringFlag(&s, "test-flag", "Test flag").Placeholder("<test_placeholder>").Env("TEST_FLAG").Default("foo")
		assert.Equal(t, "  --test-flag=<test_placeholder>\tTest flag (default: foo) [$TEST_FLAG]", f.getLongDescription())
	})
}

func TestFlagEnv(t *testing.T) {
	var v bool
	f := NewBoolFlag(&v, "test-flag", "Test flag").Env("TEST_FLAG")
	assert.Equal(t, "TEST_FLAG", f.envVarName)
}

func TestFlagPlaceholder(t *testing.T) {
	t.Run("BoolPanic", func(t *testing.T) {
		var v bool
		f := NewBoolFlag(&v, "test-flag", "Test flag")
		assert.Panics(t, func() {
			f.Placeholder("foo")
		})
	})

	t.Run("Success", func(t *testing.T) {
		var v string
		f := NewStringFlag(&v, "test-flag", "Test flag")
		assert.NotPanics(t, func() {
			f.Placeholder("<test_placeholder>")
		})
		assert.Equal(t, "<test_placeholder>", f.placeholder)
	})
}

func TestFlagDefault(t *testing.T) {
	t.Run("BoolPanic", func(t *testing.T) {
		var v bool
		f := NewBoolFlag(&v, "test-flag", "Test flag")
		assert.Panics(t, func() {
			f.Default(true)
		})
	})

	t.Run("RequiredPanic", func(t *testing.T) {
		var v string
		f := NewStringFlag(&v, "test-flag", "Test flag").Required()
		assert.Panics(t, func() {
			f.Default("foo")
		})
	})

	t.Run("Success", func(t *testing.T) {
		var v string
		f := NewStringFlag(&v, "test-flag", "Test flag")
		assert.NotPanics(t, func() {
			f.Default("foo")
		})
		assert.True(t, f.defaultValueSet)
	})
}

func TestFlagRequired(t *testing.T) {
	t.Run("BoolPanic", func(t *testing.T) {
		var v bool
		f := NewBoolFlag(&v, "test-flag", "Test flag")
		assert.Panics(t, func() {
			f.Required()
		})
	})

	t.Run("RequiredPanic", func(t *testing.T) {
		var v string
		f := NewStringFlag(&v, "test-flag", "Test flag").Default("foo")
		assert.Panics(t, func() {
			f.Required()
		})
	})

	t.Run("Success", func(t *testing.T) {
		var v string
		f := NewStringFlag(&v, "test-flag", "Test flag")
		assert.NotPanics(t, func() {
			f.Required()
		})
		assert.True(t, f.isRequired())
	})
}

func TestFlagSetValue(t *testing.T) {
	t.Run("ValidValue", func(t *testing.T) {
		var v int
		f := NewIntFlag(&v, "test-flag", "Test flag")
		err := f.setValueFromString("10")
		require.NoError(t, err)
		assert.Equal(t, 10, v)
	})

	t.Run("InvalidValue", func(t *testing.T) {
		var v int
		f := NewIntFlag(&v, "test-flag", "Test flag")
		err := f.setValueFromString("abc")
		assert.Error(t, err)
	})

	t.Run("FromEnv", func(t *testing.T) {
		t.Setenv("TEST_FLAG", "10")

		var v int
		f := NewIntFlag(&v, "test-flag", "Test flag").Env("TEST_FLAG")
		err := f.setValueFromEnv()
		require.NoError(t, err)
		assert.Equal(t, 10, v)
	})

	t.Run("FromDefault", func(t *testing.T) {
		var v int
		f := NewIntFlag(&v, "test-flag", "Test flag").Default(10)
		f.setValueFromDefault()
		assert.Equal(t, 10, v)
	})
}
