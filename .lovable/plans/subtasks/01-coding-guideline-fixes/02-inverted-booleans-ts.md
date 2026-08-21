# Subtask: Inverted Booleans (TS)
Status: pending

## Steps
1. `src/components/docs/CloneNextCommandBuilder.tsx:107`: Extract `!s.flatten` to `const isNoFlatten = !s.flatten; if (isNoFlatten) parts.push("--no-flatten");`
2. `src/components/docs/TabOrderMap.tsx:88`: Extract `!inViewport` to `const isOutsideViewport = !inViewport; if (isOutsideViewport) ...`
3. `src/components/docs/TabOrderMap.tsx:104`: Extract `!ids` to `const isIdsMissing = !ids; if (isIdsMissing) return "";`
4. `src/components/docs/TabOrderMap.tsx:211`: Extract `!aPositive` to `const isANotPositive = !aPositive; if (isANotPositive && bPositive) return 1;`
5. `src/components/docs/TabOrderMap.tsx:258`: Extract `!open` to `const isClosed = !open; if (isClosed) return;`
6. `src/components/docs/TabOrderMap.tsx:285`: Extract `!open` to `const isClosed = !open; if (isClosed) return;`
7. `src/components/docs/TabOrderMap.tsx:288`: Extract `!target` to `const isMissingTarget = !target; if (isMissingTarget) return;`
8. `src/components/docs/TabOrderMap.tsx:299`: Extract `!active` to `const isInactive = !active; if (isInactive || active === document.body) return;`
9. `src/components/projects/ProjectDetailDialog.tsx:47`: Extract `!project` to `const isMissingProject = !project; if (isMissingProject) return null;`
10. `src/components/ui/carousel.tsx:39`: Extract `!context` to `const isMissingContext = !context; if (isMissingContext) throw ...`
11. `src/components/ui/carousel.tsx:59`: Extract `!api` to `const isMissingApi = !api; if (isMissingApi) return;`
12. `src/components/ui/carousel.tsx:89`: Extract `!api || !setApi` to `const isUninitialized = !api || !setApi; if (isUninitialized) return;`
13. `src/components/ui/chart.tsx:25`: Extract `!context` to `const isMissingContext = !context; if (isMissingContext) throw ...`
14. `src/components/ui/chart.tsx:64`: Extract `!colorConfig.length` to `const isConfigEmpty = !colorConfig.length; if (isConfigEmpty) return null;`
15. `src/components/ui/chart.tsx:146`: Extract `!value` to `const isMissingValue = !value; if (isMissingValue) return null;`
16. `src/components/ui/chart.tsx:153`: Extract `!active || !payload?.length` to `const isInactive = !active || !payload?.length; if (isInactive) return null;`
17. `src/components/ui/form.tsx:40`: Extract `!fieldContext` to `const isMissingContext = !fieldContext; if (isMissingContext) throw ...`
18. `src/components/ui/form.tsx:116`: Extract `!body` to `const isMissingBody = !body; if (isMissingBody) return null;`
19. `src/components/ui/sidebar.tsx:41`: Extract `!context` to `const isMissingContext = !context; if (isMissingContext) throw ...`
20. `src/components/ui/sidebar.tsx:480`: Extract `!tooltip` to `const isMissingTooltip = !tooltip; if (isMissingTooltip) return ...`
