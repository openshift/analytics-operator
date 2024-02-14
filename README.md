# anomaly-operator
Operator that helps to deploy relevant components to detect anomalies in the cluster. 

## Description
We have created a setup that works to detect cluster anomalies using min/max as well as percentage change methods to start with. Currently, it works with metric data, in future we will incorporate logs/alerts as well. Currently, we are targeting OpenShift as the primary platform. 

Also currently we are focusing on Openshift clusters. 

## Getting Started
Login into an OpenShift cluster using the login command, below is an example of a login command. 

```sh
oc login --token=*** --server=https://example.com:6443
```

### Test It Out
1. Install the CRDs into the cluster:
```sh
make install
or
kubectl apply -f config/crd/bases/
```

2. Once installed you can verify it on the cluster by running the following command. 
```sh
kubectl get crd | grep anomaly
```

3. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):
```sh
make run
```

4. Install Instances of Custom Resources:
```sh
kubectl apply -f config/samples/observability-analytics_v1alpha1_anomalyengine.yaml
```

5. The operator should create related resources like namespace/role/serviceaccount/rolebinding/cronjob etc for the given namespace in the 4th step sample file. 

6. Based on the given configuration (Eg. configname `anomalyconfigmpname` given in the 4th step sample file), the cronjob will try to find anomalies and if anything comes up it will add anomaly data into CRD storage which can be queried like below. 
```sh
oc get anomalydata -n osa-anomaly-detection
# you can see single anomaly data with the following command 
oc describe anomalydata 2023-09-27-08-46-02-etcd-object-namespaces-namespaces -n osa-anomaly-detection
```


**NOTE:** You can also run this in one step by running: `make install run`

### Do e2e Testing
We have created e2e test script using which we can test operator is working as expected or not in an OpenShift cluster. 

**Prerequisite** : 
1. You should have access to the container registry for pushing images to the registry. 
2. You should have an OpenShift cluster, login into the cluster using the following command from your terminal. 
```sh
oc login --token=*** --server=https://example.com:6443
```
Then execute the below command to test your operator 
```sh
sh tests/run-e2e.sh 
```
This script performs the following operations to validate your operator works as expected in the OpenShift cluster.
1. Build your operator image and push it to the registry.
2. Deploy the operator into the OpenShift cluster using the created container.
3. Create CR for the anomaly engine.
4. Generate anomaly data and test whether the engine can detect an anomaly or not for Min/Max as well as Percentage Change type. 
5. Delete the operator and related resources from the cluster. 
6. Finally it should display the message "✅ All looks good :)" if everything works as expected.  


### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

### Running on the cluster using pre-published image

You can use the image from [quay.io](https://quay.io/repository/openshiftanalytics/observability-analytics-operator?tab=tags) to deploy analytics-operator.

```sh
make deploy OPERATOR_IMG=quay.io/openshiftanalytics/observability-analytics-operator:0.0.1
kubectl apply -f config/samples/
```

Alternatively, if you like to build and use your own image,

```sh
make operator-build operator-push IMG_BASE=<some-registry>
make deploy OPERATOR_IMG=<some-registry>/observability-analytics-operator:0.0.1
kubectl apply -f config/samples/
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provides a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

## License

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
