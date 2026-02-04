/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "{{.App}}",
      removal: input?.stage === "production" ? "retain" : "remove",
      protect: ["production"].includes(input?.stage),
      home: "{{.Home}}",
    };
  },
  async run() {
    new neo.{{.Home}}.Remix("MyWeb");
  },
});
