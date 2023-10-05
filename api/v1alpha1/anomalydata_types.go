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

// MetricData defined property related to metric
type MetricData struct {
	Timestamp         int64   `json:"timestamp,omitempty"`
	LatestValue       float64 `json:"latestvalue,omitempty"`
	PercentageChange  float64 `json:"percentagechange,omitempty"`
	PrevDataMeanValue float64 `json:"prevdatameanvalue,omitempty"`
	GroupedData       string  `json:"groupeddata,omitempty"`
	DataPoints        string  `json:"datapoints,omitempty"`
}

// AnomalyConfig defines the properties set while declaring anomaly defination
type AnomalyConfig struct {
	Query               string  `json:"query,omitempty"`
	Min                 int64   `json:"min,omitempty"`
	Max                 int64   `json:"max,omitempty"`
	Step                int16   `json:"step,omitempty"`
	PercentageChange    float64 `json:"percentagechange,omitempty"`
	PeriodRange         float64 `json:"periodrange,omitempty"`
	HaveMultiResultData bool    `json:"havemultiresultdata,omitempty"`
}

// AnomalyDataSpec defines the desired state of AnomalyData
type AnomalyDataSpec struct {
	AnomalyName string        `json:"anomalyname,omitempty"`
	Method      string        `json:"method,omitempty"`
	Config      AnomalyConfig `json:"config,omitempty"`
	MetricData  MetricData    `json:"metricdata,omitempty"`
}

// AnomalyDataStatus defines the observed state of AnomalyData
type AnomalyDataStatus struct {

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AnomalyData is the Schema for the anomalydata API
type AnomalyData struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AnomalyDataSpec   `json:"spec,omitempty"`
	Status AnomalyDataStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AnomalyDataList contains a list of AnomalyData
type AnomalyDataList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AnomalyData `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AnomalyData{}, &AnomalyDataList{})
}
