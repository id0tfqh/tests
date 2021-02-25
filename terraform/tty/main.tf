provider "aws" {
  region = "${var.region}"
  profile = "Light pro"
}

resource "aws_instance" "Application" {

  ami = "${var.ami}"
  instance_type = "${var.instanceType}"
  key_name = "${var.sshKeyName}"
  monitoring = "false"
  associate_public_ip_address = "true"
  availability_zone = "${var.availabilityZone}"

  root_block_device {
    volume_size = "16"
    volume_type = "gp2"
    encrypted   = "false"
    delete_on_termination = "true"
  }

  subnet_id = "${aws_subnet.main-subnet.id}"
  vpc_security_group_ids = ["${aws_security_group.application.id}", 
    "${aws_security_group.ssh-access.id}", 
    "${aws_security_group.monitoring-access.id}",
    "${aws_security_group.local-plain.id}"]

  tags = {
    environment = "${var.environment}"
    Name = "application"
  }

  volume_tags = {
    environment = "${var.environment}"
    Name = "application"
  }
}

resource "aws_instance" "mysql" {

  ami = "${var.ami}"
  instance_type = "${var.instanceType}"
  key_name = "${var.sshKeyName}"
  monitoring = "false"
  associate_public_ip_address = "true"
  availability_zone = "${var.availabilityZone}"

  root_block_device {
    volume_size = "16"
    volume_type = "gp2"
    encrypted   = "false"
    delete_on_termination = "true"
  }

  subnet_id = "${aws_subnet.main-subnet.id}"
  vpc_security_group_ids = ["${aws_security_group.database.id}", 
    "${aws_security_group.ssh-access.id}", 
    "${aws_security_group.monitoring-access.id}",
    "${aws_security_group.local-plain.id}"]

  tags = {
    environment = "${var.environment}"
    Name = "mysql"
  }

  volume_tags = {
    environment = "${var.environment}"
    Name = "mysql"
  }
}

/* Security group description */
resource "aws_security_group" "application" {
  name = "app-security"
  description = "Allow HTTP inbound traffic"
  vpc_id = "${aws_vpc.main-vpc.id}"
  tags = {
    Name = "Light pro"
    environment = "${var.environment}"
  }
}

resource "aws_security_group" "database" {
  name = "db-security"
  description = "Allow MySQL inbound traffic from office"
  vpc_id = "${aws_vpc.main-vpc.id}"
  tags = {
    Name = "Light pro"
    environment = "${var.environment}"
  }
}

resource "aws_security_group" "ssh-access" {
  name = "ssh-security"
  description = "Allow SSH inbound traffic"
  vpc_id = "${aws_vpc.main-vpc.id}"
  tags = {
    Name = "Light pro"
    environment = "${var.environment}"
  }
}

resource "aws_security_group" "monitoring-access" {
  name = "monitoring-security"
  description = "Allow for monitoring inbound traffic"
  vpc_id = "${aws_vpc.main-vpc.id}"
  tags = {
    Name = "Light pro"
    environment = "${var.environment}"
  }
}

resource "aws_security_group" "local-plain" {
  name = "local-security"
  description = "Allow all inbound traffic Ethernet"
  vpc_id = "${aws_vpc.main-vpc.id}"
  tags = {
    Name = "Light pro"
    environment = "${var.environment}"
  }
}

/* Security group rule description */
resource "aws_security_group_rule" "allow-all-local" {
    type = "ingress"
    description = "Allow all inbound traffic from local network"
    from_port = 0
    to_port = 0
    protocol = "-1"
    security_group_id = "${aws_security_group.local-plain.id}"
    cidr_blocks = ["${var.cidr}"]
}

resource "aws_security_group_rule" "allow-http" {
    type = "ingress"
    description = "Allow HTTP inbound traffic from Internet"
    from_port = 80
    to_port = 80
    protocol = "tcp"
    security_group_id = "${aws_security_group.application.id}"
    cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "allow-https" {
    type = "ingress"
    description = "Allow HTTPS inbound traffic from Internet"
    from_port = 443
    to_port = 443
    protocol = "tcp"
    security_group_id = "${aws_security_group.application.id}"
    cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "allow-mysql-office" {
    type = "ingress"
    description = "Allow MySQL inbound traffic from office"
    from_port = 3306
    to_port = 3306
    protocol = "tcp"
    security_group_id = "${aws_security_group.database.id}"
    cidr_blocks = ["${var.officeIP}"]
}

resource "aws_security_group_rule" "allow-monitoring" {
    type = "ingress"
    description = "Allow Grafana inbound traffic from office"
    from_port =9100 
    to_port = 9100
    protocol = "tcp"
    security_group_id = "${aws_security_group.monitoring-access.id}"
    cidr_blocks = ["${var.monitorIP}"]
}

resource "aws_security_group_rule" "allow-ssh" {
    type = "ingress"
    description = "Allow SSH inbound traffic from office and runners"
    from_port = 22
    to_port = 22
    protocol = "tcp"
    security_group_id = "${aws_security_group.ssh-access.id}"
    cidr_blocks = "${var.sshIPaddresses}"
}

resource "aws_security_group_rule" "allow_all_outbound" {
    description = "Allow all outbound traffic"
    type = "egress"
    from_port = 0
    to_port = 0
    protocol = "-1"
    security_group_id = "${aws_security_group.local-plain.id}"
    cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "allow_http_outbound" {
    description = "Allow tcp outbound traffic from application"
    type = "egress"
    from_port = 0
    to_port = 65535
    protocol = "tcp"
    security_group_id = "${aws_security_group.application.id}"
    cidr_blocks = ["0.0.0.0/0"]
}

/* Network description */
resource "aws_vpc" "main-vpc" {
  cidr_block = "${var.vpc_cidr}"
  tags = {
    Name = "Light pro"
    environment = "${var.environment}"
  }
}

resource "aws_subnet" "main-subnet" {
  vpc_id = "${aws_vpc.main-vpc.id}"
  cidr_block = "${var.cidr}"
  availability_zone ="${var.availabilityZone}"
  tags = {
    Name = "Light pro"
    environment = "${var.environment}"
  }
}

/* Gateway description */
resource "aws_internet_gateway" "igw-mirror" {
  vpc_id = "${aws_vpc.main-vpc.id}"
  tags = {
    Name = "Light pro"
    environment = "${var.environment}"
  }
}

resource "aws_route" "route" {
  route_table_id = "${aws_vpc.main-vpc.default_route_table_id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id = "${aws_internet_gateway.igw-mirror.id}"
}

/* S3 bucket description */
resource "aws_s3_bucket" "media-files" {
  bucket = "mirror-${var.environment}"
  acl    = "private"
  versioning {
      enabled = true
  }
  tags = {
    Name = "Light pro"
    environment = "${var.environment}"
  }
}
