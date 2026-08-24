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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"

	pipeline_yaml "go.woodpecker-ci.org/woodpecker/v3/pipeline/frontend/yaml"
	yaml_types "go.woodpecker-ci.org/woodpecker/v3/pipeline/frontend/yaml/types"
	"go.woodpecker-ci.org/woodpecker/v3/server"
	"go.woodpecker-ci.org/woodpecker/v3/server/forge"
	forge_types "go.woodpecker-ci.org/woodpecker/v3/server/forge/types"
	"go.woodpecker-ci.org/woodpecker/v3/server/model"
	"go.woodpecker-ci.org/woodpecker/v3/server/router/middleware/session"
	"go.woodpecker-ci.org/woodpecker/v3/server/store"
)

type manualPipelineInputsResponse struct {
	Inputs map[string]yaml_types.ManualInput `json:"inputs"`
}

func GetManualPipelineInputs(c *gin.Context) {
	branch := c.DefaultQuery("branch", session.Repo(c).Branch)
	if branch == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "branch is required"})
		return
	}

	inputs, status, err := resolveManualPipelineInputs(c, branch)
	if err != nil {
		_ = c.AbortWithError(status, err)
		return
	}

	c.JSON(http.StatusOK, manualPipelineInputsResponse{Inputs: inputs})
}

func ValidateManualPipelineInputs(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var opts model.PipelineOptions
	if err := json.Unmarshal(body, &opts); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if opts.Branch == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "branch is required"})
		return
	}

	inputs, status, err := resolveManualPipelineInputs(c, opts.Branch)
	if err != nil {
		_ = c.AbortWithError(status, err)
		return
	}

	variables, err := applyManualPipelineInputs(inputs, opts.Variables)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	opts.Variables = variables

	normalizedBody, err := json.Marshal(&opts)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(normalizedBody))
	c.Request.ContentLength = int64(len(normalizedBody))
}

func resolveManualPipelineInputs(c *gin.Context, branch string) (map[string]yaml_types.ManualInput, int, error) {
	_store := store.FromContext(c)
	repo := session.Repo(c)
	user := session.User(c)

	_forge, err := server.Config.Services.Manager.ForgeFromRepo(repo)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	commit, err := _forge.BranchHead(c, user, repo, branch)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("could not fetch branch head: %w", err)
	}

	repoUser, err := _store.GetUser(repo.UserID)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("could not load repository owner: %w", err)
	}
	forge.Refresh(c, _forge, _store, repoUser)

	tmpPipeline := &model.Pipeline{Event: model.EventManual, Commit: commit.SHA, Branch: branch, Ref: branch}
	configService := server.Config.Services.Manager.ConfigServiceFromRepo(repo)
	configs, fetchErr := configService.Fetch(c, _forge, repoUser, repo, tmpPipeline, nil, false)
	if errors.Is(fetchErr, &forge_types.ErrConfigNotFound{}) {
		return map[string]yaml_types.ManualInput{}, http.StatusOK, nil
	}
	if fetchErr != nil && configs == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("could not load config from forge: %w", fetchErr)
	}

	inputs := map[string]yaml_types.ManualInput{}
	for _, config := range configs {
		workflow, err := pipeline_yaml.ParseBytes(config.Data)
		if err != nil {
			return nil, http.StatusUnprocessableEntity, fmt.Errorf("could not parse %s: %w", config.Name, err)
		}
		for name, input := range workflow.Manual.Inputs {
			if existing, ok := inputs[name]; ok && !reflect.DeepEqual(existing, input) {
				return nil, http.StatusConflict, fmt.Errorf("manual input %q is defined differently by multiple workflows", name)
			}
			inputs[name] = input
		}
	}

	return inputs, http.StatusOK, nil
}

func applyManualPipelineInputs(inputs map[string]yaml_types.ManualInput, variables map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(variables)+len(inputs))
	for key, value := range variables {
		result[key] = value
	}

	for name, input := range inputs {
		value, provided := result[name]
		if !provided {
			switch defaultValue := input.Default.(type) {
			case string:
				value, provided = defaultValue, true
			case bool:
				value, provided = strconv.FormatBool(defaultValue), true
			}
		}
		if (!provided || value == "") && input.Required {
			return nil, fmt.Errorf("manual input %q is required", name)
		}
		if !provided || value == "" {
			continue
		}

		switch input.Type {
		case yaml_types.ManualInputTypeChoice:
			valid := false
			for _, option := range input.Options {
				if value == option { valid = true; break }
			}
			if !valid {
				return nil, fmt.Errorf("manual input %q must be one of %v", name, input.Options)
			}
		case yaml_types.ManualInputTypeBoolean:
			if value != "true" && value != "false" {
				return nil, fmt.Errorf("manual input %q must be true or false", name)
			}
		}
		result[name] = value
	}
	return result, nil
}
