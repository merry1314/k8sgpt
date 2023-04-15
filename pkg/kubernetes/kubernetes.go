package kubernetes

import (
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
)

type Client struct {
	Client     kubernetes.Interface
	RestClient rest.Interface
	Config     *rest.Config
}

func (c *Client) GetConfig() *rest.Config {
	return c.Config
}

func (c *Client) GetClient() kubernetes.Interface {
	return c.Client
}

func (c *Client) GetRestClient() rest.Interface {
	return c.RestClient
}

func NewClient(kubecontext string, kubeconfig string) (*Client, error) {
	var config *rest.Config
	config, err := rest.InClusterConfig()
	if err != nil {
		clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
			&clientcmd.ConfigOverrides{
				CurrentContext: kubecontext,
			})
		// create the clientset
		config, err = clientConfig.ClientConfig()
		if err != nil {
			return nil, err
		}
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	config.APIPath = "/api"
	config.GroupVersion = &scheme.Scheme.PrioritizedVersionsForGroup("")[0]
	config.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:     clientSet,
		RestClient: restClient,
		Config:     config,
	}, nil
}
