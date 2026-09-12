import { AppError } from "@/types/result";

/**
 * Generic query wrapper result structure.
 * Returns structured state with `isSuccess`, `isFail`, `isFailed`, and `hasError`.
 */
export interface QueryResult<T> {
  data: T | null;
  error: Error | AppError | null;
  isSuccess: boolean;
  isFail: boolean;
  isFailed: boolean;
  hasError: boolean;
}

function buildSuccess<T>(data: T): QueryResult<T> {
  return {
    data,
    error: null,
    isSuccess: true,
    isFail: false,
    isFailed: false,
    hasError: false,
  };
}

function buildFailure<T>(err: unknown): QueryResult<T> {
  const error = err instanceof Error ? err : new Error(String(err));
  console.error("[QueryWrapper Error]:", error);

  return {
    data: null,
    error,
    isSuccess: false,
    isFail: true,
    isFailed: true,
    hasError: true,
  };
}

/**
 * Wraps an async API or database operation, catching exceptions,
 * logging errors automatically, and returning a structured QueryResult.
 */
export async function queryWrapper<T>(
  operation: () => Promise<T>
): Promise<QueryResult<T>> {
  try {
    const data = await operation();

    return buildSuccess(data);
  } catch (err) {
    return buildFailure<T>(err);
  }
}

/**
 * Wraps a synchronous operation, catching exceptions,
 * logging errors automatically, and returning a structured QueryResult.
 */
export function queryWrapperSync<T>(
  operation: () => T
): QueryResult<T> {
  try {
    const data = operation();

    return buildSuccess(data);
  } catch (err) {
    return buildFailure<T>(err);
  }
}

export default queryWrapper;
