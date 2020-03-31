resource "tls_private_key" "ca_key" {
    algorithm = "ECDSA"
}

resource "tls_self_signed_cert" "ca_cert" {
    key_algorithm = "ECDSA"
    private_key_pem = tls_private_key.ca_key.private_key_pem

    subject {
        common_name = "ScoreStack Root CA"
        organization = "ScoreStack"
    }

    validity_period_hours = 8760
    
    allowed_uses = [
        "cert_signing",
    ]

    is_ca_certificate = true
}

resource "null_resource" "ca_cert" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "mkdir -p ${var.certificate_destination}/ca && echo '${tls_self_signed_cert.ca_cert.cert_pem}' > ${var.certificate_destination}/ca/ca.crt"
    }
}

resource "tls_private_key" "elasticsearch1_key" {
    algorithm = "ECDSA"
}

resource "null_resource" "elasticsearch1_key" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "mkdir -p ${var.certificate_destination}/elasticsearch1 && echo '${tls_private_key.elasticsearch1_key.cert_pem}' > ${var.certificate_destination}/elasticsearch1/elasticsearch1.key"
    }
}

resource "tls_cert_request" "elasticsearch1_csr" {
    key_algorithm = "ECDSA"
    private_key_pem = tls_private_key.elasticsearch1_key.private_key_pem

    subject {
        common_name = "elasticsearch1"
        organization = "ScoreStack"
    }

    dns_names = ["elasticsearch1"]
    ip_addresses = [google_compute_instance.elasticsearch1.network_interface.0.network_ip]
}

resource "tls_locally_signed_cert" "elasticsearch1_cert" {
    cert_request_pem = tls_cert_request.elasticsearch1_csr.cert_request_pem
    ca_key_algorithm = "ECDSA"
    ca_private_key_pem = tls_private_key.ca_key.private_key_pem
    ca_cert_pem = tls_private_key.ca_cert.cert_pem
    validity_period_hours = 8760

    allowed_uses = [
        "server_auth",
        "client_auth",
    ]
}

resource "null_resource" "elasticsearch1_cert" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "mkdir -p ${var.certificate_destination}/elasticsearch1 && echo '${tls_private_key.elasticsearch1_key.cert_pem}' > ${var.certificate_destination}/elasticsearch1/elasticsearch1.crt"
    }
}

resource "tls_private_key" "elasticsearch2_key" {
    algorithm = "ECDSA"
}

resource "null_resource" "elasticsearch2_key" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "mkdir -p ${var.certificate_destination}/elasticsearch2 && echo '${tls_private_key.elasticsearch2_key.cert_pem}' > ${var.certificate_destination}/elasticsearch2/elasticsearch2.key"
    }
}

resource "tls_cert_request" "elasticsearch2_csr" {
    key_algorithm = "ECDSA"
    private_key_pem = tls_private_key.elasticsearch2_key.private_key_pem

    subject {
        common_name = "elasticsearch2"
        organization = "ScoreStack"
    }

    dns_names = ["elasticsearch2"]
    ip_addresses = [google_compute_instance.elasticsearch2.network_interface.0.network_ip]
}

resource "tls_locally_signed_cert" "elasticsearch2_cert" {
    cert_request_pem = tls_cert_request.elasticsearch2_csr.cert_request_pem
    ca_key_algorithm = "ECDSA"
    ca_private_key_pem = tls_private_key.ca_key.private_key_pem
    ca_cert_pem = tls_private_key.ca_cert.cert_pem
    validity_period_hours = 8760

    allowed_uses = [
        "server_auth",
        "client_auth",
    ]
}

resource "null_resource" "elasticsearch2_cert" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "mkdir -p ${var.certificate_destination}/elasticsearch2 && echo '${tls_private_key.elasticsearch2_key.cert_pem}' > ${var.certificate_destination}/elasticsearch2/elasticsearch2.crt"
    }
}

resource "tls_private_key" "elasticsearch3_key" {
    algorithm = "ECDSA"
}

resource "null_resource" "elasticsearch3_key" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "mkdir -p ${var.certificate_destination}/elasticsearch3 && echo '${tls_private_key.elasticsearch3_key.cert_pem}' > ${var.certificate_destination}/elasticsearch3/elasticsearch3.key"
    }
}

resource "tls_cert_request" "elasticsearch3_csr" {
    key_algorithm = "ECDSA"
    private_key_pem = tls_private_key.elasticsearch3_key.private_key_pem

    subject {
        common_name = "elasticsearch3"
        organization = "ScoreStack"
    }

    dns_names = ["elasticsearch3"]
    ip_addresses = [google_compute_instance.elasticsearch3.network_interface.0.network_ip]
}

resource "tls_locally_signed_cert" "elasticsearch3_cert" {
    cert_request_pem = tls_cert_request.elasticsearch3_csr.cert_request_pem
    ca_key_algorithm = "ECDSA"
    ca_private_key_pem = tls_private_key.ca_key.private_key_pem
    ca_cert_pem = tls_private_key.ca_cert.cert_pem
    validity_period_hours = 8760

    allowed_uses = [
        "server_auth",
        "client_auth",
    ]
}

resource "null_resource" "elasticsearch3_cert" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "mkdir -p ${var.certificate_destination}/elasticsearch3 && echo '${tls_private_key.elasticsearch3_key.cert_pem}' > ${var.certificate_destination}/elasticsearch3/elasticsearch3.crt"
    }
}

resource "tls_private_key" "elasticsearch4_key" {
    algorithm = "ECDSA"
}

resource "null_resource" "elasticsearch4_key" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "mkdir -p ${var.certificate_destination}/elasticsearch4 && echo '${tls_private_key.elasticsearch4_key.cert_pem}' > ${var.certificate_destination}/elasticsearch4/elasticsearch4.key"
    }
}

resource "tls_cert_request" "elasticsearch4_csr" {
    key_algorithm = "ECDSA"
    private_key_pem = tls_private_key.elasticsearch4_key.private_key_pem

    subject {
        common_name = "elasticsearch4"
        organization = "ScoreStack"
    }

    dns_names = ["elasticsearch4"]
    ip_addresses = [google_compute_instance.elasticsearch4.network_interface.0.network_ip]
}

resource "tls_locally_signed_cert" "elasticsearch4_cert" {
    cert_request_pem = tls_cert_request.elasticsearch4_csr.cert_request_pem
    ca_key_algorithm = "ECDSA"
    ca_private_key_pem = tls_private_key.ca_key.private_key_pem
    ca_cert_pem = tls_private_key.ca_cert.cert_pem
    validity_period_hours = 8760

    allowed_uses = [
        "server_auth",
        "client_auth",
    ]
}

resource "null_resource" "elasticsearch4_cert" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "mkdir -p ${var.certificate_destination}/elasticsearch4 && echo '${tls_private_key.elasticsearch4_key.cert_pem}' > ${var.certificate_destination}/elasticsearch4/elasticsearch4.crt"
    }
}

resource "tls_private_key" "logstash_key" {
    algorithm = "ECDSA"
}

resource "null_resource" "logstash_key" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "mkdir -p ${var.certificate_destination}/logstash && echo '${tls_private_key.logstash_key.cert_pem}' > ${var.certificate_destination}/logstash/logstash.key"
    }
}

resource "tls_cert_request" "logstash_csr" {
    key_algorithm = "ECDSA"
    private_key_pem = tls_private_key.logstash_key.private_key_pem

    subject {
        common_name = "logstash"
        organization = "ScoreStack"
    }

    dns_names = ["logstash"]
    ip_addresses = [google_compute_instance.logstash.network_interface.0.network_ip]
}

resource "tls_locally_signed_cert" "logstash_cert" {
    cert_request_pem = tls_cert_request.logstash_csr.cert_request_pem
    ca_key_algorithm = "ECDSA"
    ca_private_key_pem = tls_private_key.ca_key.private_key_pem
    ca_cert_pem = tls_private_key.ca_cert.cert_pem
    validity_period_hours = 8760

    allowed_uses = [
        "server_auth",
        "client_auth",
    ]
}

resource "null_resource" "logstash_cert" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "mkdir -p ${var.certificate_destination}/logstash && echo '${tls_private_key.logstash_key.cert_pem}' > ${var.certificate_destination}/logstash/logstash.crt"
    }
}

resource "tls_private_key" "kibana_key" {
    algorithm = "ECDSA"
}

resource "null_resource" "kibana_key" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "mkdir -p ${var.certificate_destination}/kibana && echo '${tls_private_key.kibana_key.cert_pem}' > ${var.certificate_destination}/kibana/kibana.key"
    }
}

resource "tls_cert_request" "kibana_csr" {
    key_algorithm = "ECDSA"
    private_key_pem = tls_private_key.kibana_key.private_key_pem

    subject {
        common_name = "kibana"
        organization = "ScoreStack"
    }

    dns_names = ["kibana"]
    ip_addresses = [google_compute_instance.kibana.network_interface.0.network_ip]
}

resource "tls_locally_signed_cert" "kibana_cert" {
    cert_request_pem = tls_cert_request.kibana_csr.cert_request_pem
    ca_key_algorithm = "ECDSA"
    ca_private_key_pem = tls_private_key.ca_key.private_key_pem
    ca_cert_pem = tls_private_key.ca_cert.cert_pem
    validity_period_hours = 8760

    allowed_uses = [
        "server_auth",
        "client_auth",
    ]
}

resource "null_resource" "kibana_cert" {
    triggers = {
        template_rendered = data.template_file.inventory.rendered
    }

    provisioner "local-exec" {
        command = "mkdir -p ${var.certificate_destination}/kibana && echo '${tls_private_key.kibana_key.cert_pem}' > ${var.certificate_destination}/kibana/kibana.crt"
    }
}