// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arguments

import (
	"github.com/we-dcode/opentofu/pkg/tfdiags"
)

// Output represents the command-line arguments for the output command.
type Output struct {
	// Name identifies which root module output to show.  If empty, show all
	// outputs.
	Name string

	// StatePath is an optional path to a state file, from which outputs will
	// be loaded.
	StatePath string

	// ViewType specifies which output format to use: human, JSON, or "raw".
	ViewType ViewType

	Vars *Vars

	// ShowSensitive is used to display the value of variables marked as sensitive.
	ShowSensitive bool
}

// ParseOutput processes CLI arguments, returning an Output value and errors.
// If errors are encountered, an Output value is still returned representing
// the best effort interpretation of the arguments.
func ParseOutput(args []string) (*Output, tfdiags.Diagnostics) {
	var diags tfdiags.Diagnostics
	output := &Output{
		Vars: &Vars{},
	}

	var jsonOutput, rawOutput bool
	var statePath string
	cmdFlags := extendedFlagSet("output", nil, nil, output.Vars)
	cmdFlags.BoolVar(&jsonOutput, "json", false, "json")
	cmdFlags.BoolVar(&rawOutput, "raw", false, "raw")
	cmdFlags.StringVar(&statePath, "state", "", "path")
	cmdFlags.BoolVar(&output.ShowSensitive, "show-sensitive", false, "displays sensitive values")

	if err := cmdFlags.Parse(args); err != nil {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Failed to parse command-line flags",
			err.Error(),
		))
	}

	args = cmdFlags.Args()
	if len(args) > 1 {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Unexpected argument",
			"The output command expects exactly one argument with the name of an output variable or no arguments to show all outputs.",
		))
	}

	if jsonOutput && rawOutput {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Invalid output format",
			"The -raw and -json options are mutually-exclusive.",
		))

		// Since the desired output format is unknowable, fall back to default
		jsonOutput = false
		rawOutput = false
	}

	output.StatePath = statePath

	if len(args) > 0 {
		output.Name = args[0]
	}

	if rawOutput && output.Name == "" {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Output name required",
			"You must give the name of a single output value when using the -raw option.",
		))
	}

	switch {
	case jsonOutput:
		output.ViewType = ViewJSON
	case rawOutput:
		output.ViewType = ViewRaw
	default:
		output.ViewType = ViewHuman
	}

	return output, diags
}
