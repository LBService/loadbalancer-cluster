package controller

import (
	"fmt"
	glog "github.com/zoumo/logdog"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	appsinformers "k8s.io/client-go/informers/apps/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"

	lbpoolv1alpha1 "github.com/LBService/loadbalancer-cluster/pkg/apis/lbpool/v1alpha1"
	clientset "github.com/LBService/loadbalancer-cluster/pkg/client/clientset/versioned"
	lbpoolscheme "github.com/LBService/loadbalancer-cluster/pkg/client/clientset/versioned/scheme"
	informers "github.com/LBService/loadbalancer-cluster/pkg/client/informers/externalversions/lbpool/v1alpha1"
	listers "github.com/LBService/loadbalancer-cluster/pkg/client/listers/lbpool/v1alpha1"
)

const controllerAgentName = "loadbalancer-controller"

const (
	SuccessSynced         = "Synced"
	ErrResourceExists     = "ErrResourceExists"
	MessageResourceExists = "Resource %q already exists and is not managed by Loadbalancer"
	MessageResourceSynced = "Loadbalancer synced successfully"
)

type Controller struct {
	kubeclientset kubernetes.Interface
	lbclientset   clientset.Interface

	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced
	lbLister          listers.LoadBalancerLister
	lbSynced          cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
	recorder  record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	lbclientset clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	lbInformer informers.LoadBalancerInformer) *Controller {

	controller := &Controller{
		kubeclientset:     kubeclientset,
		lbclientset:       lbclientset,
		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		lbLister:          lbInformer.Lister(),
		lbSynced:          lbInformer.Informer().HasSynced,
		workqueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "LoadBalancer"),
		recorder:          recorder,
	}
	glog.Infof("Setting up event handler")
	// setup an event handler for Loadbalancer resource change
	lbInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})

	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueLoadBalancer,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueLoadBalancer(new)
		},
	})

	return controller
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the LoadBalancer resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	return nil
}

func (c *Controller) Run(worker int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	glog.Info("Starting LoadBalancer Controller")

	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.lbSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	for i := 0; i < worker; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {}
func (c *Controller) UpdateLoadbalancerStatus(lb lbpoolv1alpha1.LoadBalancer, deployment *appsv1.Deployment) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	lbCopy := lb.DeepCopy()
	lbCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas

	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the LoadBalancer resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.lbclientset.LbpoolV1alpha1().LoadBalancers(lb.Namespace).Update(lbCopy)
	return err

}

func (c *Controller) enqueueLoadBalancer(obj interface{}) {
	var key string
	var key error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the LoadBalancer resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that LoadBalancer resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool

	if object, ok := obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object tombstone, invalide type"))
			return
		}
	}
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		if ownerRef.Kind != "LoadBalancer" {
			return
		}

		lb, err := c.lbLister.LoadBalancers(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			retrun
		}
		c.enqueueLoadBalancer(lb)
		return
	}
}

func newDeployment(lb *lbpoolv1alpha1.LoadBalancer) *appsv1.Deployment {

	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      lb.Spec.DeploymentName,
			Namespace: lb.Namespace,
		},
		Spec: appsv1.DeploymentSpec{},
	}
}
