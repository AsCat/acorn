package kubernetes

import (
	"fmt"
	osapps_v1 "github.com/openshift/api/apps/v1"
	osproject_v1 "github.com/openshift/api/project/v1"
	osroutes_v1 "github.com/openshift/api/route/v1"
	apps_v1 "k8s.io/api/apps/v1"
	auth_v1 "k8s.io/api/authorization/v1"
	batch_v1 "k8s.io/api/batch/v1"
	batch_v1beta1 "k8s.io/api/batch/v1beta1"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net"
	"os"

	kialiConfig "github.com/AsCat/acorn/config"
	"github.com/AsCat/acorn/log"
)

var (
	emptyListOptions = meta_v1.ListOptions{}
	emptyGetOptions  = meta_v1.GetOptions{}
)

type PodLogs struct {
	Logs string `json:"logs,omitempty"`
}

// IstioClientInterface for mocks (only mocked function are necessary here)
type IstioClientInterface interface {
	CreateIstioObject(api, namespace, resourceType, json string) (IstioObject, error)
	DeleteIstioObject(api, namespace, resourceType, name string) error
	GetAdapter(namespace, adapterType, adapterName string) (IstioObject, error)
	GetAdapters(namespace, labelSelector string) ([]IstioObject, error)
	GetAuthorizationDetails(namespace string) (*RBACDetails, error)
	GetCronJobs(namespace string) ([]batch_v1beta1.CronJob, error)
	GetDeployment(namespace string, deploymentName string) (*apps_v1.Deployment, error)
	GetDeployments(namespace string) ([]apps_v1.Deployment, error)
	GetDeploymentsByLabel(namespace string, labelSelector string) ([]apps_v1.Deployment, error)
	GetDeploymentConfig(namespace string, deploymentconfigName string) (*osapps_v1.DeploymentConfig, error)
	GetDeploymentConfigs(namespace string) ([]osapps_v1.DeploymentConfig, error)
	GetDestinationRule(namespace string, destinationrule string) (IstioObject, error)
	GetDestinationRules(namespace string, serviceName string) ([]IstioObject, error)
	GetEndpoints(namespace string, serviceName string) (*core_v1.Endpoints, error)
	GetGateway(namespace string, gateway string) (IstioObject, error)
	GetGateways(namespace string) ([]IstioObject, error)
	GetIstioRule(namespace string, istiorule string) (IstioObject, error)
	GetIstioRules(namespace string, labelSelector string) ([]IstioObject, error)
	GetJobs(namespace string) ([]batch_v1.Job, error)
	GetNamespace(namespace string) (*core_v1.Namespace, error)
	GetNamespaces(labelSelector string) ([]core_v1.Namespace, error)
	GetPod(namespace, name string) (*core_v1.Pod, error)
	GetPodLogs(namespace, name string, opts *core_v1.PodLogOptions) (*PodLogs, error)
	GetPods(namespace, labelSelector string) ([]core_v1.Pod, error)
	GetProject(project string) (*osproject_v1.Project, error)
	GetProjects(labelSelector string) ([]osproject_v1.Project, error)
	GetQuotaSpec(namespace string, quotaSpecName string) (IstioObject, error)
	GetQuotaSpecs(namespace string) ([]IstioObject, error)
	GetQuotaSpecBinding(namespace string, quotaSpecBindingName string) (IstioObject, error)
	GetQuotaSpecBindings(namespace string) ([]IstioObject, error)
	GetReplicationControllers(namespace string) ([]core_v1.ReplicationController, error)
	GetReplicaSets(namespace string) ([]apps_v1.ReplicaSet, error)
	GetRoute(namespace string, name string) (*osroutes_v1.Route, error)
	GetSidecar(namespace string, sidecar string) (IstioObject, error)
	GetSidecars(namespace string) ([]IstioObject, error)
	GetSelfSubjectAccessReview(namespace, api, resourceType string, verbs []string) ([]*auth_v1.SelfSubjectAccessReview, error)
	GetService(namespace string, serviceName string) (*core_v1.Service, error)
	GetServices(namespace string, selectorLabels map[string]string) ([]core_v1.Service, error)
	GetServiceEntries(namespace string) ([]IstioObject, error)
	GetServiceEntry(namespace string, serviceEntryName string) (IstioObject, error)
	GetStatefulSet(namespace string, statefulsetName string) (*apps_v1.StatefulSet, error)
	GetStatefulSets(namespace string) ([]apps_v1.StatefulSet, error)
	GetTemplate(namespace, templateType, templateName string) (IstioObject, error)
	GetTemplates(namespace, labelSelector string) ([]IstioObject, error)
	GetPolicy(namespace string, policyName string) (IstioObject, error)
	GetPolicies(namespace string) ([]IstioObject, error)
	GetMeshPolicy(policyName string) (IstioObject, error)
	GetMeshPolicies() ([]IstioObject, error)
	GetClusterRbacConfig(name string) (IstioObject, error)
	GetClusterRbacConfigs() ([]IstioObject, error)
	GetRbacConfig(namespace string, name string) (IstioObject, error)
	GetRbacConfigs(namespace string) ([]IstioObject, error)
	GetServiceMeshPolicy(namespace string, name string) (IstioObject, error)
	GetServiceMeshPolicies(namespace string) ([]IstioObject, error)
	GetServiceMeshRbacConfig(namespace string, name string) (IstioObject, error)
	GetServiceMeshRbacConfigs(namespace string) ([]IstioObject, error)
	GetServiceRole(namespace string, name string) (IstioObject, error)
	GetServiceRoles(namespace string) ([]IstioObject, error)
	GetServiceRoleBinding(namespace string, name string) (IstioObject, error)
	GetServiceRoleBindings(namespace string) ([]IstioObject, error)
	GetAuthorizationPolicy(namespace string, name string) (IstioObject, error)
	GetAuthorizationPolicies(namespace string) ([]IstioObject, error)
	GetServerVersion() (*version.Info, error)
	GetToken() string
	GetVirtualService(namespace string, virtualservice string) (IstioObject, error)
	GetVirtualServices(namespace string, serviceName string) ([]IstioObject, error)
	IsMaistraApi() bool
	IsOpenShift() bool
	UpdateIstioObject(api, namespace, resourceType, name, jsonPatch string) (IstioObject, error)
}

// IstioClient is the client struct for Kubernetes and Istio APIs
// It hides the way it queries each API
type IstioClient struct {
	IstioClientInterface
	token                    string
	k8s                      *kube.Clientset
	istioConfigApi           *rest.RESTClient
	istioNetworkingApi       *rest.RESTClient
	istioAuthenticationApi   *rest.RESTClient
	istioRbacApi             *rest.RESTClient
	istioSecurityApi         *rest.RESTClient
	maistraAuthenticationApi *rest.RESTClient
	maistraRbacApi           *rest.RESTClient
	// isOpenShift private variable will check if kiali is deployed under an OpenShift cluster or not
	// It is represented as a pointer to include the initialization phase.
	// See kubernetes_service.go#IsOpenShift() for more details.
	isOpenShift *bool

	// isMaistraApi private variable will check if specific Maistra APIs for authentication and rbac are present.
	// It is represented as a pointer to include the initialization phase.
	// See kubernetes_service.go#IsMaistraApi() for more details.
	isMaistraApi *bool

	// rbacResources private variable will check which resources kiali has access to from rbac.istio.io group
	// It is represented as a pointer to include the initialization phase.
	// See istio_details_service.go#HasRbacResource() for more details.
	rbacResources *map[string]bool

	// securityResources private variable will check which resources kiali has access to from security.istio.io group
	// It is represented as a pointer to include the initialization phase.
	// See istio_details_service.go#HasSecurityResource() for more details.
	securityResources *map[string]bool
}

// GetK8sApi returns the clientset referencing all K8s rest clients
func (client *IstioClient) GetK8sApi() *kube.Clientset {
	return client.k8s
}

// GetIstioConfigApi returns the istio config rest client
func (client *IstioClient) GetIstioConfigApi() *rest.RESTClient {
	return client.istioConfigApi
}

// GetIstioNetworkingApi returns the istio config rest client
func (client *IstioClient) GetIstioNetworkingApi() *rest.RESTClient {
	return client.istioNetworkingApi
}

// GetIstioRbacApi returns the istio rbac rest client
func (client *IstioClient) GetIstioRbacApi() *rest.RESTClient {
	return client.istioRbacApi
}

// GetIstioSecurityApi returns the istio security rest client
func (client *IstioClient) GetIstioSecurityApi() *rest.RESTClient {
	return client.istioSecurityApi
}

// GetToken returns the BearerToken used from the config
func (client *IstioClient) GetToken() string {
	return client.token
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// ConfigClient return a client with the correct configuration
// Returns configuration if Kiali is in Cluster when InCluster is true
// Returns configuration if Kiali is not int Cluster when InCluster is false
// It returns an error on any problem
func ConfigClient() (*rest.Config, error) {
	if kialiConfig.Get().InCluster {
		incluster, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		incluster.QPS = kialiConfig.Get().KubernetesConfig.QPS
		incluster.Burst = kialiConfig.Get().KubernetesConfig.Burst
		return incluster, nil
	}
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		return nil, fmt.Errorf("unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined")
	}

	if !kialiConfig.Get().InCluster && kialiConfig.Get().KubeConfigPath != "" {
		config, err := clientcmd.BuildConfigFromFlags("", kialiConfig.Get().KubeConfigPath)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	return &rest.Config{
		// TODO: switch to using cluster DNS.
		Host:  "http://" + net.JoinHostPort(host, port),
		QPS:   kialiConfig.Get().KubernetesConfig.QPS,
		Burst: kialiConfig.Get().KubernetesConfig.Burst,
	}, nil
}

// NewClientFromConfig creates a new client to the Kubernetes and Istio APIs.
// It takes the assumption that Istio is deployed into the cluster.
// It hides the access to Kubernetes/Openshift credentials.
// It hides the low level use of the API of Kubernetes and Istio, it should be considered as an implementation detail.
// It returns an error on any problem.
func NewClientFromConfig(config *rest.Config) (*IstioClient, error) {
	client := IstioClient{
		token: config.BearerToken,
	}
	log.Debugf("Rest perf config QPS: %f Burst: %d", config.QPS, config.Burst)

	k8s, err := kube.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client.k8s = k8s

	// Istio is a CRD extension of Kubernetes API, so any custom type should be registered here.
	// KnownTypes registers the Istio objects we use, as soon as we get more info we will increase the number of types.
	types := runtime.NewScheme()
	schemeBuilder := runtime.NewSchemeBuilder(
		func(scheme *runtime.Scheme) error {
			// Register networking types
			for _, nt := range networkingTypes {
				scheme.AddKnownTypeWithName(NetworkingGroupVersion.WithKind(nt.objectKind), &GenericIstioObject{})
				scheme.AddKnownTypeWithName(NetworkingGroupVersion.WithKind(nt.collectionKind), &GenericIstioObjectList{})
			}
			// Register config types
			for _, cf := range configTypes {
				scheme.AddKnownTypeWithName(ConfigGroupVersion.WithKind(cf.objectKind), &GenericIstioObject{})
				scheme.AddKnownTypeWithName(ConfigGroupVersion.WithKind(cf.collectionKind), &GenericIstioObjectList{})
			}
			// Register adapter types
			for _, ad := range adapterTypes {
				scheme.AddKnownTypeWithName(ConfigGroupVersion.WithKind(ad.objectKind), &GenericIstioObject{})
				scheme.AddKnownTypeWithName(ConfigGroupVersion.WithKind(ad.collectionKind), &GenericIstioObjectList{})
			}
			// Register template types
			for _, tp := range templateTypes {
				scheme.AddKnownTypeWithName(ConfigGroupVersion.WithKind(tp.objectKind), &GenericIstioObject{})
				scheme.AddKnownTypeWithName(ConfigGroupVersion.WithKind(tp.collectionKind), &GenericIstioObjectList{})
			}
			// Register authentication types
			for _, at := range authenticationTypes {
				scheme.AddKnownTypeWithName(AuthenticationGroupVersion.WithKind(at.objectKind), &GenericIstioObject{})
				scheme.AddKnownTypeWithName(AuthenticationGroupVersion.WithKind(at.collectionKind), &GenericIstioObjectList{})
			}
			for _, at := range maistraAuthenticationTypes {
				scheme.AddKnownTypeWithName(MaistraAuthenticationGroupVersion.WithKind(at.objectKind), &GenericIstioObject{})
				scheme.AddKnownTypeWithName(MaistraAuthenticationGroupVersion.WithKind(at.collectionKind), &GenericIstioObjectList{})
			}
			// Register rbac types
			for _, rt := range rbacTypes {
				scheme.AddKnownTypeWithName(RbacGroupVersion.WithKind(rt.objectKind), &GenericIstioObject{})
				scheme.AddKnownTypeWithName(RbacGroupVersion.WithKind(rt.collectionKind), &GenericIstioObjectList{})

			}
			for _, rt := range maistraRbacTypes {
				scheme.AddKnownTypeWithName(MaistraRbacGroupVersion.WithKind(rt.objectKind), &GenericIstioObject{})
				scheme.AddKnownTypeWithName(MaistraRbacGroupVersion.WithKind(rt.collectionKind), &GenericIstioObjectList{})
			}
			for _, rt := range securityTypes {
				scheme.AddKnownTypeWithName(SecurityGroupVersion.WithKind(rt.objectKind), &GenericIstioObject{})
				scheme.AddKnownTypeWithName(SecurityGroupVersion.WithKind(rt.collectionKind), &GenericIstioObjectList{})
			}

			meta_v1.AddToGroupVersion(scheme, ConfigGroupVersion)
			meta_v1.AddToGroupVersion(scheme, NetworkingGroupVersion)
			meta_v1.AddToGroupVersion(scheme, AuthenticationGroupVersion)
			meta_v1.AddToGroupVersion(scheme, RbacGroupVersion)
			meta_v1.AddToGroupVersion(scheme, MaistraAuthenticationGroupVersion)
			meta_v1.AddToGroupVersion(scheme, MaistraRbacGroupVersion)
			meta_v1.AddToGroupVersion(scheme, SecurityGroupVersion)
			return nil
		})

	err = schemeBuilder.AddToScheme(types)
	if err != nil {
		return nil, err
	}

	// Istio needs another type as it queries a different K8S API.
	istioConfigAPI, err := newClientForAPI(config, ConfigGroupVersion, types)
	if err != nil {
		return nil, err
	}

	istioNetworkingAPI, err := newClientForAPI(config, NetworkingGroupVersion, types)
	if err != nil {
		return nil, err
	}

	istioAuthenticationAPI, err := newClientForAPI(config, AuthenticationGroupVersion, types)
	if err != nil {
		return nil, err
	}

	istioRbacApi, err := newClientForAPI(config, RbacGroupVersion, types)
	if err != nil {
		return nil, err
	}

	istioSecurityApi, err := newClientForAPI(config, SecurityGroupVersion, types)
	if err != nil {
		return nil, err
	}

	maistraAuthenticationAPI, err := newClientForAPI(config, MaistraAuthenticationGroupVersion, types)
	if err != nil {
		return nil, err
	}

	maistraRbacApi, err := newClientForAPI(config, MaistraRbacGroupVersion, types)
	if err != nil {
		return nil, err
	}

	client.istioConfigApi = istioConfigAPI
	client.istioNetworkingApi = istioNetworkingAPI
	client.istioAuthenticationApi = istioAuthenticationAPI
	client.istioRbacApi = istioRbacApi
	client.istioSecurityApi = istioSecurityApi
	client.maistraAuthenticationApi = maistraAuthenticationAPI
	client.maistraRbacApi = maistraRbacApi
	return &client, nil
}

type DirectCodecFactory struct {
	serializer.CodecFactory
}

func newClientForAPI(fromCfg *rest.Config, groupVersion schema.GroupVersion, scheme *runtime.Scheme) (*rest.RESTClient, error) {
	cfg := rest.Config{
		Host:    fromCfg.Host,
		APIPath: "/apis",
		ContentConfig: rest.ContentConfig{
			GroupVersion:         &groupVersion,
			NegotiatedSerializer: DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme)},
			ContentType:          runtime.ContentTypeJSON,
		},
		BearerToken:     fromCfg.BearerToken,
		TLSClientConfig: fromCfg.TLSClientConfig,
		QPS:             fromCfg.QPS,
		Burst:           fromCfg.Burst,
	}
	return rest.RESTClientFor(&cfg)
}