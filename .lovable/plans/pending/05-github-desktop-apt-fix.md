# 05-github-desktop-apt-fix: Fix GitHub Desktop APT Installation

## 1. Context and Problem Statement
The user reported that gitmap install github-desktop fails on Linux (pt) with exit status 100 during pt install github-desktop. 

Root cause: 
The old Shiftkey APT repository (pt.packages.shiftkey.dev) has certificate errors or is deprecated. wget -qO - fails silently, resulting in an empty GPG key. Subsequently, pt update fails to fetch the package list, and pt install fails because the package is not found.

## 2. Proposed Changes
- Replace https://apt.packages.shiftkey.dev/ubuntu/any/ANY.gpg with the new official community mirror https://mirror.mwt.me/ghd/gpgkey.
- Replace https://apt.packages.shiftkey.dev/ubuntu/ with https://mirror.mwt.me/ghd/deb/.
- Ensure wget does not fail silently by removing the -q flag if possible (optional, but good practice). Actually, it is better to leave it but the URL fix is the main solution.

## 3. Subtasks
1. **Fix URLs in installtools.go**: Update the GPG key and repo URL in unInstallGitHubDesktopLinux.
2. **Commit and Release**: Bump version, update changelog, and push.
