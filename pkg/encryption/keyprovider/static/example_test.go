// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package static_test

import (
	"fmt"
	"strings"

	"github.com/we-dcode/opentofu/pkg/configs"
	"github.com/we-dcode/opentofu/pkg/encryption"
	"github.com/we-dcode/opentofu/pkg/encryption/config"
	"github.com/we-dcode/opentofu/pkg/encryption/keyprovider/static"
	"github.com/we-dcode/opentofu/pkg/encryption/method/aesgcm"
	"github.com/we-dcode/opentofu/pkg/encryption/registry/lockingencryptionregistry"
)

var hclConfig = `key_provider "static" "foo" {
  key = "6f6f706830656f67686f6834616872756f3751756165686565796f6f72653169"
}

method "aes_gcm" "bar" {
  keys = key_provider.static.foo
}

plan {
  method = method.aes_gcm.bar
}
`

// Example is a full end-to-end example of encrypting and decrypting a plan file.
func Example() {
	registry := lockingencryptionregistry.New()
	if err := registry.RegisterKeyProvider(static.New()); err != nil {
		panic(err)
	}
	if err := registry.RegisterMethod(aesgcm.New()); err != nil {
		panic(err)
	}

	cfg, diags := config.LoadConfigFromString("test.hcl", hclConfig)
	if diags.HasErrors() {
		panic(diags)
	}

	staticEvaluator := configs.NewStaticEvaluator(nil, configs.RootModuleCallForTesting())

	enc, diags := encryption.New(registry, cfg, staticEvaluator)
	if diags.HasErrors() {
		panic(diags)
	}

	encryptor := enc.Plan()

	encryptedPlan, err := encryptor.EncryptPlan([]byte("Hello world!"))
	if err != nil {
		panic(err)
	}
	if strings.Contains(string(encryptedPlan), "Hello world!") {
		panic("The plan was not encrypted!")
	}
	decryptedPlan, err := encryptor.DecryptPlan(encryptedPlan)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", decryptedPlan)
	// Output: Hello world!
}
