// Copyright 2024 Woodpecker Authors
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

package kubernetes

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestServiceAccountTokenRetryRoundTripper(t *testing.T) {
	tokenFile := t.TempDir() + "/token"
	require.NoError(t, os.WriteFile(tokenFile, []byte("fresh-token\n"), 0o600))

	calls := 0
	rt := &serviceAccountTokenRetryRoundTripper{
		tokenFile: tokenFile,
		base: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			calls++
			if calls == 1 {
				return &http.Response{StatusCode: http.StatusUnauthorized, Body: io.NopCloser(strings.NewReader("Unauthorized")), Header: make(http.Header)}, nil
			}
			assert.Equal(t, "Bearer fresh-token", req.Header.Get("Authorization"))
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("ok")), Header: make(http.Header)}, nil
		}),
	}

	req, err := http.NewRequest(http.MethodGet, "https://kubernetes.default.svc/api/v1/pods", nil)
	require.NoError(t, err)
	resp, err := rt.RoundTrip(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, calls)
}

func TestDNSName(t *testing.T) {
	name, err := dnsName("wp_01he8bebctabr3kgk0qj36d2me_0_services_0")
	assert.NoError(t, err)
	assert.Equal(t, "wp-01he8bebctabr3kgk0qj36d2me-0-services-0", name)

	name, err = dnsName("a.0-AA")
	assert.NoError(t, err)
	assert.Equal(t, "a.0-aa", name)

	name, err = dnsName("wp-01he8bebctabr3kgk0qj36d2me-0-services-0.woodpecker-runtime.svc.cluster.local")
	assert.NoError(t, err)
	assert.Equal(t, "wp-01he8bebctabr3kgk0qj36d2me-0-services-0.woodpecker-runtime.svc.cluster.local", name)

	_, err = dnsName(".0-a")
	assert.ErrorIs(t, err, ErrDNSPatternInvalid)

	_, err = dnsName("ABC..DEF")
	assert.ErrorIs(t, err, ErrDNSPatternInvalid)

	_, err = dnsName("0.-a")
	assert.ErrorIs(t, err, ErrDNSPatternInvalid)

	_, err = dnsName("test-")
	assert.ErrorIs(t, err, ErrDNSPatternInvalid)

	_, err = dnsName("-test")
	assert.ErrorIs(t, err, ErrDNSPatternInvalid)

	_, err = dnsName("0-a.")
	assert.ErrorIs(t, err, ErrDNSPatternInvalid)

	_, err = dnsName("abc\\def")
	assert.ErrorIs(t, err, ErrDNSPatternInvalid)
}

func TestToDnsName(t *testing.T) {
	name, err := toDNSName("BUILD_AND_DEPLOY_0")
	assert.NoError(t, err)
	assert.Equal(t, "build-and-deploy-0", name)

	name, err = toDNSName("build and deploy")
	assert.NoError(t, err)
	assert.Equal(t, "build-and-deploy", name)

	name, err = toDNSName("build & deploy")
	assert.NoError(t, err)
	assert.Equal(t, "build--deploy", name)

	_, err = toDNSName("-build-and-deploy")
	assert.ErrorIs(t, err, ErrDNSPatternInvalid)
}

func TestGetHostnameOrEmpty(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Update repos", "update-repos"},
		{"MY_STEP", "my-step"},
		{"Build 🚀", ""},
	}

	for _, tt := range tests {
		got := getHostnameOrEmpty(tt.in)
		assert.Equal(t, tt.want, got, "input: %q", tt.in)
	}
}
