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

// ServiceAccountRoleBinding defines Service account role binding properties to support AnomalyEngine
type ServiceAccountRoleBinding struct {
	// Name of the Service Account
	Name string `json:"name,omitempty"`
	//  Name of the Cluster Role which have view/read access to mornitoring/thanos-api
	ClusterRoleName string `json:"clusterrolename,omitempty"`
}

// ResourceConfig defines cpu/memory resource properties required for AnomalyEngine pod
type ResourceConfig struct {
	CPURequest    string `json:"cpurequest,omitempty"`
	CPULimit      string `json:"cpulimit,omitempty"`
	MemoryRequest string `json:"memoryrequest,omitempty"`
	MemoryLimit   string `json:"memorylimit,omitempty"`
}

// CronJobConfig defines configuration required to setup cronjob
type CronJobConfig struct {
	// Schedule for the cronjob
	Schedule string `json:"schedule,omitempty"`
	// Name of the cronjob
	Name string `json:"name,omitempty"`
	// Comma-separated keys from anomalyqueryconfiguration. If not defined, the system will go through all the defined configurations. For example, if there are five configurations defined but we only want to run two for the time being, then those specific keys need to be defined here.
	AnomalyQueries string `json:"anomalyqueries,omitempty"`
	// Pod log level - DEBUG/INFO/ERROR etc
	LogLevel string         `json:"loglevel,omitempty"`
	Resource ResourceConfig `json:"resource,omitempty"`
}

// AnomalyEngineSpec defines the desired state of AnomalyEngine
type AnomalyEngineSpec struct {
	// The namespace under which Anomaly Engine cronjobs will run
	Namespace                 string                    `json:"namespace,omitempty"`
	ServiceAccountRoleBinding ServiceAccountRoleBinding `json:"serviceaccountrolebinding,omitempty"`
	CronJobConfig             CronJobConfig             `json:"cronjobconfig,omitempty"`
	// AnomalyQueryConfiguration defines the query/configuration to detect anomaly
	// You can take a look on below link to understand how to define these data structure
	// https://github.com/openshift/incluster-anomaly-detection/tree/main#understanding-anomaly-configurations
	AnomalyQueryConfiguration string `json:"anomalyqueryconfiguration,omitempty"`
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
