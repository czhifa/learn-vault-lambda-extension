resource "random_password" "password" {
  length  = 32
  special = false
}

resource "aws_db_instance" "main" {
  allocated_storage      = 20
  storage_type           = "gp2"
  engine                 = "mysql"
  engine_version         = "5.7"
  instance_class         = var.db_instance_type
  db_name                = "mysql"
  username               = "vaultadmin"
  password               = "vaultpass"
  vpc_security_group_ids = [aws_security_group.rds.id]
  skip_final_snapshot    = true
  publicly_accessible    = true
}
