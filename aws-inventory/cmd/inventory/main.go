package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	awsx "aws-inventory/internal/aws"
	"aws-inventory/internal/models"
	"aws-inventory/internal/report"
)

func main() {
	var (
		profile = flag.String("profile", "corporation", "AWS profile name")
		region  = flag.String("region", "us-east-1", "AWS region")
		cluster = flag.String("cluster", "", "ECS cluster name to inspect")
		vpcID   = flag.String("vpc-id", "", "VPC id to inspect (default: all VPCs discovered)")
		outFile = flag.String("out", "output/architecture.json", "Output JSON path")
	)
	flag.Parse()

	ctx := context.Background()

	clients, err := awsx.NewClients(ctx, *profile, *region)
	if err != nil {
		fatal("loading AWS config: %v", err)
	}

	rep := &models.Report{
		Meta: map[string]any{
			"generatedAt": time.Now().UTC().Format(time.RFC3339),
			"profile":     *profile,
			"region":      *region,
			"cluster":     *cluster,
			"vpcId":       *vpcID,
		},
		Data: map[string]map[string]any{},
	}

	// Orchestrate collectors. Each writes into its own section of the report;
	// a failure in one becomes a scoped _error and does not abort the rest.
	awsx.CollectEC2(ctx, clients, *vpcID, rep)
	roleNames := awsx.CollectECS(ctx, clients, *cluster, rep)
	awsx.CollectELB(ctx, clients, rep)
	awsx.CollectRDS(ctx, clients, rep)
	awsx.CollectIAM(ctx, clients, roleNames, rep)
	awsx.CollectS3(ctx, clients, rep)
	awsx.CollectECR(ctx, clients, rep)

	if err := report.Write(rep, *outFile); err != nil {
		fatal("%v", err)
	}
}

func fatal(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "FATAL: "+format+"\n", a...)
	os.Exit(1)
}
