/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    MyBucket: {
      name: string
      type: "neo.aws.Bucket"
    }
  }
}
export {}
