/// <reference path="./.neo/platform/config.d.ts" />
/**
 * ## Vercel domains
 *
 * Creates a router that uses domains purchased through and hosted in your Vercel account.
 * Ensure the `VERCEL_API_TOKEN` and `VERCEL_TEAM_ID` environment variables are set.
 */
export default $config({
  app(input) {
    return {
      name: "vercel-domain",
      home: "aws",
      removal: input?.stage === "production" ? "retain" : "remove",
      providers: {
        aws: true,
        "@pulumiverse/vercel": true,
      },
    };
  },
  async run() {
    const router = new neo.aws.Router("MyRouter", {
      domain: {
        name: "ion.neo.moe",
        dns: neo.vercel.dns({ domain: "neo.moe" }),
      },
      routes: {
        "/*": "https://neo.khulnasoft.com",
      },
    });
    return {
      router: router.url,
    };
  },
});
