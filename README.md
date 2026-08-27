# ReleaseFlow

ReleaseFlow is a small Go service for tracking software release packets as they move from draft to review, approval, publication, and rollback. It keeps each packet's state in a local journal so a restarted process can rebuild the current view without guessing. The CLI is meant for release coordinators who need a durable ledger of who reviewed what, when the packet was approved, and what was finally published.

## Commands

- `releaseflow init`
- `releaseflow submit`
- `releaseflow review`
- `releaseflow approve`
- `releaseflow publish`
- `releaseflow rollback`
- `releaseflow inspect`
- `releaseflow history`
- `releaseflow export`

## Storage

Release packets are stored in a JSONL journal and periodic snapshot files under the data directory. Rebuilds replay the journal into the in-memory model and then refresh the snapshot so the next startup is cheap and deterministic.

