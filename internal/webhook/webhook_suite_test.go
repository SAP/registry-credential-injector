/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and registry-credential-injector contributors
SPDX-License-Identifier: Apache-2.0
*/

package webhook_test

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"

	"github.com/sap/registry-credential-injector/internal/webhook"
)

func TestWebhook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webhook Suite")
}

var _ = Describe("Webhook", func() {
	Context("Generic Webhook", func() {
		var mutate func(pod *corev1.Pod) error
		var pod1 *corev1.Pod
		var pod2 *corev1.Pod
		var pod3 *corev1.Pod
		var pod4 *corev1.Pod

		BeforeEach(func() {
			podWebhook := &webhook.PodWebhook{
				Config: &webhook.PodWebhookConfig{
					DefaultPullSecret: "regcred-default",
					Client: kubernetesfake.NewSimpleClientset(
						&corev1.Namespace{
							ObjectMeta: metav1.ObjectMeta{
								Name: "ns1",
							},
						},
						&corev1.Namespace{
							ObjectMeta: metav1.ObjectMeta{
								Name: "ns2",
								Annotations: map[string]string{
									"regcred-injector.cs.sap.com/managed": "true",
								},
							},
						},
						&corev1.Namespace{
							ObjectMeta: metav1.ObjectMeta{
								Name: "ns3",
								Annotations: map[string]string{
									"regcred-injector.cs.sap.com/pull-secret": "regcred-ns",
								},
							},
						},
						&corev1.Namespace{
							ObjectMeta: metav1.ObjectMeta{
								Name: "ns4",
								Annotations: map[string]string{
									"regcred-injector.cs.sap.com/managed":     "true",
									"regcred-injector.cs.sap.com/pull-secret": "regcred-ns",
								},
							},
						},
					),
				},
			}

			mutate = func(pod *corev1.Pod) error {
				ctx := logr.NewContext(context.TODO(), logr.Discard())
				return podWebhook.MutateCreate(ctx, pod)
			}

			pod1 = &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod1",
				},
				Spec: corev1.PodSpec{},
			}
			pod2 = &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod2",
					Annotations: map[string]string{
						"regcred-injector.cs.sap.com/managed": "true",
					},
				},
				Spec: corev1.PodSpec{},
			}
			pod3 = &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod3",
					Annotations: map[string]string{
						"regcred-injector.cs.sap.com/pull-secret": "regcred-pod",
					},
				},
				Spec: corev1.PodSpec{},
			}
			pod4 = &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod4",
					Annotations: map[string]string{
						"regcred-injector.cs.sap.com/managed":     "true",
						"regcred-injector.cs.sap.com/pull-secret": "regcred-pod",
					},
				},
				Spec: corev1.PodSpec{},
			}
		})

		// pod1
		It("should not add a pull secret", func() {
			pod1.Namespace = "ns1"
			err := mutate(pod1)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod1.Spec.ImagePullSecrets).To(BeEmpty())
		})
		It("should add default pull secret", func() {
			pod1.Namespace = "ns2"
			err := mutate(pod1)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod1.Spec.ImagePullSecrets).To(HaveLen(1))
			Expect(pod1.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-default"}))
		})
		It("should not add a pull secret", func() {
			pod1.Namespace = "ns3"
			err := mutate(pod1)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod1.Spec.ImagePullSecrets).To(BeEmpty())
		})
		It("should add namespace pull secret", func() {
			pod1.Namespace = "ns4"
			err := mutate(pod1)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod1.Spec.ImagePullSecrets).To(HaveLen(1))
			Expect(pod1.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-ns"}))
		})

		// pod2
		It("should add default pull secret", func() {
			pod2.Namespace = "ns1"
			err := mutate(pod2)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod2.Spec.ImagePullSecrets).To(HaveLen(1))
			Expect(pod2.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-default"}))
		})
		It("should add default pull secret", func() {
			pod2.Namespace = "ns2"
			err := mutate(pod2)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod2.Spec.ImagePullSecrets).To(HaveLen(1))
			Expect(pod2.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-default"}))
		})
		It("should add namespace pull secret", func() {
			pod2.Namespace = "ns3"
			err := mutate(pod2)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod2.Spec.ImagePullSecrets).To(HaveLen(1))
			Expect(pod2.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-ns"}))
		})
		It("should add namespace pull secret", func() {
			pod2.Namespace = "ns4"
			err := mutate(pod2)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod2.Spec.ImagePullSecrets).To(HaveLen(1))
			Expect(pod2.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-ns"}))
		})

		// pod3
		It("should not add a pull secret", func() {
			pod3.Namespace = "ns1"
			err := mutate(pod3)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod3.Spec.ImagePullSecrets).To(BeEmpty())
		})
		It("should add pod pull secret", func() {
			pod3.Namespace = "ns2"
			err := mutate(pod3)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod3.Spec.ImagePullSecrets).To(HaveLen(1))
			Expect(pod3.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-pod"}))
		})
		It("should not add a pull secret", func() {
			pod3.Namespace = "ns3"
			err := mutate(pod3)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod3.Spec.ImagePullSecrets).To(BeEmpty())
		})
		It("should add pod pull secret", func() {
			pod3.Namespace = "ns4"
			err := mutate(pod3)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod3.Spec.ImagePullSecrets).To(HaveLen(1))
			Expect(pod3.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-pod"}))
		})

		// pod4
		It("should add pod pull secret", func() {
			pod4.Namespace = "ns1"
			err := mutate(pod4)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod4.Spec.ImagePullSecrets).To(HaveLen(1))
			Expect(pod4.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-pod"}))
		})
		It("should add pod pull secret", func() {
			pod4.Namespace = "ns2"
			err := mutate(pod4)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod4.Spec.ImagePullSecrets).To(HaveLen(1))
			Expect(pod4.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-pod"}))
		})
		It("should add pod pull secret", func() {
			pod4.Namespace = "ns3"
			err := mutate(pod4)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod4.Spec.ImagePullSecrets).To(HaveLen(1))
			Expect(pod4.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-pod"}))
		})
		It("should add pod pull secret", func() {
			pod4.Namespace = "ns4"
			err := mutate(pod4)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod4.Spec.ImagePullSecrets).To(HaveLen(1))
			Expect(pod4.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-pod"}))
		})

		// various
		It("should add pull secret only once", func() {
			pod4.Namespace = "ns1"
			err := mutate(pod4)
			Expect(err).NotTo(HaveOccurred())
			err = mutate(pod4)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod4.Spec.ImagePullSecrets).To(HaveLen(1))
			Expect(pod4.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-pod"}))
		})
		It("should keep existing pull secrets", func() {
			pod4.Namespace = "ns1"
			pod4.Spec.ImagePullSecrets = []corev1.LocalObjectReference{{Name: "regcred-existing"}}
			err := mutate(pod4)
			Expect(err).NotTo(HaveOccurred())
			Expect(pod4.Spec.ImagePullSecrets).To(HaveLen(2))
			Expect(pod4.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-existing"}))
			Expect(pod4.Spec.ImagePullSecrets).To(ContainElement(corev1.LocalObjectReference{Name: "regcred-pod"}))
		})
	})
})
