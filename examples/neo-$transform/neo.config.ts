/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## Default function props
 *
 * Set default props for all the functions in your app using the global [`$transform`](/docs/reference/global/#transform).
 */
export default $config({
  app(input) {
    return {
      name: "neo-transform",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    $transform(neo.aws.Function, (args) => {
      args.runtime = "nodejs14.x";
      args.environment = {
        FOO: "BAR",
      };
    });
    new neo.aws.Function("MyFunction", {
      handler: "index.ts",
    });
  },
});
