// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tofu

import (
	"github.com/we-dcode/opentofu/pkg/addrs"
	"github.com/we-dcode/opentofu/pkg/tfdiags"
)

// GraphWalker is an interface that can be implemented that when used
// with Graph.Walk will invoke the given callbacks under certain events.
type GraphWalker interface {
	EvalContext() EvalContext
	EnterPath(addrs.ModuleInstance) EvalContext
	ExitPath(addrs.ModuleInstance)
	Execute(EvalContext, GraphNodeExecutable) tfdiags.Diagnostics
}

// NullGraphWalker is a GraphWalker implementation that does nothing.
// This can be embedded within other GraphWalker implementations for easily
// implementing all the required functions.
type NullGraphWalker struct{}

func (NullGraphWalker) EvalContext() EvalContext                                     { return new(MockEvalContext) }
func (NullGraphWalker) EnterPath(addrs.ModuleInstance) EvalContext                   { return new(MockEvalContext) }
func (NullGraphWalker) ExitPath(addrs.ModuleInstance)                                {}
func (NullGraphWalker) Execute(EvalContext, GraphNodeExecutable) tfdiags.Diagnostics { return nil }
