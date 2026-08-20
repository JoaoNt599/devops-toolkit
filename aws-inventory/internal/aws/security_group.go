package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"aws-inventory/internal/models"
	"aws-inventory/internal/report"
)

// CollectSecurityGroups fetches full detail (ingress + egress rules) for a
// specific set of security-group IDs and stores them under ec2.securityGroupsDetail.
// This is a targeted complement to CollectEC2, useful when you want to expand
// the SGs actually referenced by ECS services / RDS instances.
func CollectSecurityGroups(ctx context.Context, c *Clients, groupIDs []string, rep *models.Report) {
	if len(groupIDs) == 0 {
		return
	}
	sec := rep.Section("ec2")

	var all []any
	for _, batch := range chunk(groupIDs, 50) { // API allows many; batch defensively
		out, err := c.EC2.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
			GroupIds: batch,
		})
		if err != nil {
			report.SetErr(sec, "securityGroupsDetail", err)
			return
		}
		all = append(all, report.Jsonify(out.SecurityGroups))
	}
	sec["securityGroupsDetail"] = all
}
