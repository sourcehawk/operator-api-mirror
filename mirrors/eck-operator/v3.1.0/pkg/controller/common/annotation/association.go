// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package annotation

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	commonv1 "github.com/sourcehawk/operator-api-mirrors/mirrors/eck-operator/v3.1.0/pkg/apis/common/v1"
)

// ForAssociationStatusChange constructs the annotation map for an association status change event.
func ForAssociationStatusChange(prevStatus, currStatus commonv1.AssociationStatusMap) (map[string]string, error) {
	return map[string]string{
		PrevAssocStatusAnnotation: prevStatus.String(),
		CurrAssocStatusAnnotation: currStatus.String(),
	}, nil
}

// ExtractAssociationStatusStrings extracts the association status strings from the provided meta object.
func ExtractAssociationStatusStrings(obj metav1.ObjectMeta) (prevStatus, currStatus string) {
	if obj.Annotations == nil {
		return "", ""
	}

	prevStatus = obj.Annotations[PrevAssocStatusAnnotation]
	currStatus = obj.Annotations[CurrAssocStatusAnnotation]
	return
}
