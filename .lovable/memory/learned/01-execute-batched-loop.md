The previous subagent batch execution encountered RESOURCE_EXHAUSTED (429) due to rate limits on the flash model and was subsequently interrupted by a server restart. We have switched to inherit model.
Strictly avoid using flash model for high concurrency batches that exceed rate limits.
