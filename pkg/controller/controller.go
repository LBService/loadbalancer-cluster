package controller

import (
	"time"

	"context"
	"github.com/LBService/loadbalancer-cluster/pkg/apis/lbpool/v1alpha1"
	"github.com/LBService/loadbalancer-cluster/pkg/loadbalancer"
	"github.com/sirupsen/logrus"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kwatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

var initRetryWaitTime = 30 * time.Second

type Event struct {
	Type   kwatch.EventType
	Object *v1alpha1.LoadBalancer
}

type Controller struct {
	logger *logrus.Entry
	Config

	lbs map[string]*loadbalancer.LoadBalancer
}

type Config struct {
	Namespace      string
	ClusterWide    bool
	ServiceAccount string
	KubeCli        kubernetes.Interface
	KubeExtCli     apiextensionsclient.Interface
	CreateCRD      bool
}

func New(cfg Config) *Controller {
	return &Controller{
		logger: logrus.WithField("pkg", "controller"),

		Config: cfg,
		lbs:    make(map[string]*loadbalancer.LoadBalancer),
	}
}

func (c *Controller) Start(ctx context.Context) error {
	// TODO: get rid of this init code. CRD and storage class will be managed outside of operator.
	for {
		c.logger.Infof("Not implemented")
		time.Sleep(initRetryWaitTime)
	}

	panic("unreachable")
}
