/// <reference path="./.neo/platform/config.d.ts" />

/**
 * ## AWS Bus subscriptions
 *
 * Subscribe bus events with AWS Lambda functions.
 */
export default $config({
  app(input) {
    return {
      name: "aws-bus",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    const bus = new neo.aws.Bus("Bus");

    const publisher = new neo.aws.Function("Publisher", {
      handler: "./src/publisher.handler",
      url: true,
      link: [bus],
    });

    bus.subscribe("Example", "./src/receiver.handler");
  },
});
