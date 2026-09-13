#!/usr/bin/env python3
"""
Purge GitHub Actions Artifacts Script
Autonomously discovers and deletes stored GitHub Actions artifacts via the GitHub API
across target repositories to enforce the Zero-Storage Actions Mandate and keep account
storage usage at 0.0 GB (staying well within the 0.5 GB free quota).

Usage:
  # Purge artifacts from default repositories (coding-guidelines-v24 and gitmap-v28):
  python 03-ai-scripts/34-purge-github-actions-artifacts.py

  # Purge artifacts from a specific repository:
  python 03-ai-scripts/34-purge-github-actions-artifacts.py --repo alimtvnetwork/gitmap-v28
"""

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
import json
import subprocess
import sys
import time

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

TARGET_REPOS = [
    "alimtvnetwork/coding-guidelines-v24",
    "alimtvnetwork/gitmap-v28",
]

def fetch_artifact_batch(repo: str, page: int = 1, per_page: int = 100):
    cmd = [
        "gh", "api",
        f"repos/{repo}/actions/artifacts?per_page={per_page}&page={page}",
        "--jq", "{total_count: .total_count, artifacts: [.artifacts[] | {id: .id, name: .name, size: .size_in_bytes}]}"
    ]
    res = subprocess.run(cmd, capture_output=True, text=True, encoding="utf-8", errors="replace")
    if res.returncode != 0:
        return None
    try:
        return json.loads(res.stdout)
    except Exception:
        return None

def delete_single_artifact(repo: str, artifact_id: int) -> bool:
    cmd = ["gh", "api", "-X", "DELETE", f"repos/{repo}/actions/artifacts/{artifact_id}"]
    res = subprocess.run(cmd, capture_output=True, text=True, encoding="utf-8", errors="replace")
    return res.returncode == 0

def purge_repo_artifacts(repo: str, max_workers: int = 12):
    print(f"\n========================================================", flush=True)
    print(f" Scanning artifacts for repository: {repo}", flush=True)
    print(f"========================================================", flush=True)

    first_batch = fetch_artifact_batch(repo, page=1, per_page=1)
    if not first_batch:
        print(f"❌ Failed to query artifacts for {repo} or repository not accessible.", flush=True)
        return

    total_count = first_batch.get("total_count", 0)
    print(f"Found {total_count} stored artifact(s) in {repo}.", flush=True)
    if total_count == 0:
        print(f"✅ Repository {repo} is already at 0 bytes Actions storage!", flush=True)
        return

    deleted_count = 0
    freed_bytes = 0

    while True:
        batch = fetch_artifact_batch(repo, page=1, per_page=100)
        if not batch or not batch.get("artifacts"):
            break

        artifacts = batch["artifacts"]
        if not artifacts:
            break

        batch_size = len(artifacts)
        print(f"Deleting batch of {batch_size} artifact(s) using {max_workers} worker threads...", flush=True)

        with ThreadPoolExecutor(max_workers=max_workers) as executor:
            future_to_art = {
                executor.submit(delete_single_artifact, repo, art["id"]): art
                for art in artifacts
            }
            for future in as_completed(future_to_art):
                art = future_to_art[future]
                is_success = future.result()
                if is_success:
                    deleted_count += 1
                    freed_bytes += art.get("size", 0)

        mb_freed = freed_bytes / (1024 * 1024)
        print(f"Progress: {deleted_count}/{total_count} deleted (~{mb_freed:.2f} MB freed)...", flush=True)

        # Brief pause to avoid API rate limiting
        time.sleep(0.5)

    mb_freed = freed_bytes / (1024 * 1024)
    print(f"[SUCCESS] Completed purge for {repo}: {deleted_count} artifacts deleted, ~{mb_freed:.2f} MB storage recovered.", flush=True)

def main():
    parser = argparse.ArgumentParser(description="Purge GitHub Actions artifacts to maintain 0 storage usage.")
    parser.add_argument("--repo", help="Target repository (owner/repo). If omitted, purges all standard repositories.")
    parser.add_argument("--workers", type=int, default=12, help="Number of concurrent deletion threads.")
    args = parser.parse_args()

    repos = [args.repo] if args.repo else TARGET_REPOS
    for repo in repos:
        purge_repo_artifacts(repo, max_workers=args.workers)

if __name__ == "__main__":
    main()
