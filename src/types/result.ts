/**
 * Strongly-typed Result<T> envelope and pattern matching helpers for TypeScript.
 *
 * Implements the Universal Result Envelope architecture:
 * - Discriminated union of success and failure variants
 * - Total elimination of raw throw / catch in business domain logic
 * - Exhaustive pattern matching via assertNever
 */

export interface AppError {
  readonly code: string;
  readonly message: string;
  readonly op?: string;
  readonly cause?: unknown;
}

export type Result<T> =
  | { readonly isSuccess: true; readonly isFailed: false; readonly value: T; readonly error: null }
  | { readonly isSuccess: false; readonly isFailed: true; readonly value: null; readonly error: AppError };

export function successResult<T>(value: T): Result<T> {
  return {
    isSuccess: true,
    isFailed: false,
    value,
    error: null,
  };
}

export function failureResult<T>(error: AppError): Result<T> {
  return {
    isSuccess: false,
    isFailed: true,
    value: null,
    error,
  };
}

export function assertNever(x: never): never {
  throw new Error(`Unexpected object in exhaustive check: ${JSON.stringify(x)}`);
}
