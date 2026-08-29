# Instruction (must follow): Coding Guideline Execution — React & Frontend Guidelines

Trigger Keywords & Aliases: `cg-react`, `cg-execute react`, `fix frontend guidelines`, `execute react guidelines`

```text
N = 100
```

N = total self-loop steps budget. The user may override this number when triggering the prompt.

- [ ] /goal First `N/2` steps (Phase 1) are dedicated to auditing all frontend files (`src/**/*.tsx`, `src/**/*.ts`), identifying oversized components, improper `useEffect` usages, state mutations, and tuple returns, writing the master execution plan to `.lovable/plans/pending/XX-react-guidelines-audit.md`, and decomposing it into subtasks in `.lovable/plans/subtasks/XX-react-guidelines/`.
- [ ] /goal Second `N/2` steps (Phase 2) are dedicated to executing those subtasks in an autonomous self-loop until all React components and hooks strictly comply with the frontend architecture rules.
- [ ] /goal **Linter Mandate**: If an automated component size or hook rule checker is missing in CI/CD, you MUST create an advanced ESLint/script checker and connect it directly to the CI/CD local runner and workflows.
- [ ] /learn Ingest `.lovable/coding-guidelines/coding-guidelines.md`, `.lovable/strictly-avoid.md`, and `.lovable/memory/00-index.md` before touching any code.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after the user sets them. Never change them mid-execution.

---

## 1. React Frontend Non-Negotiable Checklist

You MUST audit and strictly enforce every rule below across all frontend code:

### A. Component Size & Structure
1. **Component Length Cap**: Any React component file (`.tsx`) must NOT exceed 100 lines. Extract child components, custom hooks, and helpers into separate files.
2. **Prop & Event Types**: Component prop types must live in a dedicated `types.ts` next to the component (e.g. `type ProfileCardProps = { user: User; onSave: (u: User) => void }`). Never inline anonymous object types on component signatures.

### B. `useEffect` Rules & State Management
1. **Minimize Effects**: The default count of `useEffect` is zero. Use them ONLY for synchronization with external systems (network, timers, subscriptions, DOM). Do not use effects to derive state or transform props.
2. **Positive Guard Booleans**: Guard conditions inside `useEffect` must use positively named booleans (`isReadyToSync`, `hasFreshData`). No inline negations or nested ternaries.
3. **One Effect, One Concern**: Never combine unrelated subscriptions or data fetches in a single effect.
4. **Mandatory Cleanup**: Every effect acquiring a resource (timers, subscriptions, event listeners) MUST return a cleanup function.

### C. Immutability & Render Expressions
1. **Immutable-First State**: Never mutate state, props, or arrays in-place (`.push`, `.splice`, `.sort`, `obj.x =`). Always return fresh references using spread (`{ ...prev, field: next }`), `.map`, `.filter`, or `structuredClone`.
2. **No Raw `for` Loops in Render**: Render arrays using `.map`, `.filter`, `.reduce`, or `Array.from`. Raw `for` loops in JSX are prohibited.
3. **Stable Keys**: List items must have stable, unique keys derived from domain data. Never use array index as `key`.

### D. Custom Hooks & Public API Shapes
1. **No Tuples as Public Shapes**: Custom hooks and helper functions must return named object types (e.g. `useUser(): UserQueryResult` with `{ user, isLoading, error }`), NEVER bare tuples like `[User, boolean, Error]`.
2. **Hook Naming**: All custom hooks must begin with the `use` prefix and not be called conditionally.

---

## 2. Phase 1: Planning, Audit & Subtask Decomposition (Steps 1 .. N/2)

1. **Frontend Code Audit**: Scan `src/` for `.tsx` files exceeding 100 lines, raw state mutations, tuple-returning hooks, and `useEffect` violations.
2. **Master Plan**: Write a detailed execution plan to `.lovable/plans/pending/XX-react-guidelines-audit.md`.
3. **Subtask Files**: Decompose into subtask files in `.lovable/plans/subtasks/XX-react-guidelines/` (e.g. `01-task.md`, `02-task.md`, ...).
4. **Linter Connection**: Configure ESLint and component line-count scripts in `.lovable/ai-fix-scripts/03-cicd-local-runner.py`.

---

## 3. Phase 2: Autonomous Execution Loop (Steps N/2+1 .. N)

1. **Refactor Components**: Break oversized components into sub-components, extract custom hooks returning named objects, and enforce immutability.
2. **Verify Frontend Build**: Run `npm run build` or `vite build` to ensure all TypeScript and Vite production builds pass with zero warnings.
3. **Run CI Gates**: Verify that all CI checks pass.
4. **Update Status**: Mark completed tasks as `DONE`, move completed plans to `.lovable/plans/completed/`, and update `.lovable/plans/index.md`.

---

## 4. Pre-Commit Verification Checklist

- [ ] All `.tsx` component files are <= 100 lines.
- [ ] Custom hooks return named objects (no bare tuples).
- [ ] `useEffect` uses positive guards and provides cleanup functions.
- [ ] No in-place state mutations or raw for loops in render.
- [ ] Frontend builds pass without errors in CI/CD.
