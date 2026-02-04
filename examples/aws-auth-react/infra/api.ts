import { auth } from "./auth";

export const api = new neo.aws.Function("MyApi", {
  url: true,
  link: [auth],
  handler: "packages/functions/src/api.handler",
});
