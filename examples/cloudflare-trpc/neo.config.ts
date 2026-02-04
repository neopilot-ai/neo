/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "cloudflare-trpc",
      home: "cloudflare",
      removal: input?.stage === "production" ? "retain" : "remove",
    };
  },
  async run() {
    const trpc = new neo.cloudflare.Worker("Trpc", {
      url: true,
      handler: "index.ts",
    });

    const client = new neo.cloudflare.Worker("Client", {
      url: true,
      link: [trpc],
      handler: "client.ts",
    });

    return {
      api: trpc.url,
      client: client.url,
    };
  },
});
