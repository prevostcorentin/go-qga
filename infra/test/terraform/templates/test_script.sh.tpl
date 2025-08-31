#!/bin/sh

set -e

PROJECT_DIR="$(cd "$(dirname "$${BASH_SOURCE[0]}")" && pwd)"
GO_PROJECT_PATH="$${PROJECT_DIR}/go-qemu-project"

echo "=== Infrastructure Tests QEMU Guest Agent ==="
echo "Date: $(date)"

%{ for name, ip in vm_ips }
VM_${name}_IP="${ip}"
%{ endfor }


test_ssh_connection() {
    local vm_name=$1
    local vm_ip=$2
    
    echo "Testing SSH connection to $vm_name ($vm_ip)..."
    if ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no testuser@$vm_ip "echo 'SSH OK'"; then
        echo "✓ SSH connection successful to $vm_name"
        return 0
    else
        echo "✗ SSH connection failed to $vm_name"
        return 1
    fi
}

test_qemu_agent() {
    local vm_name=$1
    local vm_ip=$2
    
    echo "Testing QEMU Guest Agent on $vm_name..."
    if ssh testuser@$vm_ip "sudo systemctl is-active qemu-guest-agent"; then
        echo "✓ QEMU Guest Agent is running on $vm_name"
        
        echo "Testing agent communication..."
        if ssh testuser@$vm_ip "sudo qemu-ga --help >/dev/null 2>&1"; then
            echo "✓ QEMU Guest Agent communication OK on $vm_name"
            return 0
        fi
    fi
    echo "✗ QEMU Guest Agent test failed on $vm_name"
    return 1
}

test_go_project() {
    local vm_name=$1
    local vm_ip=$2
    
    echo "Deploying and testing Go project on $vm_name..."
    
    if [ -d "$GO_PROJECT_PATH" ]; then
        echo "Uploading project to $vm_name..."
        scp -r "$GO_PROJECT_PATH"/* testuser@$vm_ip:~/go-qemu-project/
        
        echo "Running Go tests on $vm_name..."
        if ssh testuser@$vm_ip "cd ~/go-qemu-project && go test -v -tags=integration ./..."; then
            echo "✓ Go tests passed on $vm_name"
            return 0
        fi
    fi
    echo "✗ Go project tests failed on $vm_name"
    return 1
}

FAILED_VMS=""

for vm_var in VM_*_IP; do
    vm_name="$${vm_var#VM_}"
    vm_name="$${vm_name%_IP}"
    vm_ip="$${!vm_var}"
    echo ""
    echo "=== Testing $vm_name ($vm_ip) ==="
   
    if [ "$vm_name" != "$${vm_name#*windows*}" ]; then
        echo "Windows VM detected - manual testing required"
        echo "RDP to $vm_ip for manual validation"
        continue
    fi
    
    if test_ssh_connection "$vm_name" "$vm_ip" && 
        test_qemu_agent "$vm_name" "$vm_ip" && 
        test_go_project "$vm_name" "$vm_ip"
    then
        echo "✓ All tests passed for $vm_name"
    else
        echo "✗ Some tests failed for $vm_name"
        FAILED_VMS="$${FAILED_VMS}$${vm_name} "
    fi
done


echo ""
echo "=== Test Summary ==="
if [ -z "$FAILED_VMS" ]
then
    echo "✓ All VM tests completed successfully!"
    exit 0
else
    echo "✗ Tests failed for: $FAILED_VMS"
    exit 1
fi
