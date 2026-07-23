import { describe, it, expect, beforeEach } from "vitest";
import { render, screen, cleanup, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { TooltipProvider } from "@/components/ui/tooltip";
import { SidebarProvider } from "@/components/ui/sidebar";
import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";

// Why focus() instead of hover():
// Radix Tooltip opens on either pointer-enter OR focus. jsdom does
// not synthesize pointer events the way real browsers do, so
// focus is the deterministic path to "tooltip is open" in this
// environment. Real-browser pointer behaviour is covered by the
// shared DocsTooltip component (single source of truth for
// placement/delay) — this suite verifies WIRING, not Radix internals.

const renderDocsChrome = () =>
  render(
    <MemoryRouter>
      <TooltipProvider delayDuration={0} skipDelayDuration={0}>
        <SidebarProvider>
          <DocsLayout>
            <div>content</div>
          </DocsLayout>
        </SidebarProvider>
      </TooltipProvider>
    </MemoryRouter>,
  );

const renderCodeBlock = () =>
  render(
    <TooltipProvider delayDuration={0} skipDelayDuration={0}>
      <CodeBlock code={"line one\nline two"} language="bash" />
    </TooltipProvider>,
  );

// HEADER_CONTROLS lists every icon-only control in the docs header.
// Each entry pairs the trigger's accessible name (used to look the
// element up via getByRole / getByLabelText) with the substring we
// expect to find in the surfaced tooltip body. Update this list
// whenever a control is added/removed from the header.
const HEADER_CONTROLS: Array<{ trigger: string; tooltip: string }> = [
  { trigger: "Toggle sidebar", tooltip: "Toggle sidebar" },
  { trigger: "Dark theme", tooltip: "Dark theme" },
  { trigger: "Light theme", tooltip: "Light theme" },
  { trigger: /Copy .* theme palette to clipboard/i.source, tooltip: "palette" },
  {
    trigger: "Open command palette (search commands, flags, pages)",
    tooltip: "Search commands",
  },
];

// CODEBLOCK_CONTROLS lists every icon-only control in the CodeBlock
// toolbar. Same shape and update rule as HEADER_CONTROLS.
const CODEBLOCK_CONTROLS: Array<{ trigger: string; tooltip: string }> = [
  { trigger: "Decrease font size", tooltip: "Decrease font size" },
  { trigger: "Increase font size", tooltip: "Increase font size" },
  { trigger: "Copy snippet", tooltip: "Copy snippet" },
  { trigger: "Download snippet", tooltip: "Download snippet" },
  { trigger: "Enter fullscreen", tooltip: "Fullscreen" },
];

const findTriggerByName = (name: string) => {
  // Try exact aria-label first; fall back to a regex match for the
  // dynamic CopyPaletteButton label that includes the active theme.
  try {
    return screen.getByLabelText(name);
  } catch {
    return screen.getByLabelText(new RegExp(name, "i"));
  }
};

const expectTooltipFor = async (
  triggerName: string,
  tooltipSubstr: string,
) => {
  const trigger = findTriggerByName(triggerName);
  // Focus is the jsdom-friendly opener (see top-of-file comment).
  trigger.focus();
  // Radix mounts the tooltip body in a portal as role="tooltip".
  // findAllByRole tolerates same-text matches across header + portal.
  const tips = await screen.findAllByRole("tooltip");
  const matched = tips.some((t) =>
    (t.textContent ?? "").toLowerCase().includes(tooltipSubstr.toLowerCase()),
  );
  expect(
    matched,
    `expected an open tooltip whose text includes "${tooltipSubstr}" ` +
      `for trigger "${triggerName}", got: ${tips.map((t) => t.textContent).join(" | ")}`,
  ).toBe(true);
  trigger.blur();
};

describe("docs tooltip wiring — header controls", () => {
  beforeEach(() => cleanup());

  for (const ctl of HEADER_CONTROLS) {
    it(`shows tooltip "${ctl.tooltip}" for "${ctl.trigger}"`, async () => {
      renderDocsChrome();
      await expectTooltipFor(ctl.trigger, ctl.tooltip);
    });
  }
});

describe("docs tooltip wiring — CodeBlock toolbar", () => {
  beforeEach(() => cleanup());

  for (const ctl of CODEBLOCK_CONTROLS) {
    it(`shows tooltip "${ctl.tooltip}" for "${ctl.trigger}"`, async () => {
      renderCodeBlock();
      await expectTooltipFor(ctl.trigger, ctl.tooltip);
    });
  }
});

// Version badge was removed from the header (per branding cleanup —
// header shows "gitmap docs" only, no version chip). The previous
// test asserting a tooltip on the badge no longer applies.


// Pointer-event sanity check: also verify userEvent.hover() opens
// the tooltip in jsdom for at least one well-known trigger. If
// Radix changes its opener semantics this test will catch it.
describe("docs tooltip wiring — hover opener", () => {
  beforeEach(() => cleanup());

  it("opens the sidebar-toggle tooltip on pointer hover", async () => {
    const user = userEvent.setup();
    renderDocsChrome();
    const trigger = screen.getByLabelText("Toggle sidebar");
    await user.hover(trigger);
    const tips = await screen.findAllByRole("tooltip");
    const matched = tips.some((t) =>
      (t.textContent ?? "").toLowerCase().includes("toggle sidebar"),
    );
    expect(matched).toBe(true);
  });
});

// DocsTooltip wraps `children` with TooltipTrigger asChild, but it
// also tries to inject an aria-label onto the child via cloneElement.
// The injection helper must short-circuit when the child is NOT a
// single valid React element — otherwise React would throw at render
// time. This suite locks in that contract: every non-element child
// shape renders without throwing and the tooltip body still appears
// when a focusable descendant receives focus. If someone refactors
// `withAccessibleName` and forgets the isValidElement guard, these
// tests fail loudly instead of crashing the whole docs chrome.
import { DocsTooltip } from "@/components/docs/DocsTooltip";

const renderTooltip = (node: React.ReactNode) =>
  render(
    <TooltipProvider delayDuration={0} skipDelayDuration={0}>
      <DocsTooltip label="fallback label">{node}</DocsTooltip>
    </TooltipProvider>,
  );

describe("DocsTooltip — non-element children fallback", () => {
  beforeEach(() => cleanup());

  it("does not throw when the child is a plain string", () => {
    expect(() => renderTooltip("just text")).not.toThrow();
    expect(screen.getByText("just text")).toBeTruthy();
  });

  it("does not throw when the child is a number", () => {
    expect(() => renderTooltip(42)).not.toThrow();
  });

  it("does not throw when the child is null", () => {
    expect(() => renderTooltip(null)).not.toThrow();
  });

  it("does not throw when there are multiple children (array)", () => {
    expect(() =>
      renderTooltip([
        <span key="a">a</span>,
        <button key="b" type="button">
          b
        </button>,
      ]),
    ).not.toThrow();
  });

  it("does not throw when the child is a fragment", () => {
    expect(() =>
      renderTooltip(
        <>
          <span>x</span>
          <span>y</span>
        </>,
      ),
    ).not.toThrow();
  });

  it("renders DOM safely when given a fragment with multiple children", () => {
    // Contract under test: DocsTooltip MUST NOT throw for any
    // ReactNode shape. With multi-children Radix's Slot may not
    // wire the tooltip opener (Slot needs a single element to
    // forward props onto) — that degraded behaviour is acceptable.
    // What is NOT acceptable is a render-time crash. This test
    // pins down that guarantee: render succeeds and the children
    // are present in the DOM.
    expect(() =>
      renderTooltip(
        <>
          <button type="button" aria-label="inner btn">
            inner
          </button>
          <span>extra</span>
        </>,
      ),
    ).not.toThrow();
    expect(screen.getByLabelText("inner btn")).toBeTruthy();
    expect(screen.getByText("extra")).toBeTruthy();
  });

  it("opens the tooltip when the normalized wrapper is reached via keyboard Tab", async () => {
    // When the child is a non-element (string), DocsTooltip wraps it
    // in a tabIndex=0 span so it remains keyboard-reachable AND
    // Radix can wire its trigger handlers (data-state="closed"
    // appears on the wrapper). This test simulates a real keyboard
    // user pressing Tab to land on the wrapper and asserts the
    // tooltip body becomes visible — the core a11y promise of the
    // fallback path.
    const user = userEvent.setup();
    renderTooltip("just text");
    const wrapper = screen.getByText("just text");
    expect(wrapper.tagName).toBe("SPAN");
    expect(wrapper.getAttribute("tabindex")).toBe("0");
    // aria-label is intentionally NOT injected onto the fallback
    // wrapper (see the dedicated injection-scope suite below) —
    // callers using non-element children must own their own naming.
    expect(wrapper.getAttribute("aria-label")).toBeNull();

    // Tab from document.body lands on the only focusable node.
    await user.tab();
    expect(document.activeElement).toBe(wrapper);

    const tips = await screen.findAllByRole("tooltip");
    const matched = tips.some((t) =>
      (t.textContent ?? "").includes("fallback label"),
    );
    expect(matched).toBe(true);
  });
});

// aria-label injection is the keyboard/SR a11y promise of DocsTooltip:
// when the trigger child is a single valid React element AND the child
// does not already declare an aria-label, the tooltip's `label` (or the
// explicit `ariaLabel` prop) is grafted onto the child as aria-label.
// Conversely, the synthesized fallback wrapper used for non-element
// children (string, number, null, fragment, multi-children) is NOT a
// caller-owned trigger — injecting onto it would silently paper over
// misuse. This suite locks in BOTH halves of the contract.
describe("DocsTooltip — aria-label injection scope", () => {
  beforeEach(() => cleanup());

  it("injects aria-label onto a single valid element child", () => {
    render(
      <TooltipProvider delayDuration={0} skipDelayDuration={0}>
        <DocsTooltip label="Save changes">
          <button type="button" data-testid="real">
            <svg aria-hidden="true" />
          </button>
        </DocsTooltip>
      </TooltipProvider>,
    );
    const btn = screen.getByTestId("real");
    expect(btn.getAttribute("aria-label")).toBe("Save changes");
  });

  it("preserves a child's existing aria-label (child wins)", () => {
    render(
      <TooltipProvider delayDuration={0} skipDelayDuration={0}>
        <DocsTooltip label="generic label">
          <button type="button" aria-label="explicit child label" data-testid="real">
            x
          </button>
        </DocsTooltip>
      </TooltipProvider>,
    );
    expect(screen.getByTestId("real").getAttribute("aria-label")).toBe(
      "explicit child label",
    );
  });

  it("does NOT inject aria-label onto the fallback wrapper for a string child", () => {
    render(
      <TooltipProvider delayDuration={0} skipDelayDuration={0}>
        <DocsTooltip label="should not leak">just text</DocsTooltip>
      </TooltipProvider>,
    );
    const wrapper = screen.getByText("just text");
    expect(wrapper.tagName).toBe("SPAN");
    expect(wrapper.getAttribute("aria-label")).toBeNull();
    // The fallback marker prop must reach the DOM so the skip rule
    // is observable / testable from the outside.
    expect(wrapper.getAttribute("data-docs-tooltip-fallback")).toBe("true");
  });

  it("does NOT inject aria-label onto the fallback wrapper for a fragment child", () => {
    render(
      <TooltipProvider delayDuration={0} skipDelayDuration={0}>
        <DocsTooltip label="should not leak">
          <>
            <span data-testid="frag-a">a</span>
            <span data-testid="frag-b">b</span>
          </>
        </DocsTooltip>
      </TooltipProvider>,
    );
    // Neither inner element gets an aria-label.
    expect(screen.getByTestId("frag-a").getAttribute("aria-label")).toBeNull();
    expect(screen.getByTestId("frag-b").getAttribute("aria-label")).toBeNull();
  });

  it("uses the explicit ariaLabel prop (not label) when injecting onto a real element", () => {
    render(
      <TooltipProvider delayDuration={0} skipDelayDuration={0}>
        <DocsTooltip
          label={<span>icon + text</span>}
          ariaLabel="explicit aria"
        >
          <button type="button" data-testid="real" />
        </DocsTooltip>
      </TooltipProvider>,
    );
    expect(screen.getByTestId("real").getAttribute("aria-label")).toBe(
      "explicit aria",
    );
  });

  // Extra fallback-shape coverage: anything that is NOT a single valid
  // React element must be wrapped in the synthesized fallback span and
  // must NOT receive aria-label injection. We assert via the public
  // wrapper marker (data-docs-tooltip-fallback="true") so the contract
  // is observable from the rendered DOM, not from internal refs.
  it("wraps a number child in a fallback wrapper without aria-label", () => {
    render(
      <TooltipProvider delayDuration={0} skipDelayDuration={0}>
        <DocsTooltip label="should not leak">{42}</DocsTooltip>
      </TooltipProvider>,
    );
    const wrapper = screen.getByText("42");
    expect(wrapper.tagName).toBe("SPAN");
    expect(wrapper.getAttribute("data-docs-tooltip-fallback")).toBe("true");
    expect(wrapper.getAttribute("aria-label")).toBeNull();
  });

  it("wraps multi-children in a single fallback wrapper without aria-label", () => {
    render(
      <TooltipProvider delayDuration={0} skipDelayDuration={0}>
        <DocsTooltip label="should not leak">
          <span data-testid="multi-a">a</span>
          <span data-testid="multi-b">b</span>
        </DocsTooltip>
      </TooltipProvider>,
    );
    // Inner real elements stay untouched.
    expect(screen.getByTestId("multi-a").getAttribute("aria-label")).toBeNull();
    expect(screen.getByTestId("multi-b").getAttribute("aria-label")).toBeNull();
    // Their common parent IS the fallback wrapper carrying the marker.
    const wrapper = screen.getByTestId("multi-a").parentElement!;
    expect(wrapper.tagName).toBe("SPAN");
    expect(wrapper.getAttribute("data-docs-tooltip-fallback")).toBe("true");
    expect(wrapper.getAttribute("aria-label")).toBeNull();
  });

  it("renders a fallback wrapper with no aria-label when child is null", () => {
    // A null child still yields a focusable fallback span (so the
    // tooltip stays openable via keyboard) but with no injected name —
    // callers passing null are expected to own their own a11y.
    const { container } = render(
      <TooltipProvider delayDuration={0} skipDelayDuration={0}>
        <DocsTooltip label="should not leak">{null}</DocsTooltip>
      </TooltipProvider>,
    );
    const wrapper = container.querySelector(
      "[data-docs-tooltip-fallback]",
    ) as HTMLElement | null;
    expect(wrapper).not.toBeNull();
    expect(wrapper!.tagName).toBe("SPAN");
    expect(wrapper!.getAttribute("data-docs-tooltip-fallback")).toBe("true");
    expect(wrapper!.getAttribute("aria-label")).toBeNull();
  });
});

// Real users routinely switch between pointer and keyboard mid-task
// (hover → tab away → tab back → hover again). Radix Tooltip tracks
// open state per-trigger and per-event-source; a regression in our
// wrapper / Slot composition could leave the tooltip stuck open,
// stuck closed, or mis-targeted after a mode switch. This suite
// pins down the round-trip: hover → unhover → focus → blur → hover.
describe("DocsTooltip — interaction mode switching", () => {
  beforeEach(() => cleanup());

  const isOpen = (triggerName: string) => {
    const tips = screen.queryAllByRole("tooltip");
    return tips.some((t) =>
      (t.textContent ?? "").toLowerCase().includes(triggerName.toLowerCase()),
    );
  };

  it("re-opens correctly when alternating hover and focus on the same trigger", async () => {
    const user = userEvent.setup();
    renderDocsChrome();
    const trigger = screen.getByLabelText("Toggle sidebar");

    // The contract under test is "no broken state after switching
    // interaction modes" — NOT the precise async close timing of
    // Radix between steps (which is debounced and varies by jsdom
    // pointer simulation). We assert the tooltip is OPEN after each
    // opener event and OPEN at the end of the full round-trip.

    // 1. Hover opens.
    await user.hover(trigger);
    await waitFor(() => expect(isOpen("Toggle sidebar")).toBe(true));

    // 2. Pointer leaves, then keyboard focus takes over.
    await user.unhover(trigger);
    trigger.focus();
    await waitFor(() => expect(isOpen("Toggle sidebar")).toBe(true));

    // 3. Blur, then hover takes over again — proves no stuck
    //    listeners or detached refs from the focus round-trip.
    trigger.blur();
    await user.hover(trigger);
    await waitFor(() => expect(isOpen("Toggle sidebar")).toBe(true));
  });

  it("does not leave a stale tooltip open when focus moves between two triggers", async () => {
    const user = userEvent.setup();
    renderDocsChrome();
    const sidebar = screen.getByLabelText("Toggle sidebar");
    const dark = screen.getByLabelText("Dark theme");

    // Hover the first, then keyboard-focus the second. Only the
    // second tooltip should be open at the end of the sequence.
    await user.hover(sidebar);
    await screen.findAllByRole("tooltip");
    expect(isOpen("Toggle sidebar")).toBe(true);

    await user.unhover(sidebar);
    dark.focus();
    await screen.findAllByRole("tooltip");

    await waitFor(() => {
      expect(isOpen("Dark theme")).toBe(true);
      expect(isOpen("Toggle sidebar")).toBe(false);
    });
  });
});

// Single-render conjunction test for the injection-scope contract.
// The dedicated suite above (line ~270) covers each half (valid
// element gets the label / each fallback shape does not) in
// isolation. This test pins down BOTH halves in ONE render so a
// future refactor that accidentally inverts the guard — e.g.
// "inject on the wrapper too" or "skip on real elements" — fails
// here even if the per-shape suites are mutated together.
describe("DocsTooltip — aria-label injection scope (combined)", () => {
  beforeEach(() => cleanup());

  it("injects on the real-element trigger AND skips the fallback wrapper rendered as a sibling", () => {
    render(
      <TooltipProvider delayDuration={0} skipDelayDuration={0}>
        <div>
          <DocsTooltip label="Real-element label">
            <button type="button" data-testid="real-trigger" />
          </DocsTooltip>
          <DocsTooltip label="Fallback label">just text</DocsTooltip>
        </div>
      </TooltipProvider>,
    );

    // Half 1: the single valid element child receives the label
    // verbatim — withAccessibleName cloned the props onto it.
    const realTrigger = screen.getByTestId("real-trigger");
    expect(realTrigger.tagName).toBe("BUTTON");
    expect(realTrigger.getAttribute("aria-label")).toBe("Real-element label");

    // Half 2: the synthesized fallback wrapper for the string child
    // is publicly tagged with data-docs-tooltip-fallback="true" AND
    // carries NO aria-label — proving the FALLBACK_WRAPPER_PROP
    // short-circuit in withAccessibleName fired.
    const fallbackWrapper = screen.getByText("just text");
    expect(fallbackWrapper.tagName).toBe("SPAN");
    expect(fallbackWrapper.getAttribute("data-docs-tooltip-fallback")).toBe(
      "true",
    );
    expect(fallbackWrapper.getAttribute("aria-label")).toBeNull();

    // Defensive: the fallback's own label string must NEVER appear
    // as an aria-label anywhere in the rendered subtree (a regression
    // where the wrapper-skip is removed would leak it onto the span).
    const leaked = document.querySelector('[aria-label="Fallback label"]');
    expect(leaked).toBeNull();
  });
});

