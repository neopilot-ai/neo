import { output } from "@pulumi/pulumi";
import { Link } from "./link";
import { Component } from "./component";
import { Input } from "./input";

export interface Definition<
  Properties extends Record<string, any> = Record<string, any>,
> {
  /**
   * Define values that the linked resource can access at runtime. These can be outputs
   * from other resources or constants.
   *
   * @example
   * ```ts
   * {
   *   properties: { foo: "bar" }
   * }
   * ```
   */
  properties: Properties;
  /**
   * Include AWS permissions or Cloudflare bindings for the linkable resource. The linked
   * resource will have these permissions or bindings.
   *
   * @example
   * Include AWS permissions.
   *
   * ```ts
   * {
   *   include: [
   *     neo.aws.permission({
   *       actions: ["lambda:InvokeFunction"],
   *       resources: ["*"]
   *     })
   *   ]
   * }
   * ```
   *
   * Include Cloudflare bindings.
   *
   * ```ts
   * {
   *   include: [
   *     neo.cloudflare.binding({
   *       type: "r2BucketBindings",
   *       properties: {
   *         bucketName: "my-bucket"
   *       }
   *     })
   *   ]
   * }
   * ```
   */
  include?: {
    type: string;
    [key: string]: any;
  }[];
}

/**
 * The `Linkable` component and the `Linkable.wrap` method lets you link any resources in your
 * app; not just the built-in NEO components. It also lets you modify the links NEO creates.
 *
 * @example
 *
 * #### Linking any value
 *
 * The `Linkable` component takes a list of properties that you want to link. These can be
 * outputs from other resources or constants.
 *
 * ```ts title="neo.config.ts"
 * new neo.Linkable("MyLinkable", {
 *   properties: { foo: "bar" }
 * });
 * ```
 *
 * You can also use this to combine multiple resources into a single linkable resource. And
 * optionally include permissions or bindings for the linked resource.
 *
 * ```ts title="neo.config.ts"
 * const bucketA = new neo.aws.Bucket("MyBucketA");
 * const bucketB = new neo.aws.Bucket("MyBucketB");
 *
 * const storage = new neo.Linkable("MyStorage", {
 *   properties: {
 *     foo: "bar",
 *     bucketA: bucketA.name,
 *     bucketB: bucketB.name
 *   },
 *   include: [
 *     neo.aws.permission({
 *       actions: ["s3:*"],
 *       resources: [bucketA.arn, bucketB.arn]
 *     })
 *   ]
 * });
 * ```
 *
 * You can now link this resource to your frontend or a function.
 *
 * ```ts title="neo.config.ts" {3}
 * new neo.aws.Function("MyApi", {
 *   handler: "src/lambda.handler",
 *   link: [storage]
 * });
 * ```
 *
 * Then use the [SDK](/docs/reference/sdk/) to access it at runtime.
 *
 * ```js title="src/lambda.ts"
 * import { Resource } from "neo";
 *
 * console.log(Resource.MyStorage.bucketA);
 * ```
 *
 * #### Linking any resource
 *
 * You can also wrap any Pulumi Resource class to make it linkable.
 *
 * ```ts title="neo.config.ts"
 * neo.Linkable.wrap(aws.dynamodb.Table, (table) => ({
 *   properties: { tableName: table.name },
 *   include: [
 *     neo.aws.permission({
 *       actions: ["dynamodb:*"],
 *       resources: [table.arn]
 *     })
 *   ]
 * }));
 * ```
 *
 * Now you create an instance of `aws.dynamodb.Table` and link it in your app like any other NEO
 * component.
 *
 * ```ts title="neo.config.ts" {7}
 * const table = new aws.dynamodb.Table("MyTable", {
 *   attributes: [{ name: "id", type: "S" }],
 *   hashKey: "id"
 * });
 *
 * new neo.aws.Nextjs("MyWeb", {
 *   link: [table]
 * });
 * ```
 *
 * And use the [SDK](/docs/reference/sdk/) to access it at runtime.
 *
 * ```js title="app/page.tsx"
 * import { Resource } from "neo";
 *
 * console.log(Resource.MyTable.tableName);
 * ```
 *
 * Your function will also have the permissions defined above.
 *
 * #### Modify built-in links
 *
 * You can also modify how NEO creates links. For example, you might want to change the
 * permissions of a linkable resource.
 *
 * ```ts title="neo.config.ts" "neo.aws.Bucket"
 *  neo.Linkable.wrap(neo.aws.Bucket, (bucket) => ({
 *    properties: { name: bucket.name },
 *    include: [
 *      neo.aws.permission({
 *        actions: ["s3:GetObject"],
 *        resources: [bucket.arn]
 *      })
 *    ]
 *  }));
 * ```
 *
 * This overrides the built-in link and lets you create your own.
 */
export class Linkable<T extends Record<string, any>>
  extends Component
  implements Link.Linkable
{
  private _name: string;
  private _definition: Definition<T>;

  public static wrappedResources = new Set<string>();

  constructor(name: string, definition: Definition<T>) {
    super("neo:sst:Linkable", name, definition, {});
    this._name = name;
    this._definition = definition;
  }

  public get name() {
    return output(this._name);
  }

  public get properties() {
    return this._definition.properties;
  }

  /** @internal */
  public getNEOLink() {
    return this._definition;
  }

  /**
   * Wrap any resource class to make it linkable. Behind the scenes this modifies the
   * prototype of the given class.
   *
   * :::tip
   * Use `Linkable.wrap` to make any resource linkable.
   * :::
   *
   * @param cls The resource class to wrap.
   * @param cb A callback that returns the definition for the linkable resource.
   *
   * @example
   *
   * Here we are wrapping the [`aws.dynamodb.Table`](https://www.pulumi.com/registry/packages/aws/api-docs/dynamodb/table/)
   * class to make it linkable.
   *
   * ```ts title="neo.config.ts"
   * Linkable.wrap(aws.dynamodb.Table, (table) => ({
   *   properties: { tableName: table.name },
   *   include: [
   *     neo.aws.permission({
   *       actions: ["dynamodb:*"],
   *       resources: [table.arn]
   *     })
   *   ]
   * }));
   * ```
   *
   * It's defining the properties that we want made accessible at runtime and the permissions
   * that the linked resource should have.
   *
   * Now you can link any `aws.dynamodb.Table` instances in your app just like any other NEO
   * component.
   *
   * ```ts title="neo.config.ts" {7}
   * const table = new aws.dynamodb.Table("MyTable", {
   *   attributes: [{ name: "id", type: "S" }],
   *   hashKey: "id",
   * });
   *
   * new neo.aws.Nextjs("MyWeb", {
   *   link: [table]
   * });
   * ```
   *
   * Since this applies to any resource, you can also use it to wrap NEO components and modify
   * how they are linked.
   *
   * ```ts title="neo.config.ts" "neo.aws.Bucket"
   * neo.Linkable.wrap(neo.aws.Bucket, (bucket) => ({
   *   properties: { name: bucket.name },
   *   include: [
   *     neo.aws.permission({
   *       actions: ["s3:GetObject"],
   *       resources: [bucket.arn]
   *     })
   *   ]
   * }));
   * ```
   *
   * This overrides the built-in link and lets you create your own.
   *
   * :::tip
   * You can modify the permissions granted by a linked resource.
   * :::
   *
   * In the above example, we're modifying the permissions to access a linked `neo.aws.Bucket`
   * in our app.
   */
  public static wrap<Resource>(
    cls: { new (...args: any[]): Resource },
    cb: (resource: Resource) => Definition,
  ) {
    // @ts-expect-error
    this.wrappedResources.add(cls.__pulumiType);

    cls.prototype.getNEOLink = function () {
      return cb(this);
    };
  }
}

/**
 * @deprecated
 * Use neo.Linkable instead.
 */
export class Resource extends Component implements Link.Linkable {
  private _properties: any;
  private _name: string;

  constructor(name: string, properties: any) {
    super(
      "neo:sst:Resource",
      name,
      {
        properties,
      },
      {},
    );
    console.warn("Resource is deprecated. Use neo.Linkable instead.");
    this._properties = properties;
    this._name = name;
  }

  public get name() {
    return output(this._name);
  }

  public get properties() {
    return this._properties;
  }

  /** @internal */
  public getNEOLink() {
    return {
      properties: this._properties,
    };
  }
}

export function env(env: Record<string, Input<string>>) {
  return {
    type: "environment" as const,
    env,
  };
}
