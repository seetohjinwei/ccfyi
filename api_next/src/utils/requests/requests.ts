import axios from "axios";

/** TODO:
 * in an interface?
 * path params
 * query params
 * body
 */

export async function get<T>(url: string): Promise<T> {
  return await axios.get<T>(url).then((res) => res.data);
}
