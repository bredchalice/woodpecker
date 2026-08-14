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
	"errors"
	"fmt"
	"net/http"
	"reflect"

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

// GetManualPipelineInputs resolves typed manual-run inputs for a repository branch.
func GetManualPipelineInputs(c *gin.Context) {
	_store := store.FromContext(c)
	repo := session.Repo(c)
	user := session.User(c)
	branch := c.DefaultQuery("branch", repo.Branch)
	if branch == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "branch is required"})
		return
	}

	_forge, err := server.Config.Services.Manager.ForgeFromRepo(repo)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	commit, err := _forge.BranchHead(c, user, repo, branch)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not fetch branch head: %w", err))
		return
	}

	repoUser, err := _store.GetUser(repo.UserID)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not load repository owner: %w", err))
		return
	}
	forge.Refresh(c, _forge, _store, repoUser)

	tmpPipeline := &model.Pipeline{
		Event:  model.EventManual,
		Commit: commit.SHA,
		Branch: branch,
		Ref:    branch,
	}

	configService := server.Config.Services.Manager.ConfigServiceFromRepo(repo)
	configs, err := configService.Fetch(c, _forge, repoUser, repo, tmpPipeline, nil, false)
	if errors.Is(err, &forge_types.ErrConfigNotFound{}) {
		c.JSON(http.StatusOK, manualPipelineInputsResponse{Inputs: map[string]yaml_types.ManualInput{}})
		return
	}
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not load config from forge: %w", err))
		return
	}

	inputs := map[string]yaml_types.ManualInput{}
	for _, config := range configs {
		workflow, err := pipeline_yaml.ParseBytes(config.Data)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("could not parse %s: %w", config.Name, err))
			return
		}

		for name, input := range workflow.Manual.Inputs {
			if existing, ok := inputs[name]; ok && !reflect.DeepEqual(existing, input) {
				c.JSON(http.StatusConflict, gin.H{
					"error": fmt.Sprintf("manual input %q is defined differently by multiple workflows", name),
				})
				return
			}
			inputs[name] = input
		}
	}

	c.JSON(http.StatusOK, manualPipelineInputsResponse{Inputs: inputs})
}
