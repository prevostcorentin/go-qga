resource "libvirt_pool" "pool" {
  name = "qemu-test-pool"
  type = "dir"
  target {
    path = "/var/lib/libvirt/images/qemu-tests"
  }
}

resource "libvirt_cloudinit_disk" "commoninit" {
  for_each = var.environments

  name = "${each.key}-cloudinit.iso"
  pool = libvirt_pool.pool.name
  user_data = templatefile("${path.module}/cloud-init/${each.key}-user-data.yml", {
    ssh_public_key = var.ssh_public_key
  })

  network_config = templatefile("${path.module}/cloud-init/network-config.yml", {})
}

resource "libvirt_volume" "base_images" {
  for_each = var.environments

  name   = "${each.key}-base.qcow2"
  pool   = libvirt_pool.pool.name
  source = each.value.image_url
  format = "qcow2"
}

resource "libvirt_volume" "vm_disks" {
  for_each = var.environments

  name           = "${each.key}-disk.qcow2"
  pool           = libvirt_pool.pool.name
  base_volume_id = libvirt_volume.base_images[each.key].id
  size           = parseint(replace(each.value.disk_size, "G", ""), 10) * 1024 * 1024 * 1024
}

resource "libvirt_network" "network" {
  name      = "qemu-test-net"
  mode      = "nat"
  domain    = "test.local"
  addresses = ["192.168.100.0/24"]

  dns {
    enabled = true
  }
}

resource "libvirt_domain" "test_vms" {
  for_each = var.environments

  name   = "qemu-test-${each.key}"
  memory = each.value.memory
  vcpu   = each.value.vcpus

  qemu_agent = true
  running    = true

  depends_on = [
    libvirt_pool.pool,
    libvirt_cloudinit_disk.commoninit,
    libvirt_volume.vm_disks,
    libvirt_network.network
  ]

  boot_device {
    dev = ["hd"]
  }

  disk {
    volume_id = libvirt_volume.vm_disks[each.key].id
    scsi      = true
  }

  disk {
    file = "${libvirt_pool.pool.target[0].path}/${each.key}-cloudinit.iso"
  }

  network_interface {
    network_name   = libvirt_network.network.name
    wait_for_lease = false
  }

  console {
    type        = "pty"
    target_port = "0"
    target_type = "serial"
  }

  graphics {
    type           = "vnc"
    listen_type    = "address"
    listen_address = "127.0.0.1"
  }

  timeouts {
    create = "5m"
  }
}

resource "local_file" "script" {
  filename = "${path.module}/run_tests.sh"
  content = templatefile("${path.module}/templates/test_script.sh.tpl", {
    vm_ips = {
      for name, vm in libvirt_domain.test_vms : name => vm.network_interface[0].addresses[0]
    }
  })
  file_permission = "0755"
}
