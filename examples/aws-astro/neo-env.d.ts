/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    MyBucket: {
      name: string
      type: "neo.aws.Bucket"
    }
    MyWeb: {
      type: "neo.aws.Astro"
      url: string
    }
  }
}
export {}
