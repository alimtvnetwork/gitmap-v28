# Error Queue Navigation

> **Parent:** [Error Modal Reference](./00-overview.md)  
> **Version:** 2.2.0  
> **Updated:** 2026-03-31

---

Multiple concurrent errors are queued and navigable:

```tsx
// In ErrorStore:
errorQueue: CapturedError[];
currentQueueIndex: number;
navigateQueue: (direction: 'prev' | 'next') => void;
getQueuedErrorsMarkdown: () => string;  // All errors as one Markdown doc

// In GlobalErrorModal header:
{hasMultipleErrors && (
  <div className="flex items-center gap-1">
    <Button onClick={() => navigateQueue('prev')}>
      <ChevronLeft />
    </Button>
    <Badge>{currentQueueIndex + 1}/{errorQueue.length}</Badge>
    <Button onClick={() => navigateQueue('next')}>
      <ChevronRight />
    </Button>
    <Button onClick={copyAllErrors}>
      <CopyPlus /> All
    </Button>
  </div>
)}
```

---

*Queue navigation — updated: 2026-03-31*
