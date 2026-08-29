import * as React from "react";

import type { ToastActionElement, ToastProps } from "@/components/ui/toast";

const TOAST_LIMIT = 1;
const TOAST_REMOVE_DELAY = 1000000;

type ToasterToast = ToastProps & {
  id: string;
  title?: React.ReactNode;
  description?: React.ReactNode;
  action?: ToastActionElement;
};

const actionTypes = {
  ADD_TOAST: "ADD_TOAST",
  UPDATE_TOAST: "UPDATE_TOAST",
  DISMISS_TOAST: "DISMISS_TOAST",
  REMOVE_TOAST: "REMOVE_TOAST",
} as const;

let count = 0;

function genId() {
  count = (count + 1) % Number.MAX_SAFE_INTEGER;

  return count.toString();
}

type ActionType = typeof actionTypes;

type Action =
  | {
      type: ActionType["ADD_TOAST"];
      toast: ToasterToast;
    }

  | {
      type: ActionType["UPDATE_TOAST"];
      toast: Partial<ToasterToast>;
    }

  | {
      type: ActionType["DISMISS_TOAST"];
      toastId?: ToasterToast["id"];
    }

  | {
      type: ActionType["REMOVE_TOAST"];
      toastId?: ToasterToast["id"];
    };

interface State {
  toasts: ToasterToast[];
}

const toastTimeouts = new Map<string, ReturnType<typeof setTimeout>>();

const addToRemoveQueue = (toastId: string) => {
  const isQueued = toastTimeouts.has(toastId);

  if (isQueued) return;

  const timeout = setTimeout(() => {
    toastTimeouts.delete(toastId);
    dispatch({ type: "REMOVE_TOAST", toastId });
  }, TOAST_REMOVE_DELAY);

  toastTimeouts.set(toastId, timeout);
};

function dismissToasts(state: State, toastId?: string): State {
  if (toastId) {
    addToRemoveQueue(toastId);
  } else {
    state.toasts.forEach((t) => addToRemoveQueue(t.id));
  }

  return {
    ...state,
    toasts: state.toasts.map((t) => {
      const isTarget = t.id === toastId || toastId === undefined;

      if (isTarget) return { ...t, open: false };

      return t;
    }),
  };
}

function removeToasts(state: State, toastId?: string): State {
  const isClearAll = toastId === undefined;

  if (isClearAll) {
    return { ...state, toasts: [] };
  }

  return {
    ...state,
    toasts: state.toasts.filter((t) => t.id !== toastId),
  };
}

export const reducer = (state: State, action: Action): State => {
  if (action.type === "ADD_TOAST") {
    return { ...state, toasts: [action.toast, ...state.toasts].slice(0, TOAST_LIMIT) };
  }

  if (action.type === "UPDATE_TOAST") {
    return {
      ...state,
      toasts: state.toasts.map((t) => (t.id === action.toast.id ? { ...t, ...action.toast } : t)),
    };
  }

  if (action.type === "DISMISS_TOAST") return dismissToasts(state, action.toastId);

  if (action.type === "REMOVE_TOAST") return removeToasts(state, action.toastId);

  return state;
};

const listeners: Array<(state: State) => void> = [];

let memoryState: State = { toasts: [] };

function dispatch(action: Action) {
  memoryState = reducer(memoryState, action);
  listeners.forEach((listener) => {
    listener(memoryState);
  });
}

type Toast = Omit<ToasterToast, "id">;

function createToastItem(props: Toast, id: string, dismiss: () => void): ToasterToast {
  return {
    ...props,
    id,
    open: true,
    onOpenChange: (open) => {
      const isClosed = !open;

      if (isClosed) dismiss();
    },
  };
}

function toast({ ...props }: Toast) {
  const id = genId();
  const update = (p: ToasterToast) => dispatch({ type: "UPDATE_TOAST", toast: { ...p, id } });
  const dismiss = () => dispatch({ type: "DISMISS_TOAST", toastId: id });

  dispatch({ type: "ADD_TOAST", toast: createToastItem(props, id, dismiss) });

  return { id, dismiss, update };
}

function useToast() {
  const [state, setState] = React.useState<State>(memoryState);

  React.useEffect(() => {
    listeners.push(setState);

    return () => {
      const index = listeners.indexOf(setState);
      const isPresent = index > -1;

      if (isPresent) listeners.splice(index, 1);
    };
  }, [state]);

  return {
    ...state,
    toast,
    dismiss: (toastId?: string) => dispatch({ type: "DISMISS_TOAST", toastId }),
  };
}

export { useToast, toast };
