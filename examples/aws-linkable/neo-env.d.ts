/* tslint:disable */
/* eslint-disable */
import "neo"
declare module "neo" {
  export interface Resource {
    ExistingResources: {
      bucketName: string
      queueName: string
      type: "neo.neo.Linkable"
    }
    Function: {
      name: string
      type: "neo.aws.Function"
      url: string
    }
    Topic: {
      arn: string
      name: string
      type: "aws.sns/topic.Topic"
    }
  }
}
export {}
