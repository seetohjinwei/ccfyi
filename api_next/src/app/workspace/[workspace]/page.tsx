"use client";

import NavBar from "@/components/workspace/navbar/NavBar";
import { Collection } from "@/models/collection";
import { Environment } from "@/models/environment";
import { Section } from "@/models/sections";
import { logger } from "@/utils/logging/logging";
import { useState } from "react";
import classes from "./page.module.css";
import None from "@/components/workspace/mainpane/none/None";

// TODO: get real data
const collections: Collection[] = [
  { workspace: "jinwei", collection: "default", apis: [] },
  { workspace: "jinwei", collection: "api_next", apis: [] },
];
const environments: Environment[] = [
  { workspace: "jinwei", environment: "default", environment_variables: [] },
  { workspace: "jinwei", environment: "test", environment_variables: [] },
];

interface Selection {
  type: Section | "none";
  key: string;
}

export default function Workspace({
  params,
}: {
  params: { workspace: string };
}) {
  const { workspace } = params;

  const [selection, setSelection] = useState<Selection>({
    type: "none",
    key: "",
  });

  const handleSelection = async (type: Section, key: string) => {
    // "use server";

    logger.debug(type, key);
    setSelection({ type, key });

    // need to fetch some stuff from database, then return to other main panel
  };

  const mainPane: JSX.Element = {
    none: <None />,
    collections: <div>TODO collections</div>,
    environments: <div>TODO environments</div>,
  }[selection.type];

  // TODO: ban creation of a user called "guest"

  return (
    <div className={classes.container}>
      <NavBar
        workspace={workspace}
        collections={collections}
        environments={environments}
        setSelection={handleSelection}
      />
      {mainPane}
    </div>
  );
}
