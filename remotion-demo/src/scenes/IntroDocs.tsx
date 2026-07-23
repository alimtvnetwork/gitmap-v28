import { AbsoluteFill, Img, useCurrentFrame, interpolate, spring, useVideoConfig, staticFile } from "remotion";
import { COLORS } from "../theme";
import { FONT_HEADING, FONT_MONO } from "../fonts";

export const IntroDocs: React.FC = () => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();

  const fadeIn = interpolate(frame, [0, 14], [0, 1], { extrapolateRight: "clamp" });
  const titleSpring = spring({ frame: frame - 6, fps, config: { damping: 22, stiffness: 100 } });
  const titleY = interpolate(titleSpring, [0, 1], [30, 0]);

  const drift = Math.sin(frame / 28) * 6;
  const ken = interpolate(frame, [0, 90], [1.0, 1.06]);

  return (
    <AbsoluteFill
      style={{
        background: `radial-gradient(circle at 30% 20%, ${COLORS.bgGradientB}, ${COLORS.bg} 70%)`,
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <div
        style={{
          opacity: fadeIn,
          transform: `translateY(${titleY}px)`,
          textAlign: "center",
          marginBottom: 36,
        }}
      >
        <div
          style={{
            fontFamily: FONT_HEADING,
            fontSize: 88,
            fontWeight: 700,
            color: COLORS.brandGold,
            letterSpacing: -2,
          }}
        >
          gitmap
        </div>
        <div
          style={{
            fontFamily: FONT_MONO,
            color: COLORS.muted,
            fontSize: 22,
            marginTop: 6,
            letterSpacing: 3,
          }}
        >
          one CLI to map every repo you own
        </div>
      </div>
      <div
        style={{
          width: 1500,
          height: 760,
          borderRadius: 18,
          overflow: "hidden",
          border: `1px solid ${COLORS.panelBorder}`,
          boxShadow: "0 50px 140px -20px rgba(0,0,0,0.8)",
          opacity: fadeIn,
          transform: `translateY(${drift}px) scale(${ken})`,
        }}
      >
        <Img
          src={staticFile("images/docs-screenshot.png")}
          style={{ width: "100%", height: "100%", objectFit: "cover", objectPosition: "top center" }}
        />
      </div>
    </AbsoluteFill>
  );
};