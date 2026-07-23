import { AbsoluteFill, useCurrentFrame, useVideoConfig, interpolate, spring } from "remotion";
import { COLORS } from "../theme";
import { FONT_HEADING } from "../fonts";
import { Terminal, TerminalLine, computeTerminalDuration } from "../Terminal";

type Props = {
  caption: string;
  cwd: string;
  lines: TerminalLine[];
};

// Each command scene: zoom-IN at start (focus on prompt), zoom-OUT at end
// (let user see the whole output before transitioning).
export const CommandScene: React.FC<Props> = ({ caption, cwd, lines }) => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();

  const total = computeTerminalDuration(lines, 1.4, 6);
  const zoomInS = spring({ frame: frame - 4, fps, config: { damping: 18, stiffness: 80 } });
  const zoomOutTrigger = total + 8;
  const zoomOutS = spring({ frame: frame - zoomOutTrigger, fps, config: { damping: 22, stiffness: 100 } });

  // Start at 1.0, zoom in to 1.18, then ease back to 1.0
  const inAdd = interpolate(zoomInS, [0, 1], [0, 0.18]);
  const outAdd = interpolate(zoomOutS, [0, 1], [0, -0.18]);
  const scale = 1 + inAdd + outAdd;

  const captionOpacity = interpolate(frame, [0, 10, total + 30, total + 50], [0, 1, 1, 0.7], {
    extrapolateRight: "clamp",
  });

  return (
    <AbsoluteFill
      style={{
        background: `linear-gradient(160deg, ${COLORS.bgGradientA}, ${COLORS.bg})`,
      }}
    >
      <div
        style={{
          position: "absolute",
          top: 60,
          width: "100%",
          textAlign: "center",
          fontFamily: FONT_HEADING,
          color: COLORS.brandGold,
          fontSize: 36,
          fontWeight: 700,
          letterSpacing: 1,
          opacity: captionOpacity,
        }}
      >
        {caption}
      </div>
      <Terminal cwd={cwd} lines={lines} scale={scale} title="gitmap — bash" />
    </AbsoluteFill>
  );
};

// Per-scene duration including zoom-out tail
export const sceneDurationFor = (lines: TerminalLine[]) =>
  Math.round(computeTerminalDuration(lines, 1.4, 6)) + 60; // tail for read + zoom-out