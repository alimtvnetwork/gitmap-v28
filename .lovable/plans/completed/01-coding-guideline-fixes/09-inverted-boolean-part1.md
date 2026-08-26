# Fix inverted boolean (Part 1)

Total items: 16

## Files to Modify

- `.\src\components\docs\TabOrderMap.tsx:93`: `const hasNoRects = rects.length === 0;`
- `.\src\components\docs\TabOrderMap.tsx:94`: `if (hasNoRects) return false;`
- `.\src\components\ui\carousel.tsx:180`: `const cannotScrollPrev = !canScrollPrev;`
- `.\src\components\ui\carousel.tsx:194`: `disabled={cannotScrollPrev}`
- `.\src\components\ui\carousel.tsx:209`: `const cannotScrollNext = !canScrollNext;`
- `.\src\components\ui\carousel.tsx:223`: `disabled={cannotScrollNext}`
- `.\src\components\ui\chart.tsx:80`: `const hasNoColors = !colorConfig.length;`
- `.\src\components\ui\chart.tsx:82`: `if (hasNoColors) {`
- `.\src\components\ui\chart.tsx:133`: `const hasNoPayload = !payload?.length;`
- `.\src\components\ui\chart.tsx:134`: `if (hideLabel `
- `.\src\components\ui\chart.tsx:159`: `const hasNoPayload = !payload?.length;`
- `.\src\components\ui\chart.tsx:160`: `if (isInactive `
- `.\src\components\ui\chart.tsx:254`: `const hasNoPayload = !payload?.length;`
- `.\src\components\ui\chart.tsx:256`: `if (hasNoPayload) {`
- `.\src\components\ui\form.tsx:88`: `const hasNoError = !error;`
- `.\src\components\ui\form.tsx:89`: `const describedBy = hasNoError ? `${formDescriptionId}` : `${formDescriptionId} ${formMessageId}`;`
