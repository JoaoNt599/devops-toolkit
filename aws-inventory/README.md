# aws-inventory

Inspects the AWS architecture in use (ECS, EC2/VPC, IAM, RDS, S3, ECR, ALB) and exports a single JSON with the raw and complete response of each `Describe*`. Idiomatic Go reorganization of the original monolithic script.

## Structure

```
aws-inventory/
├── cmd/
│   └── inventory/
│       └── main.go            # thin entrypoint: flags + orchestration
├── internal/
│   ├── aws/                   # collectors, one per domain
│   │   ├── client.go          # Clients (SDK) + helpers (chunk, roleNameFromArn)
│   │   ├── vpc.go             # CollectEC2: vpc, subnets, SGs, route tables, NAT, IGW, EIP
│   │   ├── security_group.go  # CollectSecurityGroups: SGs by ID (complement)
│   │   ├── ecs.go             # CollectECS: cluster, services, task defs, tasks
│   │   ├── alb.go             # CollectELB: LBs, listeners, rules, TGs, health
│   │   ├── iam.go             # CollectIAM: expands roles from task defs
│   │   └── misc.go            # CollectRDS / CollectS3 / CollectECR
│   ├── models/
│   │   └── models.go          # Report + Section()
│   └── report/
│       └── report.go          # Jsonify, SetErr, Write
├── output/
│   └── architecture.json      # generated at runtime
└── go.mod
```

The module is `github.com/corporation/aws-inventory` (adjust to the actual path of your organization).

## Build

```bash
go mod tidy
go build -o inventory ./cmd/inventory
```

## Usage

```bash
./inventory \
  -profile corporation \
  -region us-east-1 \
  -cluster UAT-Cluster \
  -vpc-id vpc-xxxxxxxx \
  -out output/architecture.json
```

Without `-vpc-id`, it collects all VPCs in the account. Without `-cluster`, the ECS section is skipped.
The `-cluster` chains services → complete task definitions → IAM roles.

## Design Notes

- **Independent collectors:** each `Collect*` writes to its report section; a failure becomes an `_error` in that section and does not abort the rest.
- **Raw response preserved:** `report.Jsonify` serializes the SDK struct into `map[string]any`, keeping all fields (env vars, roles, log config, SG rules, etc.).
- **ECS → IAM:** `CollectECS` returns the role names referenced by the task defs, which `main` passes to `CollectIAM` for expansion.