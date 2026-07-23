// Clipboard helper with a graceful fallback for environments where
// `navigator.clipboard` is not available.
//
// `navigator.clipboard.writeText` is undefined in three common cases:
//
//   1. Insecure contexts (plain http://, file://) — the Clipboard API is
//      gated behind a Secure Context check by every modern browser.
//   2. Older browsers (notably iOS Safari < 13.4) and embedded webviews.
//   3. Sandboxed iframes without `clipboard-write` permission.
//
// In any of those cases we fall back to the legacy `document.execCommand
// ("copy")` path, which works against a transient hidden <textarea>. It
// is officially deprecated but still implemented everywhere we care
// about, and is the canonical fallback recommended by both MDN and the
// W3C Clipboard API spec.
//
// The function ALWAYS returns a Promise<boolean>:
//   - resolves true on success
//   - resolves false on failure (never throws)
// so call sites can do `if (await copyToClipboard(text)) { ... }`
// without try/catch noise.

export async function copyToClipboard(text: string): Promise<boolean> {
  // Path 1: modern async Clipboard API. Guard for both the property
  // existing AND the document being focused — Safari rejects writes
  // from blurred documents with a NotAllowedError.
  if (
    typeof navigator !== "undefined" &&
    navigator.clipboard &&
    typeof navigator.clipboard.writeText === "function" &&
    (typeof window === "undefined" ||
      typeof window.isSecureContext === "undefined" ||
      window.isSecureContext)
  ) {
    try {
      await navigator.clipboard.writeText(text);
      return true;
    } catch {
      // Fall through to legacy path — common when the document is not
      // focused, or the Permissions API denied clipboard-write.
    }
  }

  return legacyCopy(text);
}

// legacyCopy creates an off-screen, read-only <textarea>, selects its
// contents, runs `document.execCommand("copy")`, then removes it. The
// off-screen positioning is critical: setting display:none or
// visibility:hidden makes the selection silently fail in some browsers.
function legacyCopy(text: string): boolean {
  if (typeof document === "undefined") {
    return false;
  }

  const textarea = document.createElement("textarea");
  textarea.value = text;
  // readOnly + inputMode=none keeps mobile keyboards from popping up
  // during the brief moment the textarea is in the DOM.
  textarea.setAttribute("readonly", "");
  textarea.setAttribute("aria-hidden", "true");
  textarea.style.position = "fixed";
  textarea.style.top = "0";
  textarea.style.left = "0";
  textarea.style.width = "1px";
  textarea.style.height = "1px";
  textarea.style.padding = "0";
  textarea.style.border = "none";
  textarea.style.outline = "none";
  textarea.style.boxShadow = "none";
  textarea.style.background = "transparent";
  textarea.style.opacity = "0";

  document.body.appendChild(textarea);

  // Preserve the caller's existing selection so we don't disrupt
  // whatever the user had highlighted on the page.
  const previousSelection = document.getSelection();
  const previousRange =
    previousSelection && previousSelection.rangeCount > 0
      ? previousSelection.getRangeAt(0)
      : null;

  let succeeded = false;
  try {
    textarea.focus();
    textarea.select();
    // Some browsers ignore .select() unless setSelectionRange runs too.
    textarea.setSelectionRange(0, text.length);
    succeeded = document.execCommand("copy");
  } catch {
    succeeded = false;
  } finally {
    document.body.removeChild(textarea);
    if (previousRange && previousSelection) {
      previousSelection.removeAllRanges();
      previousSelection.addRange(previousRange);
    }
  }

  return succeeded;
}
