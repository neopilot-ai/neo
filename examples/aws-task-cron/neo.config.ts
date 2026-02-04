/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## AWS Task Cron
 *
 * Use the [`Task`](/docs/component/aws/task) and [`Cron`](/docs/component/aws/cron) components
 * for long running background tasks.
 *
 * We have a node script that we want to run in `index.mjs`. It'll be deployed as a
 * Docker container using `Dockerfile`.
 *
 * It'll be invoked by a cron job that runs every 2 minutes.
 *
 * ```ts title="neo.config.ts"
 * new neo.aws.Cron("MyCron", {
 *   task,
 *   schedule: "rate(2 minutes)"
 * });
 * ```
 *
 * When this is run in `neo dev`, the task is executed locally using `dev.command`.
 *
 * ```ts title="neo.config.ts"
 * dev: {
 *   command: "node index.mjs"
 * }
 * ```
 *
 * To deploy, you need the Docker daemon running.
 */
export default $config({
  app(input) {
    return {
      name: "aws-task-cron",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const bucket = new neo.aws.Bucket("MyBucket");
    const vpc = new neo.aws.Vpc("MyVpc");

    const cluster = new neo.aws.Cluster("MyCluster", { vpc });
    const task = new neo.aws.Task("MyTask", {
      cluster,
      link: [bucket],
      dev: {
        command: "node index.mjs",
      },
    });

    new neo.aws.Cron("MyCron", {
      task,
      schedule: "rate(2 minutes)",
    });
  },
});
