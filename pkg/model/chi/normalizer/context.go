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

package normalizer

import api "github.com/minorhacks/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"

// Context specifies CHI-related normalization context
type Context struct {
	// chi specifies current CHI being normalized
	chi *api.ClickHouseInstallation
	// options specifies normalization options
	options *Options
}

// NewContext creates new Context
func NewContext(options *Options) *Context {
	return &Context{
		options: options,
	}
}

func (c *Context) GetTarget() *api.ClickHouseInstallation {
	if c == nil {
		return nil
	}
	return c.chi
}

func (c *Context) SetTarget(chi *api.ClickHouseInstallation) *api.ClickHouseInstallation {
	if c == nil {
		return nil
	}
	c.chi = chi
	return c.chi
}

func (c *Context) Options() *Options {
	if c == nil {
		return nil
	}
	return c.options
}
