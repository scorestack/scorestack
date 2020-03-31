resource "google_compute_network" "internal_network" {
    name = "internal-network"
    description = "Network for traffic in between ScoreStack nodes."
}

resource "google_compute_address" "nginx" {
    name = "nginx"
    description = "Static public IP address for the Nginx server."
}

resource "google_compute_firewall" "default" {
    name = "default"
    description = "Allow ICMP traffic."

    network = google_compute_network.internal_network.self_link

    allow {
        protocol = "icmp"
    }

    target_tags = ["default"]
}

resource "google_compute_firewall" "ssh" {
    name = "ssh"
    description = "Allow SSH traffic from the SSH jump box."

    network = google_compute_network.internal_network.self_link

    allow {
        protocol = "tcp"
        ports = ["22"]
    }

    source_tags = ["ssh-jump"]
    target_tags = ["ssh"]
}

resource "google_compute_firewall" "elasticsearch" {
    name = "elasticsearch"
    description = "Allow API traffic to Elasticsearch nodes."

    network = google_compute_network.internal_network.self_link

    allow {
        protocol = "tcp"
        ports = ["9200"]
    }

    source_tags = ["default"]
    target_tags = ["elasticsearch"]
}

resource "google_compute_firewall" "elasticsearch-transport" {
    name = "elasticsearch-transport"
    description = "Allow transport traffic in between nodes in an Elasticsearch cluster."

    network = google_compute_network.internal_network.self_link

    allow {
        protocol = "tcp"
        ports = ["9300"]
    }

    source_tags = ["elasticsearch-transport"]
    target_tags = ["elasticsearch-transport"]
}

resource "google_compute_firewall" "kibana" {
    name = "kibana"
    description = "Allow HTTP traffic to the Kibana server from the Nginx server."

    network = google_compute_network.internal_network.self_link

    allow {
        protocol = "tcp"
        ports = ["5601"]
    }

    source_tags = ["www"]
    target_tags = ["kibana"]
}

resource "google_compute_firewall" "logstash" {
    name = "logstash"
    description = "Allow traffic to the Dynamicbeat listener on the Logstash server from the Nginx server."

    network = google_compute_network.internal_network.self_link

    allow {
        protocol = "tcp"
        ports = ["5454"]
    }

    source_tags = ["www"]
    target_tags = ["logstash"]
}

resource "google_compute_firewall" "ssh-jump" {
    name = "ssh-jump"
    description = "Allow SSH traffic to the SSH jump box from the public internet."

    network = google_compute_network.internal_network.self_link

    allow {
        protocol = "tcp"
        ports = ["22"]
    }

    source_ranges = ["0.0.0.0/0"]
    target_tags = ["ssh-jump"]
}

resource "google_compute_firewall" "www" {
    name = "www"
    description = "Allow HTTP and HTTPS traffic to the Nginx server from the public internet."

    network = google_compute_network.internal_network.self_link

    allow {
        protocol = "tcp"
        ports = ["80", "443"]
    }

    source_ranges = ["0.0.0.0/0"]
    target_tags = ["www"]
}