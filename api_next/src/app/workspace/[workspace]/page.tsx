import NavBar from "@/components/workspace/navbar/NavBar";
import { Collection } from "@/models/collection";
import { Environment } from "@/models/environment";
import { logger } from "@/utils/logging/logging";
import classes from "./page.module.css";

// TODO: get real data
const collections: Collection[] = [
  { workspace: "jinwei", collection: "default", apis: [] },
  { workspace: "jinwei", collection: "api_next", apis: [] },
];
const environments: Environment[] = [
  { workspace: "jinwei", environment: "default", environment_variables: [] },
  { workspace: "jinwei", environment: "test", environment_variables: [] },
];

export default function Workspace({
  params,
}: {
  params: { workspace: string };
}) {
  const { workspace } = params;

  // TODO: ban creation of a user called "guest"

  return (
    <div className={classes.container}>
      <NavBar
        workspace={workspace}
        collections={collections}
        environments={environments}
        setSelection={async (type: string, key: string) => {
          "use server";

          logger.debug(type, key);
          // need to fetch some stuff from database, then return to other main panel
        }}
      />
      <div>RHS</div>
    </div>
  );
}
