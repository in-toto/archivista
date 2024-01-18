// Copyright 2023 The Witness Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"fmt"
	"time"
)

type Configurer interface {
	Description() string
	Name() string
	SetPrefix(string)
}

type Option interface {
	int | string | []string | bool | time.Duration
}

type ConfigOption[T any, TOption Option] struct {
	name        string
	prefix      string
	description string
	defaultVal  TOption
	setter      func(T, TOption) (T, error)
}

func (co *ConfigOption[T, TOption]) Name() string {
	if len(co.prefix) == 0 {
		return co.name
	}

	return fmt.Sprintf("%v-%v", co.prefix, co.name)
}

func (co *ConfigOption[T, TOption]) SetPrefix(prefix string) {
	co.prefix = prefix
}

func (co *ConfigOption[T, TOption]) DefaultVal() TOption {
	return co.defaultVal
}

func (co *ConfigOption[T, TOption]) Description() string {
	return co.description
}

func (co *ConfigOption[T, TOption]) Setter() func(T, TOption) (T, error) {
	return co.setter
}

func IntConfigOption[T any](name, description string, defaultVal int, setter func(T, int) (T, error)) *ConfigOption[T, int] {
	return &ConfigOption[T, int]{
		name:        name,
		description: description,
		defaultVal:  defaultVal,
		setter:      setter,
	}
}

func StringConfigOption[T any](name, description string, defaultVal string, setter func(T, string) (T, error)) *ConfigOption[T, string] {
	return &ConfigOption[T, string]{
		name:        name,
		description: description,
		defaultVal:  defaultVal,
		setter:      setter,
	}
}

func StringSliceConfigOption[T any](name, description string, defaultVal []string, setter func(T, []string) (T, error)) *ConfigOption[T, []string] {
	return &ConfigOption[T, []string]{
		name:        name,
		description: description,
		defaultVal:  defaultVal,
		setter:      setter,
	}
}

func BoolConfigOption[T any](name, description string, defaultVal bool, setter func(T, bool) (T, error)) *ConfigOption[T, bool] {
	return &ConfigOption[T, bool]{
		name:        name,
		description: description,
		defaultVal:  defaultVal,
		setter:      setter,
	}
}

func DurationConfigOption[T any](name, description string, defaultVal time.Duration, setter func(T, time.Duration) (T, error)) *ConfigOption[T, time.Duration] {
	return &ConfigOption[T, time.Duration]{
		name:        name,
		description: description,
		defaultVal:  defaultVal,
		setter:      setter,
	}
}
