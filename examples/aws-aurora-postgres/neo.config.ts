/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## AWS Aurora Postgres
 *
 * In this example, we deploy a Aurora Postgres database.
 *
 * ```ts title="neo.config.ts"
 * const postgres = new neo.aws.Aurora("MyDatabase", {
 *   engine: "postgres",
 *   vpc,
 * });
 * ```
 *
 * And link it to a Lambda function.
 *
 * ```ts title="neo.config.ts" {4}
 * new neo.aws.Function("MyApp", {
 *   handler: "index.handler",
 *   link: [postgres],
 *   url: true,
 *   vpc,
 * });
 * ```
 *
 * In the function we use the [`postgres`](https://www.npmjs.com/package/postgres) package.
 *
 * ```ts title="index.ts"
 * import postgres from "postgres";
 * import { Resource } from "neo";
 *
 * const sql = postgres({
 *   username: Resource.MyDatabase.username,
 *   password: Resource.MyDatabase.password,
 *   database: Resource.MyDatabase.database,
 *   host: Resource.MyDatabase.host,
 *   port: Resource.MyDatabase.port,
 * });
 * ```
 *
 * We also enable the `bastion` option for the VPC. This allows us to connect to the database
 * from our local machine with the `neo tunnel` CLI.
 *
 * ```bash "sudo"
 * sudo npx neo tunnel install
 * ```
 *
 * This needs _sudo_ to create a network interface on your machine. You’ll only need to do this
 * once on your machine.
 *
 * Now you can run `npx neo dev` and you can connect to the database from your local machine.
 *
 */
export default $config({
  app(input) {
    return {
      name: "aws-aurora-postgres",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const vpc = new neo.aws.Vpc("MyVpc", {
      nat: "ec2",
      bastion: true,
    });
    const postgres = new neo.aws.Aurora("MyDatabase", {
      engine: "postgres",
      vpc,
    });
    new neo.aws.Function("MyApp", {
      handler: "index.handler",
      link: [postgres],
      url: true,
      vpc,
    });

    return {
      host: postgres.host,
      port: postgres.port,
      username: postgres.username,
      password: postgres.password,
      database: postgres.database,
    };
  },
});
