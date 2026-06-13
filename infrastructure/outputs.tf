output "kms_key_ring_id" {
  description = "The fully qualified ID of the Cloud KMS KeyRing"
  value       = google_kms_key_ring.app_keyring.id
}

output "kms_crypto_key_id" {
  description = "The fully qualified ID of the Cloud KMS Biometric CryptoKey"
  value       = google_kms_crypto_key.biometric_key.id
}

output "db_instance_name" {
  description = "The name of the database instance"
  value       = google_sql_database_instance.postgres_instance.name
}

output "db_instance_connection_name" {
  description = "The connection name of the database instance (used by Cloud SQL Proxy)"
  value       = google_sql_database_instance.postgres_instance.connection_name
}

output "db_instance_public_ip" {
  description = "The public IPv4 address assigned to the Cloud SQL instance"
  value       = google_sql_database_instance.postgres_instance.public_ip_address
}
