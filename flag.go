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
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Flag[T any] struct {
	target *T
	isBool bool

	name        string
	envVarName  string
	helpMessage string
	placeholder string

	defaultValue    T
	defaultValueSet bool

	required bool
	set      bool

	parseFunc func(string) (T, error)
}

func (f *Flag[T]) Env(name string) *Flag[T] {
	f.envVarName = name
	return f
}

func (f *Flag[T]) Placeholder(placeholder string) *Flag[T] {
	if f.isBool {
		panic("setting placeholder for a bool flag is not possible")
	}

	f.placeholder = placeholder
	return f
}

func (f *Flag[T]) Default(v T) *Flag[T] {
	if f.isBool {
		panic("setting default value for a bool flag is not possible")
	}

	if f.required {
		panic("setting default value for a required flag is not possible")
	}

	f.defaultValue = v
	f.defaultValueSet = true
	return f
}

func (f *Flag[T]) Required() *Flag[T] {
	if f.isBool {
		panic("making a bool flag required is not possible")
	}

	if f.defaultValueSet {
		panic("making a flag with default value required is not possible")
	}

	f.required = true
	return f
}

func (f *Flag[T]) isRequired() bool {
	return f.required
}

func (f *Flag[T]) isSet() bool {
	return f.set
}

func (f *Flag[T]) getName() string {
	return f.name
}

func (f *Flag[T]) getShortDescription() string {
	if f.isBool {
		return fmt.Sprintf("--%s", f.name)
	}
	return fmt.Sprintf("--%s=%s", f.name, f.placeholder)
}

func (f *Flag[T]) getLongDescription() string {
	b := &strings.Builder{}

	fmt.Fprintf(b, "  %s\t%s", f.getShortDescription(), f.helpMessage)

	switch {
	case f.required:
		fmt.Fprint(b, " (required)")
	case f.defaultValueSet:
		fmt.Fprintf(b, " (default: %v)", f.defaultValue)
	}

	if f.envVarName != "" {
		fmt.Fprintf(b, " [$%s]", f.envVarName)
	}

	return b.String()
}

func (f *Flag[T]) setValue(val T) {
	*f.target = val
	f.set = true
}

func (f *Flag[T]) setValueFromString(s string) error {
	val, err := f.parseFunc(s)
	if err != nil {
		return err
	}

	f.setValue(val)

	return nil
}

func (f *Flag[T]) setValueFromEnv() error {
	val, ok := os.LookupEnv(f.envVarName)
	if !ok {
		return nil
	}

	return f.setValueFromString(val)
}

func (f *Flag[T]) setValueFromDefault() {
	if f.defaultValueSet {
		f.setValue(f.defaultValue)
	}
}

func NewBoolFlag(target *bool, name, helpMessage string) *Flag[bool] {
	return &Flag[bool]{
		target:      target,
		name:        name,
		helpMessage: helpMessage,
		isBool:      true,
		parseFunc:   strconv.ParseBool,
	}
}

func NewDurationFlag(target *time.Duration, name, helpMessage string) *Flag[time.Duration] {
	return &Flag[time.Duration]{
		target:      target,
		name:        name,
		helpMessage: helpMessage,
		placeholder: "DURATION",
		parseFunc:   time.ParseDuration,
	}
}

func NewIntFlag(target *int, name, helpMessage string) *Flag[int] {
	return &Flag[int]{
		target:      target,
		name:        name,
		helpMessage: helpMessage,
		placeholder: "INT",
		parseFunc:   strconv.Atoi,
	}
}

func NewStringFlag(target *string, name, helpMessage string) *Flag[string] {
	return &Flag[string]{
		target:      target,
		name:        name,
		helpMessage: helpMessage,
		placeholder: "STRING",
		parseFunc: func(s string) (string, error) {
			return s, nil
		},
	}
}
