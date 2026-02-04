/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    MyApi: {
      name: string
      type: "neo.aws.Function"
      url: string
    }
    MyBucket: {
      name: string
      type: "neo.aws.Bucket"
    }
    MyRouter: {
      type: "neo.aws.Router"
      url: string
    }
  }
}
export {}