import { readFileSync } from "node:fs";
import { join } from "node:path";
import { fileURLToPath } from "node:url";

const root = fileURLToPath(new URL("../", import.meta.url));
const versionPath = join(root, "version.json");

try {
  const content = readFileSync(versionPath, "utf8");
  const data = JSON.parse(content);
  if (!data.Version) {
    console.error("FAIL: version.json is missing 'Version' field.");
    process.exit(1);
  }
  console.log(`✔ version.json is valid (v${data.Version})`);
} catch (err) {
  console.error("FAIL: Could not validate version.json:", err.message);
  process.exit(1);
}
