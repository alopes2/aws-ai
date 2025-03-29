data "aws_bedrock_foundation_model" "model" {
  model_id = "anthropic.claude-3-5-sonnet-20240620-v1:0"
}

data "aws_bedrock_foundation_model" "claude" {
  model_id = "anthropic.claude-3-7-sonnet-20250219-v1:0"
}

data "aws_bedrock_inference_profile" "claude" {
  inference_profile_id = "eu.anthropic.claude-3-7-sonnet-20250219-v1:0"
}
