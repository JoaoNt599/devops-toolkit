package aws

import (
	"context"

	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"

	"aws-inventory/internal/models"
	"aws-inventory/internal/report"
)

// CollectELB gathers load balancers, their listeners and rules, target groups
// and per-target health.
func CollectELB(ctx context.Context, c *Clients, rep *models.Report) {
	sec := rep.Section("elbv2")

	lbs, err := c.ELB.DescribeLoadBalancers(ctx, &elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		report.SetErr(sec, "loadBalancers", err)
		return
	}
	sec["loadBalancers"] = report.Jsonify(lbs.LoadBalancers)

	// Listeners + rules per LB
	var listeners []any
	var rules []any
	for _, lb := range lbs.LoadBalancers {
		ls, err := c.ELB.DescribeListeners(ctx, &elbv2.DescribeListenersInput{
			LoadBalancerArn: lb.LoadBalancerArn,
		})
		if err != nil {
			continue
		}
		listeners = append(listeners, report.Jsonify(ls.Listeners))
		for _, l := range ls.Listeners {
			rs, err := c.ELB.DescribeRules(ctx, &elbv2.DescribeRulesInput{
				ListenerArn: l.ListenerArn,
			})
			if err != nil {
				continue
			}
			rules = append(rules, report.Jsonify(rs.Rules))
		}
	}
	sec["listeners"] = listeners
	sec["listenerRules"] = rules

	// Target groups + health
	tgs, err := c.ELB.DescribeTargetGroups(ctx, &elbv2.DescribeTargetGroupsInput{})
	if err != nil {
		report.SetErr(sec, "targetGroups", err)
		return
	}
	sec["targetGroups"] = report.Jsonify(tgs.TargetGroups)

	var health []any
	for _, tg := range tgs.TargetGroups {
		h, err := c.ELB.DescribeTargetHealth(ctx, &elbv2.DescribeTargetHealthInput{
			TargetGroupArn: tg.TargetGroupArn,
		})
		if err != nil {
			continue
		}
		health = append(health, map[string]any{
			"targetGroupArn":    ToString(tg.TargetGroupArn),
			"targetHealthDescs": report.Jsonify(h.TargetHealthDescriptions),
		})
	}
	sec["targetHealth"] = health
}
