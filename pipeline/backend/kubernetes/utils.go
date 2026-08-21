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

package kubernetes

import (
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	kube_core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	kube_client_cmd "k8s.io/client-go/tools/clientcmd"
)

const (
	maxDNSLabelLen                    = 63
	inClusterServiceAccountTokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

var (
	dnsPattern = regexp.MustCompile(
		`^[a-z0-9]` + // must start with
			`([-a-z0-9]*[a-z0-9])?` + // inside can als contain -
			`(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`, // allow the same pattern as before with dots in between but only one dot
	)
	dnsLabelPattern         = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	dnsDisallowedCharacters = regexp.MustCompile(`[^-^.a-z0-9]+`)
	ErrDNSPatternInvalid    = errors.New("name is not a valid kubernetes DNS name")
)

type serviceAccountTokenRetryRoundTripper struct {
	base      http.RoundTripper
	tokenFile string
}

func (r *serviceAccountTokenRetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := r.base.RoundTrip(req)
	if err != nil || resp.StatusCode != http.StatusUnauthorized || (req.Method != http.MethodGet && req.Method != http.MethodHead) {
		return resp, err
	}

	token, readErr := os.ReadFile(r.tokenFile)
	if readErr != nil {
		return resp, err
	}

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	retryReq := req.Clone(req.Context())
	retryReq.Header = req.Header.Clone()
	retryReq.Header.Set("Authorization", "Bearer "+strings.TrimSpace(string(token)))
	return r.base.RoundTrip(retryReq)
}

func getHostnameOrEmpty(name string) string {
	clean, _ := toDNSName(name)
	if clean == "" {
		clean = strings.ToLower(name)
	}
	clean = strings.ReplaceAll(clean, ".", "-")

	if len(clean) > maxDNSLabelLen {
		clean = clean[:maxDNSLabelLen]
	}

	clean = strings.Trim(clean, "-")

	if dnsLabelPattern.MatchString(clean) {
		return clean
	}
	return ""
}

func dnsName(i string) (string, error) {
	res := strings.ToLower(strings.ReplaceAll(i, "_", "-"))

	if found := dnsPattern.FindStringIndex(res); found == nil {
		return "", ErrDNSPatternInvalid
	}

	return res, nil
}

func toDNSName(in string) (string, error) {
	lower := strings.ToLower(in)
	withoutUnderscores := strings.ReplaceAll(lower, "_", "-")
	withoutSpaces := strings.ReplaceAll(withoutUnderscores, " ", "-")
	almostDNS := dnsDisallowedCharacters.ReplaceAllString(withoutSpaces, "")
	return dnsName(almostDNS)
}

func isImagePullBackOffState(pod *kube_core_v1.Pod) bool {
	for _, containerState := range pod.Status.ContainerStatuses {
		if containerState.State.Waiting != nil {
			if containerState.State.Waiting.Reason == "ImagePullBackOff" {
				return true
			}
		}
	}

	return false
}

func isInvalidImageName(pod *kube_core_v1.Pod) bool {
	for _, containerState := range pod.Status.ContainerStatuses {
		if containerState.State.Waiting != nil {
			if containerState.State.Waiting.Reason == "InvalidImageName" {
				return true
			}
		}
	}

	return false
}

// getClientOutOfCluster returns a k8s client set to the request from outside of cluster.
func getClientOutOfCluster() (kubernetes.Interface, error) {
	kubeConfigPath := os.Getenv("KUBECONFIG") // cspell:words KUBECONFIG
	if kubeConfigPath == "" {
		kubeConfigPath = os.Getenv("HOME") + "/.kube/config"
	}

	// use the current context in kube config
	config, err := kube_client_cmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

// getClientInsideOfCluster returns a k8s client set to the request from inside of cluster.
func getClientInsideOfCluster() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	previousWrap := config.WrapTransport
	config.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
		if previousWrap != nil {
			rt = previousWrap(rt)
		}
		return &serviceAccountTokenRetryRoundTripper{base: rt, tokenFile: inClusterServiceAccountTokenFile}
	}

	return kubernetes.NewForConfig(config)
}

func newBool(val bool) *bool {
	ptr := new(bool)
	*ptr = val
	return ptr
}

func newInt64(val int64) *int64 {
	ptr := new(int64)
	*ptr = val
	return ptr
}
