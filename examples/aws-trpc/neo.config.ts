/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-trpc",
      home: "aws",
      removal: input?.stage === "production" ? "retain" : "remove",
    };
  },
  async run() {
    const trpc = new neo.aws.Function("Trpc", {
      url: true,
      handler: "index.handler",
    });

    const client = new neo.aws.Function("Client", {
      url: true,
      link: [trpc],
      handler: "client.handler",
    });

    return {
      api: trpc.url,
      client: client.url,
    };
  },
});
