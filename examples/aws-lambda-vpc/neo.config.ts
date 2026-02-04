/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## AWS Lambda in a VPC
 *
 * You can use NEO to locally work on Lambda functions that are in a VPC. To do so, you'll
 * need to enable `bastion` and `nat` on the `Vpc` component.
 *
 * ```ts title="neo.config.ts"
 * new neo.aws.Vpc("MyVpc", { bastion: true, nat: "managed" });
 * ```
 *
 * The NAT gateway is necessary to allow your Lambda function to connect to the internet. While,
 * the bastion host is necessary for your local machine to be able to tunnel to the VPC.
 *
 * You'll need to install the tunnel, if you haven't done this before.
 *
 * ```bash "sudo"
 * sudo neo tunnel install
 * ```
 *
 * This needs _sudo_ to create the network interface on your machine. You'll only need to do
 * this once.
 *
 * Now you can run `neo dev`, your function can access resources in the VPC. For example, here
 * we are connecting to a Redis cluster.
 *
 * ```ts title="index.ts"
 * const redis = new Cluster(
 *   [{ host: Resource.MyRedis.host, port: Resource.MyRedis.port }],
 *   {
 *     dnsLookup: (address, callback) => callback(null, address),
 *     redisOptions: {
 *       tls: {},
 *       username: Resource.MyRedis.username,
 *       password: Resource.MyRedis.password,
 *     },
 *   }
 * );
 * ```
 *
 * The Redis cluster is in the same VPC as the function.
 */
export default $config({
  app(input) {
    return {
      name: "aws-lambda-vpc",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const vpc = new neo.aws.Vpc("MyVpc", { bastion: true, nat: "managed" });
    const redis = new neo.aws.Redis("MyRedis", { vpc });
    const api = new neo.aws.Function("MyFunction", {
      vpc,
      url: true,
      link: [redis],
      handler: "index.handler"
    });

    return {
      url: api.url,
    };
  },
});
