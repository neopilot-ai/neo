/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    MyBucket: {
      name: string
      type: "neo.aws.Bucket"
    }
    Secret1: {
      type: "neo.neo.Secret"
      value: string
    }
    Secret2: {
      type: "neo.neo.Secret"
      value: string
    }
  }
}
export {}