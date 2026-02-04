/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## Router and function URL
 *
 * Creates a router that routes all requests to a function with a URL.
 */
export default $config({
  app(input) {
    return {
      name: "aws-router",
      home: "aws",
      removal: input?.stage === "production" ? "retain" : "remove",
    };
  },
  async run() {
    const api = new neo.aws.Function("MyApi", {
      handler: "api.handler",
      url: true,
    });
    const bucket = new neo.aws.Bucket("MyBucket", {
      access: "public",
    });
    const router = new neo.aws.Router("MyRouter", {
      domain: "router.ion.dev.neo.dev",
      routes: {
        "/api/*": api.url,
        "/*": $interpolate`https://${bucket.domain}`,
      },
    });

    return {
      router: router.url,
      bucket: bucket.domain,
    };
  },
});
