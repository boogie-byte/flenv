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
	"io"
	"net/url"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
	"time"
)

type flag interface {
	isRequired() bool
	isSet() bool
	getName() string
	getLongDescription() string
	getShortDescription() string
	setValueFromDefault()
	setValueFromEnv() error
	setValueFromString(string) error
}

type Parser struct {
	envVarFormatter func(string) string
	envVarPrefix    string
	autoEnv         bool

	helpFlagName string

	appName            string
	appVersion         string
	appVersionFlagName string

	helpCalled    bool
	versionCalled bool

	flags     []flag
	flagIndex map[string]flag
}

func New(opts ...Option) *Parser {
	p := &Parser{
		flagIndex: make(map[string]flag),
		envVarFormatter: func(s string) string {
			return strings.ReplaceAll(strings.ToUpper(s), "-", "_")
		},
		autoEnv:            true,
		helpFlagName:       "help",
		appVersionFlagName: "version",
	}

	for _, opt := range opts {
		opt(p)
	}

	helpFlag := NewBoolFlag(&p.helpCalled, p.helpFlagName, "Show help message")
	p.registerFlag(p.helpFlagName, helpFlag)

	if p.appVersion != "" {
		versionFlag := NewBoolFlag(&p.versionCalled, p.appVersionFlagName, "Show application version")
		p.registerFlag(p.appVersionFlagName, versionFlag)
	}

	return p
}

func (p *Parser) Bool(target *bool, name, description string) *Flag[bool] {
	f := NewBoolFlag(target, name, description)
	p.registerFlag(name, f)

	if p.autoEnv {
		envVarName := p.envVarPrefix + p.envVarFormatter(name)
		f = f.Env(envVarName)
	}

	return f
}

func (p *Parser) Duration(target *time.Duration, name, description string) *Flag[time.Duration] {
	f := NewDurationFlag(target, name, description)
	p.registerFlag(name, f)

	if p.autoEnv {
		envVarName := p.envVarPrefix + p.envVarFormatter(name)
		f = f.Env(envVarName)
	}

	return f
}

func (p *Parser) Int(target *int, name, description string) *Flag[int] {
	f := NewIntFlag(target, name, description)
	p.registerFlag(name, f)

	if p.autoEnv {
		envVarName := p.envVarPrefix + p.envVarFormatter(name)
		f = f.Env(envVarName)
	}

	return f
}

func (p *Parser) String(target *string, name, description string) *Flag[string] {
	f := NewStringFlag(target, name, description)
	p.registerFlag(name, f)

	if p.autoEnv {
		envVarName := p.envVarPrefix + p.envVarFormatter(name)
		f = f.Env(envVarName)
	}

	return f
}

func (p *Parser) Float(target *float64, bitSize int, name, description string) *Flag[float64] {
	f := NewFloatFlag(target, bitSize, name, description)
	p.registerFlag(name, f)

	if p.autoEnv {
		envVarName := p.envVarPrefix + p.envVarFormatter(name)
		f = f.Env(envVarName)
	}

	return f
}

func (p *Parser) URL(target **url.URL, name, description string) *Flag[*url.URL] {
	f := NewURLFlag(target, name, description)
	p.registerFlag(name, f)

	if p.autoEnv {
		envVarName := p.envVarPrefix + p.envVarFormatter(name)
		f = f.Env(envVarName)
	}

	return f
}

func (p *Parser) Parse() {
	if errs := p.parse(os.Args[1:]); len(errs) != 0 {
		p.printErrs(os.Stderr, errs)
		os.Exit(1)
	}

	if p.helpCalled {
		p.printHelp(os.Stdout)
		os.Exit(0)
	}

	if p.versionCalled {
		p.printVersion(os.Stdout)
		os.Exit(0)
	}

	if errs := p.checkRequiredFlags(); len(errs) != 0 {
		p.printErrs(os.Stderr, errs)
		os.Exit(1)
	}
}

func (p *Parser) printHelp(w io.Writer) {
	slices.SortStableFunc(p.flags, func(a, b flag) int {
		return strings.Compare(a.getName(), b.getName())
	})

	appName := p.appName
	if appName == "" {
		appName = os.Args[0]
	}

	fmt.Fprintf(w, "Usage: %s", appName)
	for _, flag := range p.flags {
		if flag.isRequired() {
			fmt.Fprintf(w, " %s", flag.getShortDescription())
		}
	}
	for _, flag := range p.flags {
		if !flag.isRequired() {
			fmt.Fprintf(w, " [%s]", flag.getShortDescription())
		}
	}

	fmt.Fprint(w, "\n\n")
	fmt.Fprintln(w, "Flags:")

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, flag := range p.flags {
		fmt.Fprintln(tw, flag.getLongDescription())
	}
	tw.Flush()
}

func (p *Parser) printVersion(w io.Writer) {
	fmt.Fprintln(w, p.appVersion)
}

func (p *Parser) printErrs(w io.Writer, errs []error) {
	for _, err := range errs {
		fmt.Fprintln(w, err)
	}
	fmt.Fprintf(w, "\nUse '--%s' flag for more info.\n", p.helpFlagName)
}

func (p *Parser) registerFlag(name string, f flag) {
	if _, ok := p.flagIndex[name]; ok {
		panic(fmt.Sprintf("flag with name %s is already registered", name))
	}

	p.flags = append(p.flags, f)
	p.flagIndex[name] = f
}

func (p *Parser) set(name, value string) error {
	if f := p.flagIndex[name]; f != nil {
		return f.setValueFromString(value)
	}

	return fmt.Errorf("unknown flag: --%s", name)
}

func (p *Parser) parse(args []string) []error {
	var parseErrs []error

	for _, v := range p.flagIndex {
		v.setValueFromDefault()
		if err := v.setValueFromEnv(); err != nil {
			parseErrs = append(parseErrs, err)
		}
	}

	for len(args) > 0 {
		arg := args[0]
		args = args[1:]

		if !strings.HasPrefix(arg, "--") {
			parseErrs = append(parseErrs, fmt.Errorf("unexpected argument: %s", arg))
			return parseErrs
		}

		arg = strings.TrimPrefix(arg, "--")

		if arg == "" {
			// end of flags
			if len(args) != 0 {
				parseErrs = append(parseErrs, fmt.Errorf("unexpected arguments: %s", strings.Join(args, " ")))
				return parseErrs
			}
			break
		}

		if equalsIdx := strings.Index(arg, "="); equalsIdx != -1 {
			// --key=value
			if err := p.set(arg[:equalsIdx], arg[equalsIdx+1:]); err != nil {
				parseErrs = append(parseErrs, err)
			}
			continue
		}

		if len(args) == 0 || strings.HasPrefix(args[0], "--") {
			// --key (boolean flag)
			if err := p.set(arg, "true"); err != nil {
				parseErrs = append(parseErrs, err)
			}
			continue
		}

		// --key value
		if err := p.set(arg, args[0]); err != nil {
			parseErrs = append(parseErrs, err)
		}
		args = args[1:]
	}

	return parseErrs
}

func (p *Parser) checkRequiredFlags() []error {
	var checkErrs []error

	for _, flag := range p.flags {
		if flag.isRequired() && !flag.isSet() {
			checkErrs = append(checkErrs, fmt.Errorf("missing required flag: --%s", flag.getName()))
		}
	}

	return checkErrs
}
