// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"fmt"
	"os"

	"github.com/we-dcode/opentofu/pkg/configs"
	"github.com/we-dcode/opentofu/pkg/encryption"
	"github.com/we-dcode/opentofu/pkg/encryption/config"
	"github.com/we-dcode/opentofu/pkg/tfdiags"
)

const encryptionConfigEnvName = "TF_ENCRYPTION"

func (m *Meta) Encryption() (encryption.Encryption, tfdiags.Diagnostics) {
	path, err := os.Getwd()
	if err != nil {
		return nil, tfdiags.Diagnostics{}.Append(fmt.Errorf("Error getting pwd: %w", err))
	}

	return m.EncryptionFromPath(path)
}

func (m *Meta) EncryptionFromPath(path string) (encryption.Encryption, tfdiags.Diagnostics) {
	// This is not ideal, but given how fragmented the command package is, loading the root module here is our best option
	// See other meta commands like version check which do that same.
	module, diags := m.loadSingleModule(path, configs.SelectiveLoadEncryption)
	if diags.HasErrors() {
		return nil, diags
	}
	enc, encDiags := m.EncryptionFromModule(module)
	diags = diags.Append(encDiags)
	return enc, diags
}

func (m *Meta) EncryptionFromModule(module *configs.Module) (encryption.Encryption, tfdiags.Diagnostics) {
	cfg := module.Encryption
	var diags tfdiags.Diagnostics

	env := os.Getenv(encryptionConfigEnvName)
	if len(env) != 0 {
		envCfg, envDiags := config.LoadConfigFromString(encryptionConfigEnvName, env)
		diags = diags.Append(envDiags)
		if envDiags.HasErrors() {
			return nil, diags
		}
		cfg = cfg.Merge(envCfg)
	}

	enc, encDiags := encryption.New(encryption.DefaultRegistry, cfg, module.StaticEvaluator)
	diags = diags.Append(encDiags)

	return enc, diags
}
