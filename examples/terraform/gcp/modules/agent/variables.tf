variable "instance_count" {
  type        = string
  description = "Compute Engine instance count"
}
variable "instance_zone" {
  type        = string
  description = "Compute Engine instance zone"
}
variable "instance_template_id" {
  type        = string
  description = "Compute Engine instance template id"
}
variable "gcp_project" {
  type        = string
  description = "Google Cloud Platform project name"
}
variable "name" {
  type        = string
  description = "Deployment name"
}