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

package api

import (
	"testing"

	yaml_types "go.woodpecker-ci.org/woodpecker/v3/pipeline/frontend/yaml/types"
)

func TestApplyManualPipelineInputs(t *testing.T) {
	inputs := map[string]yaml_types.ManualInput{
		"action": {
			Type:     yaml_types.ManualInputTypeChoice,
			Required: true,
			Default:  "deploy",
			Options:  []string{"deploy", "seed"},
		},
		"force_reconfigure": {
			Type:    yaml_types.ManualInputTypeBoolean,
			Default: false,
		},
	}

	t.Run("applies defaults", func(t *testing.T) {
		variables, err := applyManualPipelineInputs(inputs, nil)
		if err != nil {
			t.Fatal(err)
		}
		if got := variables["action"]; got != "deploy" {
			t.Fatalf("expected deploy, got %q", got)
		}
		if got := variables["force_reconfigure"]; got != "false" {
			t.Fatalf("expected false, got %q", got)
		}
	})

	t.Run("accepts configured choice", func(t *testing.T) {
		variables, err := applyManualPipelineInputs(inputs, map[string]string{"action": "seed"})
		if err != nil {
			t.Fatal(err)
		}
		if got := variables["action"]; got != "seed" {
			t.Fatalf("expected seed, got %q", got)
		}
	})

	t.Run("rejects invalid choice", func(t *testing.T) {
		_, err := applyManualPipelineInputs(inputs, map[string]string{"action": "destroy"})
		if err == nil {
			t.Fatal("expected invalid choice error")
		}
	})

	t.Run("rejects invalid boolean", func(t *testing.T) {
		_, err := applyManualPipelineInputs(inputs, map[string]string{"force_reconfigure": "yes"})
		if err == nil {
			t.Fatal("expected invalid boolean error")
		}
	})

	t.Run("keeps extra variables", func(t *testing.T) {
		variables, err := applyManualPipelineInputs(inputs, map[string]string{"custom": "value"})
		if err != nil {
			t.Fatal(err)
		}
		if got := variables["custom"]; got != "value" {
			t.Fatalf("expected custom variable to be preserved, got %q", got)
		}
	})
}
