/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## Subscribe to topics
 *
 * Create an SNS topic, publish to it from a function, and subscribe to it with a function and a queue.
 */
export default $config({
  app(input) {
    return {
      name: "aws-topic",
      home: "aws",
      removal: input?.stage === "production" ? "retain" : "remove",
    };
  },
  async run() {
    const queue = new neo.aws.Queue("MyQueue");
    queue.subscribe("subscriber.handler");

    const topic = new neo.aws.SnsTopic("MyTopic");
    topic.subscribe("MySubscriber1", "subscriber.handler", {});
    topic.subscribeQueue("MySubscriber2", queue.arn);

    const app = new neo.aws.Function("MyApp", {
      handler: "publisher.handler",
      link: [topic],
      url: true,
    });

    return {
      app: app.url,
      topic: topic.name,
    };
  },
});
