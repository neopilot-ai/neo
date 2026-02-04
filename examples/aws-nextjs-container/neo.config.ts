/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-nextjs-container",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const vpc = new neo.aws.Vpc("MyVpc");
    const bucket = new neo.aws.Bucket("MyBucket", {
      access: "public",
    });

    const cluster = new neo.aws.Cluster("MyCluster", { vpc });

    new neo.aws.Service("MyService", {
      cluster,
      loadBalancer: {
        ports: [{ listen: "80/http", forward: "3000/http" }],
      },
      dev: {
        command: "npm run dev",
      },
      link: [bucket],
    });
  },
});
