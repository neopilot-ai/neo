/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-react-router",
      removal: input?.stage === "production" ? "retain" : "remove",
      protect: ["production"].includes(input?.stage),
      home: "aws",
    };
  },
  async run() {
    const bucket = new neo.aws.Bucket("MyBucket", {
      access: "public",
    });

    new neo.aws.React("MyWeb", {
      link: [bucket],
    });
  },
});
