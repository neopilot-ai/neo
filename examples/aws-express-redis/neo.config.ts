/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## AWS Express Redis
 *
 * Creates a hit counter app with Express and Redis.
 *
 * This deploys Express as a Fargate service to ECS and it's linked to Redis.
 *
 * ```ts title="neo.config.ts" {9}
 * new neo.aws.Service("MyService", {
 *   cluster,
 *   loadBalancer: {
 *     ports: [{ listen: "80/http" }],
 *   },
 *   dev: {
 *     command: "node --watch index.mjs",
 *   },
 *   link: [redis],
 * });
 * ```
 *
 * Since our Redis cluster is in a VPC, we’ll need a tunnel to connect to it from our local
 * machine.
 *
 * ```bash "sudo"
 * sudo npx neo tunnel install
 * ```
 *
 * This needs _sudo_ to create a network interface on your machine. You’ll only need to do this
 * once on your machine.
 *
 * To start your app locally run.
 *
 * ```bash
 * npx neo dev
 * ```
 *
 * Now if you go to `http://localhost:80` you’ll see a counter update as you refresh the page.
 *
 * Finally, you can deploy it using `npx neo deploy --stage production` using a `Dockerfile`
 * that's included in the example.
 */
export default $config({
  app(input) {
    return {
      name: "aws-express-redis",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const vpc = new neo.aws.Vpc("MyVpc", { bastion: true });
    const redis = new neo.aws.Redis("MyRedis", { vpc });
    const cluster = new neo.aws.Cluster("MyCluster", { vpc });

    new neo.aws.Service("MyService", {
      cluster,
      link: [redis],
      loadBalancer: {
        ports: [{ listen: "80/http" }],
      },
      dev: {
        command: "node --watch index.mjs",
      },
    });
  },
});
