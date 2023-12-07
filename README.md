# anomaly-operator
Operator that helps to deploy relevant component to detect anomaly in the cluster. 

## Description
We have created setup that works to deletct cluster anomaly using min/max as well as percentange change method to start with. Currently it works with metric data, in future we will incorporate logs/alerts as well. Currently we are targetting openshift as primary platform. 

Also currently we are focusing on Openshift cluster. 

## Getting Started
Login into openshift cluster using login command, below is the example of login command. 

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

2. Once installed you can verify it on cluster by running following command. 
```sh
kubectl get crd | grep anomaly
```

3. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):
```sh
make run
```

4. Install Instances of Custom Resources:
```sh
kubectl apply -f config/samples/backend_v1alpha1_anomalyengine.yaml
```

5. Operator should create related resources like namespace/role/serviceaccount/rolebinding/cronjob etc for given namespace in 4th step sample file. 

6. Based on given configuration (Ex : configname `anomalyconfigmpname` given in 4th step sample file), cronjob will try to find anomaly and if anything comeup it will add anomaly data into CRD storage which can be queried like below. 
```sh
oc get anomalydata -n osa-anomaly-detection
# you can see single anomaly data by below command 
oc describe anomalydata 2023-09-27-08-46-02-etcd-object-namespaces-namespaces -n osa-anomaly-detection
```


**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

### Running on the cluster using pre-published image

You can use the image from [quay.io](https://quay.io/repository/openshiftanalytics/observability-analytics-operator?tab=tags) to deploy kepler-operator.

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
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

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