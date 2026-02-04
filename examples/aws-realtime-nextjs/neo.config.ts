/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-realtime-nextjs",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const realtime = new neo.aws.Realtime("MyRealtime", {
      authorizer: "authorizer.handler",
    });

    new neo.aws.Nextjs("MyWeb", {
      link: [realtime],
    });
  },
});
