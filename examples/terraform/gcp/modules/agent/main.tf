resource "google_compute_instance_group_manager" "mg" {
  provider = google-beta
  project  = var.gcp_project
  name     = var.name

  base_instance_name = var.name
  zone               = var.instance_zone
  target_size        = var.instance_count

  wait_for_instances = true

  version {
    instance_template = var.instance_template_id
  }
}