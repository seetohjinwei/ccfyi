import TestButton from "@/components/TestButton";
import { Button } from "@mantine/core";

export default function Workspace({ params }: { params: { user: string } }) {
  const { user } = params;

  // TODO: ban creation of a user called "guest"

  return (
    <div>
      <h1>{user}&apos;s Workspace</h1>
      <TestButton />
    </div>
  );
}
