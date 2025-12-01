// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package common

import (
	commonpassword "github.com/sourcehawk/operator-api-mirror/eck-operator/v3.2.0/pkg/controller/common/password"
)

// FixedLengthRandomPasswordBytes generates a random password
func FixedLengthRandomPasswordBytes() []byte {
	return RandomBytes(24)
}

// RandomBytes generates some random bytes that can be used as a token or as a key
func RandomBytes(length int) []byte {
	return commonpassword.MustGenerate(length)
}
