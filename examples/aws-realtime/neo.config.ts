/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-realtime",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const topic = `${$app.name}/${$app.stage}/chat`;
    const realtime = new neo.aws.Realtime("MyRealtime", {
      authorizer: {
        handler: "authorizer.handler",
        environment: {
          NEO_TOPIC: topic,
        },
      },
    });
    realtime.subscribe("subscriber.handler", {
      filter: topic,
    });

    new neo.aws.StaticSite("Web", {
      path: "web",
      build: {
        command: "npm run build",
        output: "dist",
      },
      environment: {
        VITE_REALTIME_ENDPOINT: realtime.endpoint,
        VITE_TOPIC: topic,
        VITE_AUTHORIZER: realtime.authorizer,
      },
    });

    const publisher = new neo.aws.Function("MyApp", {
      handler: "publisher.handler",
      environment: {
        NEO_TOPIC: topic,
      },
      url: true,
      link: [realtime],
    });

    return {
      publisher: publisher.url,
    };
  },
});
