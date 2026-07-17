# Code Review

## Terraform + Helm

### 1 - Hardcoded values
The Terraform code contains hardcoded values.

Although variables.tf defines namespace, environment, and image_tag, none of them are used.
kubeconfig path is hardcoded as well.

Fix: Use var.namespace, var.environment, var.image_tag and add var.kubeconfig_path to variables.tf and use it

### 2 Set block

The first value is empty, so Terraform cannot parse it.
prod is unquoted, so Terraform interprets it as an undefined identifier.

Fix: Add image.tag value and use variables from variables.tf

### 3 Mismatched labels
Helm Chart Had Mismatched labels:

- Deployment label: app: myapp
- Service selector: app: myapps
- Ingress backend: homeworks

Fix: Use templated labels

### 4 Required_providers
The Terraform configuration lacked a proper required_providers block, making provider versioning non‑deterministic.

Fix: Add required_providers block into providers.tf

### 5 CI pipeline
Finish the ci pipeline, focus on the following areas:
validation, test, build, push, deploy stages
Fix: Implement the necessary stages