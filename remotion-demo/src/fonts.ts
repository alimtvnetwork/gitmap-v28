import { loadFont as loadUbuntu } from "@remotion/google-fonts/Ubuntu";
import { loadFont as loadMono } from "@remotion/google-fonts/UbuntuMono";

const ubuntu = loadUbuntu("normal", { weights: ["400", "500", "700"], subsets: ["latin"] });
const mono = loadMono("normal", { weights: ["400", "700"], subsets: ["latin"] });

export const FONT_HEADING = ubuntu.fontFamily;
export const FONT_MONO = mono.fontFamily;