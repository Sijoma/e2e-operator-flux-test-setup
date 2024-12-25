package e2e

import (
	"os"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/support/kind"
	"sigs.k8s.io/e2e-framework/third_party/flux"
)

var (
	testEnv         env.Environment
	namespace       string
	kindClusterName string
)

func TestMain(m *testing.M) {
	cfg, _ := envconf.NewFromFlags()
	testEnv = env.NewWithConfig(cfg)
	kindClusterName = envconf.RandomName("flux", 10)
	namespace = envconf.RandomName("flux", 10)
	gitRepoName := "e2e-operator-flux-test-setup"
	ksName := "hello-world"
	const gitRepo = "https://github.com/Sijoma/e2e-operator-flux-test-setup"
	testEnv.Setup(
		envfuncs.CreateCluster(kind.NewProvider(), kindClusterName),
		envfuncs.CreateNamespace(namespace),
		flux.InstallFlux(),
		flux.CreateGitRepo(gitRepoName, gitRepo, flux.WithBranch("main")),
		flux.CreateKustomization(
			ksName,
			"GitRepository/"+gitRepoName+".flux-system",
			flux.WithPath("deploy/local"),
			flux.WithArgs("--target-namespace", namespace, "--prune"),
		),
	)

	testEnv.Finish(
		flux.DeleteKustomization(ksName),
		flux.DeleteGitRepo(gitRepoName),
		flux.UninstallFlux(),
		envfuncs.DeleteNamespace(namespace),
		envfuncs.DestroyCluster(kindClusterName),
	)
	os.Exit(testEnv.Run(m))
}
