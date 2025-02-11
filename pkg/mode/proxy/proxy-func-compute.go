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

package proxy

import (
	cfg "asaintsever/open-vault-agent-injector/pkg/config"
	ctx "asaintsever/open-vault-agent-injector/pkg/context"
	"strings"
)

func proxyModeCompute(config *cfg.OVAIConfig, labels, annotations map[string]string) (ctx.ModeConfig, error) {
	proxyPort := annotations[config.VaultInjectorAnnotationsFQ[vaultInjectorAnnotationProxyPortKey]]

	if proxyPort == "" { // Default port
		proxyPort = vaultProxyDefaultPort
	}

	template := strings.Replace(config.ProxyConfig, vaultProxyPortPlaceholder, proxyPort, -1)

	return &proxyModeConfig{template}, nil
}
