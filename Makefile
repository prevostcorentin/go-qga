TERRAFORM_DIR ?= infra/test/terraform
TF_VAR_FILE ?= terraform.tfvars
SSH_KEY_PATH ?= ~/.ssh/id_rsa.pub
TERRAFORM_VERSION ?= 1.13.1

UNAME_S := $(shell uname -s)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

check-terraform:
	@which terraform > /dev/null || (echo "❌ Terraform is not installed. Run 'make install-deps'" && exit 1)
	@echo "✅ Terraform found: $(terraform version --json | jq -r '.terraform_version')"

install-deps:
	@echo "Install dependencies..."
	@if [ "$(UNAME_S)" = "Linux" ]; then \
		echo "Install on Linux..."; \
		if command -v apt-get >/dev/null 2>&1; then \
			sudo apt-get update && sudo apt-get install -y wget unzip jq libvirt-daemon-system libvirt-clients qemu-kvm qemu-system-x86-64; \
			sudo systemctl enable --now libvirtd; \
		elif command -v dnf >/dev/null 2>&1; then \
			sudo dnf install -y @virtualization virt-manager libvirt qemu-kvm jq; \
			sudo systemctl enable --now libvirtd; \
		elif command -v dnf >/dev/null 2>&1; then \
			sudo dnf install -y wget unzip jq libvirt-daemon libvirt-client qemu-kvm; \
			sudo systemctl enable --now libvirtd; \
		fi; \
		wget -q https://releases.hashicorp.com/terraform/$(TERRAFORM_VERSION)/terraform_$(TERRAFORM_VERSION)_linux_amd64.zip; \
		unzip -q terraform_$(TERRAFORM_VERSION)_linux_amd64.zip; \
		sudo mv terraform /usr/local/bin/; \
		rm terraform_$(TERRAFORM_VERSION)_linux_amd64.zip LICENSE.txt; \
		echo "Create libvirt group..."; \
		sudo addgroup libvirt; \
		sudo usermod -aG libvirt $(USER); \
		echo "✅ Added $(USER) to libvirt group"; \
	else \
		echo "Unsupported OS."; \
		echo "Install Terraform manually from https://www.terraform.io/downloads"; \
		exit 1; \
	fi
	@echo "⚠️You may need to log out and back in for group changes to take effect"
	@echo "✅ Installation complete"


init: check-terraform
	@cd $(TERRAFORM_DIR) && terraform init
	@if [ ! -f "$(TERRAFORM_DIR)/$(TF_VAR_FILE)" ]; then \
		echo "Creating default terraform.tfvars..."; \
		echo 'ssh_public_key = "'$$(cat $(SSH_KEY_PATH))'"' > $(TERRAFORM_DIR)/$(TF_VAR_FILE); \
	fi

plan: init
	cd $(TERRAFORM_DIR) && terraform plan --var-file=$(TF_VAR_FILE)

clean-libvirt:
	@echo "🧹 Cleaning up existing libvirt domains..."
	@for domain in $$(sudo virsh list --all --name | grep "qemu-test-"); do \
		echo "  Stopping domain: $$domain"; \
		sudo virsh destroy $$domain 2>/dev/null || true; \
		echo "  Removing domain: $$domain"; \
		sudo virsh undefine $$domain --remove-all-storage 2>/dev/null || true; \
	done
	@echo "🧹 Cleaning up libvirt volumes..."
	@for vol in $$(sudo virsh vol-list qemu-test-pool 2>/dev/null | grep -E "(ubuntu|centos)" | awk '{print $$1}' || true); do \
		echo "  Removing volume: $$vol"; \
		sudo virsh vol-delete $$vol --pool qemu-test-pool 2>/dev/null || true; \
	done
	@echo "✅ Libvirt cleanup complete"

clean-terraform:
	@echo "🧹 Cleaning Terraform state..."
	@cd $(TERRAFORM_DIR) && { \
		terraform state list | grep -E "(libvirt_domain|libvirt_cloudinit_disk|libvirt_volume)" | xargs -r terraform state rm; \
	} 2>/dev/null || true
	@echo "✅ Terraform state cleanup complete"

deploy: clean-libvirt clean-terraform init
	cd $(TERRAFORM_DIR) && terraform apply -var-file=$(TF_VAR_FILE) -auto-approve
	@echo "Infrastructure has been deployed! Wait 2-3 minutes for the VMs to start up completely."

deploy-linux: init
	cd $(TERRAFORM_DIR) && terraform apply -target=libvirt_domain.test_vms[ubuntu_latest] -target=libvirt_domain.test_vms[centos_stream] -var-file=$(TF_VAR_FILE)

deploy-windows: init
	cd $(TERRAFORM_DIR) && terraform apply -target=libvirt_domain.test_vms[windows_server] -var-file=$(TF_VAR_FILE)

status:
	@cd $(TERRAFORM_DIR) && terraform output -json vm_details | jq -r 'to_entries[] | "\(.key): \(.value.ip_address) (\(.value.status))"'
	@echo ""
	@cd $(TERRAFORM_DIR) &&terraform output -raw test_commands || echo "Pas de commandes de test disponibles"

test:
	@echo "Starting tests on all VMs..."
	@if [ -f "$(TERRAFORM_DIR)/run_tests.sh" ]; then \
		cd $(TERRAFORM_DIR) && ./run_tests.sh; \
	else \
		echo "Test script not found. Deploy the infrastructure first with 'make deploy'"; \
	fi

test-ubuntu: deploy
	@cd $(TERRAFORM_DIR) && vm_ip=$$(terraform output -json vm_details | jq -r '.ubuntu_latest.ip_address'); \
	if [ "$$vm_ip" != "null" ] && [ "$$vm_ip" != "pending" ]; then \
		echo "Testing Ubuntu VM at $$vm_ip..."; \
		ssh testuser@$$vm_ip "cd ~/go-qemu-project && go test -v ./..."; \
	else \
		echo "Ubuntu VM not ready or not found"; \
	fi

test-centos: deploy
	@cd $(TERRAFORM_DIR) && vm_ip=$$(terraform output -json vm_details | jq -r '.centos_stream.ip_address'); \
	if [ "$$vm_ip" != "null" ] && [ "$$vm_ip" != "pending" ]; then \
		echo "Testing CentOS VM at $$vm_ip..."; \
		ssh testuser@$$vm_ip "cd ~/go-qemu-project && go test -v ./..."; \
	else \
		echo "CentOS VM not ready or not found"; \
	fi

test-manual:
	@echo "=== Manual tests ==="
	@cd $(TERRAFORM_DIR) && terraform output test_commands

destroy:
	cd $(TERRAFORM_DIR) && terraform destroy -var-file=$(TF_VAR_FILE) -auto-approve

clean: destroy
	rm -f $(TERRAFORM_DIR)/run_tests.sh
	rm -f $(TERRAFORM_DIR)/terraform.tfstate*
	rm -f $(TERRAFORM_DIR)/.terraform.lock.hcl
	rm -rf $(TERRAFORM_DIR)/.terraform/

ssh-ubuntu:
	@cd $(TERRAFORM_DIR) && vm_ip=$$(terraform output -json vm_details | jq -r '.ubuntu_latest.ip_address'); \
	ssh testuser@$$vm_ip

ssh-centos:
	@cd $(TERRAFORM_DIR) && vm_ip=$$(terraform output -json vm_details | jq -r '.centos_stream.ip_address'); \
	ssh testuser@$$vm_ip

logs:
	cd $(TERRAFORM_DIR) && terraform show

monitor:
	@while true; do \
		clear; \
		echo "=== VM Status - $(date) ==="; \
		make status; \
		sleep 10; \
	done

.PHONY: help install-deps plan deploy test destroy clean status
