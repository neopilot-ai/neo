import "neo"
declare module "neo" {
  export interface Resource {
    Worker: {
      type: "neo.cloudflare.Worker"
      url: string
    }
  }
}
// cloudflare 
declare module "neo" {
  export interface Resource {
    MyDatabase: import("@cloudflare/workers-types").D1Database
  }
}
