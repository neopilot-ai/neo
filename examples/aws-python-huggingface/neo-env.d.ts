/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    "MyPythonFunction": {
      "name": string
      "type": "neo.aws.Function"
      "url": string
    }
  }
}
export {}
