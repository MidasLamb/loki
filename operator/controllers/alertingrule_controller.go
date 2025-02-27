package controllers

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/go-logr/logr"
	lokiv1beta1 "github.com/grafana/loki/operator/api/v1beta1"
	"github.com/grafana/loki/operator/controllers/internal/lokistack"
)

// AlertingRuleReconciler reconciles a AlertingRule object
type AlertingRuleReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=loki.grafana.com,resources=alertingrules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=loki.grafana.com,resources=alertingrules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=loki.grafana.com,resources=alertingrules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AlertingRule object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *AlertingRuleReconciler) Reconcile(ctx context.Context, _ ctrl.Request) (ctrl.Result, error) {
	err := lokistack.AnnotateForDiscoveredRules(ctx, r.Client)
	if err != nil {
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AlertingRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lokiv1beta1.AlertingRule{}).
		Watches(&source.Kind{Type: &corev1.Namespace{}}, &handler.EnqueueRequestForObject{}, builder.OnlyMetadata).
		Complete(r)
}
