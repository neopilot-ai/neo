/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "cloudflare-hono",
      home: "cloudflare",
      removal: input?.stage === "production" ? "retain" : "remove",
    };
  },
  async run() {
    const bucket = new neo.cloudflare.Bucket("MyBucket");
    const hono = new neo.cloudflare.Worker("Hono", {
      url: true,
      link: [bucket],
      handler: "index.ts",
    });

    return {
      api: hono.url,
    };
  },
});
