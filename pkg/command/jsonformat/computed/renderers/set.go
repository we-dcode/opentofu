// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package renderers

import (
	"bytes"
	"fmt"

	"github.com/we-dcode/opentofu/pkg/command/jsonformat/computed"
	"github.com/we-dcode/opentofu/pkg/plans"
)

var _ computed.DiffRenderer = (*setRenderer)(nil)

func Set(elements []computed.Diff) computed.DiffRenderer {
	return &setRenderer{
		elements: elements,
	}
}

func NestedSet(elements []computed.Diff) computed.DiffRenderer {
	return &setRenderer{
		elements:                  elements,
		overrideForcesReplacement: true,
	}
}

type setRenderer struct {
	NoWarningsRenderer

	elements []computed.Diff

	overrideForcesReplacement bool
}

func (renderer setRenderer) RenderHuman(diff computed.Diff, indent int, opts computed.RenderHumanOpts) string {
	// Sets are a bit finicky, nested sets don't render the forces replacement
	// suffix themselves, but push it onto their children. So if we are
	// overriding the forces replacement setting, we set it to true for children
	// and false for ourselves.
	displayForcesReplacementInSelf := diff.Replace && !renderer.overrideForcesReplacement
	displayForcesReplacementInChildren := diff.Replace && renderer.overrideForcesReplacement

	if len(renderer.elements) == 0 {
		return fmt.Sprintf("[]%s%s", nullSuffix(diff.Action, opts), forcesReplacement(displayForcesReplacementInSelf, opts))
	}

	elementOpts := opts.Clone()
	elementOpts.OverrideNullSuffix = true
	elementOpts.OverrideForcesReplacement = displayForcesReplacementInChildren

	unchangedElements := 0

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("[%s\n", forcesReplacement(displayForcesReplacementInSelf, opts)))
	for _, element := range renderer.elements {
		if element.Action == plans.NoOp && !opts.ShowUnchangedChildren {
			unchangedElements++
			continue
		}

		for _, warning := range element.WarningsHuman(indent+1, opts) {
			buf.WriteString(fmt.Sprintf("%s%s\n", formatIndent(indent+1), warning))
		}
		buf.WriteString(fmt.Sprintf("%s%s%s,\n", formatIndent(indent+1), writeDiffActionSymbol(element.Action, elementOpts), element.RenderHuman(indent+1, elementOpts)))
	}

	if unchangedElements > 0 {
		buf.WriteString(fmt.Sprintf("%s%s%s\n", formatIndent(indent+1), writeDiffActionSymbol(plans.NoOp, opts), unchanged("element", unchangedElements, opts)))
	}

	buf.WriteString(fmt.Sprintf("%s%s]%s", formatIndent(indent), writeDiffActionSymbol(plans.NoOp, opts), nullSuffix(diff.Action, opts)))
	return buf.String()
}
