resource "aws_apigatewayv2_api" "api" {
  name                       = "bedrock"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
  cors_configuration {
    allow_headers = ["*"]
    allow_methods = ["*"]
    allow_origins = ["*"]
  }
}

resource "aws_apigatewayv2_stage" "stage" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "live"
  auto_deploy = true
}

resource "aws_apigatewayv2_route" "example" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "$default"
  target    = "integrations/${aws_apigatewayv2_integration.default.id}"
}

resource "aws_apigatewayv2_route" "example" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "$connect"
  target    = "integrations/${aws_apigatewayv2_integration.connections.id}"
}

resource "aws_apigatewayv2_route" "example" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "$disconnect"
  target    = "integrations/${aws_apigatewayv2_integration.connections.id}"
}

resource "aws_apigatewayv2_integration" "default" {
  api_id           = aws_apigatewayv2_api.api.id
  integration_type = "AWS_PROXY"

  connection_type           = "INTERNET"
  content_handling_strategy = "CONVERT_TO_TEXT"
  description               = "Default Websocket route"
  integration_method        = "POST"
  integration_uri           = aws_lambda_function.bedrock.invoke_arn
  passthrough_behavior      = "WHEN_NO_MATCH"
}

resource "aws_apigatewayv2_integration" "connections" {
  api_id           = aws_apigatewayv2_api.api.id
  integration_type = "AWS_PROXY"

  connection_type           = "INTERNET"
  content_handling_strategy = "CONVERT_TO_TEXT"
  description               = "Connect and Disconnect Websocket route"
  integration_method        = "POST"
  integration_uri           = aws_lambda_function.connection.invoke_arn
  passthrough_behavior      = "WHEN_NO_MATCH"
}
