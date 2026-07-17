output "namespace" {
  description = "Namespace where the application is deployed"
  value       = kubernetes_namespace_v1.homework.metadata[0].name
}

output "helm_release_name" {
  description = "Name of the Helm release"
  value       = helm_release.homework.name
}
