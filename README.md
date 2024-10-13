# Description

A fees API service that allows users to create a bill, add fees to the bill, and close a bill. This application is built for local development only.

Each bill is basically a Temporal workflow that is created when a new bill is created. The bill is then updated or closed via Temporal signals.


## Running

Ensure that the Temporal dev server is running locally, before launching the Encore application.

```bash
temporal server start-dev
```

Install the dependencies using the following command:

```bash
go mod tidy
```

Then, run the application using the following command:

```bash
encore run
```

The APIs will be available at ` http://localhost:4000`.

The APIs can be easily tested via the Encore at ` http://localhost:9400/w68pi`. Or refer to the terminal output for the exact URL.

## Testing

To run the tests for the workflows and the endpoints, use the following command:

```bash
encore test ./...
```

## References

https://encore.dev/docs/how-to/temporal
https://docs.temporal.io/develop/go/testing-suite