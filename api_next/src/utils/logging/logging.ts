import { StringUnion } from "../types/types";

// type level = (typeof levels)[number];
type level = StringUnion<["info", "debug", "error"]>;

class Logger {
  // TODO: make this good
  level: level;

  constructor(level: level) {
    this.level = level;
  }

  info(...data: any) {
    console.info("INFO", ...data);
  }

  debug(...data: any) {
    console.debug("DEBUG", ...data);
  }

  error(...data: any) {
    console.error("ERROR", ...data);
  }
}

// TODO: take from env var
const level: level = "info";

export const logger = new Logger(level);
