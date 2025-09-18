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

variable "name_prefix" {
  description = "Prefix to ensure Redis resource names are unique per environment (e.g., dev|staging|prod)"
  type        = string
}

variable "cluster_num_node_groups" {
  description = "Number of node groups (shards) for Redis cluster mode"
  type        = number
  default     = 1
}

variable "cluster_replicas_per_node_group" {
  description = "Number of replicas per node group for Redis cluster mode"
  type        = number
  default     = 1
}
