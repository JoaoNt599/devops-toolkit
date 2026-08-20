package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"aws-inventory/internal/models"
	"aws-inventory/internal/report"
)

// clusterFields is the set of extra detail requested on DescribeClusters.
func clusterFields() []ecstypes.ClusterField {
	return []ecstypes.ClusterField{
		ecstypes.ClusterFieldStatistics,
		ecstypes.ClusterFieldAttachments,
		ecstypes.ClusterFieldSettings,
	}
}

// CollectECS gathers the cluster, its services, the full task definitions they
// reference, and the running tasks. It returns the IAM role names referenced by
// those task definitions so the caller can expand them via CollectIAM.
func CollectECS(ctx context.Context, c *Clients, cluster string, rep *models.Report) []string {
	sec := rep.Section("ecs")
	roleSet := map[string]struct{}{}

	if cluster == "" {
		sec["_note"] = "no cluster provided; skipping ECS"
		return nil
	}

	// Cluster detail
	cl, err := c.ECS.DescribeClusters(ctx, &ecs.DescribeClustersInput{
		Clusters: []string{cluster},
		Include:  clusterFields(),
	})
	if err != nil {
		report.SetErr(sec, "cluster", err)
	} else {
		sec["cluster"] = report.Jsonify(cl.Clusters)
	}

	// List services
	var serviceArns []string
	sList, err := c.ECS.ListServices(ctx, &ecs.ListServicesInput{Cluster: String(cluster)})
	if err != nil {
		report.SetErr(sec, "services", err)
	} else {
		serviceArns = sList.ServiceArns
	}

	// Describe services (batches of 10) and collect their task-def ARNs.
	tdSet := map[string]struct{}{}
	if len(serviceArns) > 0 {
		var allServices []any
		for _, batch := range chunk(serviceArns, 10) {
			ds, err := c.ECS.DescribeServices(ctx, &ecs.DescribeServicesInput{
				Cluster:  String(cluster),
				Services: batch,
			})
			if err != nil {
				report.SetErr(sec, "servicesDetail", err)
				break
			}
			allServices = append(allServices, report.Jsonify(ds.Services))
			for _, s := range ds.Services {
				if s.TaskDefinition != nil {
					tdSet[*s.TaskDefinition] = struct{}{}
				}
			}
		}
		sec["services"] = allServices
	}

	// Describe each referenced task definition in full.
	var taskDefs []any
	for td := range tdSet {
		out, err := c.ECS.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
			TaskDefinition: String(td),
		})
		if err != nil {
			taskDefs = append(taskDefs, map[string]any{"_error": err.Error(), "taskDefinition": td})
			continue
		}
		taskDefs = append(taskDefs, report.Jsonify(out.TaskDefinition))
		if out.TaskDefinition != nil {
			if out.TaskDefinition.TaskRoleArn != nil {
				roleSet[*out.TaskDefinition.TaskRoleArn] = struct{}{}
			}
			if out.TaskDefinition.ExecutionRoleArn != nil {
				roleSet[*out.TaskDefinition.ExecutionRoleArn] = struct{}{}
			}
		}
	}
	sec["taskDefinitions"] = taskDefs

	// Running tasks (ENIs, attachments, etc.)
	tList, err := c.ECS.ListTasks(ctx, &ecs.ListTasksInput{Cluster: String(cluster)})
	if err == nil && len(tList.TaskArns) > 0 {
		dt, err := c.ECS.DescribeTasks(ctx, &ecs.DescribeTasksInput{
			Cluster: String(cluster),
			Tasks:   tList.TaskArns,
		})
		if err == nil {
			sec["tasks"] = report.Jsonify(dt.Tasks)
		}
	}

	// Return role names for IAM expansion.
	var roleNames []string
	for arn := range roleSet {
		roleNames = append(roleNames, roleNameFromArn(arn))
	}
	return roleNames
}
