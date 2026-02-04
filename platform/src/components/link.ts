import {
  Input,
  Output,
  runtime,
  output,
  all,
  ComponentResource,
} from "@pulumi/pulumi";
import { VisibleError } from "./error.js";
import { Linkable } from "./linkable.js";

export module Link {
  export interface Definition<
    Properties extends Record<string, any> = Record<string, any>,
  > {
    properties: Properties;
    include?: {
      type: string;
      [key: string]: any;
    }[];
  }

  export class Ref extends ComponentResource {
    constructor(target: string, type: string, properties: any, include?: any) {
      super(
        "neo:neo:LinkRef",
        target + "LinkRef",
        {
          properties,
          include,
        },
        {},
      );
      this.registerOutputs({
        target: target,
        include,
        properties: {
          type: type.replaceAll(":", "."),
          ...properties,
        },
      });
    }
  }

  export function reset() {
    const links = new Set<string>();
    // Ensure component names are unique
    runtime.registerStackTransformation((args) => {
      const isLinkable =
        args.type.startsWith("neo:") ||
        Linkable.wrappedResources.has(args.type);
      if (isLinkable && !args.opts.parent) {
        const lcname = args.name.toLowerCase();

        // "App" is reserved and cannot be used as a component name.
        if (lcname === "app") {
          throw new VisibleError(
            `Component name "${args.name}" is reserved. Please choose a different name for your "${args.type}" component.`,
          );
        }

        // Ensure linkable resources have unique names. This includes all NEO components
        // and non-NEO components that are linkable.
        if (links.has(lcname)) {
          throw new VisibleError(`Component name ${args.name} is not unique.`);
        }
        links.add(lcname);
      }
      return {
        opts: args.opts,
        props: args.props,
      };
    });

    // Create link refs
    runtime.registerStackTransformation((args) => {
      const resource = args.resource;
      process.nextTick(() => {
        if (Link.isLinkable(resource) && !args.opts.parent) {
          try {
            const link = resource.getNEOLink();
            new Ref(args.name, args.type, link.properties, link.include);
          } catch (e) {}
        }
      });
      return {
        opts: args.opts,
        props: args.props,
      };
    });
  }

  export interface Linkable {
    urn: Output<string>;
    getNEOLink(): Definition;
  }

  export function isLinkable(obj: any): obj is Linkable {
    return "getNEOLink" in obj;
  }

  export function build(links: any[]) {
    return links
      .map((link) => {
        if (!link)
          throw new VisibleError(
            "An undefined link was passed into a `link` array.",
          );
        return link;
      })
      .filter((l) => isLinkable(l))
      .map((l: Linkable) => {
        const link = l.getNEOLink();
        return all([l.urn, link]).apply(([urn, link]) => ({
          name: urn.split("::").at(-1)!,
          properties: {
            ...link.properties,
            type: urn.split("::").at(-2),
          },
        }));
      });
  }

  export function getProperties(links?: Input<any[]>) {
    const linkProperties = output(links ?? []).apply((links) =>
      links
        .map((link) => {
          if (!link)
            throw new VisibleError(
              "An undefined link was passed into a `link` array.",
            );
          return link;
        })
        .filter((l) => isLinkable(l))
        .map((l: Linkable) => ({
          urn: l.urn,
          properties: l.getNEOLink().properties,
        })),
    );

    return output(linkProperties).apply((e) =>
      Object.fromEntries(
        e.map(({ urn, properties }) => {
          const name = urn.split("::").at(-1)!;
          const data = {
            ...properties,
            type: urn.split("::").at(-2),
          };
          return [name, data];
        }),
      ),
    );
  }

  export function propertiesToEnv(
    properties: ReturnType<typeof getProperties>,
  ) {
    return output(properties).apply((properties) => {
      const env = Object.fromEntries(
        Object.entries(properties).map(([key, value]) => {
          return [`NEO_RESOURCE_${key}`, JSON.stringify(value)];
        }),
      );
      env["NEO_RESOURCE_App"] = JSON.stringify({
        name: $app.name,
        stage: $app.stage,
      });
      return env;
    });
  }

  export function getInclude<T>(
    type: string,
    input?: Input<any[]>,
  ): Output<T[]> {
    if (!input) return output([]);
    return output(input).apply((links) => {
      return links.filter(isLinkable).flatMap((l: Linkable) => {
        const link = l.getNEOLink();
        return (link.include || []).filter((i) => i.type === type) as T[];
      });
    });
  }

  /** @deprecated
   * Use neo.Linkable.wrap instead.
   */
  export function linkable<T>(
    obj: { new (...args: any[]): T },
    cb: (resource: T) => Definition,
  ) {
    console.warn("neo.linkable is deprecated. Use neo.Linkable.wrap instead.");
    obj.prototype.getNEOLink = function () {
      return cb(this);
    };
  }
}
