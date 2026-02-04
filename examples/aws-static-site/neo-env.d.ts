/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    MySite: {
      type: "neo.aws.StaticSite"
      url: string
    }
  }
}
export {}