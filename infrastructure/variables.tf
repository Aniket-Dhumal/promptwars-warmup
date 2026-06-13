variable "gcp_project_id" {
  description = "The Google Cloud Project ID"
  type        = string
}

variable "region" {
  description = "The GCP region to deploy resources in"
  type        = string
  default     = "us-central1"
}

variable "db_tier" {
  description = "The machine type for the Cloud SQL database instance"
  type        = string
  default     = "db-f1-micro"
}

variable "db_version" {
  description = "The version of PostgreSQL engine to run"
  type        = string
  default     = "POSTGRES_15"
}
