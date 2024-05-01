export default function Workspace({ params }: { params: { user: string } }) {
  const { user } = params;

  // TODO: ban creation of a user called "guest"

  return <div>{user}&apos;s Workspace</div>;
}
