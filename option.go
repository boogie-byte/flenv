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

type Option func(*Parser)

func WithEnvVarPrefix(prefix string) Option {
	return func(p *Parser) {
		p.envVarPrefix = prefix
	}
}

func WithEnvVarFormatter(f func(string) string) Option {
	return func(p *Parser) {
		p.envVarFormatter = f
	}
}

func WithoutAutoEnv() Option {
	return func(p *Parser) {
		p.autoEnv = false
	}
}

func WithHelpFlagName(name string) Option {
	return func(p *Parser) {
		p.helpFlagName = name
	}
}

func WithAppVersionFlagName(name string) Option {
	return func(p *Parser) {
		p.appVersionFlagName = name
	}
}

func WithAppVersion(version string) Option {
	return func(p *Parser) {
		p.appVersion = version
	}
}

func WithAppName(name string) Option {
	return func(p *Parser) {
		p.appName = name
	}
}
