# TypeScript State Management (AI Execution Prompt)

> **/goal** Architect modular Zustand stores, enforce immutable state updates, and establish strict separation between local and global state.
> **/learn** Master state management patterns: isolate stores by domain feature, utilize typed selectors to prevent unnecessary re-renders, and preserve immutability in actions.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Split global state into small, domain-focused Zustand stores instead of a single monolithic store.
- [ ] `/learn` Never elevate component-local state to global stores unless shared across decoupled tree branches.
- [ ] `/goal` Guarantee immutable state updates in all store actions and dispatchers.
- [ ] `/learn` Validate compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

## 1. Zustand for Global State

- Prefer Zustand over Redux for global state management due to its minimal boilerplate and predictable rendering.
- Keep Zustand stores small and focused. Do not create a single monolithic store; split them by domain feature (e.g., `useAuthStore`, `useUiStore`).

## 2. Local State vs Global State

- **Rule of Thumb:** If the state is only used by a single component and its direct children, use `useState` or `useReducer`.
- Only elevate state to Zustand if it must be accessed by completely decoupled components across the app.

## 3. Immutability

- Never mutate state directly in Zustand actions. Always return a new object or use a library like Immer to produce draft mutations safely.

```typescript
// ❌ FORBIDDEN: Direct state mutation
updateProfile: (name: string) => {
  set((state) => {
    state.user.name = name;
    return state;
  });
};

// ✅ REQUIRED: Immutable state return
updateProfile: (name: string) => {
  set((state) => ({
    user: {
      ...state.user,
      name,
    },
  }));
};
```

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/02-typescript/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-TS-015: TypeScript State Management, Stores, and Reactive Architecture

**Given** TypeScript source code under implementation or review.
**When** Codebases are audited against TypeScript language standards.
**Then** Global state architectures enforce domain-bounded stores and strict immutability with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/02-typescript --check-only
```
**Expected:** exit 0. Zero violations.
