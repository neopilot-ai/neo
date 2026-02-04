import { VisibleError } from "../error";
import { FunctionPermissionArgs } from "./function";

export const URL_UNAVAILABLE = "http://url-unavailable-in-dev.mode";

/** @deprecated
 * instead try
 * ```
 * neo.Linkable.wrap(MyResource, (resource) => ({
 *   properties: { ... },
 *   with: [
 *     neo.aws.permission({ actions: ["foo:*"], resources: [resource.arn] })
 *   ]
 * }))
 * ```
 */
export function linkable<T>(
  obj: { new (...args: any[]): T },
  cb: (resource: T) => FunctionPermissionArgs[],
) {
  throw new VisibleError(
    [
      "neo.aws.linkable is deprecated. Use neo.Linkable.wrap instead.",
      "neo.Linkable.wrap(MyResource, (resource) => ({",
      "  properties: { ... },",
      "  with: [",
      '    neo.aws.permission({ actions: ["foo:*"], resources: [resource.arn] })',
      "  ]",
      "}))",
    ].join("\n"),
  );
}
