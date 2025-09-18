variable "name_prefix" {
  description = "Base name prefix for VPC resources"
  type        = string
  default     = "awesome-claude-mcp"
}

variable "cidr_block" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}
