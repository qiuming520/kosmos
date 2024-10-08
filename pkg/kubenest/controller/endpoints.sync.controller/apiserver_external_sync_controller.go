package endpointcontroller

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/kosmos.io/kosmos/pkg/apis/kosmos/v1alpha1"
	"github.com/kosmos.io/kosmos/pkg/kubenest/constants"
	"github.com/kosmos.io/kosmos/pkg/kubenest/util"
	"github.com/kosmos.io/kosmos/pkg/utils"
)

type APIServerExternalSyncController struct {
	client.Client
	EventRecorder record.EventRecorder
}

const APIServerExternalSyncControllerName string = "api-server-external-service-sync-controller"

func (e *APIServerExternalSyncController) SetupWithManager(mgr manager.Manager) error {
	skipEvent := func(obj client.Object) bool {
		return strings.Contains(obj.GetName(), "apiserver") && obj.GetNamespace() != ""
	}

	return controllerruntime.NewControllerManagedBy(mgr).
		Named(APIServerExternalSyncControllerName).
		WithOptions(controller.Options{MaxConcurrentReconciles: 5}).
		For(&v1.Endpoints{},
			builder.WithPredicates(predicate.Funcs{
				CreateFunc: func(createEvent event.CreateEvent) bool {
					return skipEvent(createEvent.Object)
				},
				UpdateFunc: func(updateEvent event.UpdateEvent) bool { return skipEvent(updateEvent.ObjectNew) },
				DeleteFunc: func(deleteEvent event.DeleteEvent) bool { return false },
			})).
		Watches(&source.Kind{Type: &v1alpha1.VirtualCluster{}}, handler.EnqueueRequestsFromMapFunc(e.newVirtualClusterMapFunc())).
		Complete(e)
}

func (e *APIServerExternalSyncController) newVirtualClusterMapFunc() handler.MapFunc {
	return func(a client.Object) []reconcile.Request {
		var requests []reconcile.Request
		vcluster := a.(*v1alpha1.VirtualCluster)

		// Join the Reconcile queue only if the status of the vcluster is Completed
		if vcluster.Status.Phase == v1alpha1.Completed {
			klog.V(4).Infof("api-server-external-sync-controller: virtualcluster change to completed: %s", vcluster.Name)
			// Add the vcluster to the Reconcile queue
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      vcluster.Name,
					Namespace: vcluster.Namespace,
				},
			})
		}
		return requests
	}
}

func (e *APIServerExternalSyncController) SyncAPIServerExternalEPS(ctx context.Context, k8sClient kubernetes.Interface) error {
	kubeEndpoints, err := k8sClient.CoreV1().Endpoints(constants.DefaultNs).Get(ctx, "kubernetes", metav1.GetOptions{})
	if err != nil {
		klog.Errorf("Error getting endpoints: %v", err)
		return err
	}
	klog.V(4).Infof("Endpoints for service 'kubernetes': %v", kubeEndpoints)
	for _, subset := range kubeEndpoints.Subsets {
		for _, address := range subset.Addresses {
			klog.V(4).Infof("IP: %s", address.IP)
		}
	}

	if len(kubeEndpoints.Subsets) != 1 {
		return fmt.Errorf("eps %s Subsets length is not 1", "kubernetes")
	}

	if kubeEndpoints.Subsets[0].Addresses == nil || len(kubeEndpoints.Subsets[0].Addresses) == 0 {
		klog.Errorf("eps %s Addresses length is nil", "kubernetes")
		return err
	}

	apiServerExternalEndpoints, err := k8sClient.CoreV1().Endpoints(constants.DefaultNs).Get(ctx, constants.APIServerExternalService, metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		klog.Errorf("failed to get endpoints for %s : %v", constants.APIServerExternalService, err)
		return err
	}

	updateEPS := apiServerExternalEndpoints.DeepCopy()

	if apiServerExternalEndpoints != nil {
		klog.V(4).Infof("apiServerExternalEndpoints: %v", apiServerExternalEndpoints)
	} else {
		klog.V(4).Info("apiServerExternalEndpoints is nil")
	}

	if updateEPS != nil {
		klog.V(4).Infof("updateEPS: %v", updateEPS)
	} else {
		klog.V(4).Info("updateEPS is nil")
	}

	if len(updateEPS.Subsets) == 1 && len(updateEPS.Subsets[0].Addresses) == 1 {
		ip := kubeEndpoints.Subsets[0].Addresses[0].IP
		klog.V(4).Infof("IP address: %s", ip)
		updateEPS.Subsets[0].Addresses[0].IP = ip

		if _, err := k8sClient.CoreV1().Endpoints(constants.DefaultNs).Update(ctx, updateEPS, metav1.UpdateOptions{}); err != nil {
			klog.Errorf("failed to update endpoints for api-server-external-service: %v", err)
			return err
		}
	} else {
		klog.ErrorS(err, "Unexpected format of endpoints for api-server-external-service", "endpoint_data", updateEPS)
		return err
	}

	return nil
}

func (e *APIServerExternalSyncController) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	klog.V(4).Infof("============ %s start to reconcile %s ============", APIServerExternalSyncControllerName, request.NamespacedName)
	defer klog.V(4).Infof("============ %s finish to reconcile %s ============", APIServerExternalSyncControllerName, request.NamespacedName)

	var virtualClusterList v1alpha1.VirtualClusterList
	if err := e.List(ctx, &virtualClusterList); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		klog.V(4).Infof("query virtualcluster failed: %v", err)
		return reconcile.Result{RequeueAfter: utils.DefaultRequeueTime}, nil
	}
	var targetVirtualCluster v1alpha1.VirtualCluster
	hasVirtualCluster := false
	for _, vc := range virtualClusterList.Items {
		if vc.Namespace == request.Namespace {
			targetVirtualCluster = vc
			klog.V(4).Infof("virtualcluster %s found", targetVirtualCluster.Name)
			hasVirtualCluster = true
			break
		}
	}
	if !hasVirtualCluster {
		klog.V(4).Infof("virtualcluster %s not found", request.Namespace)
		return reconcile.Result{}, nil
	}

	if targetVirtualCluster.Status.Phase != v1alpha1.Completed {
		return reconcile.Result{}, nil
	}

	k8sClient, err := util.GenerateKubeclient(&targetVirtualCluster)
	if err != nil {
		klog.Errorf("virtualcluster %s crd kubernetes client failed: %v", targetVirtualCluster.Name, err)
		return reconcile.Result{}, nil
	}

	if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		return e.SyncAPIServerExternalEPS(ctx, k8sClient)
	}); err != nil {
		klog.Errorf("virtualcluster %s sync apiserver external endpoints failed: %v", targetVirtualCluster.Name, err)
		return reconcile.Result{RequeueAfter: utils.DefaultRequeueTime}, nil
	}

	return reconcile.Result{}, nil
}
