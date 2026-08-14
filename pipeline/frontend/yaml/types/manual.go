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
