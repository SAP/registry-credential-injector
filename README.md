# Kubernetes Registry Credential Injector

[![REUSE status](https://api.reuse.software/badge/github.com/SAP/registry-credential-injector)](https://api.reuse.software/info/github.com/SAP/registry-credential-injector)

## About this project

This service can act as a [Mutating Kubernetes Admission Webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers) for pods, and allows to dynamically inject image pull secrets into specific pods at creation time.

Pods for which the admission webhook is called by the Kubernetes API server will be changed only if:
- there is an annotation `regcred-injector.cs.sap.com/managed: true` set at the pod's namespace or
- there is an annotation  `regcred-injector.cs.sap.com/managed: true` set at the pod itself.

In that case, the webhook  will determine the name of the pull secret to be injected as follows:
- if the inspected pod has an annotation `regcred-injector.cs.sap.com/pull-secret`, its value will be used for the pull secret
- otherwise, if the pod's namespace has the annotation `regcred-injector.cs.sap.com/pull-secret`, its value will be used for the pull secret
- otherwise, if specified, the default pull secret value (specified by command line flag) will be used
- if no pull secret value was found by the above sources, the pod will not be changed.

Note: in case this webhook has to reliably work with pods that are created or mutated by other webhooks, this one should be registered with `reinvocationPolicy: IfNeeded`.

**Command line flags**

|Flag                  |Optional|Default                            |Description                                                 |
|----------------------|--------|-----------------------------------|------------------------------------------------------------|
|-kubeconfig           |yes     |Usual kubeconfig fallback locations|Path to kubeconfig file                                     |
|-bind-address string  |yes     |:2443                              |Webhook bind address                                        |
|-tls-key-file         |no      |-                                  |File containing the TLS private key used for SSL termination|
|-tls-cert-file        |no      |-                                  |File containing the TLS certificate matching the private key|
|-default-pull-secret  |yes     |-                                  |Name of the default pull secret to be injected              |


## Requirements and Setup

The recommended deployment method is to use the [Helm chart](https://github.com/sap/registry-credential-injector-helm):

```bash
helm upgrade -i registry-credential-injector oci://ghcr.io/sap/registry-credential-injector-helm/registry-credential-injector
```

## Documentation
 
The API reference is here: [https://pkg.go.dev/github.com/sap/registry-credential-injector](https://pkg.go.dev/github.com/sap/registry-credential-injector).

## Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/SAP/registry-credential-injector/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/SAP/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright 2024 SAP SE or an SAP affiliate company and registry-credential-injector contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/SAP/registry-credential-injector).
