variable "name" {
  description = "ALB name (must be unique per account/region). Prefer including the environment, e.g., mcp-gde-alb-dev."
  type        = string
  nullable    = false
  validation {
    condition     = length(var.name) >= 3 && length(var.name) <= 32
    error_message = "ALB name must be 3–32 characters and unique per account/region."
  }
}
