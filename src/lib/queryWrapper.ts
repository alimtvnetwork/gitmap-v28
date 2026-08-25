/**
 * Generic query wrapper result structure.
 * Returns structured state with `isSuccess: false,
      isFail: true,` or `isSuccess: true,
      isFail: false,`.
 * DO NOT use `isSuccess` inline.
 */
export interface QueryResult<T> {
  data: T | null;
  error: Error | null;
  isSuccess: boolean;
  isFail: boolean;
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
    return {
      data,
      error: null,
      isSuccess: true,
      isFail: false,
    };
  } catch (err) {
    const error = err instanceof Error ? err : new Error(String(err));
    console.error("[QueryWrapper Error]:", error);
    return {
      data: null,
      error,
      isSuccess: false,
      isFail: true,
    };
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
    return {
      data,
      error: null,
      isSuccess: true,
      isFail: false,
    };
  } catch (err) {
    const error = err instanceof Error ? err : new Error(String(err));
    console.error("[QueryWrapper Error]:", error);
    return {
      data: null,
      error,
      isSuccess: false,
      isFail: true,
    };
  }
}

export default queryWrapper;



