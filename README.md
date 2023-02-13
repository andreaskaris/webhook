## About

This is a test mutating webhook. The code was copied and modified from
https://github.com/kubernetes-sigs/controller-runtime/tree/master/examples/builtins

This webhook's only purpose is to inject capabilities into each pod that is requesting them as an annotation.

This webhook relies on cert-manager for certificate management and was only tested on OpenShift.

## Installation

Install the cert-manager on OpenShift:
~~~
make install-openshift-cert-manager
~~~

Make and push the image:
~~~
make podman-push
~~~

Deploy the webhook:
~~~
make deploy
~~~

Remove the webhook:
~~~
make undeploy
~~~
