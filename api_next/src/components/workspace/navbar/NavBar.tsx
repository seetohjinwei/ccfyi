"use client";

import { Collection } from "@/models/collection";
import { Environment } from "@/models/environment";
import { toTitleCase } from "@/utils/strings/cases";
import { SegmentedControl, Text, TextInput, rem } from "@mantine/core";
import {
  IconLogout,
  IconSearch,
  IconSwitchHorizontal,
} from "@tabler/icons-react";
import { ChangeEventHandler, useState } from "react";
import classes from "./NavBar.module.css";
import { Section, sections } from "@/models/sections";

interface Item {
  label: string;
  active: boolean;
}

interface Props {
  workspace: string;
  collections: Collection[];
  environments: Environment[];
  setSelection: (type: Section, key: string) => void;
}

export default function NavBar({
  workspace,
  collections,
  environments,
  setSelection,
}: Props) {
  const [section, setSection] = useState<Section>("collections");
  const [activeItem, setActiveItem] = useState<string | undefined>(undefined);
  const [searchValue, setSearchValue] = useState<string>("");

  const setChosen = (key: string) => {
    setActiveItem(key);
    setSelection(section, key);
  };

  const handleSearchValue: ChangeEventHandler<HTMLInputElement> = (e) => {
    setSearchValue(e.currentTarget.value);
  };

  function ItemComponent({ item }: { item: Item }) {
    return (
      <a
        className={classes.link}
        data-active={item.label === activeItem}
        onClick={(event) => {
          event.preventDefault();
          setChosen(item.label);
        }}
      >
        <span>{item.label}</span>
      </a>
    );
  }

  const items: Item[] = {
    collections: collections.map((c) => c.collection),
    environments: environments.map((e) => e.environment),
  }[section]
    .map((item: string) => ({
      label: item,
      active: true,
    }))
    .filter((item) => searchValue === "" || item.label.includes(searchValue));

  return (
    <nav className={classes.navbar}>
      <div>
        <Text fw={500} size="sm" className={classes.title} c="dimmed" mb="xs">
          Workspace {workspace}
        </Text>

        <TextInput
          placeholder="Search"
          size="xs"
          value={searchValue}
          onChange={handleSearchValue}
          leftSection={
            <IconSearch
              style={{ width: rem(12), height: rem(12) }}
              stroke={1.5}
            />
          }
          styles={{ section: { pointerEvents: "none" } }}
          mb="sm"
        />

        <SegmentedControl
          value={section}
          onChange={(value: any) => setSection(value)}
          transitionTimingFunction="ease"
          fullWidth
          data={sections.map((s) => ({
            label: toTitleCase(s),
            value: s,
          }))}
        />
      </div>

      <div className={classes.navbarMain}>
        {items.map((item, key) => (
          <ItemComponent key={key} item={item} />
        ))}
        {items.length === 0 && searchValue === "" && (
          <p>Nothing found in your {section}!</p>
        )}
        {items.length === 0 && searchValue !== "" && (
          <p>Nothing matching your search was found in your {section}!</p>
        )}
      </div>

      <div className={classes.footer}>
        <a
          href="#"
          className={classes.link}
          onClick={(event) => event.preventDefault()}
        >
          <IconSwitchHorizontal className={classes.linkIcon} stroke={1.5} />
          <span>Change account</span>
        </a>

        <a
          href="#"
          className={classes.link}
          onClick={(event) => event.preventDefault()}
        >
          <IconLogout className={classes.linkIcon} stroke={1.5} />
          <span>Logout</span>
        </a>
      </div>
    </nav>
  );
}

// UI modified from: https://ui.mantine.dev/component/navbar-segmented/
