#!/bin/bash

# Script to wait for Terraform state lock to be released
# Usage: ./scripts/wait-for-terraform-unlock.sh [environment]

set -e

ENVIRONMENT=${1:-"dev"}
TERRAFORM_DIR="infrastructure/terraform"
MAX_WAIT_TIME=1800  # 30 minutes max wait
POLL_INTERVAL=10    # Check every 10 seconds

echo "🔒 Waiting for Terraform state lock to be released for environment: $ENVIRONMENT"
echo "⏰ Max wait time: $((MAX_WAIT_TIME / 60)) minutes"
echo "🔄 Polling every $POLL_INTERVAL seconds"
echo ""

cd "$TERRAFORM_DIR"

# Initialize if needed
if [ ! -d ".terraform" ]; then
    echo "🔧 Initializing Terraform..."
    terraform init -backend-config="backend-${ENVIRONMENT}.hcl"
fi

start_time=$(date +%s)

while true; do
    current_time=$(date +%s)
    elapsed=$((current_time - start_time))

    if [ $elapsed -ge $MAX_WAIT_TIME ]; then
        echo "❌ Timeout: State lock was not released within $((MAX_WAIT_TIME / 60)) minutes"
        echo "💡 You may need to manually investigate and force-unlock if necessary"
        exit 1
    fi

    echo -n "⏳ Checking state lock status... (elapsed: $((elapsed / 60))m $((elapsed % 60))s)"

    # Try to acquire a state lock by running a harmless plan with detailed exit code
    if terraform plan -detailed-exitcode -var-file="environments/${ENVIRONMENT}.tfvars" -out=/dev/null >/dev/null 2>&1; then
        echo ""
        echo "✅ State lock is available! Terraform can now proceed."
        exit 0
    else
        # Check if the error is specifically about state lock
        error_output=$(terraform plan -detailed-exitcode -var-file="environments/${ENVIRONMENT}.tfvars" -out=/dev/null 2>&1 || true)

        if echo "$error_output" | grep -q "state lock"; then
            echo " 🔒 Still locked"
        else
            echo ""
            echo "❌ Different error encountered (not a lock issue):"
            echo "$error_output"
            exit 1
        fi
    fi

    sleep $POLL_INTERVAL
done
