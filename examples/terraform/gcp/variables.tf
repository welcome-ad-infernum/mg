variable "instance_count" {
  type        = string
  description = "Compute Engine instance count"
}
variable "gcp_project" {
  type        = string
  description = "Google Cloud Platform project name"
}
variable "instance_zones" {
  type        = set(string)
  description = "List of compute zones"
}
variable "instance_type" {
  type        = string
  description = "Compute Engine instance type"
}
variable "name" {
  type        = string
  description = "Deployment name"
}