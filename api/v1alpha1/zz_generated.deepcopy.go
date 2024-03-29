//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2024 Redhat.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnomalyConfig) DeepCopyInto(out *AnomalyConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnomalyConfig.
func (in *AnomalyConfig) DeepCopy() *AnomalyConfig {
	if in == nil {
		return nil
	}
	out := new(AnomalyConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnomalyData) DeepCopyInto(out *AnomalyData) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnomalyData.
func (in *AnomalyData) DeepCopy() *AnomalyData {
	if in == nil {
		return nil
	}
	out := new(AnomalyData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AnomalyData) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnomalyDataList) DeepCopyInto(out *AnomalyDataList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AnomalyData, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnomalyDataList.
func (in *AnomalyDataList) DeepCopy() *AnomalyDataList {
	if in == nil {
		return nil
	}
	out := new(AnomalyDataList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AnomalyDataList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnomalyDataSpec) DeepCopyInto(out *AnomalyDataSpec) {
	*out = *in
	out.Config = in.Config
	out.MetricData = in.MetricData
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnomalyDataSpec.
func (in *AnomalyDataSpec) DeepCopy() *AnomalyDataSpec {
	if in == nil {
		return nil
	}
	out := new(AnomalyDataSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnomalyDataStatus) DeepCopyInto(out *AnomalyDataStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnomalyDataStatus.
func (in *AnomalyDataStatus) DeepCopy() *AnomalyDataStatus {
	if in == nil {
		return nil
	}
	out := new(AnomalyDataStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnomalyEngine) DeepCopyInto(out *AnomalyEngine) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnomalyEngine.
func (in *AnomalyEngine) DeepCopy() *AnomalyEngine {
	if in == nil {
		return nil
	}
	out := new(AnomalyEngine)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AnomalyEngine) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnomalyEngineList) DeepCopyInto(out *AnomalyEngineList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AnomalyEngine, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnomalyEngineList.
func (in *AnomalyEngineList) DeepCopy() *AnomalyEngineList {
	if in == nil {
		return nil
	}
	out := new(AnomalyEngineList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AnomalyEngineList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnomalyEngineSpec) DeepCopyInto(out *AnomalyEngineSpec) {
	*out = *in
	out.ServiceAccountRoleBinding = in.ServiceAccountRoleBinding
	out.CronJobConfig = in.CronJobConfig
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnomalyEngineSpec.
func (in *AnomalyEngineSpec) DeepCopy() *AnomalyEngineSpec {
	if in == nil {
		return nil
	}
	out := new(AnomalyEngineSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnomalyEngineStatus) DeepCopyInto(out *AnomalyEngineStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnomalyEngineStatus.
func (in *AnomalyEngineStatus) DeepCopy() *AnomalyEngineStatus {
	if in == nil {
		return nil
	}
	out := new(AnomalyEngineStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CronJobConfig) DeepCopyInto(out *CronJobConfig) {
	*out = *in
	out.Resource = in.Resource
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CronJobConfig.
func (in *CronJobConfig) DeepCopy() *CronJobConfig {
	if in == nil {
		return nil
	}
	out := new(CronJobConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricData) DeepCopyInto(out *MetricData) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricData.
func (in *MetricData) DeepCopy() *MetricData {
	if in == nil {
		return nil
	}
	out := new(MetricData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceConfig) DeepCopyInto(out *ResourceConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceConfig.
func (in *ResourceConfig) DeepCopy() *ResourceConfig {
	if in == nil {
		return nil
	}
	out := new(ResourceConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceAccountRoleBinding) DeepCopyInto(out *ServiceAccountRoleBinding) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceAccountRoleBinding.
func (in *ServiceAccountRoleBinding) DeepCopy() *ServiceAccountRoleBinding {
	if in == nil {
		return nil
	}
	out := new(ServiceAccountRoleBinding)
	in.DeepCopyInto(out)
	return out
}
