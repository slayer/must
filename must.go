package must

import (
	"fmt"
	"os"
	"slices"
	"sync"
	"unsafe"
)

// OnFailure is a function type that defines the signature for functions to be called on assertion failures.
type OnFailure func(message string, details string)

var (
	failureHandlers      []OnFailure = []OnFailure{}
	failureHandlersMutex sync.Mutex
)

// RegisterFailureHandler registers a function to be called when an assertion fails.
// This allows for custom handling of assertion failures, such as logging or sending errors to a monitoring service.
// After calling all registered functions, the program will panic with the failure message.
// The function will be called with the failure message and any additional details.
func RegisterFailureHandler(f OnFailure) {
	failureHandlersMutex.Lock()
	defer failureHandlersMutex.Unlock()

	failureHandlers = append(failureHandlers, f)
}

// abort is a helper function that panics with a message and details.
// It is used internally by the assertion functions to handle assertion failures.
func abort(message string, details string) {
	for _, f := range failureHandlers {
		f(message, details)
	}
	panic(message + ": " + fmt.Sprint(details))
}

// NotNil checks if the given value is nil and panics if it is.
// It validates if the interface of the value is not nil and if the underlying value is not nil.
// This handles cases like nil pointers where the interface is not nil but the underlying value is.
// Uses unsafe to check the internal representation of the interface without reflection.
func NotNil(value any, message string) {
	// First, check if the interface itself is nil
	if value == nil {
		abort(message, "expected a non-nil value, got nil")
	}

	// Check if the data pointer inside the interface is nil (e.g., *string(nil))
	// In Go, an interface value is represented by a header with two words:
	// a pointer to type information and a pointer to the actual data
	type eface struct {
		_type unsafe.Pointer
		data  unsafe.Pointer
	}

	// Convert the interface to our internal representation to check the data pointer
	valuePtr := (*eface)(unsafe.Pointer(&value)) // #nosec G103

	// If the data pointer is nil, it means the interface contains a nil pointer value
	if valuePtr.data == nil && valuePtr._type != nil {
		abort(message, fmt.Sprintf("expected a non-nil value, got nil pointer of type %T", value))
	}
}

// NoError checks if the given error is nil and panics if it is not.
func NoError(err error, message string) {
	if err != nil {
		abort(message, fmt.Sprintf("expected no error, got: %v", err))
	}
}

// Error checks if the given error is not nil and panics if it is.
func Error(err error, message string) {
	if err == nil {
		abort(message, "expected an error, got nil")
	}
}

// NotEqual checks if the given value is not equal to the expected value and panics if it is.
// It is used to ensure that two values are not equal before proceeding with further operations.
func NotEqual[T comparable](expected, value T, message string) {
	if expected == value {
		abort(message, fmt.Sprintf("expected %v to not be equal to %v", expected, value))
	}
}

// Equal checks if the given value is equal to the expected value and panics if it is not.
// It is used to ensure that two values are equal before proceeding with further operations.
func Equal[T comparable](expected, value T, message string) {
	if expected != value {
		abort(message, fmt.Sprintf("expected %v to be equal to %v", expected, value))
	}
}

// True checks if the given value is true and panics if it is not.
// It is used to ensure that a boolean condition is true before proceeding with further operations.
func True(value bool, message string) {
	if !value {
		abort(message, "expected true, got false")
	}
}

// False checks if the given value is false and panics if it is.
// It is used to ensure that a boolean condition is false before proceeding with further operations.
func False(value bool, message string) {
	if value {
		abort(message, "expected false, got true")
	}
}

// NotZero checks if the given value is zero and panics if it is.
// It is used to ensure that a numeric value is not zero before proceeding with further operations.
func NotZero[T ~int | float64](value T, message string) {
	if value == 0 {
		abort(message, "expected non-zero value, got zero")
	}
}

func GreaterThan[T ~int | float64](value, threshold T, message string) {
	if value <= threshold {
		abort(message, fmt.Sprintf("expected %v to be greater than %v", value, threshold))
	}
}
func LessThan[T ~int | float64](value, threshold T, message string) {
	if value >= threshold {
		abort(message, fmt.Sprintf("expected %v to be less than %v", value, threshold))
	}
}
func GreaterThanOrEqual[T ~int | float64](value, threshold T, message string) {
	if value < threshold {
		abort(message, fmt.Sprintf("expected %v to be greater than or equal to %v", value, threshold))
	}
}
func LessThanOrEqual[T ~int | float64](value, threshold T, message string) {
	if value > threshold {
		abort(message, fmt.Sprintf("expected %v to be less than or equal to %v", value, threshold))
	}
}

// NotEmpty checks if the given value (map, slice or string) is empty and panics if it is.
func NotEmpty(value any, message string) {
	switch v := value.(type) {
	case map[any]any:
		if len(v) == 0 {
			abort(message, "expected a non-empty map, got empty")
		}
	case []any:
		if len(v) == 0 {
			abort(message, "expected a non-empty slice, got empty")
		}
	case string:
		if v == "" {
			abort(message, "expected a non-empty string, got empty")
		}
	default:
		abort(message, fmt.Sprintf("expected a map, slice or string, got %T", v))
	}
}

// Empty checks if the given value (map, slice or string) is not empty and panics if it is not.
func Empty(value any, message string) {
	switch v := value.(type) {
	case map[any]any:
		if len(v) != 0 {
			abort(message, "expected an empty map, got non-empty")
		}
	case []any:
		if len(v) != 0 {
			abort(message, "expected an empty slice, got non-empty")
		}
	case string:
		if v != "" {
			abort(message, "expected an empty string, got non-empty")
		}
	default:
		abort(message, fmt.Sprintf("expected a map, slice or string, got %T", v))
	}
}

// Contains checks if the given slice contains the specified value and panics if it does not.
func Contains[T comparable](slice []T, value T, message string) {
	if slices.Contains(slice, value) {
		return
	}
	abort(message, fmt.Sprintf("expected slice to contain %v, but it does not", value))
}

// NotContains checks if the given slice does not contain the specified value and panics if it does.
func NotContains[T comparable](slice []T, value T, message string) {
	if !slices.Contains(slice, value) {
		return
	}
	abort(message, fmt.Sprintf("expected slice to not contain %v, but it does", value))
}

// IsNil checks if the given value is nil and panics if it is not.
func IsNil(value any, message string) {
	if value != nil {
		abort(message, "expected nil, got non-nil")
	}
}

// IsNotNil checks if the given value is not nil and panics if it is.
func IsNotNil(value any, message string) {
	if value == nil {
		abort(message, "expected non-nil, got nil")
	}
}

// FileExists checks if the given file path exists and panics if it does not.
// It is used to ensure that a file exists before proceeding with further operations.
func FileExists(path string, message string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		abort(message, fmt.Sprintf("expected file %s to exist, but it does not", path))
	}
}

// DirExists checks if the given directory path exists and panics if it does not.
// It is used to ensure that a directory exists before proceeding with further operations.
func DirExists(path string, message string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		abort(message, fmt.Sprintf("expected directory %s to exist, but it does not", path))
	}
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		abort(message, fmt.Sprintf("expected %s to be a directory, but it is not", path))
	}
}

// TypeOf checks if the given value is of the expected type and panics if it is not.
// It is used to ensure that a value is of a specific type before proceeding with further operations.
func TypeOf[T any](value any, message string) {
	if _, ok := value.(T); !ok {
		abort(message, fmt.Sprintf("expected value of type %T, got %T", (*T)(nil), value))
	}
}

// TypeOfNot checks if the given value is not of the expected type and panics if it is.
// It is used to ensure that a value is not of a specific type before proceeding with further operations.
func TypeOfNot[T any](value any, message string) {
	if _, ok := value.(T); ok {
		abort(message, fmt.Sprintf("expected value not of type %T, got %T", (*T)(nil), value))
	}
}

// PointsToSame checks if two pointers point to the same value and panics if they do not.
func PointsToSame[T comparable](a, b *T, message string) {
	if a == nil || b == nil { // nolint:staticcheck
		abort(message, "expected non-nil pointers, got nil")
	}
	if *a != *b { // nolint:staticcheck
		abort(message, fmt.Sprintf("expected pointers to point to the same value, got %v and %v", *a, *b))
	}
}

func PointsToNotSame[T comparable](a, b *T, message string) {
	if a == nil || b == nil { // nolint:staticcheck
		abort(message, "expected non-nil pointers, got nil")
	}
	if *a == *b { // nolint:staticcheck
		abort(message, fmt.Sprintf("expected pointers to point to different values, got %v and %v", *a, *b))
	}
}

func SliceHas[T comparable](slice []T, value T, message string) {
	if !slices.Contains(slice, value) {
		abort(message, fmt.Sprintf("expected slice to have %v, but it does not", value))
	}
}
func SliceNotHas[T comparable](slice []T, value T, message string) {
	if slices.Contains(slice, value) {
		abort(message, fmt.Sprintf("expected slice to not have %v, but it does", value))
	}
}

func MapHas[K comparable, V any](m map[K]V, key K, message string) {
	if _, ok := m[key]; !ok {
		abort(message, fmt.Sprintf("expected map to have key %v, but it does not", key))
	}
}
func MapNotHas[K comparable, V any](m map[K]V, key K, message string) {
	if _, ok := m[key]; ok {
		abort(message, fmt.Sprintf("expected map to not have key %v, but it does", key))
	}
}

func MapNotEmpty[K comparable, V any](m map[K]V, message string) {
	if len(m) == 0 {
		abort(message, "expected map to be non-empty, but it is empty")
	}
}

func MapEmpty[K comparable, V any](m map[K]V, message string) {
	if len(m) != 0 {
		abort(message, "expected map to be empty, but it is not")
	}
}

func IsEmpty[T comparable](slice []T, message string) {
	if len(slice) != 0 {
		abort(message, "expected slice to be empty, but it is not")
	}
}
