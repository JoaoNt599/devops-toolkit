package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"aws-inventory/internal/models"
	"aws-inventory/internal/report"
)

// CollectEC2 gathers VPCs, subnets, security groups (with full rules),
// route tables, NAT gateways, internet gateways and Elastic IPs.
// If vpcID is empty, resources across all VPCs are returned.
func CollectEC2(ctx context.Context, c *Clients, vpcID string, rep *models.Report) {
	sec := rep.Section("ec2")

	var idFilter []ec2types.Filter
	if vpcID != "" {
		idFilter = []ec2types.Filter{{Name: String("vpc-id"), Values: []string{vpcID}}}
	}

	// VPCs
	vpcs, err := c.EC2.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{Filters: idFilter})
	if err != nil {
		report.SetErr(sec, "vpcs", err)
	} else {
		sec["vpcs"] = report.Jsonify(vpcs.Vpcs)
	}

	// Subnets
	subnets, err := c.EC2.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{Filters: idFilter})
	if err != nil {
		report.SetErr(sec, "subnets", err)
	} else {
		sec["subnets"] = report.Jsonify(subnets.Subnets)
	}

	// Security groups (with full ingress/egress rules)
	sgs, err := c.EC2.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{Filters: idFilter})
	if err != nil {
		report.SetErr(sec, "securityGroups", err)
	} else {
		sec["securityGroups"] = report.Jsonify(sgs.SecurityGroups)
	}

	// Route tables
	rts, err := c.EC2.DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{Filters: idFilter})
	if err != nil {
		report.SetErr(sec, "routeTables", err)
	} else {
		sec["routeTables"] = report.Jsonify(rts.RouteTables)
	}

	// NAT gateways
	nats, err := c.EC2.DescribeNatGateways(ctx, &ec2.DescribeNatGatewaysInput{Filter: idFilter})
	if err != nil {
		report.SetErr(sec, "natGateways", err)
	} else {
		sec["natGateways"] = report.Jsonify(nats.NatGateways)
	}

	// Internet gateways (attached to the vpc)
	var igwFilter []ec2types.Filter
	if vpcID != "" {
		igwFilter = []ec2types.Filter{{Name: String("attachment.vpc-id"), Values: []string{vpcID}}}
	}
	igws, err := c.EC2.DescribeInternetGateways(ctx, &ec2.DescribeInternetGatewaysInput{Filters: igwFilter})
	if err != nil {
		report.SetErr(sec, "internetGateways", err)
	} else {
		sec["internetGateways"] = report.Jsonify(igws.InternetGateways)
	}

	// Elastic IPs (outbound whitelisting / NAT)
	eips, err := c.EC2.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{})
	if err != nil {
		report.SetErr(sec, "elasticIPs", err)
	} else {
		sec["elasticIPs"] = report.Jsonify(eips.Addresses)
	}
}
