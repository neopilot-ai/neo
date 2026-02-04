/// <reference path="./.neo/platform/config.d.ts" />
export default $config({
  app(input) {
    return {
      name: "scrap",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
      version: "3.10.13",
      providers: { command: "1.0.2" },
    };
  },
  async run() {
    const bucket = new neo.aws.Bucket("MyBucket");
    new neo.aws.Function("Hono", {
      url: true,
      link: [bucket],
      handler: "src/index.handler",
    });
    // new neo.aws.Router("Router");
  },
});
