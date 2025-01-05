// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/we-dcode/opentofu/pkg/grpcwrap"
	"github.com/we-dcode/opentofu/pkg/plugin"
	simple "github.com/we-dcode/opentofu/pkg/provider-simple"
	"github.com/we-dcode/opentofu/pkg/tfplugin5"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		GRPCProviderFunc: func() tfplugin5.ProviderServer {
			return grpcwrap.Provider(simple.Provider())
		},
	})
}
