/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-bundle",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    new neo.aws.Function("Function", {
      bundle: "./src",
      handler: "index.handler",
      url: true,
    });
  },
});
