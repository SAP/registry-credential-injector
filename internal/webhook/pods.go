/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and registry-credential-injector contributors
SPDX-License-Identifier: Apache-2.0
*/

package webhook

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodWebhookConfig struct {
	DefaultManaged    bool
	DefaultPullSecret string
	Client            kubernetes.Interface
}

type PodWebhook struct {
	Config *PodWebhookConfig
}

func (w *PodWebhook) MutateCreate(ctx context.Context, pod *corev1.Pod) error {
	return w.mutate(ctx, pod)
}

func (w *PodWebhook) MutateUpdate(ctx context.Context, oldPod *corev1.Pod, newPod *corev1.Pod) error {
	return w.mutate(ctx, newPod)
}

func (w *PodWebhook) mutate(ctx context.Context, pod *corev1.Pod) error {
	log, err := logr.FromContext(ctx)
	if err != nil {
		panic(err)
	}

	namespace := pod.Namespace
	name := pod.Name
	if name == "" {
		name = pod.GenerateName + "<new>"
	}

	ns, err := w.Config.Client.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		log.Error(err, "error reading pod namespace", "namespace", namespace, "pod", name)
		return errors.Wrapf(err, "error reading namespace of pod %s/%s", namespace, name)
	}

	managed := w.Config.DefaultManaged
	if v, ok := ns.Annotations["regcred-injector.x4.sap.com/managed"]; ok && v == "true" {
		managed = true
	}
	if v, ok := ns.Annotations["regcred-injector.cs.sap.com/managed"]; ok && v == "true" {
		managed = true
	}
	if v, ok := pod.Annotations["regcred-injector.x4.sap.com/managed"]; ok && v == "true" {
		managed = true
	}
	if v, ok := pod.Annotations["regcred-injector.cs.sap.com/managed"]; ok && v == "true" {
		managed = true
	}

	pullSecret := w.Config.DefaultPullSecret
	if v, ok := ns.Annotations["regcred-injector.x4.sap.com/pull-secret"]; ok {
		pullSecret = v
	}
	if v, ok := ns.Annotations["regcred-injector.cs.sap.com/pull-secret"]; ok {
		pullSecret = v
	}
	if v, ok := pod.Annotations["regcred-injector.x4.sap.com/pull-secret"]; ok {
		pullSecret = v
	}
	if v, ok := pod.Annotations["regcred-injector.cs.sap.com/pull-secret"]; ok {
		pullSecret = v
	}

	if managed {
		if pullSecret == "" {
			log.Info("unable to determine pull secret for managed pod", "namespace", namespace, "pod", name)
		} else {
			exists := false
			for _, ps := range pod.Spec.ImagePullSecrets {
				if ps.Name == pullSecret {
					exists = true
					break
				}
			}
			if !exists {
				log.Info("adding pull secret to pod", "namespace", namespace, "pod", name, "pullSecret", pullSecret)
				pod.Spec.ImagePullSecrets = append(pod.Spec.ImagePullSecrets, corev1.LocalObjectReference{Name: pullSecret})
			}
		}
	}

	return nil
}
