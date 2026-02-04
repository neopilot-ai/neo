/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-step-functions",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    // Create a function to be invoked by the state machine
    const app = new neo.aws.Function("MyApp", {
      handler: "index.handler",
    });

    // Define all the states of the state machine
    const lambdaInvoke = neo.aws.StepFunctions.lambdaInvoke({
      name: "LambdaInvoke",
      function: app,
      payload: {
        foo: "bar",
      },
    });
    const pass = neo.aws.StepFunctions.pass({ name: "Pass" });
    const wait = neo.aws.StepFunctions.wait({
      name: "Wait",
      time: "2 seconds",
    });
    const choice = neo.aws.StepFunctions.choice({ name: "Choice" });
    const parallel = neo.aws.StepFunctions.parallel({ name: "Parallel" });
    const parallelA = neo.aws.StepFunctions.pass({ name: "ParallelA" });
    const parallelB = neo.aws.StepFunctions.pass({ name: "ParallelB" });
    const parallelC = neo.aws.StepFunctions.pass({ name: "ParallelC" });
    const mapA = neo.aws.StepFunctions.pass({ name: "MapA" });
    const mapB = neo.aws.StepFunctions.pass({ name: "MapB" });
    const map = neo.aws.StepFunctions.map({
      name: "Map",
      processor: mapA.next(mapB),
      items: ["a", "b", "c"],
    });
    const success = neo.aws.StepFunctions.succeed({ name: "Succeed" });
    const fail = neo.aws.StepFunctions.fail({ name: "Fail" });
    const last = neo.aws.StepFunctions.pass({ name: "Last" });

    // Create the state machine
    new neo.aws.StepFunctions("MyStateMachine", {
      definition: lambdaInvoke
        .catch(fail)
        .next(pass)
        .next(wait)
        .next(parallel.branch(parallelA.next(parallelB)).branch(parallelC))
        .next(map)
        .next(
          choice
            .when("{% 1+1 = 2 %}", success)
            .when("{% 1+1 = 3 %}", fail)
            .otherwise(last)
        ),
    });
  },
});
