variable "environments" {
  description = "Test environments to deploy"
  type = map(object({
    image_url = string
    vcpus     = string
    memory    = number
    disk_size = string
  }))
  default = {
    ubuntu = {
      image_url = "https://cloud-images.ubuntu.com/releases/noble/release/ubuntu-24.04-server-cloudimg-amd64.img"
      vcpus     = 2
      memory    = 4096
      disk_size = "40G"
    }
    centos = {
      image_url = "https://cloud.centos.org/centos/10-stream/x86_64/images/CentOS-Stream-GenericCloud-10-latest.x86_64.qcow2"
      vcpus     = 2
      memory    = 4096
      disk_size = "40G"
    }
  }
}

variable "ssh_public_key" {
  description = "SSH public key to access VMs"
  type        = string
  default     = "~/.ssh/id_rsa.pub"
}
