/**
 * The AWS Permission Linkable helper is used to define the AWS permissions included with the
 * [`neo.Linkable`](/docs/component/linkable/) component.
 *
 * @example
 *
 * ```ts
 * neo.aws.permission({
 *   actions: ["lambda:InvokeFunction"],
 *   resources: ["*"]
 * })
 * ```
 *
 * @packageDocumentation
 */

import { Prettify } from "../component.js";
import { FunctionPermissionArgs } from "./function.js";

export interface InputArgs extends Prettify<FunctionPermissionArgs> {}

/**
 * The AWS Permission Linkable helper is used to define the AWS permissions included with the
 * [`neo.Linkable`](/docs/component/linkable/) component.
 *
 * @example
 *
 * ```ts
 * neo.aws.permission({
 *   actions: ["lambda:InvokeFunction"],
 *   resources: ["*"]
 * })
 * ```
 */
export function permission(input: InputArgs) {
  return {
    type: "aws.permission" as const,
    ...input,
  };
}

export type Permission = ReturnType<typeof permission>;
