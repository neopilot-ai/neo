/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "cloudflare-vite",
      removal: input?.stage === "production" ? "retain" : "remove",
      providers: {
        aws: {},
        cloudflare: {},
      },
      home: "aws",
    };
  },
  async run() {
    new neo.cloudflare.StaticSite("Web", {
      build: {
        command: "pnpm run build",
        output: "dist",
      },
      domain: "vite.neoion.com",
    });
  },
});
