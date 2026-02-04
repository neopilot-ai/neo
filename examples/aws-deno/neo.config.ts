/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-deno",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const vpc = new neo.aws.Vpc("MyVpc");
    const bucket = new neo.aws.Bucket("MyBucket");

    const cluster = new neo.aws.Cluster("MyCluster", { vpc });

    new neo.aws.Service("MyService", {
      cluster,
      link: [bucket],
      loadBalancer: {
        ports: [{ listen: "80/http", forward: "8000/http" }],
      },
      dev: {
        command: "deno task dev",
      },
    });
  },
});
