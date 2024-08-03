// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
//
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

package k8s

import (
	"fmt"

	core "k8s.io/api/core/v1"

	api "github.com/minorhacks/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"
)

// ServiceSpecVerifyPorts verifies core.ServiceSpec to have reasonable ports specified
func ServiceSpecVerifyPorts(spec *core.ServiceSpec) error {
	for i := range spec.Ports {
		servicePort := &spec.Ports[i]
		if api.IsPortInvalid(servicePort.Port) {
			return fmt.Errorf(fmt.Sprintf("incorrect port :%d", servicePort.Port))
		}
	}
	return nil
}
