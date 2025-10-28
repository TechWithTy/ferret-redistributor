# Workflows (Go)

Runtime orchestration helpers for scheduling and posting.

- `scheduler_workflow.go`: wraps atomic claim logic.

## Claim API
```
rows, err := workflows.Scheduler{DB: db}.Claim(ctx, 10*time.Minute, 50)
```
- Moves rows to `processing` atomically and returns them.
- Prevents duplicate posting when multiple runners are active.

Integrate with `go/cmd/scheduler` or your own service.

