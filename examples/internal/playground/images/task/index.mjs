import { Resource } from "neo";

console.log({
  sdk: Resource.MyBucket.name,
  env: process.env.NEO_RESOURCE_MyBucket,
});
