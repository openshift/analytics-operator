/*
Copyright 2023 Redhat.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceAccountRoleBinding defines Service account role binding properties
type ServiceAccountRoleBinding struct {
	Name            string `json:"name,omitempty"`
	ClusterRoleName string `json:"clusterrolename,omitempty"`
	SATokenName     string `json:"satokenname,omitempty"`
}

// ResourceConfig defines cpu/memory resource properties
type ResourceConfig struct {
	CPURequest    string `json:"cpurequest,omitempty"`
	CPULimit      string `json:"cpulimit,omitempty"`
	MemoryRequest string `json:"memoryrequest,omitempty"`
	MemoryLimit   string `json:"memorylimit,omitempty"`
}

// CronJobConfig defines configuration required to setup cronjob
type CronJobConfig struct {
	Schedule       string         `json:"schedule,omitempty"`
	Name           string         `json:"name,omitempty"`
	Image          string         `json:"image,omitempty"`
	AnomalyQueries string         `json:"anomalyqueries,omitempty"`
	LogLevel       string         `json:"loglevel,omitempty"`
	Resource       ResourceConfig `json:"resource,omitempty"`
}

// AnomalyEngineSpec defines the desired state of AnomalyEngine
type AnomalyEngineSpec struct {
	Namespace                 string                    `json:"namespace,omitempty"`
	ServiceAccountRoleBinding ServiceAccountRoleBinding `json:"serviceaccountrolebinding,omitempty"`
	CronJobConfig             CronJobConfig             `json:"cronjobconfig,omitempty"`
	AnomalyConfigmpName       string                    `json:"anomalyconfigmpname,omitempty"`
}

// AnomalyEngineStatus defines the observed state of AnomalyEngine
type AnomalyEngineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AnomalyEngine is the Schema for the anomalyengines API
type AnomalyEngine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AnomalyEngineSpec   `json:"spec,omitempty"`
	Status AnomalyEngineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AnomalyEngineList contains a list of AnomalyEngine
type AnomalyEngineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AnomalyEngine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AnomalyEngine{}, &AnomalyEngineList{})
}
