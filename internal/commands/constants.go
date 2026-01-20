package commands

// Output format constants
const (
	OutputFormatJSON  = "json"
	OutputFormatYAML  = "yaml"
	OutputFormatTable = "table"
)

// Status constants
const (
	StatusReady    = "Ready"
	StatusNotReady = "NotReady"
)

// Common string constants
const (
	DefaultNamespace = "default"
	NoneValue        = "<none>"
)

// Resource type shortcuts
const (
	ResourcePod         = "pod"
	ResourcePods        = "pods"
	ResourcePo          = "po"
	ResourceDeploy      = "deploy"
	ResourceDeployment  = "deployment"
	ResourceDeployments = "deployments"
	ResourceService     = "service"
	ResourceServices    = "services"
	ResourceSvc         = "svc"
	ResourceIngress    = "ingress"
	ResourceIngresses  = "ingresses"
	ResourceIng         = "ing"
	ResourceServiceAccount = "serviceaccount"
	ResourceServiceAccounts = "serviceaccounts"
	ResourceSa          = "sa"
	ResourceConfigMap   = "configmap"
	ResourceConfigMaps  = "configmaps"
	ResourceCm          = "cm"
	ResourceSecret      = "secret"
	ResourceSecrets     = "secrets"
	ResourceSec         = "sec"
)
