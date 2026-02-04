/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-flutter-web",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    new neo.aws.StaticSite("MySite", {
      build: {
        command: "flutter build web",
        output: "build/web",
      },
    });
  },
});
