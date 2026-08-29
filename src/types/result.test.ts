import { describe, expect, it } from "vitest";
import { failureResult, newFailure, successResult, type Result } from "./result";

describe("TypeScript Result<T> Envelope", () => {
  it("creates a successful Result envelope", () => {
    const res: Result<string> = successResult("gitmap");
    expect(res.isSuccess).toBe(true);
    expect(res.isFailed).toBe(false);
    expect(res.hasError).toBe(false);
    expect(res.value).toBe("gitmap");
    expect(res.data).toBe("gitmap");
    expect(res.error).toBeNull();
    expect(res.unwrap()).toEqual(["gitmap", null]);
    expect(res.unwrapOr("fallback")).toBe("gitmap");
  });

  it("creates a failed Result envelope with AppError", () => {
    const res: Result<number> = failureResult({
      code: "E_NOT_FOUND",
      message: "Item not found",
    });

    expect(res.isSuccess).toBe(false);
    expect(res.isFailed).toBe(true);
    expect(res.hasError).toBe(true);
    expect(res.value).toBeNull();
    expect(res.data).toBeNull();
    expect(res.error?.code).toBe("E_NOT_FOUND");
    expect(res.unwrap()).toEqual([null, { code: "E_NOT_FOUND", message: "Item not found" }]);
    expect(res.unwrapOr(42)).toBe(42);
  });

  it("creates a failed Result using newFailure helper", () => {
    const res = newFailure<string>("E_VALIDATION", "Invalid input", "testOp");
    expect(res.isFailed).toBe(true);
    expect(res.error?.code).toBe("E_VALIDATION");
    expect(res.error?.op).toBe("testOp");
  });
});
