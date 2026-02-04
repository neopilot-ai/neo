import { bus } from "neo/aws/bus";
import { z } from "zod";
import { Resource } from "neo";
import { MyEvent } from "./events";

export async function handler() {
  await bus.publish(Resource.Bus, MyEvent, { foo: "hello" });

  return {
    statusCode: 200,
  };
}
