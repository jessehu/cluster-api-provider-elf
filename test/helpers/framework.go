/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helpers

import (
	goctx "context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime"
	capie2e "sigs.k8s.io/cluster-api/test/e2e"
	"sigs.k8s.io/cluster-api/test/framework"
	"sigs.k8s.io/cluster-api/test/framework/bootstrap"
	"sigs.k8s.io/cluster-api/test/framework/clusterctl"
)

// Util functions to interact with the clusterctl e2e framework

func LoadE2EConfig(configPath string) (*clusterctl.E2EConfig, error) {
	config := clusterctl.LoadE2EConfig(goctx.TODO(), clusterctl.LoadE2EConfigInput{ConfigPath: configPath})
	if config == nil {
		return nil, fmt.Errorf("cannot load E2E config found at %s", configPath)
	}

	return config, nil
}

func CreateClusterctlLocalRepository(config *clusterctl.E2EConfig, repositoryFolder string, cniEnabled bool) (string, error) {
	createRepositoryInput := clusterctl.CreateRepositoryInput{
		E2EConfig:        config,
		RepositoryFolder: repositoryFolder,
	}

	if cniEnabled {
		// Ensuring a CNI file is defined in the config and register a FileTransformation to inject the referenced file as in place of the CNI_RESOURCES envSubst variable.
		cniPath, ok := config.Variables[capie2e.CNIPath]
		if !ok {
			return "", fmt.Errorf("missing %s variable in the config", capie2e.CNIPath)
		}

		if _, err := os.Stat(cniPath); err != nil {
			return "", fmt.Errorf("the %s variable should resolve to an existing file", capie2e.CNIPath)
		}

		createRepositoryInput.RegisterClusterResourceSetConfigMapTransformation(cniPath, capie2e.CNIResources)
	}

	clusterctlConfig := clusterctl.CreateRepository(goctx.TODO(), createRepositoryInput)
	if _, err := os.Stat(clusterctlConfig); err != nil {
		return "", fmt.Errorf("the clusterctl config file does not exists in the local repository %s", repositoryFolder)
	}

	return clusterctlConfig, nil
}

func SetupBootstrapCluster(config *clusterctl.E2EConfig, scheme *runtime.Scheme, useExistingCluster bool) (bootstrap.ClusterProvider, framework.ClusterProxy, error) {
	var clusterProvider bootstrap.ClusterProvider
	kubeconfigPath := ""
	if !useExistingCluster {
		clusterProvider = bootstrap.CreateKindBootstrapClusterAndLoadImages(goctx.TODO(), bootstrap.CreateKindBootstrapClusterAndLoadImagesInput{
			Name:               config.ManagementClusterName,
			RequiresDockerSock: config.HasDockerProvider(),
			Images:             config.Images,
		})

		kubeconfigPath = clusterProvider.GetKubeconfigPath()
		if _, err := os.Stat(kubeconfigPath); err != nil {
			return nil, nil, errors.New("failed to get the kubeconfig file for the bootstrap cluster")
		}
	}

	clusterProxy := framework.NewClusterProxy("bootstrap", kubeconfigPath, scheme)

	return clusterProvider, clusterProxy, nil
}

func InitBootstrapCluster(bootstrapClusterProxy framework.ClusterProxy, config *clusterctl.E2EConfig, clusterctlConfig, artifactFolder string) {
	clusterctl.InitManagementClusterAndWatchControllerLogs(goctx.TODO(), clusterctl.InitManagementClusterAndWatchControllerLogsInput{
		ClusterProxy:            bootstrapClusterProxy,
		ClusterctlConfigPath:    clusterctlConfig,
		InfrastructureProviders: config.InfrastructureProviders(),
		LogFolder:               filepath.Join(artifactFolder, "clusters", bootstrapClusterProxy.GetName()),
	}, config.GetIntervals(bootstrapClusterProxy.GetName(), "wait-controllers")...)
}

func TearDown(bootstrapClusterProvider bootstrap.ClusterProvider, bootstrapClusterProxy framework.ClusterProxy) {
	if bootstrapClusterProxy != nil {
		bootstrapClusterProxy.Dispose(goctx.TODO())
	}
	if bootstrapClusterProvider != nil {
		bootstrapClusterProvider.Dispose(goctx.TODO())
	}
}
