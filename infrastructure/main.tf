provider "google" {
  project = var.gcp_project_id
  region  = var.region
}

# Cryptographic Key Storage Layer (Cloud KMS CMEK)
resource "google_kms_key_ring" "app_keyring" {
  name     = "digital-twin-keyring-mumbai"
  location = "global"
}

resource "google_kms_crypto_key" "biometric_key" {
  name            = "metabolic-secure-key"
  key_ring        = google_kms_key_ring.app_keyring.id
  rotation_period = "7776000s"

  lifecycle {
    prevent_destroy = true
  }
}

# Secured PostgreSQL Instance Boundary with CMEK Encryption
resource "google_sql_database_instance" "postgres_instance" {
  name                = "twin-secure-ledger"
  database_version    = var.db_version
  region              = var.region
  encryption_key_name = google_kms_crypto_key.biometric_key.id

  settings {
    tier = var.db_tier

    ip_configuration {
      ipv4_enabled = true
      require_ssl  = true
    }
  }
}
