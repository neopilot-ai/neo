/// <reference path="./.neo/platform/config.d.ts" />
export default $config({
  app(input) {
    return {
      name: "aws-rails",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const bucket = new neo.aws.Bucket("MyBucket", {
      access: "public",
    });
    const vpc = new neo.aws.Vpc("MyVpc");

    const cluster = new neo.aws.Cluster("MyCluster", { vpc });
    new neo.aws.Service("MyService", {
      cluster,
      loadBalancer: {
        ports: [{ listen: "80/http", forward: "3000/http" }],
      },
      environment: {
        RAILS_MASTER_KEY: (await import("fs")).readFileSync(
          "config/master.key",
          "utf8"
        ),
      },
      dev: {
        command: "bin/rails server",
      },
      link: [bucket],
    });
    return { vpc: vpc.id };
  },
});
