import { Composition } from "remotion";
import { MainVideo, VIDEO_FPS, VIDEO_WIDTH, VIDEO_HEIGHT, TOTAL_FRAMES } from "./MainVideo";

export const RemotionRoot = () => (
  <Composition
    id="main"
    component={MainVideo}
    durationInFrames={TOTAL_FRAMES}
    fps={VIDEO_FPS}
    width={VIDEO_WIDTH}
    height={VIDEO_HEIGHT}
  />
);