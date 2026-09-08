# Multi-agent workflow

For non-trivial tasks, the primary agent owns requirements, decomposition, architecture, complex debugging, cross-module decisions, result synthesis, and final decisions.

- Delegate codebase search, call-chain tracing, and dependency mapping to `explorer` first.
- Delegate bounded implementation and clearly diagnosed fixes to `developer`.
- After non-trivial changes, delegate an independent read-only review to `reviewer`.
- When work reaches verification and commit, delegate to `tester`. Say explicitly whether the task includes committing; otherwise tester must only test and report.
- Parallelize independent explorer or reviewer work when it materially improves speed or quality, but do not split trivial work solely to use subagents.
- The primary agent must evaluate all subagent results and remains responsible for final quality and decisions.
