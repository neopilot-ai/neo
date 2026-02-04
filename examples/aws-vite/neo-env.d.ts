/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    "Web": {
      "type": "neo.aws.StaticSite"
      "url": string
    }
  }
}
export {}
