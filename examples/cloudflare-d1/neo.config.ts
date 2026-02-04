/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "cloudflare-d1",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "cloudflare",
    };
  },
  async run() {
    const db = new neo.cloudflare.D1("MyDatabase");
    const worker = new neo.cloudflare.Worker("Worker", {
      link: [db],
      url: true,
      handler: "index.ts",
    });

    return {
      url: worker.url,
    };
  },
});
