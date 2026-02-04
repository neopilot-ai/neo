package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/neopilot-ai/neo/v3/sdk/golang/resource"
)

func handler() (string, error) {
	bucket, err := resource.Get("MyBucket", "name")
	if err != nil {
		return "", err
	}
	return bucket.(string), nil
}

func main() {
	lambda.Start(handler)
}
