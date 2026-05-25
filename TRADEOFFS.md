# Architectural Tradeoffs

## Concurrency Strategy

I chose sequential batch processing over a worker pool or full parallelism.

**Why:** Simpler to reason about, easier to debug, and safer for a first implementation. A worker pool would be faster but introduces goroutine coordination complexity.

**Tradeoff:** Slower throughput. A 500-tweet archive takes longer than it would with parallel workers. This is acceptable for a personal audit tool that runs infrequently.

## Batching Decision

Tweets are processed in batches of 20 per Gemini API call.

**Why:** Sending tweets one by one would be expensive in API calls and slow. Sending all at once risks hitting token limits. 20 is a reasonable middle ground.

**Tradeoff:** Batch size is not dynamically adjusted. A tweet with very long text could push a batch over the token limit.

## Error Handling

On a failed batch, the tool retries up to 3 times with a 5-second delay for rate limit errors (429). For any other error, it logs and skips the batch.

**Why:** Retrying on 429s is necessary since rate limits are temporary. Skipping non-recoverable errors keeps the pipeline moving rather than halting the entire audit.

**Tradeoff:** Skipped batches mean some tweets are never evaluated. The user is not explicitly notified which tweets were skipped.

## Prompt Design

The Gemini prompt encodes user criteria as structured rules — forbidden words, forbidden phrases, outdated opinions, and custom plain-English rules. Gemini is instructed to return a JSON array of flagged tweets.

**Why:** Plain-English custom rules give the user flexibility without requiring them to write code. Requesting JSON output makes parsing deterministic.

**Tradeoff:** Gemini can still return malformed JSON or deviate from the schema. The parser handles this gracefully but some responses may be lost.

## Performance vs Safety

The tool prioritizes safety over speed — sequential processing, retries before skipping, and incremental CSV writes ensure no data is lost mid-run.

**Tradeoff:** The tool is slower than it could be with concurrency.
