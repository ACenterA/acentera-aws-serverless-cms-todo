package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/myplugin/gofaas"
)

func main() {
	lambda.Start(gofaas.CognitoClientDomains)
}
