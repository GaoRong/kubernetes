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
	"io"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/admission"
	api "k8s.io/kubernetes/pkg/apis/core"
)

// PluginName indicates name of admission plugin.
const PluginName = "DisableServiceLinks"

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register(PluginName, func(config io.Reader) (admission.Interface, error) {
		return NewDisableServiceLinks(), nil
	})
}

// DisableServiceLinks is an implementation of admission.Interface.
// It looks at all new pods and overrides each pod's enableServiceLinks to false.
type DisableServiceLinks struct {
	*admission.Handler
}

var _ admission.MutationInterface = &DisableServiceLinks{}

// Admit makes an admission decision based on the request attributes
func (a *DisableServiceLinks) Admit(attributes admission.Attributes, o admission.ObjectInterfaces) (err error) {
	// Ignore all calls to subresources or resources other than pods.
	if shouldIgnore(attributes) {
		return nil
	}
	pod, ok := attributes.GetObject().(*api.Pod)
	if !ok {
		return apierrors.NewBadRequest("Resource was marked with kind Pod but was unable to be converted")
	}
	enableServiceLinks := false
	pod.Spec.EnableServiceLinks = &enableServiceLinks
	return nil
}

func shouldIgnore(attributes admission.Attributes) bool {
	// Ignore all calls to subresources or resources other than pods.
	if len(attributes.GetSubresource()) != 0 || attributes.GetResource().GroupResource() != api.Resource("pods") {
		return true
	}

	return false
}

// NewDisableServiceLinks creates a new always pull images admission control handler
func NewDisableServiceLinks() *DisableServiceLinks {
	return &DisableServiceLinks{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}
}
