# Must!

[![Go Reference](https://pkg.go.dev/badge/github.com/slayer/must.svg)](https://pkg.go.dev/github.com/slayer/must)

**must** – for when you're shocked that the impossible happened… again.

**must** — a gentle reminder that the universe owes you nothing.

**must** — verifying reality, because QA can’t cover everything.

**must** — because “should never happen” happens all the time.

**must** — because sometimes you need to be reminded that the sky is blue and the grass is green.

This library is designed to ensure that the world hasn't gone mad — that things behave as expected and the laws of physics still apply.
It's not meant for handling routine errors like "Connection timeout" or "File not found."
In other words, it serves as a runtime equivalent of `assert`, helping you verify that fundamental assumptions hold true during execution.

# Why Use This Library?

This library is useful for situations where you want to enforce certain conditions in your code and ensure that they are met. It can be particularly helpful in debugging and testing scenarios, where you want to catch unexpected behavior early on.
It can also be used in production code to enforce invariants and ensure that your code behaves as expected. This is especially useful in cases where you have complex logic or dependencies that could lead to unexpected behavior.

For example, if you have a function that should never return `nil`, you can use this library to check that it doesn't.

## Usage

To use the `must` library, simply import it into your Go project and call the functions provided. The library will panic if the conditions you specify are not met.
Here's a simple example:

```go
package main

import (
  "github.com/slayer/must"
)

func main() {
  // This will panic if the condition is false
  must.True(2*2 == 4, "Math is broken!")

  // This will panic if the value is nil
  var x *int
  must.NotNil(x, "x should not be nil")

  // This will panic if the slice, map or string is empty
  must.NotEmpty([]int{1, 2, 3}, "Array should not be empty")
  must.NotEmpty(map[string]int{"a": 1, "b": 2}, "Map should not be empty")
  must.NotEmpty("Hello, world!", "String should not be empty")

  // This will panic if map does not contain the key
  must.SliceHas([]int{1, 2, 3}, 4, "Slice should contain 4")

  // This will panic if map does not contain the key
  must.MapHas(map[string]int{"a": 1, "b": 2}, "c", "Map should contain key 'c'")

  must.FileExists("test.txt", "File should exist")
  must.DirExists("test_dir", "Directory should exist")

}
```

You can also register custom callbacks to handle panics gracefully, allowing you to log the error or take other actions before the program exits.

```go
package main

import (
  "fmt"
  "github.com/slayer/must"
)

must.RegisterFailureHandler(func(message string, details ) {
  fmt.Println("World is mad! Error:", message)
  if len(details) > 0 {
    fmt.Println("Details:", details)
  }
})

func main() {
  // This will panic if the condition is false
  must.True(2*2 == 4, "Math is broken!")
}

```

## Documentation

For more detailed documentation, including all available functions and their usage, please refer to the [GoDoc](https://pkg.go.dev/github.com/slayer/must) page.

## Installation

To install the `must` library, use the following command:

```bash
go get github.com/slayer/must
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! If you have suggestions for improvements or new features, please open an issue or submit a pull request.
