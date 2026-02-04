import { RealtimeAuthHandler } from "neo";

export const handler = RealtimeAuthHandler(async () => {
  return {
    subscribe: [process.env.NEO_TOPIC],
    publish: [process.env.NEO_TOPIC],
  };
});
