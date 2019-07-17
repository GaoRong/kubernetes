/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package disableservicelinks

import (
	"testing"

	"k8s.io/apiserver/pkg/admission"
	admissiontesting "k8s.io/apiserver/pkg/admission/testing"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/core/helper"
)

func TestDisableServiceLinks(t *testing.T) {
	handler := admissiontesting.WithReinvocationTesting(t, NewDisableServiceLinks())

	enableServiceLinks := true
	disableServiceLinks := false
	tests := []struct {
		description  string
		requestedPod api.Pod
		expectedPod  api.Pod
	}{
		{
			description: "pod has no tolerations, expect add tolerations for `not-ready:NoExecute` and `unreachable:NoExecute`",
			requestedPod: api.Pod{
				Spec: api.PodSpec{
					EnableServiceLinks: &disableServiceLinks,
				},
			},
			expectedPod: api.Pod{
				Spec: api.PodSpec{
					EnableServiceLinks: &enableServiceLinks,
				},
			},
		},
	}

	for _, test := range tests {
		err := handler.Admit(admission.NewAttributesRecord(&test.requestedPod, nil, api.Kind("Pod").WithVersion("version"), "foo", "name", api.Resource("pods").WithVersion("version"), "", "ignored", nil, false, nil), nil)
		if err != nil {
			t.Errorf("[%s]: unexpected error %v for pod %+v", test.description, err, test.requestedPod)
		}

		if !helper.Semantic.DeepEqual(test.expectedPod.Spec.Tolerations, test.requestedPod.Spec.Tolerations) {
			t.Errorf("[%s]: expected %#v got %#v", test.description, test.expectedPod.Spec.Tolerations, test.requestedPod.Spec.Tolerations)
		}
	}
}

func TestHandles(t *testing.T) {
	handler := NewDisableServiceLinks()
	tests := map[admission.Operation]bool{
		admission.Update:  true,
		admission.Create:  true,
		admission.Delete:  false,
		admission.Connect: false,
	}
	for op, expected := range tests {
		result := handler.Handles(op)
		if result != expected {
			t.Errorf("Unexpected result for operation %s: %v\n", op, result)
		}
	}
}
