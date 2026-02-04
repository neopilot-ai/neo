import { defineConfig } from "astro/config";
import aws from "astro-neo";

export default defineConfig({
  output: "server",
  adapter: aws(),
});
