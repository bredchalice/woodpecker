// Copyright 2026 Woodpecker Authors
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

package types

import "fmt"

const (
	ManualInputTypeString  = "string"
	ManualInputTypeChoice  = "choice"
	ManualInputTypeBoolean = "boolean"
)

// Manual defines options exposed when a workflow is started manually.
type Manual struct {
	Inputs map[string]ManualInput `yaml:"inputs,omitempty" json:"inputs,omitempty"`
}

// ManualInput defines one typed input in the manual pipeline form.
// Values are submitted through the existing pipeline variables mechanism.
type ManualInput struct {
	Type        string   `yaml:"type,omitempty" json:"type,omitempty"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Required    bool     `yaml:"required,omitempty" json:"required,omitempty"`
	Default     any      `yaml:"default,omitempty" json:"default,omitempty"`
	Options     []string `yaml:"options,omitempty" json:"options,omitempty"`
}

// NormalizeAndValidate applies defaults and validates a manual input schema.
func (m *Manual) NormalizeAndValidate() error {
	for name, input := range m.Inputs {
		if name == "" {
			return fmt.Errorf("manual input name must not be empty")
		}

		if input.Type == "" {
			input.Type = ManualInputTypeString
		}

		switch input.Type {
		case ManualInputTypeString:
			if input.Default != nil {
				if _, ok := input.Default.(string); !ok {
					return fmt.Errorf("manual input %q: string default must be a string", name)
				}
			}
		case ManualInputTypeChoice:
			if len(input.Options) == 0 {
				return fmt.Errorf("manual input %q: choice requires at least one option", name)
			}
			if input.Default != nil {
				defaultValue, ok := input.Default.(string)
				if !ok {
					return fmt.Errorf("manual input %q: choice default must be a string", name)
				}
				found := false
				for _, option := range input.Options {
					if option == defaultValue {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("manual input %q: default %q is not one of the configured options", name, defaultValue)
				}
			}
		case ManualInputTypeBoolean:
			if len(input.Options) != 0 {
				return fmt.Errorf("manual input %q: boolean inputs cannot define options", name)
			}
			if input.Default != nil {
				if _, ok := input.Default.(bool); !ok {
					return fmt.Errorf("manual input %q: boolean default must be true or false", name)
				}
			}
		default:
			return fmt.Errorf("manual input %q: unsupported type %q", name, input.Type)
		}

		m.Inputs[name] = input
	}

	return nil
}
