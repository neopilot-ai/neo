/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-angular",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const bucket = new neo.aws.Bucket("MyBucket", {
      access: "public",
    });

    const pre = new neo.aws.Function("MyFunction", {
      url: true,
      link: [bucket],
      handler: "functions/presigned.handler",
    });

    new neo.aws.StaticSite("MyWeb", {
      dev: {
        command: "npm run start",
      },
      build: {
        output: "dist/browser",
        command: "ng build --output-path dist",
      },
      environment: {
        NG_APP_PRESIGNED_API: pre.url
      }
    });
  },
});

