resource "aws_dynamodb_table" "connections" {
  name           = "connections"
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "ConnectionID"

  attribute {
    name = "ConnectionID"
    type = "S"
  }
}
