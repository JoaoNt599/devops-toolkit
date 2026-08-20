package aws

import (
	"context"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Clients bundles every SDK client the collectors need.
type Clients struct {
	EC2 *ec2.Client
	ECS *ecs.Client
	ELB *elbv2.Client
	RDS *rds.Client
	IAM *iam.Client
	S3  *s3.Client
	ECR *ecr.Client
}

// NewClients loads the shared AWS config for the given profile/region and
// builds one client per service.
func NewClients(ctx context.Context, profile, region string) (*Clients, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	return &Clients{
		EC2: ec2.NewFromConfig(cfg),
		ECS: ecs.NewFromConfig(cfg),
		ELB: elbv2.NewFromConfig(cfg),
		RDS: rds.NewFromConfig(cfg),
		IAM: iam.NewFromConfig(cfg),
		S3:  s3.NewFromConfig(cfg),
		ECR: ecr.NewFromConfig(cfg),
	}, nil
}

// re-export a couple of aws helpers so collectors don't import the SDK aws pkg directly.
var (
	String   = awssdk.String
	ToString = awssdk.ToString
	ToBool   = awssdk.ToBool
	ToInt32  = awssdk.ToInt32
)

// chunk splits a slice into batches of at most size elements.
func chunk[T any](s []T, size int) [][]T {
	var out [][]T
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		out = append(out, s[i:end])
	}
	return out
}

// roleNameFromArn extracts the role name from an IAM role ARN.
// arn:aws:iam::acct:role/NAME -> NAME
func roleNameFromArn(arn string) string {
	if i := strings.LastIndex(arn, "/"); i != -1 {
		return arn[i+1:]
	}
	return arn
}
