# Code Review

## Terraform

The current implementation relies on hardcoded values, e.g. namespace,environment or image tag. Variables already defined in variables.tf, but not in use. Hard coded values reduce flexibility and generate extra overhead when we maintanance our code.
The recommendation is that use these variables, like var.namespace