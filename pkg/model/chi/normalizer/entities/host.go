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

package entities

import (
	core "k8s.io/api/core/v1"

	api "github.com/minorhacks/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"
	model "github.com/minorhacks/clickhouse-operator/pkg/model/chi"
)

// NormalizeHostPorts ensures api.ChiReplica.Port is reasonable
func NormalizeHostPorts(host *api.ChiHost) {
	// Walk over all assigned ports of the host and append each port to the list of service's ports
	model.HostWalkInvalidPorts(
		host,
		func(name string, port *int32, protocol core.Protocol) bool {
			*port = api.PortUnassigned()
			// Do not abort, continue iterating
			return false
		},
	)
}
