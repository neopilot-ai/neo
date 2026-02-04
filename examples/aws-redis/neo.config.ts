/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-redis",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    // NAT Gateways are required for Lambda functions
    const vpc = new neo.aws.Vpc("MyVpc", { nat: "managed" });
    const redis = new neo.aws.Redis("MyRedis", { vpc });
    new neo.aws.Function("MyApp", {
      handler: "index.handler",
      url: true,
      vpc,
      link: [redis],
    });
  },
});
