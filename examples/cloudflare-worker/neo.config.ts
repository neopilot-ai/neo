/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "cloudflare-worker",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "cloudflare",
    };
  },
  async run() {
    const bucket = new neo.cloudflare.Bucket("MyBucket");
    const worker = new neo.cloudflare.Worker("MyWorker", {
      handler: "./index.ts",
      link: [bucket],
      url: true,
    });

    return {
      api: worker.url,
    };
  },
});
