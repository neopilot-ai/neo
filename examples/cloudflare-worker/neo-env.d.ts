/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    MyWorker: {
      type: "neo.cloudflare.Worker"
      url: string
    }
  }
}
// cloudflare 
declare module "neo" {
  export interface Resource {
    MyBucket: import("@cloudflare/workers-types").R2Bucket
  }
}
