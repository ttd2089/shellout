# Shellout

The shellout package provides a thin wrapper around os/exec to improve testability and ergonomics
when executing simple shell commands programmatically.

The API exports a type `Shell` with a single method `Run` that takes the start options and returns
the result. Using the `Shell` interface facilitates mocking results for testing so you don't have
to execute actual commands.

The package also exports a `Run` function that proxies the call to a default instead of `Shell`.

The `Run` method returns two kinds of errors; `ErrCommandNotFound` and `ErrCommandProcessFailed`.

See the tests for detailed documentation.

## Differences from os/exec

- Shellout only exposes a subset of the options that ox/exec does.

- Non-zero exit codes are not considered errors, they are simply included on the `Result` type.

## Considerations

- Stdout and stderr are read into buffers and returned as part of the result. This means the
  package is not suitable for running processes that require dynamic interaction with the output
  nor for running processes that will produce a volume of output that would be undesirable to hold
  in memory.
