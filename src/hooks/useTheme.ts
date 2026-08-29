// Shared theme hook — single source of truth used by every component that
// needs to read or change the active light/dark palette. All instances stay
// in sync via a custom DOM event, and the hook also reacts to OS-level
// prefers-color-scheme changes when the user has no saved preference.

import { useCallback, useEffect, useState } from "react";
import {
  getCurrentTheme,
  setTheme as persistTheme,
  THEME_STORAGE_KEY,
  type Theme,
  ThemeType,
} from "@/lib/theme";
import { queryWrapperSync } from "@/lib/queryWrapper";

const THEME_CHANGE_EVENT = "gitmap:theme-change";
const THEME_SOURCE_EVENT = "gitmap:theme-source-change";

export enum ThemeSourceType {
  System = "system",
  User = "user",
}

export type ThemeSource = ThemeSourceType;

interface UseThemeResult {
  theme: Theme;
  isDark: boolean;
  source: ThemeSourceType;
  isSystem: boolean;
  setTheme: (theme: Theme) => void;
  toggleTheme: () => void;
}

function readSource(): ThemeSourceType {
  const res = queryWrapperSync(() => localStorage.getItem(THEME_STORAGE_KEY));

  if (res.isFail) return ThemeSourceType.System;
  const stored = res.data;

  return stored === "light" || stored === "dark" ? ThemeSourceType.User : ThemeSourceType.System;
}

export function useTheme(): UseThemeResult {
  const [theme, setThemeState] = useState<Theme>(() => getCurrentTheme());
  const [source, setSourceState] = useState<ThemeSourceType>(() => readSource());

  // Sync across hook instances + cross-tab + OS-level changes.
  useEffect(() => {
    const handleThemeEvent = (event: Event) => {
      const next = (event as CustomEvent<Theme>).detail;

      if (next === "light" || next === "dark") setThemeState(next);
    };

    const handleSourceEvent = (event: Event) => {
      const next = (event as CustomEvent<ThemeSourceType>).detail;

      if (next === "system" || next === "user") setSourceState(next);
    };

    const handleStorage = (event: StorageEvent) => {
      if (event.key !== THEME_STORAGE_KEY) return;

      if (event.newValue === "light" || event.newValue === "dark") {
        setThemeState(event.newValue as ThemeType);
        setSourceState(ThemeSourceType.User);
        persistTheme(event.newValue as ThemeType);
      } else if (event.newValue === null) {
        setSourceState(ThemeSourceType.System);
      }
    };

    const mediaQuery = window.matchMedia("(prefers-color-scheme: light)");
    const handleSystemChange = (event: MediaQueryListEvent) => {
      // Only follow the OS when the user hasn't explicitly chosen a theme.
      const checkRes = queryWrapperSync(() => localStorage.getItem(THEME_STORAGE_KEY));
      const hasUserPreference = checkRes.isSuccess && Boolean(checkRes.data);

      if (hasUserPreference) {
        return;
      }

      const next: ThemeType = event.matches ? ThemeType.Light : ThemeType.Dark;
      setThemeState(next);
      setSourceState(ThemeSourceType.System);
      document.documentElement.classList.toggle("dark", next === ThemeType.Dark);
      document.documentElement.classList.toggle("light", next === ThemeType.Light);
    };

    window.addEventListener(THEME_CHANGE_EVENT, handleThemeEvent);
    window.addEventListener(THEME_SOURCE_EVENT, handleSourceEvent);
    window.addEventListener("storage", handleStorage);
    mediaQuery.addEventListener("change", handleSystemChange);

    return () => {
      window.removeEventListener(THEME_CHANGE_EVENT, handleThemeEvent);
      window.removeEventListener(THEME_SOURCE_EVENT, handleSourceEvent);
      window.removeEventListener("storage", handleStorage);
      mediaQuery.removeEventListener("change", handleSystemChange);
    };
  }, []);

  const setTheme = useCallback((next: Theme) => {
    persistTheme(next);
    setThemeState(next);
    setSourceState(ThemeSourceType.User);
    window.dispatchEvent(new CustomEvent<Theme>(THEME_CHANGE_EVENT, { detail: next }));
    window.dispatchEvent(
      new CustomEvent<ThemeSourceType>(THEME_SOURCE_EVENT, { detail: ThemeSourceType.User }),
    );
  }, []);

  const toggleTheme = useCallback(() => {
    setTheme(theme === ThemeType.Dark ? ThemeType.Light : ThemeType.Dark);
  }, [theme, setTheme]);

  return {
    theme,
    isDark: theme === ThemeType.Dark,
    source,
    isSystem: source === ThemeSourceType.System,
    setTheme,
    toggleTheme,
  };
}
