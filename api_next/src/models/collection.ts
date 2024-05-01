import { Settings } from "./settings";

export interface Collection {
  workspace: string;
  collection: string;
  apis: API[];
}

export interface API {
  api: string;
  method: string;
  url: string; // does **not** store the query params
  description: string;
  pre_req: string; // stores javascript code
  post_req: string; // stores javascript code
  query_params: QueryParams;
  path_params: PathParams;
  headers: Headers;
  body: any;
  settings: Settings;
}

export interface QueryParams {}
export interface PathParams {}
export interface Headers {}
