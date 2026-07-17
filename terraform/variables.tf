variable "namespace" {
  type        = string
  description = "Kubernetes namespace where the application will be deployed"
  default     = "production"
}

variable "environment" {
  type        = string
  description = "Environment name"
  default     = "prod"
}

variable "image_tag" {
  type        = string
  description = "Docker image tag to deploy (e.g. CI_COMMIT_SHA)"
}

variable "kubeconfig_path" {
  type        = string
  description = "Path to kubeconfig file"
}