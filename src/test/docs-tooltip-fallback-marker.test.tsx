import { describe, it, expect, afterEach } from "vitest";
import { render, cleanup } from "@testing-library/react";
import { TooltipProvider } from "@/components/ui/tooltip";
import { DocsTooltip } from "@/components/docs/DocsTooltip";

// This suite verifies the data-docs-tooltip-fallback marker contract
// on DocsTooltip:
//   - present and "true" on the synthesized wrapper when the child
//     is NOT a single valid React element (string, fragment, multiple
//     children, null) — normalizeTrigger had to invent markup.
//   - absent everywhere when the child IS a single valid React
//     element — no wrapper was synthesized, so no marker should leak.
// The marker drives withAccessibleName's skip logic, so regressing
// it would silently re-introduce aria-label injection on invented
// markup. Keeping a dedicated test makes that contract loud.

const FALLBACK_ATTR = "data-docs-tooltip-fallback";

const renderWithProvider = (ui: React.ReactNode) =>
  render(
    <TooltipProvider delayDuration={0} skipDelayDuration={0}>
      {ui}
    </TooltipProvider>,
  );

afterEach(() => cleanup());

describe("DocsTooltip — fallback wrapper marker", () => {
  it("sets data-docs-tooltip-fallback=\"true\" on synthesized wrapper for string child", () => {
    const { container } = renderWithProvider(
      <DocsTooltip label="Hint">just text</DocsTooltip>,
    );
    const wrapper = container.querySelector(`[${FALLBACK_ATTR}]`);
    expect(wrapper).not.toBeNull();
    expect(wrapper?.getAttribute(FALLBACK_ATTR)).toBe("true");
  });

  it("sets the marker when children are multiple elements", () => {
    const { container } = renderWithProvider(
      <DocsTooltip label="Hint">
        <span>a</span>
        <span>b</span>
      </DocsTooltip>,
    );
    const wrapper = container.querySelector(`[${FALLBACK_ATTR}]`);
    expect(wrapper).not.toBeNull();
    expect(wrapper?.getAttribute(FALLBACK_ATTR)).toBe("true");
  });

  it("sets the marker when child is null", () => {
    const { container } = renderWithProvider(
      <DocsTooltip label="Hint">{null}</DocsTooltip>,
    );
    const wrapper = container.querySelector(`[${FALLBACK_ATTR}]`);
    expect(wrapper).not.toBeNull();
    expect(wrapper?.getAttribute(FALLBACK_ATTR)).toBe("true");
  });

  it("does NOT set the marker when child is a single valid element (button)", () => {
    const { container } = renderWithProvider(
      <DocsTooltip label="Hint">
        <button type="button">Click</button>
      </DocsTooltip>,
    );
    const wrapper = container.querySelector(`[${FALLBACK_ATTR}]`);
    expect(wrapper).toBeNull();
    // Sanity: the real button is rendered and has the injected aria-label.
    const button = container.querySelector("button");
    expect(button).not.toBeNull();
    expect(button?.getAttribute("aria-label")).toBe("Hint");
  });

  it("does NOT set the marker when child is a single valid element (span)", () => {
    const { container } = renderWithProvider(
      <DocsTooltip label="Hint">
        <span>icon</span>
      </DocsTooltip>,
    );
    const wrapper = container.querySelector(`[${FALLBACK_ATTR}]`);
    expect(wrapper).toBeNull();
  });
});

// Fragments are a special shape: React.Children.count treats a
// top-level fragment as a single child and isValidElement returns
// true for it. That means normalizeTrigger does NOT synthesize a
// fallback wrapper for fragment children — the fragment itself is
// passed straight to TooltipTrigger asChild. Document that:
//   - the fallback marker must NOT appear anywhere in the subtree
//   - aria-label must NOT land on the fragment's children, because
//     cloneElement on a fragment can't forward props onto inner
//     DOM nodes (and we never inject onto the fallback path).
// This pins behavior: callers passing a fragment must own a11y
// on the inner element themselves.
describe("DocsTooltip — fragment children", () => {
  it("fragment with multiple inner nodes: no fallback marker, no aria-label leakage", () => {
    const { container } = renderWithProvider(
      <DocsTooltip label="Hint">
        <>
          <span data-testid="frag-a">a</span>
          <span data-testid="frag-b">b</span>
        </>
      </DocsTooltip>,
    );
    expect(container.querySelector(`[${FALLBACK_ATTR}]`)).toBeNull();
    const all = container.querySelectorAll("[aria-label]");
    all.forEach((el) => expect(el.getAttribute("aria-label")).not.toBe("Hint"));
  });

  it("fragment wrapping a single valid element: no fallback marker, no aria-label on the inner element", () => {
    const { container } = renderWithProvider(
      <DocsTooltip label="Hint">
        <>
          <button type="button">Click</button>
        </>
      </DocsTooltip>,
    );
    expect(container.querySelector(`[${FALLBACK_ATTR}]`)).toBeNull();
    const button = container.querySelector("button");
    expect(button).not.toBeNull();
    // Fragments don't forward cloned props to inner DOM nodes —
    // pin that the caller (not DocsTooltip) owns a11y here.
    expect(button?.hasAttribute("aria-label")).toBe(false);
  });
});
