import axios from "axios";

export interface RequestParameters {
  query_params: QueryParam[]; // need to perform url encoding
  path_params: PathParam[];
  body: any; // need to stringify
}

const empty_request_parameters: RequestParameters = {
  query_params: [],
  path_params: [],
  body: null,
};

export interface QueryParam {
  key: string;
  value: string | null;
}

export interface PathParam {
  key: string;
  value: string;
}

export async function get<T>(
  url: string,
  params: RequestParameters = empty_request_parameters,
): Promise<T> {
  // TODO: do stuff with said params
  return await axios.get<T>(url).then((res) => res.data);
}
