import { bus } from "neo/aws/bus";
import { testEvent } from "./event";

export const handler = bus.subscriber([testEvent], async (evt, raw) => {
  console.log("event", evt, raw, process.env);
  const message = evt.properties.message;

  if (message !== "hello") {
    throw new Error("🚨 bus subscriber failed");
  }
});
