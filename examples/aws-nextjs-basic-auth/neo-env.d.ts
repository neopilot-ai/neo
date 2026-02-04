/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    "MyWeb": {
      "type": "neo.aws.Nextjs"
      "url": string
    }
    "PASSWORD": {
      "type": "neo.neo.Secret"
      "value": string
    }
    "USERNAME": {
      "type": "neo.neo.Secret"
      "value": string
    }
  }
}
export {}
