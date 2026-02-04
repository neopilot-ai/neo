/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-auth-nextjs",
      removal: input?.stage === "production" ? "retain" : "remove",
      protect: ["production"].includes(input?.stage),
      home: "aws",
    };
  },
  async run() {
    const auth = new neo.aws.Auth("MyAuth", {
      issuer: "auth/index.handler",
    });

    new neo.aws.Nextjs("MyWeb", {
      link: [auth],
    });
  },
});
