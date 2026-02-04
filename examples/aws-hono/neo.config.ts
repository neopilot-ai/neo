/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-hono",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
      version: "3.10.13",
    };
  },
  async run() {
    const bucket = new neo.aws.Bucket("MyBucket");
    new neo.aws.Function("Hono", {
      url: true,
      link: [bucket],
      handler: "src/index.handler",
    });
  },
});
