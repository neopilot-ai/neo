// @ts-check
import aws from "astro-neo";
import { defineConfig } from 'astro/config';

// https://astro.build/config
export default defineConfig({
  output: "server",
  adapter: aws({
    responseMode: "stream",
  }),
});
