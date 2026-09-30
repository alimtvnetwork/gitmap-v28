# TypeScript State Management

## 1. Zustand for Global State

- Prefer Zustand over Redux for global state management due to its minimal boilerplate and predictable rendering.
- Keep Zustand stores small and focused. Do not create a single monolithic store; split them by domain feature (e.g., `useAuthStore`, `useUiStore`).

## 2. Local State vs Global State

- **Rule of Thumb:** If the state is only used by a single component and its direct children, use `useState` or `useReducer`.
- Only elevate state to Zustand if it must be accessed by completely decoupled components across the app.

## 3. Immutability

- Never mutate state directly in Zustand actions. Always return a new object or use a library like Immer to produce draft mutations safely.
