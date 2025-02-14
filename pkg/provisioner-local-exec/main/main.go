// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	localexec "github.com/we-dcode/opentofu/pkg/builtin/provisioners/local-exec"
	"github.com/we-dcode/opentofu/pkg/grpcwrap"
	"github.com/we-dcode/opentofu/pkg/plugin"
	"github.com/we-dcode/opentofu/pkg/tfplugin5"
)

func main() {
	// Provide a binary version of the internal terraform provider for testing
	plugin.Serve(&plugin.ServeOpts{
		GRPCProvisionerFunc: func() tfplugin5.ProvisionerServer {
			return grpcwrap.Provisioner(localexec.New())
		},
	})
}
