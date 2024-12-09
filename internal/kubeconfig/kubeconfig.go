/*
SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and registry-credential-injector contributors
SPDX-License-Identifier: Apache-2.0
*/

package kubeconfig

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Kubeconfig struct {
	path string
}

func New(path string) *Kubeconfig {
	return &Kubeconfig{
		path: path,
	}
}

func (k *Kubeconfig) Load() (*rest.Config, error) {
	if k.path == "" {
		return rest.InClusterConfig()
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(&clientcmd.ClientConfigLoadingRules{ExplicitPath: k.path}, nil).ClientConfig()
}
