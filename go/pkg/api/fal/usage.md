# Usage Snippet

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()

client := fal.NewClient(os.Getenv("FAL_API_KEY"))

call, err := client.Call(ctx, "fal-ai/image-tool", fal.RequestOptions{
	Method:  http.MethodPost,
	Payload: map[string]any{"prompt": "friendly robot on mars"},
})
if err != nil {
	log.Fatalf("fal call failed: %v", err)
}

resp, err := client.PollUntilCompletion(ctx, call)
if err != nil {
	if errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("Fal request timed out")
	}
	log.Fatalf("Fal request failed: %v", err)
}

fmt.Println(resp.Images[0].URL)
```


