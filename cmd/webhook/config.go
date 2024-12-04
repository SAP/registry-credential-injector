/*
SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and registry-credential-injector contributors
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"flag"

	"github.com/pkg/errors"

	"k8s.io/client-go/kubernetes"

	"github.com/sap/registry-credential-injector/internal/kubeconfig"
	"github.com/sap/registry-credential-injector/internal/webhook"
)

var kubeconfigPath string
var defaultManaged bool
var defaultPullSecret string

func defineFlags() {
	flag.StringVar(&kubeconfigPath, "kubeconfig", kubeconfigPath, "File containing kubeconfig to be used by the webhook")
	flag.StringVar(&defaultPullSecret, "default-pull-secret", defaultPullSecret, "Name of the default pull secret to supply to managed pods")
}

func buildConfigFromFlags() (*webhook.PodWebhookConfig, error) {
	config := &webhook.PodWebhookConfig{}

	restConfig, err := kubeconfig.New(kubeconfigPath).Load()
	if err != nil {
		return nil, errors.Wrapf(err, "error loading kubeconfig %s", kubeconfigPath)
	}
	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error building kubernetes client")
	}
	config.Client = client

	config.DefaultManaged = defaultManaged
	config.DefaultPullSecret = defaultPullSecret

	return config, nil
}
