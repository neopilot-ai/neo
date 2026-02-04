/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    "MyBucket": {
      "name": string
      "type": "neo.aws.Bucket"
    }
    "MyFunction": {
      "name": string
      "type": "neo.aws.Function"
      "url": string
    }
    "MyWeb": {
      "type": "neo.aws.StaticSite"
      "url": string
    }
  }
}
export {}
