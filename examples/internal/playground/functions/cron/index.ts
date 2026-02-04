import { Resource } from "neo";

export const handler = async (event) => {
  console.log("event", event);
  console.log("Resource.MyBucket", Resource.MyBucket);
};
