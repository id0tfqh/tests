provider "aws" {
  region  = "${var.region}"
  profile = "Test profile"
}

data "aws_ami" "latest" {
    owners = ["137112412989"]
    most_recent = true
    filter {
        name = "name"
        values = ["amzn2-ami-hvm-2.0.*-x86_64-gp2"]
    }
}

/* EC2 description */
resource "aws_instance" "d99df" {

  ami                         = "${data.aws_ami.latest.id}"
  instance_type               = "${var.instanceType}"
  key_name                    = "${var.sshKeyName}"
  monitoring                  = "true"
  associate_public_ip_address = "true"
  availability_zone           = "${var.aZone}"
  subnet_id                   = "${aws_subnet.bi-subnet.id}"
  vpc_security_group_ids      = ["${aws_security_group.target.id}",
    "${aws_security_group.control.id}",
    "${aws_security_group.public.id}"]
  count                       = "${var.instanceCount}"
  user_data                   = file("bootstrap.sh")
  
  root_block_device {
   volume_size            = "${var.volumeSize}"
    volume_type           = "gp2"
    encrypted             = "false"
    delete_on_termination = "true"
  }

  tags = {
    Name = "d99df"
    Project = "Test profile"
  }
}

/* Network description */
resource "aws_vpc" "bi-vpc" {
  cidr_block = "${var.vpcCidr}"
  tags = {
    Name = "d99df"
    Project = "Test profile"
  }
}

resource "aws_subnet" "bi-subnet" {
  vpc_id = "${aws_vpc.bi-vpc.id}"
  cidr_block = "${var.cidr}"
  availability_zone ="${var.aZone}"
  tags = {
    Name = "d99df"
    Project = "Test profile"
  }
}

/* Gateway description */
resource "aws_internet_gateway" "bi-igw" {
  vpc_id = "${aws_vpc.bi-vpc.id}"
  tags = {
    Name = "d99df"
    Project = "Test profile"
  }
}

resource "aws_route" "bi-route" {
  route_table_id = "${aws_vpc.bi-vpc.default_route_table_id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id = "${aws_internet_gateway.bi-igw.id}"
}

/* Security group description */
resource "aws_security_group" "target" {
  name = "target-access"
  description = "Allow HTTP inbound traffic"
  vpc_id = "${aws_vpc.bi-vpc.id}"
  tags = {
    Name = "d99df"
    Project = "Test profile"
  }
}

resource "aws_security_group" "control" {
  name = "ssh-access"
  description = "Allow SSH inbound traffic"
  vpc_id = "${aws_vpc.bi-vpc.id}"
  tags = {
    Name = "d99df"
    Project = "Test profile"
  }
}

resource "aws_security_group" "public" {
  name = "outgoing-traffic"
  description = "Allow all outbound traffic"
  vpc_id = "${aws_vpc.bi-vpc.id}"
  tags = {
    Name = "d99df"
    Project = "Test profile"
  }
}

/* Security group rule description */
resource "aws_security_group_rule" "http" {
    description = "Allow HTTP inbound traffic"
    type = "ingress"
    from_port = 8080
    to_port = 8080
    protocol = "tcp"
    security_group_id = "${aws_security_group.target.id}"
    cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "ssh" {
    description = "Allow SSH inbound traffic"
    type = "ingress"
    from_port = 22
    to_port = 22
    protocol = "tcp"
    security_group_id = "${aws_security_group.control.id}"
    cidr_blocks = ["${var.mainIP}"]
}

resource "aws_security_group_rule" "allow-outbound-all" {
    type = "egress"
    from_port = "0"
    to_port = "0"
    protocol = "-1"
    security_group_id = "${aws_security_group.public.id}"
    cidr_blocks = ["0.0.0.0/0"]
}

/* Regestry */
resource "aws_ecr_repository" "hello-world" {
  name                 = "test_profile"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "d99df"
    Project = "Test profile"
  }
}
