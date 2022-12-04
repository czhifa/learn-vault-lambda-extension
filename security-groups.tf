resource "aws_default_vpc" "default" {
}

resource "aws_security_group" "vault-server" {
  name        = "${var.environment_name}-vault-server-sg"
  description = "SSH and Internal Traffic"
  vpc_id      = aws_default_vpc.default.id

  tags = {
    Name = var.environment_name
  }

  # SSH
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Vault API traffic
  ingress {
    from_port   = 8200
    to_port     = 8200
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Vault cluster traffic
  ingress {
    from_port   = 8201
    to_port     = 8201
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Internal Traffic
  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    self      = true
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "rds" {
  name        = "${var.environment_name}-rds-sg"
  description = "MySQL traffic"
  vpc_id      = aws_default_vpc.default.id

  tags = {
    Name = var.environment_name
  }

  # MySQL traffic
  ingress {
    from_port   = 3309
    to_port     = 3309
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

