variable "node_type" {
  description = "Redis node instance type"
  type        = string
  default     = "cache.t3.micro"
}

variable "engine_version" {
  description = "Redis engine version"
  type        = string
  default     = "7.0"
}

variable "vpc_id" {
  description = "VPC ID for Redis SG"
  type        = string
}

variable "private_subnet_ids" {
  description = "Private subnet IDs for Redis subnet group"
  type        = list(string)
}

variable "allowed_sg_ids" {
  description = "List of security group IDs allowed to connect to Redis (port 6379)"
  type        = list(string)
}

variable "num_cache_clusters" {
  description = "Number of cache nodes (1 for single node)"
  type        = number
  default     = 1
}
