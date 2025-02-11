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
	m "asaintsever/open-vault-agent-injector/pkg/mode"
)

func init() {
	// Register mode
	m.RegisterMode(
		m.VaultInjectorModeInfo{
			Key:                  m.VaultInjectorModeProxy,
			DefaultMode:          false,
			EnableDefaultMode:    false,
			Annotations:          []string{vaultInjectorAnnotationProxyPortKey},
			ComputeTemplatesFunc: proxyModeCompute,
			InjectContainerFunc:  proxyModeInject,
		},
	)
}

func (proxyModeCfg *proxyModeConfig) GetTemplate() string {
	return proxyModeCfg.template
}
