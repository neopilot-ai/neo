/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-api",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const bucket = new neo.aws.Bucket("MyBucket");
    const api = new neo.aws.ApiGatewayV2("MyApi");
    api.route("GET /", {
      link: [bucket],
      handler: "index.upload",
    });
    api.route("GET /latest", {
      link: [bucket],
      handler: "index.latest",
    });
  },
});
