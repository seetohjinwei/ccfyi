import axios from "axios";

export interface RequestParameters {
  path_params: PathParam[];
  query_params: QueryParam[]; // need to perform url encoding
  body: any; // need to stringify
}

const empty_request_parameters: RequestParameters = {
  path_params: [],
  query_params: [],
  body: null,
};

export interface PathParam {
  key: string;
  value: string;
}

export interface QueryParam {
  key: string;
  value: string | null;
}

export async function get<T>(
  url: string,
  params: RequestParameters = empty_request_parameters,
): Promise<T> {
  // TODO: do stuff with said params
  return await axios.get<T>(url).then((res) => res.data);
}
