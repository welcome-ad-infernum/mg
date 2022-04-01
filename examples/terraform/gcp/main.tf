provider "google" {
  project = var.gcp_project
}

resource "google_service_account" "mg" {
  project      = var.gcp_project
  account_id   = "compute-sa"
  display_name = "SA for Google Compute Engine"
}

resource "google_compute_instance_template" "mg" {
  provider     = google-beta
  project      = var.gcp_project
  name         = "${var.name}-template"
  machine_type = var.instance_type
  tags         = ["default-allow-ssh"]

  metadata = {
    ssh-keys = "user:${file("~/.ssh/id_rsa.pub")}"
  }

  metadata_startup_script = <<SCRIPT
    sleep 10
    docker run --detach --restart always --pull always --net host --name mg-agent vladstarr/mg-agent:latest
  SCRIPT

  disk {
    source_image = "cos-cloud/cos-stable"
    auto_delete  = true
    boot         = true
  }

  service_account {
    email  = google_service_account.mg.email
    scopes = ["cloud-platform"]
  }

  network_interface {
    network = "default"
    access_config {}
  }

  scheduling {
    preemptible        = true
    automatic_restart  = false
    provisioning_model = "SPOT"
  }
}

module "agent" {
  source = "./modules/agent"

  for_each = var.instance_zones

  gcp_project          = var.gcp_project
  name                 = var.name
  instance_count       = var.instance_count
  instance_zone        = each.key
  instance_template_id = google_compute_instance_template.mg.id
}