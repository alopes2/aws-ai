resource "aws_apigatewayv2_api" "api" {
  name                       = "bedrock"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_stage" "stage" {
  api_id = aws_apigatewayv2_api.api.id
  name   = "live"
}

resource "aws_apigatewayv2_route" "example" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "$default"
}

resource "aws_apigatewayv2_integration" "example" {
  api_id           = aws_apigatewayv2_api.api.id
  integration_type = "AWS_PROXY"

  connection_type           = "INTERNET"
  content_handling_strategy = "CONVERT_TO_TEXT"
  description               = "Default Websocket route"
  integration_method        = "POST"
  integration_uri           = aws_lambda_function.bedrock.invoke_arn
  passthrough_behavior      = "WHEN_NO_MATCH"
}
