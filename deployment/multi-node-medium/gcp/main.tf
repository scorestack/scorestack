provider "google" {
  version     = "3.5.0"
  credentials = file(var.credentials_file)
  project     = var.project
  region      = var.region
  zone        = var.zone
}