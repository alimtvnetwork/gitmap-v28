# Subtask 03: Live Ubuntu Host Upgrade to Antigravity IDE 2.22.0

- **Subtask ID:** `SUB-03`
- **Parent Plan:** `.ai-memory/plans/root-cause-analysis-and-update-instructions.md`
- **Status:** `COMPLETED`
- **Target Target:** Local Ubuntu Host (`~/.local/share/antigravity-ide/`)

---

## Technical Tasks

1. **Safety Backup:**
   - Verify `~/.local/share/antigravity-ide` safety backup is created if updating in-place.
   - Strictly avoid touching `~/.gemini/` or `~/.config/Antigravity/`.

2. **Execute In-Place Update via GitMap:**
   - Execute `gitmap agy update` or `gitmap update ide`.
   - Verify download of 2.22.0 tarball from `https://storage.googleapis.com/antigravity-public/antigravity-hub/2.22.0-5446056071266304/linux-x64/Antigravity.tar.gz`.
   - Verify extraction into `~/.local/share/antigravity-ide`.

3. **Verify Launcher Wrapper & Symlinks:**
   - Verify `~/.local/share/antigravity-ide/antigravity-ide.run` is executable and passes `--no-sandbox`.
   - Verify `~/.local/bin/{antigravity, agy, antigravity-ide}` point to `antigravity-ide.run`.
   - Verify `file ~/.local/share/antigravity-ide/antigravity` confirms 64-bit ELF executable.
