# Fal.ai Go SDK (Ferret Edition)

Thin wrapper around Fal.ai's queue API so Ferret services can launch AI workflows, poll for completion, and stream progress without re-implementing the HTTP plumbing each time.

## Features

- Minimal client with pluggable `http.Client`, queue base URL, and poll intervals
- Friendly helpers: `Call`, `PollUntilCompletion`, `PollWithProgress`, `Call.Cancel`
- Typed responses for images, timings, queue position, etc.
- Context-aware requests for cancellation/timeouts

## Usage

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()

cli := fal.NewClient(
  os.Getenv("FAL_API_KEY"),
  fal.WithPollInterval(750*time.Millisecond),
)

call, err := cli.Call(ctx, "fal-ai/image-tool", fal.RequestOptions{
  Method: http.MethodPost,
  Payload: map[string]any{
    "prompt": "retro robot surfing on mars",
  },
})
if err != nil {
  log.Fatalf("call failed: %v", err)
}

resp, err := cli.PollUntilCompletion(ctx, call)
if err != nil {
  log.Fatalf("polling failed: %v", err)
}

fmt.Println("image url:", resp.Images[0].URL)
```

To surface progress in a UI:

```go
progressCh := make(chan *fal.Response, 4)
go func() {
  for update := range progressCh {
    fmt.Printf("status=%s, queue=%d\n", update.Status, update.QueuePosition)
  }
}()

resp, err := cli.PollWithProgress(ctx, call, progressCh)
```

## Testing

`client_test.go` uses `httptest.Server` to simulate Fal endpoints (`/status`, `/response`, `/cancel`). No network calls are made.


