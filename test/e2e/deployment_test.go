package e2e

import (
	"context"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestFluxRepoWorkflow(t *testing.T) {
	feature := features.New("Install resources by flux").
		Assess("check if deployment was successful", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			deployment := &appsv1.Deployment{
				ObjectMeta: v1.ObjectMeta{
					Name:      "hello-app",
					Namespace: c.Namespace(),
				},
				Spec: appsv1.DeploymentSpec{},
			}

			err := wait.For(conditions.New(c.Client().Resources()).
				DeploymentConditionMatch(
					deployment,
					appsv1.DeploymentAvailable,
					corev1.ConditionStatus(v1.ConditionTrue),
				), wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal("Error deployment not found", err)
			}

			return ctx
		}).Feature()

	testEnv.Test(t, feature)
}
