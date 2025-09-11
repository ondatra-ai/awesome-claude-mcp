variable "repositories" {
  description = "List of repository names to create"
  type        = list(string)
  default     = ["frontend", "backend", "mcp-service"]
}
