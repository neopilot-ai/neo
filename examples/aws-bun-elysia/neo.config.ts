/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## AWS Bun Elysia container
 *
 * Deploys a Bun [Elysia](https://elysiajs.com/) API to AWS.
 *
 * You can get started by running.
 *
 * ```bash
 * bun create elysia aws-bun-elysia
 * cd aws-bun-elysia
 * bunx neo init
 * ```
 *
 * Now you can add a service.
 *
 * ```ts title="neo.config.ts"
 * new neo.aws.Service("MyService", {
 *   cluster,
 *   loadBalancer: {
 *     ports: [{ listen: "80/http", forward: "3000/http" }],
 *   },
 *   dev: {
 *     command: "bun dev",
 *   },
 * });
 * ```
 *
 * Start your app locally.
 *
 * ```bash
 * bun neo dev
 * ```
 *
 * This example lets you upload a file to S3 and then download it.
 *
 * ```bash
 * curl -F file=@elysia.png http://localhost:3000/
 * curl http://localhost:3000/latest
 * ```
 *
 * Finally, you can deploy it using `bun neo deploy --stage production`.
 */
export default $config({
  app(input) {
    return {
      name: "aws-bun-elysia",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const bucket = new neo.aws.Bucket("MyBucket");
    const vpc = new neo.aws.Vpc("MyVpc");

    const cluster = new neo.aws.Cluster("MyCluster", { vpc });
    new neo.aws.Service("MyService", {
      cluster,
      loadBalancer: {
        ports: [{ listen: "80/http", forward: "3000/http" }],
      },
      dev: {
        command: "bun dev",
      },
      link: [bucket],
    });
  },
});
