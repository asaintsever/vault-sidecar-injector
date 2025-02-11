// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package secrets

import (
	ctx "asaintsever/open-vault-agent-injector/pkg/context"
	m "asaintsever/open-vault-agent-injector/pkg/mode"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

func secretsModeInject(containerBasePath string, podContainers []corev1.Container, containerName string, env []corev1.EnvVar, context *ctx.InjectionContext) (bool, error) {
	for _, cntName := range secretsContainerNames[containerBasePath] {
		if cntName == containerName {
			// Look type of secrets: inject init container only for static secrets
			if (IsSecretsStatic(context) && (containerBasePath == ctx.JsonPathInitContainers)) ||
				(!IsSecretsStatic(context) && (containerBasePath == ctx.JsonPathContainers)) {

				if (cntName == secretsEnvInitContainerName) && !IsSecretsInjectionEnv(context) {
					// Do not inject env init container if injection method is not 'env'
					return false, nil
				}

				klog.Infof("[%s] Injecting container %s (path: %s)", m.VaultInjectorModeSecrets, containerName, containerBasePath)

				// Resolve secrets env vars
				for envIdx := range env {
					if env[envIdx].Name == secretsTemplatesPlaceholderEnv {
						env[envIdx].Value = context.ModesConfig[m.VaultInjectorModeSecrets].GetTemplate()
					}
				}

				return true, nil
			}
		}
	}

	return false, nil
}
