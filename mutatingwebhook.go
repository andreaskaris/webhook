/*
Copyright 2018 The Kubernetes Authors.
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

package main

import (
	"context"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io

const (
	webhookAnnotation = "webhook/capabilities"
)

// podAnnotator annotates Pods
type podAnnotator struct{}

func (a *podAnnotator) Default(ctx context.Context, obj runtime.Object) error {
	log := logf.FromContext(ctx)
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod but got a %T", obj)
	}

	log.Info("Inspecting pod",
		"namespace", pod.ObjectMeta.Namespace,
		"name", pod.ObjectMeta.Name,
		"ownerReferences", pod.ObjectMeta.OwnerReferences,
		"labels", pod.ObjectMeta.Labels)

	var capabilities []string
	if pod.Annotations != nil {
		for key, content := range pod.Annotations {
			if key == webhookAnnotation {
				if err := json.Unmarshal([]byte(content), &capabilities); err != nil {
					return fmt.Errorf("could not unmarshal content %s, err: %q", content, err)
				}
			}
		}
	}
	if len(capabilities) > 0 {
		for _, container := range pod.Spec.Containers {
			if container.SecurityContext == nil {
				container.SecurityContext = &corev1.SecurityContext{}
			}
			if container.SecurityContext.Capabilities == nil {
				container.SecurityContext.Capabilities = &corev1.Capabilities{}
			}
			add := container.SecurityContext.Capabilities.Add
			for _, c := range capabilities {
				log.Info("Adding capabilities to pod", "capabilities", c)
				add = append(add, corev1.Capability(c))
			}
			container.SecurityContext.Capabilities.Add = add
		}
	}

	return nil
}
