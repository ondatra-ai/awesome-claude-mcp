variable "repositories" {
  description = "List of repository names to create"
  type        = list(string)
  default     = ["frontend", "backend", "mcp-service"]
  validation {
    condition = alltrue([
      for r in var.repositories : (
        length(r) >= 2 && length(r) <= 256 &&
        can(regex("^(?:[a-z0-9]+(?:[._-][a-z0-9]+)*/)*[a-z0-9]+(?:[._-][a-z0-9]+)*$", r))
      )
    ])
    error_message = "Each ECR repository name must be 2–256 chars and match AWS ECR naming (lowercase letters, digits, '.', '_', '-', optional '/' namespaces)."
  }
}
