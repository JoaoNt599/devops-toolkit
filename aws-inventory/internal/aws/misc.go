package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"aws-inventory/internal/models"
	"aws-inventory/internal/report"
)

// CollectRDS gathers DB instances and DB subnet groups.
func CollectRDS(ctx context.Context, c *Clients, rep *models.Report) {
	sec := rep.Section("rds")

	dbs, err := c.RDS.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		report.SetErr(sec, "dbInstances", err)
	} else {
		sec["dbInstances"] = report.Jsonify(dbs.DBInstances)
	}

	subs, err := c.RDS.DescribeDBSubnetGroups(ctx, &rds.DescribeDBSubnetGroupsInput{})
	if err != nil {
		report.SetErr(sec, "dbSubnetGroups", err)
	} else {
		sec["dbSubnetGroups"] = report.Jsonify(subs.DBSubnetGroups)
	}
}

// CollectS3 lists buckets and, for each, its region, encryption, public-access
// block and versioning status.
func CollectS3(ctx context.Context, c *Clients, rep *models.Report) {
	sec := rep.Section("s3")

	list, err := c.S3.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		report.SetErr(sec, "buckets", err)
		return
	}

	var buckets []any
	for _, b := range list.Buckets {
		entry := map[string]any{"name": ToString(b.Name), "creationDate": b.CreationDate}

		if loc, err := c.S3.GetBucketLocation(ctx, &s3.GetBucketLocationInput{Bucket: b.Name}); err == nil {
			entry["locationConstraint"] = string(loc.LocationConstraint)
		}
		if enc, err := c.S3.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{Bucket: b.Name}); err == nil {
			entry["encryption"] = report.Jsonify(enc.ServerSideEncryptionConfiguration)
		}
		if pab, err := c.S3.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{Bucket: b.Name}); err == nil {
			entry["publicAccessBlock"] = report.Jsonify(pab.PublicAccessBlockConfiguration)
		}
		if ver, err := c.S3.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{Bucket: b.Name}); err == nil {
			entry["versioning"] = string(ver.Status)
		}

		buckets = append(buckets, entry)
	}
	sec["buckets"] = buckets
}

// CollectECR gathers repositories and their lifecycle policies.
func CollectECR(ctx context.Context, c *Clients, rep *models.Report) {
	sec := rep.Section("ecr")

	repos, err := c.ECR.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{})
	if err != nil {
		report.SetErr(sec, "repositories", err)
		return
	}
	sec["repositories"] = report.Jsonify(repos.Repositories)

	var lifecycles []any
	for _, r := range repos.Repositories {
		lp, err := c.ECR.GetLifecyclePolicy(ctx, &ecr.GetLifecyclePolicyInput{
			RepositoryName: r.RepositoryName,
		})
		if err != nil {
			continue
		}
		lifecycles = append(lifecycles, map[string]any{
			"repositoryName":      ToString(r.RepositoryName),
			"lifecyclePolicyText": ToString(lp.LifecyclePolicyText),
		})
	}
	sec["lifecyclePolicies"] = lifecycles
}
