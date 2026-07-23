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
} from "@/lib/theme";

const THEME_CHANGE_EVENT = "gitmap:theme-change";
const THEME_SOURCE_EVENT = "gitmap:theme-source-change";

export type ThemeSource = "system" | "user";

interface UseThemeResult {
  theme: Theme;
  isDark: boolean;
  source: ThemeSource;
  isSystem: boolean;
  setTheme: (theme: Theme) => void;
  toggleTheme: () => void;
}

function readSource(): ThemeSource {
  try {
    const stored = localStorage.getItem(THEME_STORAGE_KEY);
    return stored === "light" || stored === "dark" ? "user" : "system";
  } catch {
    return "system";
  }
}

export function useTheme(): UseThemeResult {
  const [theme, setThemeState] = useState<Theme>(() => getCurrentTheme());
  const [source, setSourceState] = useState<ThemeSource>(() => readSource());

  // Sync across hook instances + cross-tab + OS-level changes.
  useEffect(() => {
    const handleThemeEvent = (event: Event) => {
      const next = (event as CustomEvent<Theme>).detail;
      if (next === "light" || next === "dark") setThemeState(next);
    };

    const handleSourceEvent = (event: Event) => {
      const next = (event as CustomEvent<ThemeSource>).detail;
      if (next === "system" || next === "user") setSourceState(next);
    };

    const handleStorage = (event: StorageEvent) => {
      if (event.key !== THEME_STORAGE_KEY) return;
      if (event.newValue === "light" || event.newValue === "dark") {
        setThemeState(event.newValue);
        setSourceState("user");
        persistTheme(event.newValue);
      } else if (event.newValue === null) {
        setSourceState("system");
      }
    };

    const mediaQuery = window.matchMedia("(prefers-color-scheme: light)");
    const handleSystemChange = (event: MediaQueryListEvent) => {
      // Only follow the OS when the user hasn't explicitly chosen a theme.
      try {
        if (localStorage.getItem(THEME_STORAGE_KEY)) return;
      } catch {
        return;
      }
      const next: Theme = event.matches ? "light" : "dark";
      setThemeState(next);
      setSourceState("system");
      document.documentElement.classList.toggle("dark", next === "dark");
      document.documentElement.classList.toggle("light", next === "light");
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
    setSourceState("user");
    window.dispatchEvent(new CustomEvent<Theme>(THEME_CHANGE_EVENT, { detail: next }));
    window.dispatchEvent(
      new CustomEvent<ThemeSource>(THEME_SOURCE_EVENT, { detail: "user" }),
    );
  }, []);

  const toggleTheme = useCallback(() => {
    setTheme(theme === "dark" ? "light" : "dark");
  }, [theme, setTheme]);

  return {
    theme,
    isDark: theme === "dark",
    source,
    isSystem: source === "system",
    setTheme,
    toggleTheme,
  };
}
