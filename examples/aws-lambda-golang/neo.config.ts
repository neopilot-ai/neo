/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## AWS Lambda Go
 *
 * This example shows how to use the [`go`](https://golang.org/) runtime in your Lambda
 * functions.
 *
 * Our Go function is in the `src` directory and we point to it in our function.
 *
 * ```ts title="neo.config.ts" {5}
 * new neo.aws.Function("MyFunction", {
 *   url: true,
 *   runtime: "go",
 *   link: [bucket],
 *   handler: "./src",
 * });
 * ```
 *
 * We are also linking it to an S3 bucket. We can reference the bucket in our function.
 *
 * ```go title="src/main.go" {2}
 * func handler() (string, error) {
 *   bucket, err := resource.Get("MyBucket", "name")
 *   if err != nil {
 *     return "", err
 *   }
 *   return bucket.(string), nil
 * }
 * ```
 *
 * The `resource.Get` function is from the NEO Go SDK.
 *
 * ```go title="src/main.go" {2}
 * import (
 *   "github.com/neopilot-ai/neo/v3/sdk/golang/resource"
 * )
 * ```
 *
 * The `neo dev` CLI also supports running your Go function [_Live_](/docs/live).
 */
export default $config({
  app(input) {
    return {
      name: "aws-lambda-golang",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const bucket = new neo.aws.Bucket("MyBucket");

    new neo.aws.Function("MyFunction", {
      url: true,
      runtime: "go",
      link: [bucket],
      handler: "./src",
    });
  },
});
