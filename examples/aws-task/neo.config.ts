/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## AWS Task
 *
 * Use the [`Task`](/docs/component/aws/task) component to run background tasks.
 *
 * We have a node script that we want to run in `image/index.mjs`. It'll be deployed as a
 * Docker container using `image/Dockerfile`.
 *
 * We also have a function that the task is linked to. It uses the [SDK](/docs/reference/sdk/)
 * to start the task.
 *
 * ```ts title="index.ts" {5}
 * import { Resource } from "neo";
 * import { task } from "neo/aws/task";
 *
 * export const handler = async () => {
 *   const ret = await task.run(Resource.MyTask);
 *   return {
 *     statusCode: 200,
 *     body: JSON.stringify(ret, null, 2),
 *   };
 * };
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
      name: "aws-task",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const bucket = new neo.aws.Bucket("MyBucket");
    const vpc = new neo.aws.Vpc("MyVpc", { nat: "ec2" });

    const cluster = new neo.aws.Cluster("MyCluster", { vpc });

    const task = new neo.aws.Task("MyTask", {
      cluster,
      link: [bucket],
      image: {
        context: "image",
      },
      dev: {
        command: "node index.mjs",
      },
    });

    new neo.aws.Function("MyApp", {
      vpc,
      url: true,
      link: [task],
      handler: "index.handler",
    });
  },
});
