// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package rbac

import (
	"context"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/sourcehawk/operator-api-mirrors/mirrors/otel-operator/v0.138.0/internal/autodetect/autodetectutils"
	"github.com/sourcehawk/operator-api-mirrors/mirrors/otel-operator/v0.138.0/internal/rbac"
)

// CheckRBACPermissions checks if the operator has the needed permissions to create RBAC resources automatically.
// If the RBAC is there, no errors nor warnings are returned.
func CheckRBACPermissions(ctx context.Context, reviewer *rbac.Reviewer) (admission.Warnings, error) {
	namespace, err := autodetectutils.GetOperatorNamespace()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "not possible to check RBAC rules", err)
	}

	serviceAccount, err := autodetectutils.GetOperatorServiceAccount()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "not possible to check RBAC rules", err)
	}

	rules := []*rbacv1.PolicyRule{
		{
			APIGroups: []string{"rbac.authorization.k8s.io"},
			Resources: []string{"clusterrolebindings", "clusterroles"},
			Verbs:     []string{"create", "delete", "get", "list", "patch", "update"},
		},
	}

	if subjectAccessReviews, err := reviewer.CheckPolicyRules(ctx, serviceAccount, namespace, rules...); err != nil {
		return nil, fmt.Errorf("%s: %w", "unable to check rbac rules", err)
	} else if allowed, deniedReviews := rbac.AllSubjectAccessReviewsAllowed(subjectAccessReviews); !allowed {
		return rbac.WarningsGroupedByResource(deniedReviews), nil
	}
	return nil, nil
}
