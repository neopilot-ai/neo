// @ts-check
import { defineConfig } from "astro/config";
import aws from "astro-neo";
import cloudflare from "@astrojs/cloudflare";
//import aws from "../../../../../../astro-neo/packages/astro-neo/dist/adapter";

// https://astro.build/config
export default defineConfig({
  image: {
    domains: ["neo.dev"],
  },
  output: "server",
  adapter: cloudflare(),
  redirects: {
    "/redirect-to-route": "/prerendered",
    "/redirect-to-url": "https://www.google.com",
    "/redirect/[slug]": "/sub/[slug]",
  },
  base: undefined,
});
