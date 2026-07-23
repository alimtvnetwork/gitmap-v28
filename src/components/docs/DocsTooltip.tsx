import {
  Children,
  cloneElement,
  isValidElement,
  ReactElement,
  ReactNode,
} from "react";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  DOCS_TOOLTIP_ALIGN,
  DOCS_TOOLTIP_SIDE,
  DOCS_TOOLTIP_SIDE_OFFSET,
} from "@/components/docs/docsTooltipConfig";

interface DocsTooltipProps {
  // The control the tooltip describes. Pass the trigger element
  // directly; this wrapper hands it to TooltipTrigger asChild so
  // the trigger keeps its own ref/keyboard semantics.
  children: ReactNode;
  // Tooltip body. Keep it short (one short phrase) — long text
  // belongs in inline help, not a hover tooltip.
  label: ReactNode;
  // Optional explicit accessible name. Defaults to `label` when it
  // is a string. Pass this when `label` is JSX (e.g. with an icon)
  // so screen-reader users still get a clean spoken name.
  ariaLabel?: string;
}

// Resolve the accessible name we want to expose on the trigger:
// 1. explicit `ariaLabel` prop (wins)
// 2. `label` when it is a plain string
// 3. otherwise undefined — the child must already supply its own.
const resolveAccessibleName = (
  label: ReactNode,
  ariaLabel: string | undefined,
): string | undefined => {
  if (ariaLabel) return ariaLabel;
  if (typeof label === "string") return label;
  return undefined;
};

// Marker prop set on the fallback wrapper produced by
// normalizeTrigger. withAccessibleName uses this to skip aria-label
// injection on synthesized wrappers — the a11y contract only applies
// when the caller passed a real, single React element. Callers using
// non-element children must own their own accessible naming.
const FALLBACK_WRAPPER_PROP = "data-docs-tooltip-fallback" as const;

// Inject aria-label onto the trigger child so keyboard / screen-reader
// users get the same information sighted users get from the tooltip.
// We never overwrite an aria-label the child already set — the child
// wins so callers can be more specific when they need to.
// We also skip the synthesized fallback wrapper (see normalizeTrigger):
// injecting aria-label onto invented markup would silently paper over
// a misuse where the caller passed a non-element child.
const withAccessibleName = (
  child: ReactNode,
  accessibleName: string | undefined,
): ReactNode => {
  if (!accessibleName) return child;
  if (!isValidElement(child)) return child;
  const childProps = child.props as Record<string, unknown>;
  if (childProps[FALLBACK_WRAPPER_PROP]) return child;
  if (childProps["aria-label"]) return child;
  return cloneElement(child as ReactElement, { "aria-label": accessibleName });
};

// Radix's TooltipTrigger uses Slot+React.Children.only under the hood,
// which throws when given a string, number, fragment, null, or
// multiple children. To make DocsTooltip safe for ANY child shape
// (defensive rendering — never let a tooltip crash the docs chrome)
// we normalize non-single-element children into a focusable <span>.
// The wrapper keeps tabIndex=0 so keyboard users can still focus
// the trigger and surface the tooltip body. The wrapper is tagged
// with FALLBACK_WRAPPER_PROP so withAccessibleName skips it.
const normalizeTrigger = (child: ReactNode): ReactElement => {
  const count = Children.count(child);
  if (count === 1 && isValidElement(child)) return child;
  return (
    <span
      tabIndex={0}
      className="inline-flex"
      {...{ [FALLBACK_WRAPPER_PROP]: true }}
    >
      {child}
    </span>
  );
};

export const DocsTooltip = ({ children, label, ariaLabel }: DocsTooltipProps) => {
  const accessibleName = resolveAccessibleName(label, ariaLabel);
  const normalized = normalizeTrigger(children);
  const trigger = withAccessibleName(normalized, accessibleName);
  return (
    <Tooltip>
      <TooltipTrigger asChild>{trigger}</TooltipTrigger>
      <TooltipContent
        side={DOCS_TOOLTIP_SIDE}
        sideOffset={DOCS_TOOLTIP_SIDE_OFFSET}
        align={DOCS_TOOLTIP_ALIGN}
      >
        {label}
      </TooltipContent>
    </Tooltip>
  );
};

export default DocsTooltip;
