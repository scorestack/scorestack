variable "project" {}
variable "credentials_file" {}
variable "ssh_pub_key_file" {}
variable "ssh_priv_key_file" {}
variable "domain" {}

variable "region" {
    default = "us-central1"
}

variable "zone1" {
    default = "us-central1-a"
}

variable "zone2" {
    default = "us-central1-b"
}

variable "zone3" {
    default = "us-central1-c"
}

variable "ubuntu" {
    default = "family/ubuntu-1804-lts"
}

variable "ssh_user" {
    default = "ubuntu"
}

variable "inventory_destination" {
    default = "../ansible/inventory.ini"
}

variable "certificate_destination" {
    default = "../ansible/certificates"
}