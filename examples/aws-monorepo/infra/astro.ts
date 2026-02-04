import { api } from "./api";

const bucket = new neo.aws.Bucket("MyBucket");
export const astro = new neo.aws.Astro("Astro", {
  path: "packages/astro",
  link: [bucket],
  environment: {
    VITE_API_URL: api.url,
  },
});
