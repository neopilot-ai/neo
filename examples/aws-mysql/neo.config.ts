/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## AWS MySQL
 *
 * In this example, we deploy an RDS MySQL database.
 *
 * ```ts title="neo.config.ts"
 * const mysql = new neo.aws.Mysql("MyDatabase", {
 *   vpc,
 * });
 * ```
 *
 * And link it to a Lambda function.
 *
 * ```ts title="neo.config.ts" {3}
 * new neo.aws.Function("MyApp", {
 *   handler: "index.handler",
 *   link: [mysql],
 *   url: true,
 *   vpc,
 * });
 * ```
 *
 * Now in the function we can access the database.
 *
 * ```ts title="index.ts"
 * const connection = await mysql.createConnection({
 *   database: Resource.MyDatabase.database,
 *   host: Resource.MyDatabase.host,
 *   port: Resource.MyDatabase.port,
 *   user: Resource.MyDatabase.username,
 *   password: Resource.MyDatabase.password,
 * });
 * ```
 *
 * We also enable the `bastion` option for the VPC. This allows us to connect to
 * the database from our local machine with the `neo tunnel` CLI.
 *
 * ```bash "sudo"
 * sudo npx neo tunnel install
 * ```
 *
 * This needs _sudo_ to create a network interface on your machine. You’ll only
 * need to do this once on your machine.
 *
 * Now you can run `npx neo dev` and you can connect to the database from your local machine.
 *
 */
export default $config({
  app(input) {
    return {
      name: "aws-mysql",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const vpc = new neo.aws.Vpc("MyVpc", { nat: "ec2", bastion: true });
    const mysql = new neo.aws.Mysql("MyDatabase", {
      vpc,
    });
    const app = new neo.aws.Function("MyApp", {
      handler: "index.handler",
      link: [mysql],
      url: true,
      vpc,
    });

    return {
      app: app.url,
      host: mysql.host,
      port: mysql.port,
      username: mysql.username,
      password: mysql.password,
      database: mysql.database,
    };
  },
});
