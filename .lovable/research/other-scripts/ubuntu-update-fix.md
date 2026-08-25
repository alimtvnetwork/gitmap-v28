# Ubuntu 26.04 LTS — `apt` fails with `403 Forbidden` from `my.archive.ubuntu.com`

> Self-contained troubleshooting document. Everything needed to understand and fix
> the problem is in this file — you can paste the whole thing into another AI or
> hand it to a colleague with no extra context.

---

## 1. Problem statement

On a **freshly installed Ubuntu 26.04 LTS** machine (codename `resolute`) located in
**Malaysia**, `apt update` works fine but `apt upgrade` **cannot download one package**
and aborts the whole upgrade:

```
Err:1 http://my.archive.ubuntu.com/ubuntu resolute-updates/main amd64 linux-firmware-amd-misc all 20260319.git217ca6e4-0ubuntu1.1
  403  Forbidden [IP: 202.79.180.254 80]
Error: Failed to fetch http://my.archive.ubuntu.com/ubuntu/pool/main/l/linux-firmware-amd-misc/linux-firmware-amd-misc_20260319.git217ca6e4-0ubuntu1.1_all.deb  403  Forbidden [IP: 202.79.180.254 80]
Error: Unable to fetch some archives, maybe run apt update or try with --fix-missing?
```

Because `apt` downloads everything before configuring anything, **a single 403 blocks
all 132 pending upgrades**.

### Full verbatim terminal output

```console
$ apt update -y --fix-missing && apt-get update -y && apt upgrade -y
Hit:1 http://my.archive.ubuntu.com/ubuntu resolute InRelease
Hit:2 http://my.archive.ubuntu.com/ubuntu resolute-updates InRelease
Hit:3 http://my.archive.ubuntu.com/ubuntu resolute-backports InRelease
Hit:4 http://security.ubuntu.com/ubuntu resolute-security InRelease
Hit:5 https://ppa.launchpadcontent.net/apt-fast/stable/ubuntu resolute InRelease
Hit:6 https://ppa.launchpadcontent.net/git-core/ppa/ubuntu resolute InRelease
138 packages can be upgraded. Run 'apt list --upgradable' to see them.

Hit:1 http://my.archive.ubuntu.com/ubuntu resolute InRelease
Hit:2 http://my.archive.ubuntu.com/ubuntu resolute-updates InRelease
Hit:3 http://my.archive.ubuntu.com/ubuntu resolute-backports InRelease
Hit:4 http://security.ubuntu.com/ubuntu resolute-security InRelease
Hit:5 https://ppa.launchpadcontent.net/apt-fast/stable/ubuntu resolute InRelease
Hit:6 https://ppa.launchpadcontent.net/git-core/ppa/ubuntu resolute InRelease
Reading package lists... Done

The following packages were automatically installed and are no longer required:
  linux-headers-7.0.0-14                 linux-main-modules-zfs-7.0.0-14-generic  linux-tools-7.0.0-14-generic
  linux-headers-7.0.0-14-generic         linux-modules-7.0.0-14-generic
  linux-image-unsigned-7.0.0-14-generic  linux-tools-7.0.0-14
Use 'sudo apt autoremove' to remove them.

Get more security updates through Ubuntu Pro with 'esm-apps' enabled:
  libavcodec62 libavfilter11 libavutil60 libswscale9 libswresample6
  libavformat62
Learn more about Ubuntu Pro at https://ubuntu.com/pro

Upgrading:
  alsa-ucm-conf  apparmor  apport  ...  (132 packages, GNOME/mesa/linux-firmware/sssd/grub/snapd stack)

Installing dependencies:
  gst-audio-thumbnailer  gst-video-thumbnailer

Not upgrading yet due to phasing:
  base-files              nautilus       python3-software-properties
  libnautilus-extension4  nautilus-data  software-properties-common

Summary:
  Upgrading: 132, Installing: 2, Removing: 0, Not Upgrading: 6
  Download size: 584 kB / 890 MB
  Space needed: 3,958 kB / 86.5 GB available

Err:1 http://my.archive.ubuntu.com/ubuntu resolute-updates/main amd64 linux-firmware-amd-misc all 20260319.git217ca6e4-0ubuntu1.1
  403  Forbidden [IP: 202.79.180.254 80]
Error: Failed to fetch http://my.archive.ubuntu.com/ubuntu/pool/main/l/linux-firmware-amd-misc/linux-firmware-amd-misc_20260319.git217ca6e4-0ubuntu1.1_all.deb  403  Forbidden [IP: 202.79.180.254 80]
Error: Unable to fetch some archives, maybe run apt update or try with --fix-missing?
```

---

## 2. Environment facts

| Item | Value |
| --- | --- |
| Distro | Ubuntu 26.04 LTS, codename `resolute` |
| Architecture | amd64 |
| State | fresh install, first big upgrade |
| Physical location / ISP | Malaysia |
| Archive mirror in use | `http://my.archive.ubuntu.com/ubuntu` (Malaysia country redirect) |
| Mirror IP seen in the error | `202.79.180.254`, port `80`, plain HTTP |
| Security suite | `http://security.ubuntu.com/ubuntu resolute-security` (works) |
| Third-party sources | PPA `apt-fast/stable`, PPA `git-core/ppa` (both `Hit:` fine) |
| Extra tooling installed | `apt-fast` (aria2-based parallel downloader) |
| Pending work | 138 upgradable, 132 to upgrade, 2 new deps, 890 MB total, 584 kB already cached |
| Disk | 86.5 GB free — not a space problem |
| Held back | 6 packages "due to phasing" (`base-files`, `nautilus`, …) |
| Orphans | `linux-*-7.0.0-14*` flagged by `autoremove` |
| Failing package | `linux-firmware-amd-misc 20260319.git217ca6e4-0ubuntu1.1` (`_all.deb`) |

---

## 3. Diagnosis

**What the symptoms tell us:**

1. **All `InRelease` files download fine.** So DNS, routing, and the mirror's
   `dists/` tree are reachable. This is *not* a network outage and *not* an
   expired-key / signature problem.
2. **The failure is HTTP `403 Forbidden`, on a single file in `pool/`.**
   - `404 Not Found` would mean "mirror is stale / file not synced yet".
   - Timeout / `Could not resolve` would mean network.
   - **`403` means something in the path actively refused to serve the file** —
     either the mirror node itself (misconfigured vhost, partial sync where the
     directory exists but permissions are wrong, rate limiting / geo or UA
     blocking) or a **transparent HTTP proxy / ISP cache** sitting on port 80
     between you and `202.79.180.254`.
3. **Everything is plain HTTP (port 80).** That is exactly the traffic a
   transparent ISP cache can intercept and rewrite. HTTPS traffic to the PPAs
   works perfectly — a strong hint that the problem is HTTP-specific.

**Conclusion:** this is a **mirror-selection / mirror-health issue tied to the
Malaysia country mirror**, not a broken package and not an Ubuntu bug. Yes — it is
effectively a "country issue" in the sense that `my.archive.ubuntu.com` is the
Malaysia redirect and that particular node (or your ISP's HTTP path to it) is
refusing pool downloads. The primary fix is to **stop using that mirror over
plain HTTP.**

### Where the sources live on Ubuntu 26.04 (important)

Ubuntu 24.04 and later use the **deb822** format. There is **no useful
`/etc/apt/sources.list`** anymore. Your archive definition is here:

```
/etc/apt/sources.list.d/ubuntu.sources
```

and looks roughly like this:

```deb822
Types: deb
URIs: http://my.archive.ubuntu.com/ubuntu/
Suites: resolute resolute-updates resolute-backports
Components: main restricted universe multiverse
Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg

Types: deb
URIs: http://security.ubuntu.com/ubuntu/
Suites: resolute-security
Components: main restricted universe multiverse
Signed-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg
```

Every fix below edits the **first** stanza's `URIs:` line. Leave the
`security.ubuntu.com` stanza pointing at Canonical — security updates should come
straight from the source, and it is already working.

### About "Not upgrading yet due to phasing"

`base-files`, `nautilus`, `software-properties-common`, etc. are **not an error**.
Ubuntu rolls updates out to a percentage of machines at a time (phased updates).
Those packages will appear on their own in a few days. Do **not** force them.

### About the `linux-*-7.0.0-14` orphans

Those are an older kernel and its modules/headers. Safe to remove **only after**
you have booted successfully into the newer kernel. See §8.

---

## 4. Solutions — 12 options, cheapest first

Run everything as root (`sudo -i`) or prefix with `sudo`.

Before touching anything, **back up your apt config once**:

```bash
sudo cp -a /etc/apt /etc/apt.bak.$(date +%Y%m%d-%H%M%S)
ls -d /etc/apt.bak.*
```

---

### Solution 1 — Clear the partial cache and retry (30 seconds, try first)

Sometimes a truncated/poisoned partial file or a momentary mirror hiccup is all it is.

```bash
sudo apt clean
sudo rm -rf /var/lib/apt/lists/partial/* /var/cache/apt/archives/partial/*
sudo apt update
sudo apt upgrade -y
```

- **Use when:** always, as the first attempt.
- **Trade-off:** re-downloads metadata (a few MB). Won't help if the mirror is
  genuinely refusing.

---

### Solution 2 — Switch to the global `archive.ubuntu.com` ⭐ recommended

This is the fix that resolves the overwhelming majority of these 403s.

```bash
sudo cp /etc/apt/sources.list.d/ubuntu.sources /etc/apt/sources.list.d/ubuntu.sources.bak
sudo sed -i 's|http://my\.archive\.ubuntu\.com/ubuntu|http://archive.ubuntu.com/ubuntu|g' \
  /etc/apt/sources.list.d/ubuntu.sources
grep -n 'URIs' /etc/apt/sources.list.d/ubuntu.sources
sudo apt clean && sudo apt update && sudo apt upgrade -y
```

Rollback:

```bash
sudo mv /etc/apt/sources.list.d/ubuntu.sources.bak /etc/apt/sources.list.d/ubuntu.sources
sudo apt update
```

- **Use when:** the country mirror is the suspect (your case).
- **Trade-off:** `archive.ubuntu.com` is UK/global-anycast; from Malaysia it may be
  slower than a good local mirror, but it is authoritative and reliable.

---

### Solution 3 — Switch to a specific known-good regional mirror (test first)

Good Asia-Pacific mirrors, all HTTPS-capable:

| Mirror | URL |
| --- | --- |
| 0x.sg (Singapore) | `https://mirror.0x.sg/ubuntu/` |
| NUS/Singapore | `https://download.nus.edu.sg/mirror/ubuntu/` |
| JAIST (Japan) | `https://ftp.jaist.ac.jp/pub/Linux/ubuntu/` |
| TUNA (China) | `https://mirrors.tuna.tsinghua.edu.cn/ubuntu/` |
| Kartolo (Indonesia) | `http://kartolo.sby.datautama.net.id/ubuntu/` |
| Rackspace (global) | `https://mirror.rackspace.com/ubuntu/` |
| UPM (Malaysia, alt) | `https://mirror.upm.edu.my/ubuntu/` |

**Test reachability against the exact file that failed** before committing:

```bash
FILE="pool/main/l/linux-firmware-amd-misc/linux-firmware-amd-misc_20260319.git217ca6e4-0ubuntu1.1_all.deb"
for M in \
  https://mirror.0x.sg/ubuntu/ \
  https://download.nus.edu.sg/mirror/ubuntu/ \
  https://ftp.jaist.ac.jp/pub/Linux/ubuntu/ \
  https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ \
  https://mirror.rackspace.com/ubuntu/ \
  http://archive.ubuntu.com/ubuntu/ ; do
  CODE=$(curl -o /dev/null -sIL -w '%{http_code}' --max-time 10 "${M}${FILE}")
  printf '%-55s %s\n' "$M" "$CODE"
done
```

Then apply the one that printed `200`:

```bash
NEW="https://mirror.0x.sg/ubuntu"     # <-- change to your winner
sudo cp /etc/apt/sources.list.d/ubuntu.sources /etc/apt/sources.list.d/ubuntu.sources.bak
sudo sed -i "s|http://my\.archive\.ubuntu\.com/ubuntu|${NEW}|g" \
  /etc/apt/sources.list.d/ubuntu.sources
sudo apt clean && sudo apt update && sudo apt upgrade -y
```

- **Use when:** you want local speed and the country redirect is bad.
- **Trade-off:** third-party mirrors can lag a few hours behind. Packages are still
  GPG-verified, so a stale mirror is safe, just occasionally out of date.

---

### Solution 4 — Force HTTPS instead of HTTP (defeats transparent proxies)

If an ISP/corporate transparent cache on port 80 is injecting the 403, HTTPS bypasses it.

```bash
sudo apt install -y ca-certificates
sudo sed -i 's|http://my\.archive\.ubuntu\.com/ubuntu|https://my.archive.ubuntu.com/ubuntu|g' \
  /etc/apt/sources.list.d/ubuntu.sources
sudo apt clean && sudo apt update && sudo apt upgrade -y
```

If `my.archive.ubuntu.com` has no valid TLS, use HTTPS on `archive.ubuntu.com`:

```bash
sudo sed -i 's|https\?://my\.archive\.ubuntu\.com/ubuntu|https://archive.ubuntu.com/ubuntu|g' \
  /etc/apt/sources.list.d/ubuntu.sources
sudo apt clean && sudo apt update
```

- **Use when:** metadata works but pool files 403 — the classic MITM-cache signature.
- **Trade-off:** slightly more CPU, no downside otherwise. Recommended permanently.

---

### Solution 5 — Force IPv4 / IPv6, or blackhole the bad mirror node

The 403 named a specific node (`202.79.180.254`). If the mirror is DNS round-robin,
forcing the other address family, or making that IP unreachable so DNS falls
through, can be enough.

```bash
# Try IPv4 only
sudo apt -o Acquire::ForceIPv4=true update
sudo apt -o Acquire::ForceIPv4=true upgrade -y

# Or IPv6 only
sudo apt -o Acquire::ForceIPv6=true upgrade -y
```

Make it permanent:

```bash
echo 'Acquire::ForceIPv4 "true";' | sudo tee /etc/apt/apt.conf.d/99force-ipv4
```

See which addresses the mirror has, and route the bad one nowhere:

```bash
getent ahosts my.archive.ubuntu.com
# temporary, non-persistent:
sudo ip route add blackhole 202.79.180.254/32
# undo with:  sudo ip route del blackhole 202.79.180.254/32
```

- **Use when:** the mirror hostname resolves to several IPs and only one is broken.
- **Trade-off:** blackhole is a blunt instrument; remember to remove it.

---

### Solution 6 — Let apt pick a mirror automatically (`mirror://`)

```bash
sudo cp /etc/apt/sources.list.d/ubuntu.sources /etc/apt/sources.list.d/ubuntu.sources.bak
sudo sed -i 's|http://my\.archive\.ubuntu\.com/ubuntu|mirror+http://mirrors.ubuntu.com/mirrors.txt|g' \
  /etc/apt/sources.list.d/ubuntu.sources
sudo apt clean && sudo apt update && sudo apt upgrade -y
```

- **Use when:** you want apt to fail over between mirrors by itself.
- **Trade-off:** the geo list can hand you the same bad Malaysian node. Also makes
  debugging harder because the effective URL changes per run.

---

### Solution 7 — Benchmark mirrors with `netselect-apt` / `apt-select`

```bash
# Option A
sudo apt install -y netselect-apt
sudo netselect-apt -n -c MY -t 20 resolute        # writes ./sources.list
cat sources.list

# Option B (Python, understands the Launchpad mirror list)
sudo apt install -y python3-pip
pip3 install --user apt-select
~/.local/bin/apt-select --country MY --top-number 10 --min-status up-to-date
```

Then translate the winning URL into `ubuntu.sources` using the `sed` from Solution 3.

- **Use when:** you want the objectively fastest healthy mirror.
- **Trade-off:** these tools emit the *legacy* `sources.list` format — you must copy
  only the URL into the deb822 file, not the whole generated file.

---

### Solution 8 — Fetch the one failing `.deb` by hand, then continue

Useful when you don't want to change mirrors at all.

```bash
cd /tmp
PKG="linux-firmware-amd-misc_20260319.git217ca6e4-0ubuntu1.1_all.deb"
# try several sources until one downloads
wget "http://archive.ubuntu.com/ubuntu/pool/main/l/linux-firmware-amd-misc/$PKG" \
  || wget "https://mirror.0x.sg/ubuntu/pool/main/l/linux-firmware-amd-misc/$PKG" \
  || wget "https://ftp.jaist.ac.jp/pub/Linux/ubuntu/pool/main/l/linux-firmware-amd-misc/$PKG"

# stage it into apt's cache so apt uses it instead of downloading
sudo cp "$PKG" /var/cache/apt/archives/
sudo apt upgrade -y
```

Or install directly:

```bash
sudo apt install -y "/tmp/$PKG"
sudo apt --fix-broken install -y
sudo apt upgrade -y
```

You can also let apt tell you every URL it wants, and feed them to `wget`:

```bash
sudo apt-get --print-uris -y upgrade | awk -F"'" '/^\x27http/ {print $2}' > /tmp/urls.txt
wc -l /tmp/urls.txt
cd /var/cache/apt/archives && sudo wget -c -i /tmp/urls.txt
sudo apt upgrade -y
```

- **Use when:** exactly one or two files 403 and you need to move on now.
- **Trade-off:** manual, doesn't fix the root cause; the next upgrade may 403 again.

---

### Solution 9 — Isolate third-party sources (PPAs)

Rule out `apt-fast` and `git-core` PPAs, and anything else in `sources.list.d`.

```bash
ls -l /etc/apt/sources.list.d/
# disable everything except the official ubuntu.sources
for f in /etc/apt/sources.list.d/*.list /etc/apt/sources.list.d/*.sources; do
  case "$(basename "$f")" in ubuntu.sources) continue;; esac
  [ -e "$f" ] && sudo mv "$f" "$f.disabled"
done
sudo apt update && sudo apt upgrade -y
```

Re-enable afterwards:

```bash
for f in /etc/apt/sources.list.d/*.disabled; do sudo mv "$f" "${f%.disabled}"; done
sudo apt update
```

- **Use when:** you suspect a PPA is shadowing or redirecting packages.
- **Trade-off:** packages installed from those PPAs won't see updates while disabled.

---

### Solution 10 — Rule out `apt-fast` / aria2

`apt-fast` is installed on this machine. It downloads with aria2 in parallel, which
changes User-Agent, uses many connections, and can trip mirror rate limits — a very
common cause of sudden `403`s.

```bash
# 1. Confirm plain apt behaves differently
sudo /usr/bin/apt update
sudo /usr/bin/apt upgrade -y

# 2. Check for an apt-fast/aria2 apt hook
grep -rn 'aria2\|apt-fast' /etc/apt/apt.conf.d/ 2>/dev/null
cat /etc/apt-fast.conf 2>/dev/null | grep -v '^#' | grep -v '^$'

# 3. Reduce parallelism in /etc/apt-fast.conf
#    _MAXNUM=5   -> _MAXNUM=2
#    _MAXCONPERSRV=10 -> _MAXCONPERSRV=2
sudo sed -i 's/^_MAXNUM=.*/_MAXNUM=2/; s/^_MAXCONPERSRV=.*/_MAXCONPERSRV=2/' /etc/apt-fast.conf

# 4. Or remove it entirely while debugging
sudo apt remove -y apt-fast
```

Also cap plain apt's own parallelism and set a normal UA:

```bash
sudo tee /etc/apt/apt.conf.d/99tuning >/dev/null <<'EOF'
Acquire::Queue-Mode "access";
Acquire::http::Pipeline-Depth "0";
Acquire::Retries "3";
EOF
sudo apt update && sudo apt upgrade -y
```

- **Use when:** `apt-fast` is present (it is) — always worth eliminating.
- **Trade-off:** slower downloads with pipelining off; re-enable once fixed.

---

### Solution 11 — Route around the ISP: proxy, VPN, or a local apt cache

If HTTPS also 403s and multiple mirrors fail identically, your ISP or a network
appliance is filtering.

```bash
# a) explicit HTTP proxy for apt only
sudo tee /etc/apt/apt.conf.d/95proxy >/dev/null <<'EOF'
Acquire::http::Proxy "http://PROXY_HOST:PORT/";
Acquire::https::Proxy "http://PROXY_HOST:PORT/";
EOF
sudo apt update
# remove with: sudo rm /etc/apt/apt.conf.d/95proxy

# b) VPN / WireGuard: bring the tunnel up, then retry
sudo apt update && sudo apt upgrade -y

# c) run your own caching proxy on the LAN (also great for multiple machines)
sudo apt install -y apt-cacher-ng
# then on clients:
echo 'Acquire::http::Proxy "http://CACHE_HOST:3142/";' | sudo tee /etc/apt/apt.conf.d/01proxy
```

- **Use when:** every mirror fails the same way from this network.
- **Trade-off:** requires infrastructure; a VPN changes your effective geo, which is
  itself a diagnostic (if it works over VPN, the ISP path was the culprit).

---

### Solution 12 — Ubuntu Pro / ESM (informational, unrelated to the 403)

The output suggested `esm-apps` for `libavcodec62` and friends. This is a separate,
free-for-personal-use offer and has nothing to do with the 403.

```bash
sudo pro status
# free personal token from https://ubuntu.com/pro/dashboard
sudo pro attach <YOUR_TOKEN>
sudo pro enable esm-apps
sudo apt update && sudo apt upgrade -y
```

- **Use when:** you want extended security maintenance for universe packages.
- **Trade-off:** requires a Canonical account; pulls from `esm.ubuntu.com`, another
  network dependency.

---

## 5. One do-everything script — `fix-apt-403.sh`

Backs up your sources, probes candidate mirrors against a real pool file, rewrites
`ubuntu.sources` to the first mirror that answers `200`, cleans the cache, and runs
the upgrade. Idempotent and non-destructive (it never removes packages).

```bash
#!/usr/bin/env bash
# fix-apt-403.sh — recover from "403 Forbidden" apt pool download failures on
# Ubuntu 24.04+ (deb822 sources). Safe to re-run.
set -euo pipefail

[[ ${EUID} -eq 0 ]] || { echo "Run with sudo: sudo bash $0" >&2; exit 1; }

SRC="/etc/apt/sources.list.d/ubuntu.sources"
STAMP="$(date +%Y%m%d-%H%M%S)"
BACKUP="${SRC}.bak.${STAMP}"
CODENAME="$(. /etc/os-release && echo "${VERSION_CODENAME}")"

# A file that is guaranteed to exist in every complete mirror.
PROBE="dists/${CODENAME}/Release"

CANDIDATES=(
  "https://mirror.0x.sg/ubuntu"
  "https://download.nus.edu.sg/mirror/ubuntu"
  "https://ftp.jaist.ac.jp/pub/Linux/ubuntu"
  "https://mirrors.tuna.tsinghua.edu.cn/ubuntu"
  "https://mirror.rackspace.com/ubuntu"
  "https://archive.ubuntu.com/ubuntu"
  "http://archive.ubuntu.com/ubuntu"
)

echo "==> Codename: ${CODENAME}"
echo "==> Backing up apt config"
cp -a /etc/apt "/etc/apt.bak.${STAMP}"
cp -a "${SRC}" "${BACKUP}"
echo "    ${BACKUP}"

CURRENT="$(awk '/^URIs:/ {print $2; exit}' "${SRC}")"
echo "==> Current archive URI: ${CURRENT}"

echo "==> Probing candidate mirrors (HTTP status for /${PROBE})"
WINNER=""
for M in "${CANDIDATES[@]}"; do
  CODE="$(curl -o /dev/null -sIL --max-time 12 -w '%{http_code}' "${M}/${PROBE}" || echo 000)"
  printf '    %-52s %s\n' "${M}" "${CODE}"
  if [[ -z "${WINNER}" && "${CODE}" == "200" ]]; then WINNER="${M}"; fi
done

[[ -n "${WINNER}" ]] || { echo "!! No candidate mirror answered 200. Network or DNS problem — see Solution 11." >&2; exit 2; }
echo "==> Selected mirror: ${WINNER}"

echo "==> Rewriting archive URI (security.ubuntu.com is left untouched)"
python3 - "$SRC" "$WINNER" <<'PY'
import re, sys
path, new = sys.argv[1], sys.argv[2].rstrip('/') + '/'
text = open(path).read()
def fix(m):
    uri = m.group(1)
    return m.group(0) if 'security.ubuntu.com' in uri else f'URIs: {new}'
open(path, 'w').write(re.sub(r'URIs:\s*(\S+)', fix, text))
PY
grep -n 'URIs:' "${SRC}"

echo "==> Sensible acquire defaults"
cat >/etc/apt/apt.conf.d/99fix-403 <<'EOF'
Acquire::Retries "3";
Acquire::http::Pipeline-Depth "0";
Acquire::Queue-Mode "access";
EOF

echo "==> Clearing caches"
apt clean
rm -rf /var/lib/apt/lists/partial/* /var/cache/apt/archives/partial/* 2>/dev/null || true

echo "==> apt update"
apt update

echo "==> apt full-upgrade"
DEBIAN_FRONTEND=noninteractive apt -y -o Dpkg::Options::=--force-confdef \
  -o Dpkg::Options::=--force-confold full-upgrade

echo
echo "==> Done. Rollback if needed:"
echo "    sudo cp ${BACKUP} ${SRC} && sudo rm -f /etc/apt/apt.conf.d/99fix-403 && sudo apt update"
[[ -f /var/run/reboot-required ]] && echo "==> A reboot is required."
exit 0
```

Save and run:

```bash
mkdir -p ~/bin && nano ~/bin/fix-apt-403.sh   # paste the script
chmod +x ~/bin/fix-apt-403.sh
sudo bash ~/bin/fix-apt-403.sh
```

---

## 6. Verification checklist

```bash
# 1. Metadata refreshes with no Err: lines
sudo apt update

# 2. Dependency tree is consistent
sudo apt-get check

# 3. Nothing half-configured
sudo dpkg --audit
sudo dpkg --configure -a

# 4. The upgrade actually completes
sudo apt full-upgrade -y

# 5. What (if anything) is still pending
apt list --upgradable

# 6. The specific package that used to fail
apt policy linux-firmware-amd-misc
dpkg -l linux-firmware-amd-misc | tail -n 3

# 7. Which mirror is in effect
grep -n 'URIs:' /etc/apt/sources.list.d/ubuntu.sources

# 8. Reboot needed?
[ -f /var/run/reboot-required ] && cat /var/run/reboot-required*
```

Success looks like: `apt update` with only `Hit:`/`Get:` lines, `apt full-upgrade`
finishing with `0 upgraded, 0 newly installed`, and `apt list --upgradable` showing
only phased packages (or nothing).

---

## 7. Routine maintenance script — `ubuntu-update.sh`

For day-to-day use on the fresh install.

```bash
#!/usr/bin/env bash
# ubuntu-update.sh — routine full system update for Ubuntu 24.04+
set -euo pipefail
[[ ${EUID} -eq 0 ]] || { echo "Run with sudo: sudo bash $0" >&2; exit 1; }

export DEBIAN_FRONTEND=noninteractive
APT_OPTS=(-y -o Dpkg::Options::=--force-confdef -o Dpkg::Options::=--force-confold)

echo "==> [1/7] Refreshing package lists"
apt update

echo "==> [2/7] Upgrading packages"
apt "${APT_OPTS[@]}" full-upgrade

echo "==> [3/7] Fixing anything half-installed"
dpkg --configure -a || true
apt "${APT_OPTS[@]}" --fix-broken install || true

echo "==> [4/7] Removing orphans (kernels kept per policy below)"
apt "${APT_OPTS[@]}" autoremove --purge
apt clean

echo "==> [5/7] Refreshing snaps"
command -v snap >/dev/null && snap refresh || echo "    snap not installed"

echo "==> [6/7] Firmware"
if command -v fwupdmgr >/dev/null; then
  fwupdmgr refresh --force || true
  fwupdmgr get-updates || echo "    no firmware updates"
fi

echo "==> [7/7] Status"
echo "    Kernel running : $(uname -r)"
echo "    Kernels present:"; dpkg -l 'linux-image-*' 2>/dev/null | awk '/^ii/ {print "      " $2}'
if [[ -f /var/run/reboot-required ]]; then
  echo "    *** REBOOT REQUIRED ***"
  cat /var/run/reboot-required.pkgs 2>/dev/null | sed 's/^/      /'
else
  echo "    No reboot required."
fi
echo "==> Complete."
```

Optional weekly timer:

```bash
sudo cp ubuntu-update.sh /usr/local/sbin/ubuntu-update.sh
sudo chmod +x /usr/local/sbin/ubuntu-update.sh
sudo tee /etc/systemd/system/ubuntu-update.service >/dev/null <<'EOF'
[Unit]
Description=Routine Ubuntu update
[Service]
Type=oneshot
ExecStart=/usr/local/sbin/ubuntu-update.sh
EOF
sudo tee /etc/systemd/system/ubuntu-update.timer >/dev/null <<'EOF'
[Unit]
Description=Weekly Ubuntu update
[Timer]
OnCalendar=Sun 03:00
Persistent=true
[Install]
WantedBy=timers.target
EOF
sudo systemctl daemon-reload && sudo systemctl enable --now ubuntu-update.timer
```

---

## 8. Safety notes

- **Never `autoremove` kernel packages for the kernel you are currently booted into.**
  Check first:
  ```bash
  uname -r
  dpkg -l 'linux-image-*' | awk '/^ii/ {print $2}'
  ```
  Keep **at least two** kernels (current + one previous) so you can recover from a
  bad boot. Reboot into the new kernel *before* purging the old one.
- **Always back up `/etc/apt` before editing** — `sudo cp -a /etc/apt /etc/apt.bak.$(date +%F)`.
  Every solution above includes a rollback line.
- **Do not disable GPG verification.** `--allow-unauthenticated`,
  `[trusted=yes]`, and `Acquire::AllowInsecureRepositories` are *not* fixes for a
  403 and open you to package tampering.
- **Don't force phased packages.** Wait or `apt install <pkg>` explicitly only if
  you have a specific reason.
- **`grub-pc` / `grub-common` upgrades may prompt** for the install device. Answer
  it interactively at least once rather than blindly using `--force-confold` on a
  boot-critical package.
- **Snap and apt are separate.** `apt upgrade` does not update snaps; use `snap refresh`.
- Prefer `full-upgrade` over `upgrade` on a fresh install — `upgrade` refuses to
  remove packages and can stall on transitions.

---

## 9. Copy-paste block for another AI

```text
CONTEXT
I have a freshly installed Ubuntu 26.04 LTS (codename "resolute"), amd64, located
in Malaysia, on a residential/office ISP. apt uses the deb822 format; my archive is
defined in /etc/apt/sources.list.d/ubuntu.sources with
URIs: http://my.archive.ubuntu.com/ubuntu (the Malaysia country mirror redirect).
Security comes from http://security.ubuntu.com/ubuntu resolute-security.
I also have two PPAs enabled: ppa:apt-fast/stable and ppa:git-core/ppa, and the
apt-fast (aria2 parallel downloader) package is installed.

PROBLEM
`sudo apt update` succeeds completely — every repository returns "Hit:", including
the PPAs over HTTPS. But `sudo apt upgrade` (132 packages, 890 MB) aborts while
downloading with:

  Err:1 http://my.archive.ubuntu.com/ubuntu resolute-updates/main amd64
        linux-firmware-amd-misc all 20260319.git217ca6e4-0ubuntu1.1
    403  Forbidden [IP: 202.79.180.254 80]
  Error: Failed to fetch http://my.archive.ubuntu.com/ubuntu/pool/main/l/
    linux-firmware-amd-misc/linux-firmware-amd-misc_20260319.git217ca6e4-0ubuntu1.1_all.deb
    403  Forbidden [IP: 202.79.180.254 80]
  Error: Unable to fetch some archives, maybe run apt update or try with --fix-missing?

Other facts from the same run:
- 138 upgradable, 132 to upgrade, 2 new deps, 0 removals, 6 held back "due to phasing"
  (base-files, nautilus, nautilus-data, libnautilus-extension4,
   python3-software-properties, software-properties-common).
- 86.5 GB free disk; only 3.9 MB needed. Not a space issue.
- Old kernel packages linux-*-7.0.0-14* are flagged as autoremovable.
- Ubuntu Pro esm-apps is suggested for libavcodec62/libavfilter11/libavutil60/
  libswscale9/libswresample6/libavformat62.
- I ran: apt update -y --fix-missing && apt-get update -y && apt upgrade -y

WHAT I THINK IS HAPPENING
Metadata (dists/.../InRelease) downloads fine, so DNS/routing/keys are OK. The
failure is HTTP 403 on a single file under pool/ over plain HTTP port 80 — an
active refusal, not a 404 (stale mirror) and not a timeout (network). Likely the
Malaysia mirror node 202.79.180.254 is misconfigured/rate-limiting, or a
transparent ISP HTTP proxy is intercepting port 80. HTTPS to the PPAs works.

WHAT I WANT
Concrete shell commands to (a) confirm the root cause, and (b) fix it permanently,
using the deb822 file /etc/apt/sources.list.d/ubuntu.sources — not the legacy
/etc/apt/sources.list. Please keep security.ubuntu.com pointing at Canonical, do
not disable GPG verification, and warn me before anything touches kernel packages.
```

---

## 10. Quick reference — the 60-second fix

```bash
sudo cp -a /etc/apt /etc/apt.bak.$(date +%F)
sudo sed -i 's|http://my\.archive\.ubuntu\.com/ubuntu|https://archive.ubuntu.com/ubuntu|g' \
  /etc/apt/sources.list.d/ubuntu.sources
sudo apt clean
sudo apt update
sudo apt full-upgrade -y
[ -f /var/run/reboot-required ] && echo "Reboot required" || echo "All good"
```
