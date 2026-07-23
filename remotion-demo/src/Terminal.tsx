import React from "react";
import { AbsoluteFill, useCurrentFrame, interpolate, spring, useVideoConfig } from "remotion";
import { COLORS, TRAFFIC } from "./theme";
import { FONT_MONO, FONT_HEADING } from "./fonts";

export type Token = { t: string; c?: string };
export type TerminalLine =
  | { kind: "prompt"; cwd: string; tokens: Token[] }
  | { kind: "out"; tokens: Token[] }
  | { kind: "blank" };

type Props = {
  title?: string;
  cwd?: string;
  lines: TerminalLine[];
  // frames per character while typing the prompt line
  typeSpeed?: number;
  // frames to wait between lines (after typing prompt)
  linePause?: number;
  // overall scale (zoom) applied via transform; can animate
  scale?: number;
  // optional caption shown above terminal
  caption?: string;
};

// Layout the lines onto a virtual frame timeline so we can render
// progressive typing of prompts and immediate output reveals.
const buildTimeline = (lines: TerminalLine[], typeSpeed: number, linePause: number) => {
  const starts: number[] = [];
  const ends: number[] = [];
  let cursor = 0;
  for (const ln of lines) {
    starts.push(cursor);
    if (ln.kind === "prompt") {
      const chars = ln.tokens.reduce((n, t) => n + t.t.length, 0);
      cursor += Math.max(8, chars * typeSpeed) + linePause;
    } else if (ln.kind === "out") {
      cursor += linePause;
    } else {
      cursor += Math.floor(linePause / 2);
    }
    ends.push(cursor);
  }
  return { starts, ends, total: cursor };
};

const renderTokens = (tokens: Token[], reveal: number) => {
  let used = 0;
  const out: React.ReactNode[] = [];
  for (let i = 0; i < tokens.length; i++) {
    const tk = tokens[i];
    const remaining = Math.max(0, reveal - used);
    const slice = tk.t.slice(0, remaining);
    if (slice.length > 0) {
      out.push(
        <span key={i} style={{ color: tk.c ?? COLORS.text }}>{slice}</span>,
      );
    }
    used += tk.t.length;
    if (used >= reveal) break;
  }
  return out;
};

export const Terminal: React.FC<Props> = ({
  title = "gitmap — bash",
  cwd = "~/repos",
  lines,
  typeSpeed = 1.4,
  linePause = 6,
  scale = 1,
  caption,
}) => {
  const frame = useCurrentFrame();
  const { width, height } = useVideoConfig();
  const { starts } = buildTimeline(lines, typeSpeed, linePause);

  const termWidth = Math.min(1500, width - 240);
  const termHeight = Math.min(820, height - 240);

  const cursorBlink = Math.floor(frame / 12) % 2 === 0;

  return (
    <AbsoluteFill style={{ alignItems: "center", justifyContent: "center" }}>
      {caption && (
        <div
          style={{
            position: "absolute",
            top: 70,
            fontFamily: FONT_HEADING,
            color: COLORS.muted,
            fontSize: 28,
            letterSpacing: 4,
            textTransform: "uppercase",
          }}
        >
          {caption}
        </div>
      )}
      <div
        style={{
          width: termWidth,
          height: termHeight,
          transform: `scale(${scale})`,
          transformOrigin: "center center",
          background: COLORS.panel,
          border: `1px solid ${COLORS.panelBorder}`,
          borderRadius: 14,
          overflow: "hidden",
          boxShadow: "0 40px 120px -20px rgba(0,0,0,0.7), 0 0 0 1px rgba(232,179,74,0.06)",
          display: "flex",
          flexDirection: "column",
        }}
      >
        {/* Title bar */}
        <div
          style={{
            height: 44,
            background: COLORS.titleBar,
            display: "flex",
            alignItems: "center",
            paddingLeft: 18,
            borderBottom: `1px solid ${COLORS.panelBorder}`,
            position: "relative",
          }}
        >
          <div style={{ display: "flex", gap: 8 }}>
            <Dot color={TRAFFIC.red} />
            <Dot color={TRAFFIC.yellow} />
            <Dot color={TRAFFIC.green} />
          </div>
          <div
            style={{
              position: "absolute",
              left: 0,
              right: 0,
              textAlign: "center",
              color: COLORS.muted,
              fontFamily: FONT_MONO,
              fontSize: 16,
              letterSpacing: 0.5,
            }}
          >
            {title}
          </div>
        </div>
        {/* Body */}
        <div
          style={{
            flex: 1,
            padding: "22px 28px",
            fontFamily: FONT_MONO,
            fontSize: 22,
            lineHeight: 1.55,
            color: COLORS.text,
            overflow: "hidden",
          }}
        >
          {lines.map((ln, idx) => {
            const start = starts[idx];
            const local = frame - start;
            if (local < 0) return <div key={idx} style={{ height: 34 }} />;

            if (ln.kind === "blank") {
              return <div key={idx} style={{ height: 14 }} />;
            }

            if (ln.kind === "prompt") {
              const totalChars = ln.tokens.reduce((n, t) => n + t.t.length, 0);
              const reveal = Math.min(totalChars, Math.floor(local / typeSpeed));
              const isTyping = reveal < totalChars;
              return (
                <div key={idx} style={{ display: "flex", flexWrap: "wrap" }}>
                  <span style={{ color: COLORS.success }}>➜</span>
                  <span style={{ color: COLORS.flag, marginLeft: 10 }}>{ln.cwd}</span>
                  <span style={{ color: COLORS.prompt, marginLeft: 10, marginRight: 10 }}>$</span>
                  <span>{renderTokens(ln.tokens, reveal)}</span>
                  {isTyping && cursorBlink && (
                    <span
                      style={{
                        display: "inline-block",
                        width: 12,
                        height: 24,
                        background: COLORS.prompt,
                        marginLeft: 2,
                        verticalAlign: "middle",
                      }}
                    />
                  )}
                </div>
              );
            }

            // output line: fade in
            const opacity = interpolate(local, [0, 4], [0, 1], { extrapolateRight: "clamp" });
            return (
              <div key={idx} style={{ opacity }}>
                {renderTokens(ln.tokens, 9999)}
              </div>
            );
          })}
        </div>
      </div>
    </AbsoluteFill>
  );
};

const Dot: React.FC<{ color: string }> = ({ color }) => (
  <div style={{ width: 14, height: 14, borderRadius: 7, background: color }} />
);

// Helper: compute the natural duration (in frames) of a given lines array
export const computeTerminalDuration = (
  lines: TerminalLine[],
  typeSpeed = 1.4,
  linePause = 6,
) => buildTimeline(lines, typeSpeed, linePause).total;

// Spring-driven scale helper used by scenes for zoom-in / zoom-out moments
export const useZoom = (inFrame: number, outFrame: number, fps: number) => {
  const f = useCurrentFrame();
  const inS = spring({ frame: f - inFrame, fps, config: { damping: 18, stiffness: 90 } });
  const outS = spring({ frame: f - outFrame, fps, config: { damping: 22, stiffness: 110 } });
  const baseToZoom = interpolate(inS, [0, 1], [1, 1.18]);
  const zoomToBase = interpolate(outS, [0, 1], [0, -0.18]);
  return baseToZoom + zoomToBase;
};