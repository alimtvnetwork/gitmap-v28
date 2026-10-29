# GitMap UI/UX and Icon Improvement

Version: 6.630.0
Updated: 2026-10-30

## Purpose

The first mark in the shell is the GitMap icon from `assets/01-icons/01-gitmap.svg`. The same mark is the page icon and sits beside the page title. Navigation keeps a yellow bar on the active item while hover only darkens the background. Save sits at the top of Settings as a green button. Installer profiles render as a tree. Prompt rows name the instance they come from. An SSH timeout stays on the node as "Timed out" instead of falling back to a blank offline list.

## GitLab plan

No GitLab client code in this change. The follow-on adds multiple GitLab accounts beside the existing GitHub accounts and exposes the same basic GitMap actions there: clone, pull, commit, push, status, and release listing. Account pull is tested on the current GitHub path first. GitLab stays a plan until that account model is specified.

## Time-restriction documents

The restriction behavior stays in the commit path. The restriction spec and plan files are removed from the tree. Rewriting published git history on main is not part of this change.
