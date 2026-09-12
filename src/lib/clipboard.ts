import { queryWrapper, queryWrapperSync } from "./queryWrapper";

function isModernClipboardSupported(): boolean {
  if (typeof navigator === "undefined" || !navigator.clipboard) {
    return false;
  }

  if (typeof navigator.clipboard.writeText !== "function") {
    return false;
  }

  if (typeof window !== "undefined" && typeof window.isSecureContext !== "undefined" && !window.isSecureContext) {
    return false;
  }

  return true;
}

async function tryModernClipboardCopy(text: string): Promise<boolean> {
  if (!isModernClipboardSupported()) {
    return false;
  }

  const res = await queryWrapper(async () => {
    await navigator.clipboard.writeText(text);

    return true;
  });

  if (res.isFail || !res.data) {
    return false;
  }

  return true;
}

export async function copyToClipboard(text: string): Promise<boolean> {
  const isCopied = await tryModernClipboardCopy(text);
  if (isCopied) {
    return true;
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
  const res = queryWrapperSync(() => {
    textarea.focus();
    textarea.select();
    // Some browsers ignore .select() unless setSelectionRange runs too.
    textarea.setSelectionRange(0, text.length);

    return document.execCommand("copy");
  });

  if (res.isFail) {
    succeeded = false;
  } else {
    succeeded = Boolean(res.data);
  }

  document.body.removeChild(textarea);

  if (previousRange && previousSelection) {
    previousSelection.removeAllRanges();
    previousSelection.addRange(previousRange);
  }

  return succeeded;
}
