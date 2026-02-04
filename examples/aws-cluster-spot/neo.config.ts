/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## AWS Cluster Spot capacity
 *
 * This example, shows how to use the Fargate Spot capacity provider for your services.
 *
 * We have it set to use only Fargate Spot instances for all non-production stages. Learn more
 * about the [`capacity`](/docs/component/aws/cluster#capacity) prop.
 */
export default $config({
  app(input) {
    return {
      name: "aws-cluster-spot",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const vpc = new neo.aws.Vpc("MyVpc");

    const cluster = new neo.aws.Cluster("MyCluster", { vpc });
    new neo.aws.Service("MyService", {
      cluster,
      loadBalancer: {
        ports: [{ listen: "80/http" }],
      },
      capacity: $app.stage === "production" ? undefined : "spot",
    });
  },
});
