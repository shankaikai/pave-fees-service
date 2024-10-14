# Description

A fees API service that allows users to create a bill, add fees to the bill, and close a bill. This application is built for local development only.

Each bill is basically a Temporal workflow that is created when a new bill is created. The bill is then updated or closed via Temporal signals.

Functions available:

1. Create a bill
2. Add a fee to a bill
3. Close a bill
4. List all bills
5. Get a bill by ID

## Running

Ensure that the Temporal dev server is running locally, before launching the Encore application. This project assumes default ports.

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

The APIs can be easily tested via the Encore at `http://localhost:9400/w68pi`. Refer to the terminal output for the exact URL, if it happens to be different.

## Testing

To run the tests for the workflows and the endpoints, use the following command:

```bash
encore test ./...
```

## Assumptions

- Currency is set when the bill is created and cannot be changed.
  - Fees are assumed to be in the same currency as the bill.
- Bills have no limits on the number of fees that can be added.
- Fees can only be positive values.
- Bills can only have two states: open and closed.
- Bills cannot be reopened once closed.

## Future Improvements

- Authentication and authorization, if the service is to be exposed to the public and consumed by a FE.
- Temporal documentation mentioned that frequent use of the list workflows API might affect persistence performance. Could possibly introduce a cache layer to reduce the number of calls to the Temporal server or store the bills in our own database using activities.

## References

https://encore.dev/docs/how-to/temporal
https://docs.temporal.io/develop/go/testing-suite
