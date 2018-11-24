package main

import (
	"context"
	"flag"
	"os"
	"runtime"
	"time"

	controller "github.com/LBService/loadbalancer-cluster/pkg/controller"

	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
	"github.com/LBService/loadbalancer-cluster/pkg/util/constants"
	"github.com/LBService/loadbalancer-cluster/pkg/version"
	"github.com/LBService/loadbalancer-cluster/pkg/k8sutil"
)

var (
	createCRD bool
)

func init() {
	flag.BoolVar(&createCRD, "create-crd", true, "The lbaas operator will not create the CRD when this flag is set to false.")
	flag.Parse()
}

func main() {
	namespace := os.Getenv(constants.EnvOperatorPodNamespace)
	if len(namespace) == 0 {
		logrus.Fatalf("must set env %s", constants.EnvOperatorPodNamespace)
	}
	name := os.Getenv(constants.EnvOperatorPodName)
	if len(name) == 0 {
		logrus.Fatalf("must set env %s", constants.EnvOperatorPodName)
	}
	id, err := os.Hostname()
	if err != nil {
		logrus.Fatalf("failed to get hostname: %v", err)
	}

	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("lbaas-operator Version: %v", version.Get())


	kubecli := k8sutil.MustNewKubeClient()
	rl, err := resourcelock.New(
		resourcelock.EndpointsResourceLock,
		namespace,
		"lbaas-operator",
		kubecli.Core(),
		resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: createRecorder(kubecli, name, namespace),
		},
	)
	if err != nil {
		logrus.Fatalf("error creating lock: %v", err)
	}

	leaderelection.RunOrDie(context.TODO(), leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: 15 * time.Second,
		RenewDeadline: 10 * time.Second,
		RetryPeriod:   2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				logrus.Fatalf("leader election lost")
			},
		},
	})
}

func createRecorder(kubecli kubernetes.Interface, name, namespace string) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logrus.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(kubecli.Core().RESTClient()).Events(namespace)})
	return eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: name})
}

func run(ctx context.Context) {
	cfg := newControllerConfig()

	c := controller.New(cfg)
	err := c.Start(ctx)
	if err != nil {
		logrus.Fatalf("operator stopped with error: %v", err)
	}
}



func newControllerConfig() controller.Config {
	kubecli := k8sutil.MustNewKubeClient()


	cfg := controller.Config{
		KubeCli:        kubecli,
		CreateCRD:      createCRD,
	}

	return cfg
}