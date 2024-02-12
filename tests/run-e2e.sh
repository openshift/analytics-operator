set -e -u -o pipefail

trap cleanup INT

PROJECT_ROOT="$(git rev-parse --show-toplevel)"
declare -r PROJECT_ROOT

source "$PROJECT_ROOT/hack/utils.bash"

declare -r LOCAL_BIN="$PROJECT_ROOT/tmp/bin"
export PATH="$LOCAL_BIN:$PATH"


declare IMG_BASE='quay.io/openshiftanalytics'
declare NO_DEPLOY=false
declare NO_BUILD=false
declare VERSION=e2e-test
declare BATCH_SIZE=30

# Deploy operator to openshift cluster
deploy_operator() {

    header "Deploy Operator"


	$NO_DEPLOY && {
		info "skipping deploying of operator"
		return 0
	}

    run make deploy IMG_BASE="$IMG_BASE" VERSION="$VERSION" 
    wait_for_operators_ready "analytics-operator-system"
}

# Build operator image and push to container registry 
build_and_push() {
	header "Build Operator Images"

	$NO_BUILD && {
		info "skipping building of images"
		return 0
	}

    run make manifests generate operator-build operator-push IMG_BASE="$IMG_BASE" VERSION="$VERSION"

}

# Wait till operator becomes ready 
wait_for_operators_ready() {
	local ns="$1"
	shift

	header "Wait for analytics operator to be ready"

	local tries=30
	while [[ $tries -gt 0 ]] &&
		! kubectl -n "$ns" rollout status deploy/analytics-operator-controller-manager; do
		sleep 10
		((tries--))
	done

	kubectl wait -n "$ns" --for=condition=Available deploy/analytics-operator-controller-manager --timeout=300s

	ok "Analytics operator is up and running"
}

# Create CRD instance for anomaly engine
create_cr_for_anomaly_engine(){

    header "Configuring Anomaly Engine by creating CR"

    # temp code-----start------
    ls .
    echo "----------"
    ls ./
    echo "----------"
    ls ./config/samples/
    echo "----------"
    # temp code-----end------

    info "get anomalyengine"
    kubectl get anomalyengine

    info "create anomalyengine instance"
    kubectl apply -f ./config/samples/observability-analytics_v1alpha1_anomalyengine.yaml
    
    if ! kubectl -n "osa-anomaly-detection" get cronjob | grep "osa-anomaly-detection"; then
        fail "Cronjob not present to detect Anomaly"
        die
    fi

    info "Cronjob to detect anomaly engine created"

}

# Check pod for anomaly engine running fine or not
check_pod_for_anomaly_engine(){

    header "See pod is getting created from cronjob and running as expected."

    kubectl_command="kubectl -n "osa-anomaly-detection" get pods --field-selector=status.phase=Succeeded"

    # Set the total duration to wait and interval for checking pod
    total_duration=600  # 10 minutes
    check_interval=60  # 1 minute

    # Calculate the number of iterations needed
    iterations=$((total_duration / check_interval))

    # Counter for loop
    count=0

    # Run the loop
    while [ $count -lt $iterations ]; do
        # Execute the kubectl command
        pod_status=$($kubectl_command)

        # Check if the command was successful and the pod with "Succeeded" status exists
        if [[ $? -eq 0 && -n "$pod_status" ]]; then
            ok "Pod with Succeeded status found."
            break
        else
            warn "Pod with Succeeded status not found. Checking again in 1 minute..."
            sleep $check_interval
            count=$((count + 1))
        fi
    done
    # If the loop completes without finding the pod, print a final message
    if [ $count -eq $iterations ]; then
        fail "Pod with Succeeded status not found within 10 minutes."
        die
    fi
}

# Ingest namespaces into cluster
ingest_namespace(){

    info "Ingest Namespaces"
    commands_list=()

    for i in {1..100}; do
        namespace_name="osa-e2e-ns-${i}"
        command="kubectl create namespace "$namespace_name""
        commands_list+=("$command")
    done

    # Run the commands in parallel
    run_commands "${commands_list[@]}"  
    
    info "Namesapces created successfully"
}

# Delete ingested namesapces from cluster
delete_namespaces(){

    info "Delete Namespaces"
    commands_list=()

    for i in {1..100}; do
        namespace_name="osa-e2e-ns-${i}"
        command="kubectl delete namespace "$namespace_name""
        commands_list+=("$command")
    done

    # Run the commands in parallel
    run_commands "${commands_list[@]}"  

    info "Namesapces deleted successfully"
}

# Ingest Congifmaps into cluster
ingest_configmaps(){

    info "Ingest Configmaps"
    commands_list=()

    for i in {1..600}; do
        configmap_name="osa-e2e-cm-${i}"
        command="kubectl -n "osa-anomaly-detection" create configmap "$configmap_name" --from-literal=key1=value1"
        commands_list+=("$command")
    done

    # Run the commands in parallel
    run_commands "${commands_list[@]}" 

    info "Configmaps created successfully"

}

# Delete ingested configmaps from cluster
delete_configmaps(){

    info "Delete Configmaps"

    commands_list=()

    for i in {1..600}; do
        configmap_name="osa-e2e-cm-${i}"
        command="kubectl -n "osa-anomaly-detection" delete configmap "$configmap_name""
        commands_list+=("$command")
    done

    # Run the commands in parallel
    run_commands "${commands_list[@]}" 

    info "Configmaps deleted successfully"
}

# Run given commands parallel 
run_commands() {
    local commands=("$@")

    for((i=0; i < ${#commands[@]}; i+=$BATCH_SIZE))
    do
        part=( "${commands[@]:i:$BATCH_SIZE}" )
        # echo "Elements in this group: ${part[*]}"
         for cmd in "${part[@]}"; do
            # Run the command in the background
            eval "$cmd" &
        done
        
        # Wait for all background processes to finish
        wait
    done
}

# Inspect namespace anomaly 
inpsect_namespace_anomaly(){
    
    header "Inspect namespace anomaly (Min/Max Configuration)."
    # Set the total duration to wait and interval for checking anomaly data exist or not
    total_duration=600  # 10 minutes
    check_interval=60  # 1 minute

    # Calculate the number of iterations needed
    iterations=$((total_duration / check_interval))

    # Counter for loop
    count=0

    info "Check for anomlay data for 10 minutes in a while loop with 1 minute of interval."
    # Run the loop
    while [ $count -lt $iterations ]; do
        # Execute the kubectl command
        anomaly_status=$(kubectl get anomalydata -n osa-anomaly-detection)
        anomaly_present=false
        for elm in $(kubectl get anomalydata -n osa-anomaly-detection | grep "etcd-object-namespaces-namespaces"); do 
            if [[ "$elm" == *"etcd-object-namespaces-namespaces" ]]; then
                anomaly_present=true
                break
            fi
        done

        # ...do something interesting...
        if [ "$anomaly_present" = true ] ; then
            ok "Anomaly for the namespace found."
            break
        else
            warn "Anomaly for the namespace not found. Checking again in 1 minute..."
            sleep $check_interval
            count=$((count + 1))
        fi
    done
    # If the loop completes without finding the anomaly, print a final message and break the code. 
    if [ $count -eq $iterations ]; then
        fail "Anomaly for the namespace not found within 10 minutes."
        die
    fi
}

# Inspect configmap anomaly
inpsect_configmap_anomaly(){
    
    header "Inspect configmap anomaly (Percentange Change Configuration)."
    # Set the total duration to wait and interval for checking anomaly data exist or not
    total_duration=600  # 10 minutes
    check_interval=60  # 1 minute

    # Calculate the number of iterations needed
    iterations=$((total_duration / check_interval))

    # Counter for loop
    count=0

    info "Check for anomlay data for 10 minutes in a while loop with 1 minute of interval."
    # Run the loop
    while [ $count -lt $iterations ]; do
        # Execute the kubectl command
        anomaly_status=$(kubectl get anomalydata -n osa-anomaly-detection)
        anomaly_present=false
        for elm in $(kubectl get anomalydata -n osa-anomaly-detection | grep "etcd-object-secrets-config-maps-configmaps"); do 
            if [[ "$elm" == *"etcd-object-secrets-config-maps-configmaps" ]]; then
                anomaly_present=true
                break
            fi
        done

        # ...do something interesting...
        if [ "$anomaly_present" = true ] ; then
            ok "Anomaly for the configmap found."
            break
        else
            warn "Anomaly for the configmap not found. Checking again in 1 minute..."
            sleep $check_interval
            count=$((count + 1))
        fi
    done
    # If the loop completes without finding the anomaly, print a final message and break the code. 
    if [ $count -eq $iterations ]; then
        fail "Anomaly for the configmap not found within 10 minutes."
        die
    fi
}

# Test Namespace Anomaly 
test_namespace_anomaly(){
    header "Testing for Namespace Anomaly with 'Min/Max' configuration."
    ingest_namespace
    inpsect_namespace_anomaly
    delete_namespaces
}

# Test Configmap Anomaly 
test_configmap_anomaly(){
    header "Testing for Configmap Anomaly with 'Percentage Change' configuration."
    ingest_configmaps
    inpsect_configmap_anomaly
    delete_configmaps
}

# Delete operator and created resources. 
delete_operator_releated_resources(){
    header "Delete operator and releated resources"
    run make undeploy IMG_BASE="$IMG_BASE" VERSION="$VERSION"
    run kubectl delete namespace "osa-anomaly-detection"
}

# print configuration/parameters
print_config() {
	header "Test Configuration"
	cat <<-EOF
		  Image base:      $IMG_BASE
		  Skip Builds:     $NO_BUILD
		  Skip Deploy:     $NO_DEPLOY
		  Version:         $VERSION
		  Batch Size:      $BATCH_SIZE
		  Project root:    $PROJECT_ROOT
	EOF
	line 50
}

# parse arguments
parse_args() {
	### while there are args parse them
	while [[ -n "${1+xxx}" ]]; do
		case $1 in
		--no-deploy)
			NO_DEPLOY=true
			shift
			;;
		--no-build)
			NO_BUILD=true
			shift
			;;
		--image-base)
			shift
			IMG_BASE="$1"
			shift
			;;
		--version)
			shift
			VERSION="$1"
			shift
			;;
		*) return 1 ;; # show usage on everything else
		esac
	done

	return 0
}

main() {

	info "Inside Main"
    info "PROJECT_ROOT : ${PROJECT_ROOT}"

    parse_args "$@" || die "parse args failed"
    print_config

    cd "$PROJECT_ROOT"
    info "operator-sdk version $(operator-sdk version)" 

    build_and_push
    deploy_operator

    create_cr_for_anomaly_engine
    check_pod_for_anomaly_engine

    test_namespace_anomaly
    test_configmap_anomaly

    delete_operator_releated_resources

    ok "All looks good :)"

}

main "$@"
