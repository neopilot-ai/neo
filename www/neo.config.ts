/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "www",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
      version: "3.13.20",
    };
  },
  console: {
    autodeploy: {
      target(event) {
        if (
          event.type === "branch" &&
          event.branch === "dev" &&
          event.action === "pushed"
        ) {
          return { stage: "production" };
        }
      },
      async workflow({ $, event }) {
        await $`bun i`;
        await $`goenv install 1.21.3 && goenv global 1.21.3`;
        await $`cd ../platform && ./scripts/build`;
        await $`bun i neo-linux-x64`;
        event.action === "removed"
          ? await $`bun neo remove`
          : await $`bun neo deploy`;
      },
    },
  },
  async run() {
    const domain =
      {
        production: "neo.dev",
        dev: "dev.neo.dev",
      }[$app.stage] || $app.stage + "dev.neo.dev";

    // Redirect /examples to guide.neo.dev/examples
    // Redirect /chapters to guide.neo.dev/chapters
    // Redirect /archives to guide.neo.dev/archives
    const redirectToGuideBehavior = {
      targetOriginId: "redirect",
      viewerProtocolPolicy: "redirect-to-https",
      allowedMethods: ["GET", "HEAD", "OPTIONS"],
      cachedMethods: ["GET", "HEAD"],
      functionAssociations: [
        {
          eventType: "viewer-request",
          functionArn: new aws.cloudfront.Function("AstroRedirect", {
            runtime: "cloudfront-js-2.0",
            code: [
              `async function handler(event) {`,
              `  const request = event.request;`,
              // ie. request.uri is /examples/foo
              `  return {`,
              `    statusCode: 302,`,
              `    statusDescription: 'Found',`,
              `    headers: {`,
              `      location: { value: "https://guide.neo.dev" + request.uri }`,
              `    },`,
              `  };`,
              `}`,
            ].join("\n"),
          }).arn,
        },
      ],
      forwardedValues: {
        queryString: true,
        headers: ["Origin"],
        cookies: { forward: "none" },
      },
    };

    // Redirect /u/* to api.console.neo.dev/link/*
    const redirectToConsoleBehavior = {
      targetOriginId: "redirect",
      viewerProtocolPolicy: "redirect-to-https",
      allowedMethods: ["GET", "HEAD", "OPTIONS"],
      cachedMethods: ["GET", "HEAD"],
      functionAssociations: [
        {
          eventType: "viewer-request",
          functionArn: new aws.cloudfront.Function("ConsoleRedirect", {
            runtime: "cloudfront-js-2.0",
            code: [
              `async function handler(event) {`,
              `  const request = event.request;`,
              // ie. request.uri is /u/123
              `  return {`,
              `    statusCode: 302,`,
              `    statusDescription: 'Found',`,
              `    headers: {`,
              `      location: { value: "https://api.console.neo.dev/link" + request.uri }`,
              `    },`,
              `  };`,
              `}`,
            ].join("\n"),
          }).arn,
        },
      ],
      forwardedValues: {
        queryString: true,
        headers: ["Origin"],
        cookies: { forward: "none" },
      },
    };

    // Redirect /install to https://raw.githubusercontent.com/neo/neo/dev/install
    const redirectToInstallBehavior = {
      targetOriginId: "redirect",
      viewerProtocolPolicy: "redirect-to-https",
      allowedMethods: ["GET", "HEAD", "OPTIONS"],
      cachedMethods: ["GET", "HEAD"],
      functionAssociations: [
        {
          eventType: "viewer-request",
          functionArn: new aws.cloudfront.Function("InstallRedirect", {
            runtime: "cloudfront-js-2.0",
            code: [
              `async function handler(event) {`,
              `  const request = event.request;`,
              `  return {`,
              `    statusCode: 302,`,
              `    statusDescription: 'Found',`,
              `    headers: {`,
              `      location: { value: "https://raw.githubusercontent.com/neo/neo/dev/install" }`,
              `    },`,
              `  };`,
              `}`,
            ].join("\n"),
          }).arn,
        },
      ],
      forwardedValues: {
        queryString: true,
        headers: ["Origin"],
        cookies: { forward: "none" },
      },
    };

    // Strip .html from /blog
    const stripHtmlBehavior = {
      targetOriginId: "redirect",
      viewerProtocolPolicy: "redirect-to-https",
      allowedMethods: ["GET", "HEAD", "OPTIONS"],
      cachedMethods: ["GET", "HEAD"],
      functionAssociations: [
        {
          eventType: "viewer-request",
          functionArn: new aws.cloudfront.Function("StripHtml", {
            runtime: "cloudfront-js-2.0",
            code: [
              `async function handler(event) {`,
              `  return {`,
              `    statusCode: 308,`,
              `    headers: {`,
              `      location: { value: event.request.uri.replace(/\.html$/, "") }`,
              `    },`,
              `  };`,
              `}`,
            ].join("\n"),
          }).arn,
        },
      ],
      forwardedValues: {
        queryString: true,
        headers: ["Origin"],
        cookies: { forward: "none" },
      },
    };

    new neo.aws.Astro("Astro", {
      domain:
        $app.stage === "production"
          ? {
              name: domain,
              redirects: [
                "www.neo.dev",
                "ion.neo.dev",
                "serverless-stack.com",
                "www.serverless-stack.com",
              ],
            }
          : domain,
      transform: {
        cdn: (args) => {
          args.origins = $output(args.origins).apply((origins) => [
            ...origins,
            {
              domainName: "guide.neo.dev",
              originId: "redirect",
              customOriginConfig: {
                httpPort: 80,
                httpsPort: 443,
                originProtocolPolicy: "https-only",
                originReadTimeout: 20,
                originSslProtocols: ["TLSv1.2"],
              },
            },
          ]);
          args.orderedCacheBehaviors = $output(
            args.orderedCacheBehaviors,
          ).apply((cacheBehaviors) => [
            ...(cacheBehaviors || []),
            { pathPattern: "/blog/*.html", ...stripHtmlBehavior },
            { pathPattern: "/install", ...redirectToInstallBehavior },
            { pathPattern: "/examples*", ...redirectToGuideBehavior },
            { pathPattern: "/chapters*", ...redirectToGuideBehavior },
            { pathPattern: "/archives*", ...redirectToGuideBehavior },
            { pathPattern: "/u/*", ...redirectToConsoleBehavior },
          ]);
        },
      },
    });

    // Redirect docs.neo.dev to neo.dev/docs
    if ($app.stage === "production") {
      new neo.aws.Router("DocsRouter", {
        domain: {
          name: "docs.neo.dev",
          aliases: ["docs.serverless-stack.com"],
        },
        routes: {
          "/*": {
            url: `https://neo.dev/docs`,
            edge: {
              viewerRequest: {
                injection: `
return {
  statusCode: 301,
  statusDescription: 'Moved Permanently',
  headers: {
    location: { value: "https://neo.dev/docs" }
  }
};
              `,
              },
            },
          },
        },
      });
    }

    // Redirect telemetry.ion.neo.dev to us.i.posthog.com
    new neo.aws.Router("TelemetryRouter", {
      domain: {
        name: "telemetry.ion." + domain,
      },
      routes: {
        "/*": "https://us.i.posthog.com",
      },
    });
  },
});
