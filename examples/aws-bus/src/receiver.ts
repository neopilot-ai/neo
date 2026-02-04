import { bus } from "neo/aws/bus";
import { MyEvent } from "./events";

export const handler = bus.subscriber([MyEvent], async (event) => {
  console.log({ event });
});
