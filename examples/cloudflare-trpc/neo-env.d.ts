import "neo"
declare module "neo" {
  export interface Resource {
    Trpc: import("@cloudflare/workers-types").Service
  }
}
export {}