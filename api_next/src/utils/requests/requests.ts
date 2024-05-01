import axios from "axios";

export async function get<T>(url: string): Promise<T> {
  return await axios.get<T>(url).then((res) => res.data);
}
