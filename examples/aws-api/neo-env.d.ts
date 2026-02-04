/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    MyApi: {
      type: "neo.aws.ApiGatewayV2"
      url: string
    }
    MyBucket: {
      name: string
      type: "neo.aws.Bucket"
    }
  }
}
export {}
