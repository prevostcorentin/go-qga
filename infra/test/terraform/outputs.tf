output "vm_details" {
  description = "Détails des VMs créées"
  value = {
    for name, vm in libvirt_domain.test_vms : name => {
      ip_address = try(vm.network_interface[0].addresses[0], "pending")
      vnc_port   = vm.graphics[0].listen_address != "" ? "VNC disponible sur 127.0.0.1" : "VNC non configuré"
      status     = "running"
    }
  }
}

output "test_commands" {
  description = "Commandes pour lancer les tests"
  value = {
    test_script = "./run_tests.sh"
    manual_test = "go test -tags=integration ./tests/..."
    ssh_access = {
      for name, vm in libvirt_domain.test_vms : name => "ssh testuser@${try(vm.network_interface[0].addresses[0], "pending")}"
    }
  }
}
