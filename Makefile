.PHONY: up down destroy-all

up:
	cd terraform/platform-compute && terraform init -input=false && terraform apply -auto-approve
	cd terraform/workload && terraform init -input=false && terraform apply -auto-approve

down:
	cd terraform/workload && terraform destroy -auto-approve
	cd terraform/platform-compute && terraform destroy -auto-approve

destroy-all: down
	cd terraform/platform-data && terraform destroy -auto-approve
