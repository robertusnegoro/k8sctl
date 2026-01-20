# Kubernetes Sample Manifests

This directory contains sample Kubernetes manifest files for various Kubernetes objects.

## Files

### Core Workloads
- **deployment.yaml** - Deployment with replicas, resource limits, and volume mounts
- **statefulset.yaml** - StatefulSet with persistent volumes and volume claim templates
- **daemonset.yaml** - DaemonSet for node-level logging
- **pod.yaml** - Standalone Pod and multi-container Pod examples
- **job.yaml** - Job and CronJob examples for batch processing

### Services & Networking
- **service-clusterip.yaml** - ClusterIP service (default)
- **service-nodeport.yaml** - NodePort service for external access
- **service-loadbalancer.yaml** - LoadBalancer service for cloud providers
- **service-headless.yaml** - Headless service for StatefulSets
- **ingress.yaml** - Ingress with TLS and multiple hosts
- **networkpolicy.yaml** - Network policies for pod-to-pod communication

### Configuration
- **configmap.yaml** - ConfigMaps for application configuration
- **secret.yaml** - Secrets for sensitive data (Opaque, TLS, Docker registry)
- **namespace.yaml** - Namespace examples

### Storage
- **persistentvolumeclaim.yaml** - PVC examples (ReadWriteOnce and ReadWriteMany)

### RBAC
- **serviceaccount.yaml** - ServiceAccount with Role, RoleBinding, ClusterRole, and ClusterRoleBinding

## Usage

Apply these manifests to your Kubernetes cluster:

```bash
# Apply all manifests
kubectl apply -f examples/

# Apply specific manifest
kubectl apply -f examples/deployment.yaml

# Apply with namespace
kubectl apply -f examples/deployment.yaml -n production
```

## Notes

- Replace placeholder values (passwords, certificates, etc.) with your actual values
- Adjust resource requests/limits based on your requirements
- Modify storage classes based on your cluster configuration
- Update image tags to use specific versions in production
