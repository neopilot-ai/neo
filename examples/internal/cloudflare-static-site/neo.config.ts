/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "cloudflare-static-site",
      removal: input?.stage === "production" ? "retain" : "remove",
      providers: {
        cloudflare: true,
      },
      home: "aws",
    };
  },
  async run() {
    new neo.cloudflare.StaticSite("MySite", {
      domain: "static.neoion.com",
    });
  },
});
