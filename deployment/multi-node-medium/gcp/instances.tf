resource "google_compute_instance" "elasticsearch1" {
    name = "elasticsearch1"
    description = "An Elasticsearch master-eligible data node."

    boot_disk {
        initialize_params {
            size = 20
            image = var.ubuntu
        }
    }

    machine_type = "n1-standard-2"
    zone = var.zone1

    network_interface {
        network = google_compute_network.internal_network.self_link
        access_config {}
    }

    metadata = {
        ssh-keys = "${var.ssh_user}:${file(var.ssh_pub_key_file)}"
    }

    tags = ["default", "ssh", "elasticsearch", "elasticsearch-transport"]
}

resource "google_compute_instance" "elasticsearch2" {
    name = "elasticsearch2"
    description = "An Elasticsearch master-eligible data node."

    boot_disk {
        initialize_params {
            size = 20
            image = var.ubuntu
        }
    }

    machine_type = "n1-standard-2"
    zone = var.zone2

    network_interface {
        network = google_compute_network.internal_network.self_link
        access_config {}
    }

    metadata = {
        ssh-keys = "${var.ssh_user}:${file(var.ssh_pub_key_file)}"
    }

    tags = ["default", "ssh", "elasticsearch", "elasticsearch-transport"]
}

resource "google_compute_instance" "elasticsearch3" {
    name = "elasticsearch3"
    description = "An Elasticsearch master-eligible data node."

    boot_disk {
        initialize_params {
            size = 20
            image = var.ubuntu
        }
    }

    machine_type = "n1-standard-2"
    zone = var.zone3

    network_interface {
        network = google_compute_network.internal_network.self_link
        access_config {}
    }

    metadata = {
        ssh-keys = "${var.ssh_user}:${file(var.ssh_pub_key_file)}"
    }

    tags = ["default", "ssh", "elasticsearch", "elasticsearch-transport"]
}

resource "google_compute_instance" "elasticsearch4" {
    name = "elasticsearch4"
    description = "An Elasticsearch coordinating-only node."

    boot_disk {
        initialize_params {
            size = 20
            image = var.ubuntu
        }
    }

    machine_type = "n1-standard-1"
    zone = var.zone1

    network_interface {
        network = google_compute_network.internal_network.self_link
        access_config {}
    }

    metadata = {
        ssh-keys = "${var.ssh_user}:${file(var.ssh_pub_key_file)}"
    }

    tags = ["default", "ssh", "elasticsearch", "elasticsearch-transport"]
}

resource "google_compute_instance" "kibana" {
    name = "kibana"
    description = "A Kibana node."

    boot_disk {
        initialize_params {
            image = var.ubuntu
        }
    }

    machine_type = "n1-standard-2"
    zone = var.zone1

    network_interface {
        network = google_compute_network.internal_network.self_link
        access_config {}
    }

    metadata = {
        ssh-keys = "${var.ssh_user}:${file(var.ssh_pub_key_file)}"
    }

    tags = ["default", "ssh", "kibana"]
}

resource "google_compute_instance" "logstash" {
    name = "logstash"
    description = "A Logstash node."

    boot_disk {
        initialize_params {
            image = var.ubuntu
        }
    }

    machine_type = "n1-standard-1"
    zone = var.zone1

    network_interface {
        network = google_compute_network.internal_network.self_link
        access_config {}
    }

    metadata = {
        ssh-keys = "${var.ssh_user}:${file(var.ssh_pub_key_file)}"
    }

    tags = ["default", "ssh", "logstash"]
}

resource "google_compute_instance" "nginx" {
    name = "nginx"
    description = "An Nginx node that also serves as an SSH jump box."

    boot_disk {
        initialize_params {
            image = var.ubuntu
        }
    }

    machine_type = "g1-small"
    zone = var.zone1

    network_interface {
        network = google_compute_network.internal_network.self_link
        access_config {
            nat_ip = google_compute_address.nginx.address
        }
    }

    metadata = {
        ssh-keys = "${var.ssh_user}:${file(var.ssh_pub_key_file)}"
    }

    tags = ["default", "ssh-jump", "proxy", "logstash", "elasticsearch"]
}