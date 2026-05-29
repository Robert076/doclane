.PHONY: up down destroy-all

# Bring up the full environment (~15-20 min cold start).
up:
	cd terraform/platform-compute && terraform init -input=false && terraform apply -auto-approve
	cd terraform/workload && terraform init -input=false && terraform apply -auto-approve

# Tear down everything compute-related ($0 cost after this).
down:
	cd terraform/workload && terraform destroy -auto-approve
	cd terraform/platform-compute && terraform destroy -auto-approve

# Nuclear option: destroy everything including persistent data (ECR, Cognito, certs).
destroy-all: down
	cd terraform/platform-data && terraform destroy -auto-approve
