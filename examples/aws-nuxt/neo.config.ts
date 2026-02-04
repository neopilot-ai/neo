/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-nuxt",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const bucket = new neo.aws.Bucket("MyBucket", {
      access: "public",
    });
    new neo.aws.Nuxt("MyWeb", {
      link: [bucket],
    });
  },
});
