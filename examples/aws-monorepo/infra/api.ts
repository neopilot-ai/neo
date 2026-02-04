import { database } from "./database";

export const api = new neo.aws.Function("Api", {
  url: true,
  link: [database],
  handler: "./packages/functions/src/api.handler",
  environment: {
    foo: "9",
  },
});
