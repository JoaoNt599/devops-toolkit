package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"

	"aws-inventory/internal/models"
	"aws-inventory/internal/report"
)

// CollectIAM expands the given role names (typically discovered from ECS task
// definitions) with their role detail, attached managed policies and inline
// policy names.
func CollectIAM(ctx context.Context, c *Clients, roleNames []string, rep *models.Report) {
	sec := rep.Section("iam")
	if len(roleNames) == 0 {
		sec["_note"] = "no task-def roles discovered"
		return
	}

	var roles []any
	for _, name := range roleNames {
		entry := map[string]any{"roleName": name}

		gr, err := c.IAM.GetRole(ctx, &iam.GetRoleInput{RoleName: String(name)})
		if err != nil {
			entry["_error"] = err.Error()
			roles = append(roles, entry)
			continue
		}
		entry["role"] = report.Jsonify(gr.Role)

		if ap, err := c.IAM.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{
			RoleName: String(name),
		}); err == nil {
			entry["attachedPolicies"] = report.Jsonify(ap.AttachedPolicies)
		}

		if ip, err := c.IAM.ListRolePolicies(ctx, &iam.ListRolePoliciesInput{
			RoleName: String(name),
		}); err == nil {
			entry["inlinePolicyNames"] = report.Jsonify(ip.PolicyNames)
		}

		roles = append(roles, entry)
	}
	sec["roles"] = roles
}
