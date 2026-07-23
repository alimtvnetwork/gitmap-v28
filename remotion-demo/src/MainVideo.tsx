import { AbsoluteFill } from "remotion";
import React from "react";
import { TransitionSeries, springTiming } from "@remotion/transitions";
import { fade } from "@remotion/transitions/fade";
import { slide } from "@remotion/transitions/slide";

import { IntroDocs } from "./scenes/IntroDocs";
import { IntroTUI } from "./scenes/IntroTUI";
import { CommandScene, sceneDurationFor } from "./scenes/CommandScene";
import { COMMANDS } from "./commands";

export const VIDEO_FPS = 30;
export const VIDEO_WIDTH = 1920;
export const VIDEO_HEIGHT = 1080;

const INTRO_DOCS_FRAMES = 110; // ~3.7s
const INTRO_TUI_FRAMES = 100;  // ~3.3s
const TRANSITION_FRAMES = 18;

const cmdDurations = COMMANDS.map((c) => sceneDurationFor(c.lines));
const cmdTotal = cmdDurations.reduce((a, b) => a + b, 0);
const transitionsCount = 1 /*docs→tui*/ + 1 /*tui→cmd0*/ + (COMMANDS.length - 1);

export const TOTAL_FRAMES = Math.round(
  INTRO_DOCS_FRAMES + INTRO_TUI_FRAMES + cmdTotal - transitionsCount * TRANSITION_FRAMES,
);

export const MainVideo: React.FC = () => {
  return (
    <AbsoluteFill style={{ background: "#0b0d10" }}>
      <TransitionSeries>
        <TransitionSeries.Sequence durationInFrames={INTRO_DOCS_FRAMES}>
          <IntroDocs />
        </TransitionSeries.Sequence>

        <TransitionSeries.Transition
          presentation={fade()}
          timing={springTiming({ config: { damping: 200 }, durationInFrames: TRANSITION_FRAMES })}
        />

        <TransitionSeries.Sequence durationInFrames={INTRO_TUI_FRAMES}>
          <IntroTUI />
        </TransitionSeries.Sequence>

        <TransitionSeries.Transition
          presentation={slide({ direction: "from-right" })}
          timing={springTiming({ config: { damping: 200 }, durationInFrames: TRANSITION_FRAMES })}
        />

        {COMMANDS.flatMap((c, i) => {
          const nodes: React.ReactNode[] = [];
          if (i > 0) {
            nodes.push(
              <TransitionSeries.Transition
                key={`tr-${i}`}
                presentation={fade()}
                timing={springTiming({ config: { damping: 200 }, durationInFrames: TRANSITION_FRAMES })}
              />,
            );
          }
          nodes.push(
            <TransitionSeries.Sequence key={`seq-${i}`} durationInFrames={cmdDurations[i]}>
              <CommandScene caption={c.caption} cwd={c.cwd} lines={c.lines} />
            </TransitionSeries.Sequence>,
          );
          return nodes;
        })}
      </TransitionSeries>
    </AbsoluteFill>
  );
};