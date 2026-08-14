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

package yaml

import "testing"

func TestParseManualInputs(t *testing.T) {
	workflow, err := ParseString(`
when:
  - event: manual
manual:
  inputs:
    action:
      type: choice
      options: [deploy, seed]
      default: deploy
      required: true
    force_reconfigure:
      type: boolean
      default: false
steps:
  - name: run
    image: alpine
    commands:
      - echo run
`)
	if err != nil {
		t.Fatal(err)
	}

	if got := workflow.Manual.Inputs["action"].Default; got != "deploy" {
		t.Fatalf("expected deploy default, got %#v", got)
	}
	if got := workflow.Manual.Inputs["force_reconfigure"].Default; got != false {
		t.Fatalf("expected false default, got %#v", got)
	}
}

func TestParseManualInputsRejectsInvalidChoiceDefault(t *testing.T) {
	_, err := ParseString(`
manual:
  inputs:
    action:
      type: choice
      options: [deploy, seed]
      default: destroy
steps:
  - name: run
    image: alpine
    commands:
      - echo run
`)
	if err == nil {
		t.Fatal("expected invalid manual input schema to fail parsing")
	}
}
