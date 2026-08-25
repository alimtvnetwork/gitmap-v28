# Swapp Off


## Swap Off

To permanently disable swap on Ubuntu, follow these steps:

1. **Turn off swap immediately**:
   ```bash
   sudo swapoff -a
   ```

2. **Remove the swap file** (if you are using a swap file):
   ```bash
   sudo rm /swapfile
   ```

3. **Edit the `/etc/fstab` file** to prevent swap from being re-enabled at boot:
   ```bash
   sudo nano /etc/fstab
   ```
   Find the line that references the swap file or partition (e.g., `/swapfile` or `/dev/sdX`) and comment it out by adding a `#` at the beginning of the line. Save and exit the editor.

4. **Verify that swap is disabled**:
   ```bash
   sudo swapon --show
   ```
   There should be no output if swap is disabled.

5. **Reboot your system** to ensure the changes take effect:
   ```bash
   sudo reboot
   ```
