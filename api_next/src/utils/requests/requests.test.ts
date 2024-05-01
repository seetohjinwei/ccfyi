import { get } from "./requests";
import { test, expect } from "@jest/globals";

test("simple get", async () => {
  expect(
    await get("https://jsonplaceholder.typicode.com/todos/1")
  ).toBeTruthy();
});
