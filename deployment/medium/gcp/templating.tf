resource "random_password" "bootstrap" {
    length = 32
    special = false
}

resource "random_password" "elastic" {
    length = 32
    special = false
}

resource "random_password" "kibana" {
    length = 32
    special = false
}

resource "random_password" "remote_monitoring_user" {
    length = 32
    special = false
}

data "template_file" "inventory" {
    template = file("${path.module}/inventory_template.ini")

    vars = {
        nginx_ip = google_compute_address.nginx.address
        nginx_internal_ip = google_compute_instance.nginx.network_interface.0.network_ip
        kibana_ip = google_compute_instance.kibana.network_interface.0.network_ip
        elasticsearch1_ip = google_compute_instance.elasticsearch1.network_interface.0.network_ip
        elasticsearch2_ip = google_compute_instance.elasticsearch2.network_interface.0.network_ip
        elasticsearch3_ip = google_compute_instance.elasticsearch3.network_interface.0.network_ip
        elasticsearch4_ip = google_compute_instance.elasticsearch4.network_interface.0.network_ip
        ssh_user = var.ssh_user
        ssh_priv_key_file = var.ssh_priv_key_file
        bootstrap_password = random_password.bootstrap.result
        elastic_password = random_password.elastic.result
        kibana_password = random_password.kibana.result
        remote_monitoring_user_password = random_password.remote_monitoring_user.result
        fqdn = var.fqdn
    }
}

resource "null_resource" "inventory" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "echo '${data.template_file.inventory.rendered}' > ${var.inventory_destination}"
    }
}