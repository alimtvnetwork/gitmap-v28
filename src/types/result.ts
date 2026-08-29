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
  | {
      readonly isSuccess: true;
      readonly isFailed: false;
      readonly hasError: false;
      readonly value: T;
      readonly data: T;
      readonly error: null;
      unwrap(): [T, null];
      unwrapOr(defaultVal: T): T;
    }

  | {
      readonly isSuccess: false;
      readonly isFailed: true;
      readonly hasError: true;
      readonly value: null;
      readonly data: null;
      readonly error: AppError;
      unwrap(): [null, AppError];
      unwrapOr(defaultVal: T): T;
    };

export function successResult<T>(value: T): Result<T> {
  return {
    isSuccess: true,
    isFailed: false,
    hasError: false,
    value,
    data: value,
    error: null,
    unwrap: () => [value, null],
    unwrapOr: () => value,
  };
}

export function failureResult<T>(error: AppError): Result<T> {
  return {
    isSuccess: false,
    isFailed: true,
    hasError: true,
    value: null,
    data: null,
    error,
    unwrap: () => [null, error],
    unwrapOr: (defaultVal: T) => defaultVal,
  };
}

export function newFailure<T>(
  code: string,
  message: string,
  op?: string,
): Result<T> {
  return failureResult<T>({ code, message, op });
}

export function assertNever(x: never): never {
  throw new Error(`Unexpected object in exhaustive check: ${JSON.stringify(x)}`);
}
