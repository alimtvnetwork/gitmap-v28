import { AbsoluteFill, useCurrentFrame, interpolate, spring, useVideoConfig } from "remotion";
import { COLORS, TRAFFIC } from "../theme";
import { FONT_HEADING, FONT_MONO } from "../fonts";

const ROWS = [
  { name: "gitmap-v27", branch: "main", status: "clean", color: COLORS.success },
  { name: "lovable-cloud", branch: "main", status: "behind 3", color: COLORS.warn },
  { name: "remotion-demos", branch: "feat/intro", status: "dirty", color: COLORS.danger },
  { name: "design-system", branch: "main", status: "clean", color: COLORS.success },
  { name: "scripts-shared", branch: "main", status: "clean", color: COLORS.success },
  { name: "infra-iac", branch: "release", status: "behind 1", color: COLORS.warn },
];

export const IntroTUI: React.FC = () => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();

  const fade = interpolate(frame, [0, 12], [0, 1], { extrapolateRight: "clamp" });
  const cursorIdx = Math.min(ROWS.length - 1, Math.floor(frame / 14));

  return (
    <AbsoluteFill
      style={{
        background: `linear-gradient(140deg, ${COLORS.bgGradientA}, ${COLORS.bg})`,
        alignItems: "center",
        justifyContent: "center",
        opacity: fade,
      }}
    >
      <div style={{ width: 1500, height: 820, background: COLORS.panel, border: `1px solid ${COLORS.panelBorder}`, borderRadius: 14, overflow: "hidden", boxShadow: "0 40px 120px -20px rgba(0,0,0,0.7)" }}>
        <div style={{ height: 44, background: COLORS.titleBar, display: "flex", alignItems: "center", paddingLeft: 18, borderBottom: `1px solid ${COLORS.panelBorder}` }}>
          <div style={{ display: "flex", gap: 8 }}>
            <Dot c={TRAFFIC.red} /><Dot c={TRAFFIC.yellow} /><Dot c={TRAFFIC.green} />
          </div>
          <div style={{ flex: 1, textAlign: "center", color: COLORS.muted, fontFamily: FONT_MONO, fontSize: 16 }}>gitmap interactive</div>
        </div>
        <div style={{ padding: 28, fontFamily: FONT_MONO, color: COLORS.text }}>
          <div style={{ fontFamily: FONT_HEADING, fontSize: 30, color: COLORS.brandGold, marginBottom: 6 }}>Repos</div>
          <div style={{ color: COLORS.muted, fontSize: 18, marginBottom: 22 }}>↑↓ navigate · enter open · / search · q quit</div>
          {/* tabs */}
          <div style={{ display: "flex", gap: 18, marginBottom: 18, fontSize: 18 }}>
            {["Repos", "Actions", "Groups", "Status", "Releases", "Logs"].map((t, i) => (
              <span key={t} style={{ color: i === 0 ? COLORS.brandGold : COLORS.muted, borderBottom: i === 0 ? `2px solid ${COLORS.brandGold}` : "none", paddingBottom: 4 }}>{t}</span>
            ))}
          </div>
          <div style={{ borderTop: `1px solid ${COLORS.panelBorder}`, paddingTop: 12 }}>
            <Row name="REPOSITORY" branch="BRANCH" status="STATUS" color={COLORS.muted} header />
            {ROWS.map((r, i) => (
              <Row
                key={r.name}
                name={r.name}
                branch={r.branch}
                status={r.status}
                color={r.color}
                highlighted={i === cursorIdx}
              />
            ))}
          </div>
        </div>
      </div>
    </AbsoluteFill>
  );
};

const Dot: React.FC<{ c: string }> = ({ c }) => <div style={{ width: 14, height: 14, borderRadius: 7, background: c }} />;

const Row: React.FC<{ name: string; branch: string; status: string; color: string; highlighted?: boolean; header?: boolean }> = ({ name, branch, status, color, highlighted, header }) => (
  <div
    style={{
      display: "grid",
      gridTemplateColumns: "60px 1fr 1fr 1fr",
      padding: "10px 14px",
      fontSize: header ? 16 : 22,
      letterSpacing: header ? 2 : 0,
      background: highlighted ? "rgba(232,179,74,0.1)" : "transparent",
      borderLeft: highlighted ? `3px solid ${COLORS.brandGold}` : "3px solid transparent",
      color: header ? COLORS.muted : COLORS.text,
    }}
  >
    <span>{header ? "" : highlighted ? "▶" : ""}</span>
    <span style={{ color: header ? COLORS.muted : highlighted ? COLORS.brandGold : COLORS.text }}>{name}</span>
    <span style={{ color: header ? COLORS.muted : COLORS.flag }}>{branch}</span>
    <span style={{ color }}>{status}</span>
  </div>
);