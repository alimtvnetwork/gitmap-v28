import { queryWrapperSync } from "./queryWrapper";

export const THEME_STORAGE_KEY = "gitmap-theme";

export enum ThemeType {
  Light = "light",
  Dark = "dark",
}
export type Theme = ThemeType;

/** Read the currently applied theme from the <html> element. */
export function getCurrentTheme(): ThemeType {
  if (typeof document === "undefined") return ThemeType.Dark;
  if (document.documentElement.classList.contains("light")) return ThemeType.Light;
  if (document.documentElement.classList.contains("dark")) return ThemeType.Dark;

  const res = queryWrapperSync(() => localStorage.getItem(THEME_STORAGE_KEY));
  if (!res.isFail && (res.data === "light" || res.data === "dark")) {
    return res.data as ThemeType;
  }
  return ThemeType.Dark;
}

/** Apply a theme to the <html> element AND persist it to localStorage. */
export function setTheme(theme: ThemeType): void {
  if (typeof document === "undefined") return;
  document.documentElement.classList.toggle("dark", theme === ThemeType.Dark);
  document.documentElement.classList.toggle("light", theme === ThemeType.Light);
  queryWrapperSync(() => localStorage.setItem(THEME_STORAGE_KEY, theme));
}

/** Toggle between light and dark, persisting the new value. */
export function toggleTheme(): ThemeType {
  const next: ThemeType = getCurrentTheme() === ThemeType.Dark ? ThemeType.Light : ThemeType.Dark;
  setTheme(next);
  return next;
}
