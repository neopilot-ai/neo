/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    MyApi: {
      type: "neo.aws.ApiGatewayV1"
      url: string
    }
    MyApiAuthorizerMyAuthorizerFunction: {
      name: string
      type: "neo.aws.Function"
    }
  }
}
export {}