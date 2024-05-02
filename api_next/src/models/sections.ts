export const sections = ["collections", "environments"] as const;
export type Section = (typeof sections)[number];
