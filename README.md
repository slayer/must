# Must!


**must** – for when you're shocked that the impossible happened… again.
**must** — a gentle reminder that the universe owes you nothing.
**must** — it's like assert, but it doesn't just quietly cry in dev mode.
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

  // This will panic if the value is not equal to the expected value
  must.Equal(42, 43, "The answer to life, the universe, and everything is wrong!")
}
```

You can also register custom callbacks to handle panics gracefully, allowing you to log the error or take other actions before the program exits.

```go
package main

import (
  "fmt"
  "github.com/slayer/must"
)

must.RegisterFailureHandler(func(message string, details ...string) {
  fmt.Println("World is mad! Error:", message)
  if len(details) > 0 {
    fmt.Println("Details:", details)
  }
})

func main() {
  // This will panic if the condition is false
  must.True(2*2 == 4, "Math is broken!")

  // This will panic if the value is not equal to the expected value
  must.Equal(2*2, 5, "Math is broken!")

  // This will panic if the value is nil
  var x *int
  must.NotNil(x, "x should not be nil")
}

```



## Installation