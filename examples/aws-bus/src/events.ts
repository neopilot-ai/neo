import { event } from "neo/event";
import { ZodValidator } from "neo/event/validator";
import { z } from "zod";

const defineEvent = event.builder({
  validator: ZodValidator,
  metadata: () => {
    return {
      timestamp: Date.now(),
    };
  },
});

export const MyEvent = defineEvent(
  "app.myevent",
  z.object({
    foo: z.string(),
  }),
);
