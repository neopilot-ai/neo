/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    Api: {
      name: string
      type: "neo.aws.Function"
      url: string
    }
  }
}
export {}