/* tslint:disable */
/* eslint-disable */
import "neo";
declare module "neo" {
  export interface Resource {
    Database: {
      database: string;
      host: string;
      password: string;
      port: number;
      type: "supabase.index/project.Project";
      user: string;
    };
  }
}
export {};
