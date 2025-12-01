// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package v1

import (
	commonv1 "github.com/sourcehawk/operator-api-mirrors/mirrors/eck-operator/pkg/apis/common/v1"
)

// GetIdentityLabels will return the common Elastic assigned labels for the APM Server.
func (as *ApmServer) GetIdentityLabels() map[string]string {
	return map[string]string{
		commonv1.TypeLabelName:    "apm-server",
		"apm.k8s.elastic.co/name": as.Name,
	}
}
