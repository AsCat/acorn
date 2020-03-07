package business

import (
	"sync"

	"github.com/AsCat/acorn/config"
	"github.com/AsCat/acorn/jaeger"
	"github.com/AsCat/acorn/kubernetes"
	"github.com/AsCat/acorn/kubernetes/cache"
	"github.com/AsCat/acorn/log"
	"github.com/AsCat/acorn/prometheus"
)

// Layer is a container for fast access to inner services
type Layer struct {
	Namespace      NamespaceService
	OpenshiftOAuth OpenshiftOAuthService
}

// Global clientfactory and prometheus clients.
var clientFactory kubernetes.ClientFactory
var prometheusClient prometheus.ClientInterface
var once sync.Once
var kialiCache cache.KialiCache

func initKialiCache() {
	if config.Get().KubernetesConfig.CacheEnabled {
		if cache, err := cache.NewKialiCache(); err != nil {
			log.Errorf("Error initializing Kiali Cache. Details: %s", err)
		} else {
			kialiCache = cache
		}
	}
}

func GetUnauthenticated() (*Layer, error) {
	return Get("")
}

// Get the business.Layer
func Get(token string) (*Layer, error) {
	// Kiali Cache will be initialized once at first use of Business layer
	once.Do(initKialiCache)

	// Use an existing client factory if it exists, otherwise create and use in the future
	if clientFactory == nil {
		userClient, err := kubernetes.GetClientFactory()
		if err != nil {
			return nil, err
		}
		clientFactory = userClient
	}

	// Creates a new k8s client based on the current users token
	k8s, err := clientFactory.GetClient(token)
	if err != nil {
		return nil, err
	}

	// Use an existing Prometheus client if it exists, otherwise create and use in the future
	if prometheusClient == nil {
		prom, err := prometheus.NewClient()
		if err != nil {
			return nil, err
		}
		prometheusClient = prom
	}

	// Create Jaeger client
	jaegerLoader := func() (jaeger.ClientInterface, error) {
		return jaeger.NewClient(token)
	}

	return NewWithBackends(k8s, prometheusClient, jaegerLoader), nil
}

// SetWithBackends allows for specifying the ClientFactory and Prometheus clients to be used.
// Mock friendly. Used only with tests.
func SetWithBackends(cf kubernetes.ClientFactory, prom prometheus.ClientInterface) {
	clientFactory = cf
	prometheusClient = prom
}

// NewWithBackends creates the business layer using the passed k8s and prom clients
func NewWithBackends(k8s kubernetes.IstioClientInterface, prom prometheus.ClientInterface, jaegerClient JaegerLoader) *Layer {
	temporaryLayer := &Layer{}

	temporaryLayer.Namespace = NewNamespaceService(k8s)

	return temporaryLayer
}

func Stop() {
	if kialiCache != nil {
		kialiCache.Stop()
	}
}