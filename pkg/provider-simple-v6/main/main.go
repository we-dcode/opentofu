// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/we-dcode/opentofu/pkg/grpcwrap"
	plugin "github.com/we-dcode/opentofu/pkg/plugin6"
	simple "github.com/we-dcode/opentofu/pkg/provider-simple-v6"
	"github.com/we-dcode/opentofu/pkg/tfplugin6"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		GRPCProviderFunc: func() tfplugin6.ProviderServer {
			return grpcwrap.Provider6(simple.Provider())
		},
	})
}
