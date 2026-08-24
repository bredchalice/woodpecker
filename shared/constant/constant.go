// Copyright 2022 Woodpecker Authors
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

package constant

import "time"

// DefaultConfigOrder represent the priority in which woodpecker searches for a pipeline config by default
// folders are indicated by supplying a trailing slash.
var DefaultConfigOrder = []string{
	".woodpecker/",
	".woodpecker.yaml",
	".woodpecker.yml",
}

const (
	// DefaultClonePlugin can be changed by 'WOODPECKER_DEFAULT_CLONE_PLUGIN' at runtime.
	// renovate: datasource=docker depName=woodpeckerci/plugin-git
	DefaultClonePlugin = "docker.io/woodpeckerci/plugin-git:2.9.2"
)

// TrustedClonePlugins can be changed by 'WOODPECKER_PLUGINS_TRUSTED_CLONE' at runtime.
var TrustedClonePlugins = []string{
	DefaultClonePlugin,
	"docker.io/woodpeckerci/plugin-git",
	"quay.io/woodpeckerci/plugin-git",
}

// TaskTimeout is the grace period after the last successful agent lease refresh
// before a running task is considered abandoned and returned to the queue.
// Agents refresh active workflow leases every 15 seconds. Five minutes gives
// transient gRPC, API-server and control-plane interruptions room to recover
// without reclaiming a workflow whose build pod is still running.
var TaskTimeout = 5 * time.Minute
