variable "accessKey" {
  type = "string"
  default = ""
  description = ""
}

variable "secretKey" {
  type = "string"
  default = ""
  description = ""
}

variable "region" {
    type = "string"
    default = "us-east-1"
    description = "US East (N. Virginia)"
}

variable "availabilityZone" {
    type = "string"
    default = "us-east-1a"
    description = "US East (N. Virginia)"
}

variable "ami" {
    type = "string"
    default = "ami-068663a3c619dd892"
    description = "Ubuntu Server 20.04 LTS (HVM)"
}

variable "sshKeyName" {
    type = "string"
    default = "t-light"
    description = "prod key"
}

variable "environment" {
  type = "string"
  default = "production"
}

variable "vpc_cidr" {
    type = "string"
    description = ""
    default = "192.168.100.0/26"
}

variable "cidr" {
    type = "string"
    description = ""
    default = "192.168.100.0/26"
}

variable "instanceType" {
    type = "string"
    default = "t3.small"
    description = ""
}

variable "officeIP" {
    type = "string"
    default = "29.155.215.5/32"
}

variable "monitorIP" {
    type = "string"
    default = "31.9.22.253/32"
}
variable "sshIPaddresses" {
    type = "list"
    default = [""]
}
