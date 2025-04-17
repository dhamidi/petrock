# Worker Inspection

The Petrock framework now supports inspection of worker components through the `self inspect` command. This guide explains how to use this feature and how to provide custom information for your workers.

## Using the Self-Inspect Command

The `self inspect` command returns information about all registered application components, including workers:

```bash
./yourapp self inspect
```

This command returns a JSON object including worker information under the `workers` key.

## Worker Information

The worker inspection provides the following information for each worker:

- `name`: The name of the worker (from WorkerInfo or derived from type)
- `description`: A description of the worker's purpose (if provided)
- `type`: The Go type name of the worker
- `methods`: A list of available methods on the worker

## Making Your Workers Self-Describing

You can make your workers provide custom inspection information by implementing the optional `WorkerInfo()` method:

```go
// WorkerInfo provides self-description information
func (w *MyWorker) WorkerInfo() *core.WorkerInfo {
    return &core.WorkerInfo{
        Name:        "MyFeature Worker",
        Description: "Processes background tasks for my feature",
    }
}
```

If you don't implement this method, the framework will use reflection to extract basic information about your worker.

## Example Output

Here's an example of worker information in the self-inspect output:

```json
{
  "workers": [
    {
      "name": "Posts Worker",
      "description": "Handles background processing for the posts feature",
      "type": "*posts.Worker",
      "methods": ["Start", "Stop", "Work", "WorkerInfo"]
    }
  ]
}
```

## Benefits of Worker Inspection

- **Debugging**: Quickly see all registered workers in your application
- **Documentation**: Self-documenting code through worker descriptions
- **Tooling**: Build tools that can understand your application's worker structure
- **Monitoring**: Potential foundation for worker monitoring and health checking