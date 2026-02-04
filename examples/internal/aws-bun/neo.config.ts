/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-bun",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const vpc = new neo.aws.Vpc("MyVpc");
    const redis = new neo.aws.Redis("MyRedis", { vpc });
    const cluster = new neo.aws.Cluster("MyCluster", { vpc });

    new neo.aws.Service("MyService", {
      cluster,
      link: [redis],
      loadBalancer: {
        ports: [{ listen: "80/http", forward: "3000/http" }],
      },
      dev: {
        command: "bun dev",
      },
    });
  },
});
