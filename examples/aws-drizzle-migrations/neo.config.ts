/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## Drizzle migrations in CI/CD
 *
 * An example on how to run Drizzle migrations as a part of your CI/CD.
 *
 * Start by creating a function that runs migrations.
 *
 * ```ts title="neo.config.ts"
 * const migrator = new neo.aws.Function("DatabaseMigrator", {
 *   handler: "src/migrator.handler",
 *   link: [rds],
 *   vpc,
 *   copyFiles: [
 *     {
 *       from: "migrations",
 *       to: "./migrations",
 *     },
 *   ],
 * });
 * ```
 *
 * Where `src/migrator.ts` looks like.
 *
 * ```ts title="src/migrator.ts"
 * import { db } from "./drizzle";
 * import { migrate } from "drizzle-orm/postgres-js/migrator";
 *
 * export const handler = async (event: any) => {
 *   await migrate(db, {
 *     migrationsFolder: "./migrations",
 *   });
 * };
 * ```
 *
 * And we can set it up to run on every deploy.
 *
 * ```ts title="neo.config.ts"
 * if (!$dev){
 *   new aws.lambda.Invocation("DatabaseMigratorInvocation", {
 *     input: Date.now().toString(),
 *     functionName: migrator.name,
 *   });
 * }
 * ```
 *
 * We use the current time to make sure the function runs on every deploy.
 */
export default $config({
  app(input) {
    return {
      name: "aws-drizzle-migrations",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const vpc = new neo.aws.Vpc("MyVpc", { bastion: true, nat: "ec2" });
    const rds = new neo.aws.Postgres("MyPostgres", { vpc, proxy: true });

    new neo.aws.Function("MyApi", {
      vpc,
      url: true,
      link: [rds],
      handler: "src/api.handler",
    });

    const migrator = new neo.aws.Function("DatabaseMigrator", {
      handler: "src/migrator.handler",
      link: [rds],
      vpc,
      copyFiles: [
        {
          from: "migrations",
          to: "./migrations",
        },
      ],
    });

    if (!$dev) {
      new aws.lambda.Invocation("DatabaseMigratorInvocation", {
        input: Date.now().toString(),
        functionName: migrator.name,
      });
    }

    new neo.x.DevCommand("Studio", {
      link: [rds],
      dev: {
        command: "npx drizzle-kit studio",
      },
    });
  },
});
