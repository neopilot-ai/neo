/// <reference path="./.neo/platform/config.d.ts" />

export default $config({
  app(input) {
    return {
      name: "aws-svelte-kit-remote-functions",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "aws",
    };
  },
  async run() {
    new neo.aws.SvelteKit("MyWeb", {
      dev: {
        command: "bun run dev",
      },
    });
  },
});
