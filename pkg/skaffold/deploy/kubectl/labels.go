/*
Copyright 2019 The Skaffold Authors

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

package kubectl

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// SetLabels add labels to a list of Kubernetes manifests.
func (l *ManifestList) SetLabels(labels map[string]string) (ManifestList, error) {
	replacer := newLabelsSetter(labels)

	updated, err := l.Visit(replacer)
	if err != nil {
		return nil, errors.Wrap(err, "setting labels")
	}

	logrus.Debugln("manifests with labels", updated.String())

	return updated, nil
}

type labelsSetter struct {
	labels map[string]string
}

func newLabelsSetter(labels map[string]string) *labelsSetter {
	return &labelsSetter{
		labels: labels,
	}
}

func (r *labelsSetter) Matches(key string) bool {
	return key == "metadata"
}

func (r *labelsSetter) NewValue(old interface{}) (bool, interface{}) {
	if len(r.labels) == 0 {
		return false, nil
	}

	metadata, ok := old.(map[interface{}]interface{})
	if !ok {
		return false, nil
	}

	l, present := metadata["labels"]
	if !present {
		metadata["labels"] = r.labels
		return true, metadata
	}

	labels, ok := l.(map[interface{}]interface{})
	if !ok {
		return false, nil
	}

	for k, v := range r.labels {
		labels[k] = v
	}

	return true, metadata
}
