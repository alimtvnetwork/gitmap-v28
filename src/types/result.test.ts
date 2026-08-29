import { describe, expect, it } from "vitest";
import { failureResult, successResult, type Result } from "./result";

describe("TypeScript Result<T> Envelope", () => {
  it("creates a successful Result envelope", () => {
    const res: Result<string> = successResult("gitmap");
    expect(res.isSuccess).toBe(true);
    expect(res.isFailed).toBe(false);
    expect(res.value).toBe("gitmap");
    expect(res.error).toBeNull();
  });

  it("creates a failed Result envelope with AppError", () => {
    const res: Result<number> = failureResult({
      code: "E_NOT_FOUND",
      message: "Item not found",
    });

    expect(res.isSuccess).toBe(false);
    expect(res.isFailed).toBe(true);
    expect(res.value).toBeNull();
    expect(res.error?.code).toBe("E_NOT_FOUND");
  });
});
