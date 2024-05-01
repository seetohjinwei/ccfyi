"use client";

import { logger } from "@/utils/logging/logging";
import { get } from "@/utils/requests/requests";
import { Button } from "@mantine/core";

export default function TestButton() {
  return (
    <Button
      onClick={async () => {
        const res = await get("https://jsonplaceholder.typicode.com/todos/1");
        logger.info(res);
      }}
    >
      Le Test Button
    </Button>
  );
}
