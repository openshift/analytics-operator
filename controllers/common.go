package controllers

import (
	"context"

	backendv1alpha1 "github.com/k8s-analytics/anomaly-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var anomalyConfigFileName = "anomaly_config.yaml"
var storageAdminRoleName = "osa-crd-admin"
var storageServiceAccountName = "osa-crd-sa-user"
var configmapName = "osa-anomaly-config"

func (r *AnomalyEngineReconciler) ensureNamespace(request reconcile.Request, instance *backendv1alpha1.AnomalyEngine) (*reconcile.Result, error) {

	var log = logf.Log.WithName("anomay-detection -> Namespace")

	// check namespace exist
	namespaceName := instance.Spec.Namespace
	namespaceFound := &corev1.Namespace{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name: namespaceName,
	}, namespaceFound)
	log.Info("namespace check completed", "name", namespaceName)

	if err != nil && errors.IsNotFound(err) {
		// Create the namespace
		log.Info("Creating a new namespace", "name", namespaceName)
		err = r.Client.Create(context.TODO(), &corev1.Namespace{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Namespace",
				APIVersion: corev1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: namespaceName,
			},
		})

		if err != nil {
			// namespace creation failed
			log.Error(err, "Failed to create new namespace", "name", namespaceName)
			return &reconcile.Result{}, err
		}
		log.Info("Sucessfully created Namespace", "name", namespaceName)
	} else if err != nil {
		// Error that isn't due to the namespace not exist
		log.Error(err, "Failed to get namespace data", "name", namespaceName)
		return &reconcile.Result{}, err
	} else {
		log.Info("Namespace already exist", "name", namespaceName)
	}
	// Namespace craetion was successful
	return nil, nil
}

func (r *AnomalyEngineReconciler) getServiceAccountForStorage(instance *backendv1alpha1.AnomalyEngine) *corev1.ServiceAccount {
	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "osa-crd-sa-user",
			Namespace: instance.Spec.Namespace,
		},
	}
	controllerutil.SetControllerReference(instance, sa, r.Scheme)
	return sa
}

func (r *AnomalyEngineReconciler) getServiceAccountForAnomalyEngine(instance *backendv1alpha1.AnomalyEngine) *corev1.ServiceAccount {
	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Spec.ServiceAccountRoleBinding.Name,
			Namespace: instance.Spec.Namespace,
		},
		AutomountServiceAccountToken: &[]bool{true}[0],
	}
	controllerutil.SetControllerReference(instance, sa, r.Scheme)
	return sa
}

func (r *AnomalyEngineReconciler) ensureServiceAccount(request reconcile.Request,
	instance *backendv1alpha1.AnomalyEngine,
	sa *corev1.ServiceAccount,
	createSecret bool,
) (*reconcile.Result, error) {

	var log = logf.Log.WithName("anomay-detection -> SA ->")
	found := &corev1.ServiceAccount{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      sa.Name,
		Namespace: sa.Namespace,
	}, found)
	log.Info("SA check completed", "Name", sa.Name)
	if err != nil && errors.IsNotFound(err) {

		// Create the Role
		log.Info("Creating a new SA", "namespace", sa.Namespace, "name", sa.Name)
		err = r.Client.Create(context.TODO(), sa)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new SA", "namespace", sa.Namespace, "name", sa.Name)
			return &reconcile.Result{}, err
		}
		log.Info("SA creation completed", "namespace", sa.Namespace, "name", sa.Name)

		if createSecret {

			// create secret for the Service Account
			err = r.Client.Create(context.TODO(), &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      sa.Name + "-token",
					Namespace: instance.Spec.Namespace,
					Annotations: map[string]string{
						corev1.ServiceAccountNameKey: sa.Name,
					},
				},
				Type: corev1.SecretTypeServiceAccountToken,
			})

			if err != nil {
				// service account creation failed
				log.Error(err, "Failed to create secret for Service Account", "name", sa.Name)
				return &reconcile.Result{}, err
			}
			log.Info("Sucessfully created secret for Service Account", "name", sa.Name)
		}

	} else if err != nil {
		// Error that isn't due to the SA not existing
		log.Error(err, "Failed to get SA")
		return &reconcile.Result{}, err
	} else {
		log.Info("SA already exist", "namespace", sa.Namespace, "name", sa.Name)
	}
	// Creation was successful
	return nil, nil
}

func (r *AnomalyEngineReconciler) getStorageAdminRole(instance *backendv1alpha1.AnomalyEngine) *rbac.Role {
	role := &rbac.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "osa-crd-admin",
			Namespace: instance.Spec.Namespace,
		},
		Rules: []rbac.PolicyRule{
			{
				Verbs:     []string{"get", "watch", "list", "create", "update", "delete"},
				Resources: []string{"anomalydata"},
				APIGroups: []string{"backend.anomaly.io"},
			},
		},
	}

	controllerutil.SetControllerReference(instance, role, r.Scheme)
	return role
}

func (r *AnomalyEngineReconciler) ensureRole(request reconcile.Request,
	instance *backendv1alpha1.AnomalyEngine,
	role *rbac.Role,
) (*reconcile.Result, error) {

	var log = logf.Log.WithName("anomay-detection -> Role ->")
	found := &rbac.Role{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      role.Name,
		Namespace: role.Namespace,
	}, found)
	log.Info("Role check completed", "name", role.Name)
	if err != nil && errors.IsNotFound(err) {

		// Create the Role
		log.Info("Creating a new Role", "namespace", role.Namespace, "name", role.Name)
		err = r.Client.Create(context.TODO(), role)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Role", "namespace", role.Namespace, "name", role.Name)
			return &reconcile.Result{}, err
		}
		log.Info("Role creation completed", "namespace", role.Namespace, "name", role.Name)
	} else if err != nil {
		// Error that isn't due to the Role not existing
		log.Error(err, "Failed to get Role")
		return &reconcile.Result{}, err
	} else {
		log.Info("Role already exist", "namespace", role.Namespace, "name", role.Name)
	}
	// Creation was successful
	return nil, nil
}

func (r *AnomalyEngineReconciler) ensureRoleBinding(request reconcile.Request, instance *backendv1alpha1.AnomalyEngine,
	rb *rbac.RoleBinding) (*reconcile.Result, error) {

	var log = logf.Log.WithName("anomay-detection -> RoleBinding ->")

	//  check role binding exist
	roleBindingFound := &rbac.RoleBinding{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      rb.Name,
		Namespace: rb.Namespace,
	}, roleBindingFound)
	log.Info("RoleBinding check completed", "name", rb.Name)

	if err != nil && errors.IsNotFound(err) {
		// Create the RoleBinding
		log.Info("Creating a new RoleBinding", "name", rb.Name)
		err = r.Client.Create(context.TODO(), rb)

		if err != nil {
			// RoleBinding creation failed
			log.Error(err, "Failed to create role binding", "name", rb.Name)
			return &reconcile.Result{}, err
		}
		log.Info("Sucessfully created RoleBinding", "name", rb.Name)
	} else if err != nil {
		// Error that isn't due to the role binding not exist
		log.Error(err, "Failed to get role binding data", "name", rb.Name)
		return &reconcile.Result{}, err
	} else {
		log.Info("RoleBinding already exist", "name", rb.Name)
	}
	// RoleBinding creation was successful
	return nil, nil
}

func (r *AnomalyEngineReconciler) ensureClusterRoleBinding(request reconcile.Request, instance *backendv1alpha1.AnomalyEngine,
	rb *rbac.ClusterRoleBinding) (*reconcile.Result, error) {

	var log = logf.Log.WithName("anomay-detection -> ClusterRoleBinding ->")

	// check ClusterRoleBinding exist
	cluterRoleBindingFound := &rbac.ClusterRoleBinding{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      rb.Name,
		Namespace: rb.Namespace,
	}, cluterRoleBindingFound)
	log.Info("ClusterRoleBinding check completed", "name", rb.Name)

	if err != nil && errors.IsNotFound(err) {
		// Create the ClusterRoleBinding
		log.Info("Creating a new ClusterRoleBinding", "name", rb.Name)
		err = r.Client.Create(context.TODO(), rb)

		if err != nil {
			// ClusterRoleBinding creation failed
			log.Error(err, "Failed to create ClusterRoleBinding", "name", rb.Name)
			return &reconcile.Result{}, err
		}
		log.Info("Sucessfully created ClusterRoleBinding", "name", rb.Name)
	} else if err != nil {
		// Error that isn't due to the ClusterRoleBinding not exist
		log.Error(err, "Failed to get ClusterRoleBinding data", "name", rb.Name)
		return &reconcile.Result{}, err
	} else {
		log.Info("ClusterRoleBinding already exist", "name", rb.Name)
	}
	// ClusterRoleBinding creation was successful
	return nil, nil
}

func (r *AnomalyEngineReconciler) getRoleBindingForStorageServiceAccount(instance *backendv1alpha1.AnomalyEngine) *rbac.RoleBinding {
	rb := &rbac.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: rbac.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "osa-crd-admin-sa-rolebinding",
			Namespace: instance.Spec.Namespace,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "Role",
			Name:     storageAdminRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      storageServiceAccountName,
				Namespace: instance.Spec.Namespace,
			},
		},
	}
	controllerutil.SetControllerReference(instance, rb, r.Scheme)
	return rb
}

func (r *AnomalyEngineReconciler) getRoleBindingForAnomalyEngine(instance *backendv1alpha1.AnomalyEngine) *rbac.ClusterRoleBinding {
	rb := &rbac.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: rbac.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Spec.ServiceAccountRoleBinding.Name,
			Namespace: instance.Spec.Namespace,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     instance.Spec.ServiceAccountRoleBinding.ClusterRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      instance.Spec.ServiceAccountRoleBinding.Name,
				Namespace: instance.Spec.Namespace,
			},
		},
	}
	controllerutil.SetControllerReference(instance, rb, r.Scheme)
	return rb
}

func (r *AnomalyEngineReconciler) ensureConfigMap(request reconcile.Request, instance *backendv1alpha1.AnomalyEngine) (*reconcile.Result, error) {

	var log = logf.Log.WithName("anomay-detection -> Configmap ->")

	//  check configmap exist
	configmapFound := &corev1.ConfigMap{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      configmapName,
		Namespace: instance.Spec.Namespace,
	}, configmapFound)
	log.Info("Comfigmap check completed", "name", configmapName)

	if err != nil && errors.IsNotFound(err) {

		// Create the configmap
		log.Info("Creating a new comfigmap", "name", configmapName)
		data := string(instance.Spec.AnomalyQueryConfiguration)
		err = r.Client.Create(context.TODO(), &corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				APIVersion: corev1.SchemeGroupVersion.String(),
				Kind:       "ConfigMap",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      configmapName,
				Namespace: instance.Spec.Namespace,
			},
			Data: map[string]string{
				anomalyConfigFileName: data,
			},
		})

		if err != nil {
			// Configmap creation failed
			log.Error(err, "Failed to create Configmap")
			return &reconcile.Result{}, err
		}
		log.Info("Sucessfully created Configmap")
	} else if err != nil {
		// Error that isn't due to the configmap not exist
		log.Error(err, "Failed to get Configmap data")
		return &reconcile.Result{}, err
	}
	// Configmap creation was successful
	return nil, nil
}

func (r *AnomalyEngineReconciler) ensureCronJob(request reconcile.Request, instance *backendv1alpha1.AnomalyEngine) (*reconcile.Result, error) {

	var log = logf.Log.WithName("anomay-detection -> CronJob ->")

	//  check cronjob exist
	cronjobName := instance.Spec.CronJobConfig.Name
	cronjobFound := &batchv1.CronJob{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      cronjobName,
		Namespace: instance.Spec.Namespace,
	}, cronjobFound)
	log.Info("Cronjob check completed", "name", cronjobName)

	if err != nil && errors.IsNotFound(err) {
		// Create the cronjob
		log.Info("Creating a new cronjob", "name", cronjobName, "image", instance.Spec.CronJobConfig.Image)
		logLevel := instance.Spec.CronJobConfig.LogLevel
		if logLevel == "" {
			logLevel = "DEBUG"
		}
		err = r.Client.Create(context.TODO(), &batchv1.CronJob{
			TypeMeta: metav1.TypeMeta{
				APIVersion: batchv1.SchemeGroupVersion.String(),
				Kind:       "CronJob",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      cronjobName,
				Namespace: instance.Spec.Namespace,
			},
			Spec: batchv1.CronJobSpec{
				Schedule:                   instance.Spec.CronJobConfig.Schedule,
				ConcurrencyPolicy:          "Forbid",
				StartingDeadlineSeconds:    &[]int64{100}[0],
				SuccessfulJobsHistoryLimit: &[]int32{20}[0],
				FailedJobsHistoryLimit:     &[]int32{20}[0],
				JobTemplate: batchv1.JobTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name: cronjobName,
					},
					Spec: batchv1.JobSpec{
						BackoffLimit: &[]int32{3}[0],
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								RestartPolicy: corev1.RestartPolicyNever,
								Volumes: []corev1.Volume{
									{
										Name: "config-volume",
										VolumeSource: corev1.VolumeSource{
											ConfigMap: &corev1.ConfigMapVolumeSource{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: configmapName,
												},
											},
										},
									},
								},
								ServiceAccountName: storageServiceAccountName,
								Containers: []corev1.Container{
									{
										Name:            cronjobName,
										Image:           instance.Spec.CronJobConfig.Image,
										ImagePullPolicy: corev1.PullPolicy(corev1.PullAlways),
										Command:         []string{"python"},
										Args:            []string{"-m", "src.driver", "-aq", instance.Spec.CronJobConfig.AnomalyQueries},
										Env: []corev1.EnvVar{
											{
												Name:  "LOG_LEVEL",
												Value: logLevel,
											},
											{
												Name: "ANOMALY_CONFIG_FILE",
												// Value: "src/data_asset/anomaly_config.yaml",
												Value: "/etc/config/" + anomalyConfigFileName,
											},
											{
												Name: "ACCESS_TOKEN",
												ValueFrom: &corev1.EnvVarSource{
													SecretKeyRef: &corev1.SecretKeySelector{
														LocalObjectReference: corev1.LocalObjectReference{
															Name: instance.Spec.ServiceAccountRoleBinding.Name + "-token",
														},
														Key: "token",
													},
												},
											},
										},
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												"cpu":    resource.MustParse(instance.Spec.CronJobConfig.Resource.CPULimit),
												"memory": resource.MustParse(instance.Spec.CronJobConfig.Resource.MemoryLimit),
											},
											Requests: corev1.ResourceList{
												"cpu":    resource.MustParse(instance.Spec.CronJobConfig.Resource.CPURequest),
												"memory": resource.MustParse(instance.Spec.CronJobConfig.Resource.MemoryRequest),
											},
										},
										VolumeMounts: []corev1.VolumeMount{
											{
												Name:      "config-volume",
												MountPath: "/etc/config",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		})

		if err != nil {
			// cronjob creation failed
			log.Error(err, "Failed to create Cronjob")
			return &reconcile.Result{}, err
		}
		log.Info("Sucessfully created Cronjob")
	} else if err != nil {
		// Error that isn't due to the cronjob not exist
		log.Error(err, "Failed to get Cronjob data")
		return &reconcile.Result{}, err
	}
	// Cronjob creation was successful
	return nil, nil
}
