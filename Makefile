build:
	go build -o bin/main . 

podman-build:
	podman build . -t quay.io/akaris/webhook:latest

podman-push: podman-build
	podman push quay.io/akaris/webhook:latest

install-openshift-cert-manager:
	oc apply -f resources/cert-manager.yaml

uninstall-openshift-cert-manager:
	oc delete -f resources/cert-manager.yaml

create-cert-manager-cert:
	oc apply -f resources/webhook-cert.yaml

deploy-webhook:
	oc apply -f resources/service.yaml
	oc apply -f resources/webhook.yaml

deploy: install-openshift-cert-manager create-cert-manager-cert deploy-webhook

undeploy:
	oc delete -f resources/webhook.yaml
	oc delete -f resources/service.yaml
	oc delete -f resources/webhook-cert.yaml

logs:
	oc logs -n webhook -l deployment=webhook -f
