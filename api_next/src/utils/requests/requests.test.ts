import { sum } from "./requests";
import { test, expect } from "@jest/globals";

test("1+2=3", () => {
  expect(sum(1, 2)).toBe(3);
});
