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
	"errors"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/klog"
)

func getSecretsModeConfig(config ctx.ModeConfig) (*secretsModeConfig, error) {
	if config != nil {
		secModeCfg, ok := config.(*secretsModeConfig) // here we use type assertion (https://golang.org/ref/spec#Type_assertions)

		if ok {
			return secModeCfg, nil
		}

		err := errors.New("Provided type cannot be casted to 'secretsModeConfig'")
		klog.Error(err.Error())
		return nil, err
	}

	err := errors.New("Secrets mode config is null")
	klog.Warning(err.Error())
	return nil, err
}

func IsSecretsStatic(context *ctx.InjectionContext) bool {
	if secretsModeCfg, err := getSecretsModeConfig(context.ModesConfig[m.VaultInjectorModeSecrets]); err == nil {
		return secretsModeCfg.secretsType == vaultInjectorSecretsTypeStatic
	}

	return false
}

func IsSecretsInjectionEnv(context *ctx.InjectionContext) bool {
	if secretsModeCfg, err := getSecretsModeConfig(context.ModesConfig[m.VaultInjectorModeSecrets]); err == nil {
		return secretsModeCfg.secretsInjectionMethod == vaultInjectorSecretsInjectionMethodEnv
	}

	return false
}

func getMountPathOfSecretsVolume(cnt corev1.Container) string {
	var secretsVolMountPath string

	for _, volMount := range cnt.VolumeMounts {
		if volMount.Name == SecretsVolName {
			secretsVolMountPath = volMount.MountPath
			break
		}
	}

	return secretsVolMountPath
}
