package must

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRegisterFailureHandler tests the RegisterFailureHandler function
func TestRegisterFailureHandler(t *testing.T) {
	t.Parallel()

	// Save the original failure handlers and restore them after the test
	originalHandlers := failureHandlers
	defer func() { failureHandlers = originalHandlers }()

	// Reset failure handlers for this test
	failureHandlers = []OnFailure{}

	// Create a test failure handler that records whether it was called
	var handlerCalled bool
	testHandler := func(message, details string) {
		handlerCalled = true
		assert.Equal(t, "test message", message)
		assert.Contains(t, details, "test details")
	}

	// Register the test handler
	RegisterFailureHandler(testHandler)

	// Verify the handler was registered
	require.Len(t, failureHandlers, 1)

	// Test that the handler is called when abort is called
	defer func() {
		r := recover()
		assert.NotNil(t, r, "Expected abort to panic")
		assert.True(t, handlerCalled, "Expected failure handler to be called")
	}()

	// Call abort which should call the registered handler and panic
	abort("test message", "test details")
}

// TestNotNil tests the NotNil function
func TestNotNil(t *testing.T) {
	t.Parallel()

	// Test success case - should not panic
	NotNil("not nil", "should not panic")
	NotNil(123, "should not panic")
	NotNil([]string{"foo"}, "should not panic")

	// Test failure case - should panic
	t.Run("nil value", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotNil(t, r, "Expected NotNil to panic on nil value")
		}()

		var nilValue any
		NotNil(nilValue, "should panic")
	})

	t.Run("nil pointer to string", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotNil(t, r, "Expected NotNil to panic on nil value")
		}()

		var nilString *string
		NotNil(nilString, "should panic")
	})

	// Test if interface is not nil but value is nil
	t.Run("interface with nil value", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotNil(t, r, "Expected NotNil to panic on nil interface")
		}()

		var nilError error // interface with nil value
		NotNil(nilError, "should panic")
	})

	t.Run("interface not nil", func(t *testing.T) {
		NotNil("not nil", "should not panic")
	})
}

// TestNotNilWithUnsafe specifically tests the improved NotNil function
// that uses unsafe to detect nil pointers inside non-nil interfaces
func TestNotNilWithUnsafe(t *testing.T) {
	t.Parallel()

	// Test with regular non-nil values (should not panic)
	normalString := "test"
	NotNil(normalString, "regular string should not panic")

	normalInt := 42
	NotNil(normalInt, "regular int should not panic")

	// Test with pointers to non-nil values (should not panic)
	strPtr := new(string)
	*strPtr = "test pointer"
	NotNil(strPtr, "pointer to string should not panic")

	intPtr := new(int)
	*intPtr = 42
	NotNil(intPtr, "pointer to int should not panic")

	// Test with nil interface (should panic)
	t.Run("nil interface", func(t *testing.T) {
		t.Parallel()
		assert.Panics(t, func() {
			var nilInterface any
			NotNil(nilInterface, "should panic with nil interface")
		})
	})

	// Test with nil pointer to string (should panic)
	t.Run("nil string pointer", func(t *testing.T) {
		t.Parallel()
		assert.Panics(t, func() {
			var nilStrPtr *string
			NotNil(nilStrPtr, "should panic with nil string pointer")
		})
	})

	// Test with nil pointer to int (should panic)
	t.Run("nil int pointer", func(t *testing.T) {
		t.Parallel()
		assert.Panics(t, func() {
			var nilIntPtr *int
			NotNil(nilIntPtr, "should panic with nil int pointer")
		})
	})

	// Test with nil error interface (should panic)
	t.Run("nil error interface", func(t *testing.T) {
		t.Parallel()
		assert.Panics(t, func() {
			var nilErr error
			NotNil(nilErr, "should panic with nil error")
		})
	})

	// Test with typed nil (should panic)
	t.Run("typed nil", func(t *testing.T) {
		t.Parallel()
		type customStruct struct{}
		assert.Panics(t, func() {
			var nilCustomPtr *customStruct
			NotNil(nilCustomPtr, "should panic with typed nil")
		})
	})
}

// TestNoError tests the NoError function
func TestNoError(t *testing.T) {
	t.Parallel()

	// Test success case - should not panic
	NoError(nil, "should not panic")

	// Test failure case - should panic
	t.Run("error value", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotNil(t, r, "Expected NoError to panic on error value")
		}()

		err := errors.New("test error")
		NoError(err, "should panic")
	})
}

// TestError tests the Error function
func TestError(t *testing.T) {
	t.Parallel()

	// Test success case - should not panic
	err := errors.New("test error")
	Error(err, "should not panic")

	// Test failure case - should panic
	t.Run("nil error", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotNil(t, r, "Expected Error to panic on nil error")
		}()

		Error(nil, "should panic")
	})
}

// TestEqualNotEqual tests the Equal and NotEqual functions
func TestEqualNotEqual(t *testing.T) {
	t.Parallel()

	// Test Equal success case - should not panic
	Equal(42, 42, "should not panic")
	Equal("foo", "foo", "should not panic")

	// Test Equal failure case - should panic
	t.Run("not equal values", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotNil(t, r, "Expected Equal to panic when values are not equal")
		}()

		Equal(42, 43, "should panic")
	})

	// Test NotEqual success case - should not panic
	NotEqual(42, 43, "should not panic")
	NotEqual("foo", "bar", "should not panic")

	// Test NotEqual failure case - should panic
	t.Run("equal values", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotNil(t, r, "Expected NotEqual to panic when values are equal")
		}()

		NotEqual(42, 42, "should panic")
	})
}

// TestBooleanAssertions tests the True and False functions
func TestBooleanAssertions(t *testing.T) {
	t.Parallel()

	// Test True success case - should not panic
	True(true, "should not panic")

	// Test True failure case - should panic
	t.Run("false value for True", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotNil(t, r, "Expected True to panic when value is false")
		}()

		True(false, "should panic")
	})

	// Test False success case - should not panic
	False(false, "should not panic")

	// Test False failure case - should panic
	t.Run("true value for False", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotNil(t, r, "Expected False to panic when value is true")
		}()

		False(true, "should panic")
	})
}

// TestNumericAssertions tests the numeric comparison functions
func TestNumericAssertions(t *testing.T) {
	t.Parallel()

	// Test NotZero success and failure cases
	t.Run("NotZero", func(t *testing.T) {
		t.Parallel()

		// Success case
		NotZero(42, "should not panic")

		// Failure case
		t.Run("zero value", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected NotZero to panic on zero value")
			}()

			NotZero(0, "should panic")
		})
	})

	// Test GreaterThan success and failure cases
	t.Run("GreaterThan", func(t *testing.T) {
		t.Parallel()

		// Success case
		GreaterThan(42, 41, "should not panic")

		// Failure cases
		t.Run("equal values", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected GreaterThan to panic when values are equal")
			}()

			GreaterThan(42, 42, "should panic")
		})

		t.Run("less than", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected GreaterThan to panic when value is less")
			}()

			GreaterThan(41, 42, "should panic")
		})
	})

	// Test LessThan success and failure cases
	t.Run("LessThan", func(t *testing.T) {
		t.Parallel()

		// Success case
		LessThan(41, 42, "should not panic")

		// Failure cases
		t.Run("equal values", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected LessThan to panic when values are equal")
			}()

			LessThan(42, 42, "should panic")
		})

		t.Run("greater than", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected LessThan to panic when value is greater")
			}()

			LessThan(42, 41, "should panic")
		})
	})

	// Test GreaterThanOrEqual success and failure cases
	t.Run("GreaterThanOrEqual", func(t *testing.T) {
		t.Parallel()

		// Success cases
		GreaterThanOrEqual(42, 42, "should not panic")
		GreaterThanOrEqual(43, 42, "should not panic")

		// Failure case
		t.Run("less than", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected GreaterThanOrEqual to panic when value is less")
			}()

			GreaterThanOrEqual(41, 42, "should panic")
		})
	})

	// Test LessThanOrEqual success and failure cases
	t.Run("LessThanOrEqual", func(t *testing.T) {
		t.Parallel()

		// Success cases
		LessThanOrEqual(42, 42, "should not panic")
		LessThanOrEqual(41, 42, "should not panic")

		// Failure case
		t.Run("greater than", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected LessThanOrEqual to panic when value is greater")
			}()

			LessThanOrEqual(43, 42, "should panic")
		})
	})
}

// TestEmptyAssertions tests the Empty and NotEmpty functions
func TestEmptyAssertions(t *testing.T) {
	t.Parallel()

	// Test NotEmpty success and failure cases for different types
	t.Run("NotEmpty", func(t *testing.T) {
		t.Parallel()

		// Success cases
		NotEmpty(map[any]any{"key": "value"}, "should not panic")
		NotEmpty([]any{42}, "should not panic")
		NotEmpty("foo", "should not panic")

		// Failure cases
		t.Run("empty map", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected NotEmpty to panic on empty map")
			}()

			NotEmpty(map[any]any{}, "should panic")
		})

		t.Run("empty slice", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected NotEmpty to panic on empty slice")
			}()

			NotEmpty([]any{}, "should panic")
		})

		t.Run("empty string", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected NotEmpty to panic on empty string")
			}()

			NotEmpty("", "should panic")
		})

		t.Run("unsupported type", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected NotEmpty to panic on unsupported type")
			}()

			NotEmpty(42, "should panic")
		})
	})

	// Test Empty success and failure cases for different types
	t.Run("Empty", func(t *testing.T) {
		t.Parallel()

		// Success cases
		Empty(map[any]any{}, "should not panic")
		Empty([]any{}, "should not panic")
		Empty("", "should not panic")

		// Failure cases
		t.Run("non-empty map", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected Empty to panic on non-empty map")
			}()

			Empty(map[any]any{"key": "value"}, "should panic")
		})

		t.Run("non-empty slice", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected Empty to panic on non-empty slice")
			}()

			Empty([]any{42}, "should panic")
		})

		t.Run("non-empty string", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected Empty to panic on non-empty string")
			}()

			Empty("foo", "should panic")
		})

		t.Run("unsupported type", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected Empty to panic on unsupported type")
			}()

			Empty(42, "should panic")
		})
	})
}

// TestContainsAssertions tests the Contains and NotContains functions
func TestContainsAssertions(t *testing.T) {
	t.Parallel()

	// Test Contains success and failure cases
	t.Run("Contains", func(t *testing.T) {
		t.Parallel()

		// Success case
		Contains([]string{"foo", "bar", "baz"}, "bar", "should not panic")

		// Failure case
		t.Run("missing value", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected Contains to panic when value is missing")
			}()

			Contains([]string{"foo", "bar", "baz"}, "qux", "should panic")
		})
	})

	// Test NotContains success and failure cases
	t.Run("NotContains", func(t *testing.T) {
		t.Parallel()

		// Success case
		NotContains([]string{"foo", "bar", "baz"}, "qux", "should not panic")

		// Failure case
		t.Run("contained value", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected NotContains to panic when value is contained")
			}()

			NotContains([]string{"foo", "bar", "baz"}, "bar", "should panic")
		})
	})
}

// TestNilAssertions tests the IsNil and IsNotNil functions
func TestNilAssertions(t *testing.T) {
	t.Parallel()

	// Test IsNil success and failure cases
	t.Run("IsNil", func(t *testing.T) {
		t.Parallel()

		// Success case
		var nilValue any
		IsNil(nilValue, "should not panic")

		// Failure case
		t.Run("non-nil value", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected IsNil to panic on non-nil value")
			}()

			IsNil("not nil", "should panic")
		})
	})

	// Test IsNotNil success and failure cases
	t.Run("IsNotNil", func(t *testing.T) {
		t.Parallel()

		// Success case
		IsNotNil("not nil", "should not panic")

		// Failure case
		t.Run("nil value", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected IsNotNil to panic on nil value")
			}()

			var nilValue any
			IsNotNil(nilValue, "should panic")
		})
	})
}

// TestFileAssertions tests the FileExists and DirExists functions
func TestFileAssertions(t *testing.T) {
	// Don't use t.Parallel() since we're dealing with file system operations

	// Create a temporary file and directory for testing
	tempFile, err := os.CreateTemp("", "must-test-file-*.txt")
	require.NoError(t, err)

	// Write some data to the file to ensure it exists properly
	_, err = tempFile.Write([]byte("test content"))
	require.NoError(t, err)

	// Sync to ensure content is written to disk
	require.NoError(t, tempFile.Sync())

	// Get the file path before we close it
	tempFilePath := tempFile.Name()

	// Close the file but don't delete it yet
	require.NoError(t, tempFile.Close())

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "must-test-dir-*")
	require.NoError(t, err)

	// Set up cleanup to happen after all subtests complete
	defer func() {
		os.Remove(tempFilePath)
		os.RemoveAll(tempDir)
	}()

	nonExistentPath := "/path/that/does/not/exist"

	// Test FileExists success and failure cases
	t.Run("FileExists", func(t *testing.T) {
		// Don't use t.Parallel() here as it references shared resources

		// Success case
		FileExists(tempFilePath, "should not panic")

		// Failure case
		t.Run("non-existent file", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected FileExists to panic on non-existent file")
			}()

			FileExists(nonExistentPath, "should panic")
		})
	})

	// Test DirExists success and failure cases
	t.Run("DirExists", func(t *testing.T) {
		// Don't use t.Parallel() here as it references shared resources

		// Success case
		DirExists(tempDir, "should not panic")

		// Failure cases
		t.Run("non-existent directory", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected DirExists to panic on non-existent directory")
			}()

			DirExists(nonExistentPath, "should panic")
		})

		t.Run("file instead of directory", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected DirExists to panic when path is a file")
			}()

			DirExists(tempFilePath, "should panic")
		})
	})
}

// TestTypeAssertions tests the TypeOf and TypeOfNot functions
func TestTypeAssertions(t *testing.T) {
	t.Parallel()

	// Test TypeOf success and failure cases
	t.Run("TypeOf", func(t *testing.T) {
		t.Parallel()

		// Success case
		TypeOf[string]("string value", "should not panic")

		// Failure case
		t.Run("wrong type", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected TypeOf to panic on wrong type")
			}()

			TypeOf[string](42, "should panic")
		})
	})

	// Test TypeOfNot success and failure cases
	t.Run("TypeOfNot", func(t *testing.T) {
		t.Parallel()

		// Success case
		TypeOfNot[int]("string value", "should not panic")

		// Failure case
		t.Run("matching type", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected TypeOfNot to panic on matching type")
			}()

			TypeOfNot[int](42, "should panic")
		})
	})
}

// TestPointerAssertions tests the PointsToSame and PointsToNotSame functions
func TestPointerAssertions(t *testing.T) {
	t.Parallel()

	// Test PointsToSame success and failure cases
	t.Run("PointsToSame", func(t *testing.T) {
		t.Parallel()

		// Success case - pointers to same value
		value1 := "same value"
		ptr1 := &value1
		ptr2 := &value1
		PointsToSame(ptr1, ptr2, "should not panic")

		// Failure cases
		t.Run("different values", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected PointsToSame to panic on different values")
			}()

			value2 := "different value"
			ptr2 := &value2
			PointsToSame(ptr1, ptr2, "should panic")
		})

		t.Run("nil pointers", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected PointsToSame to panic on nil pointers")
			}()

			var nilPtr *string
			PointsToSame(ptr1, nilPtr, "should panic")
		})
	})

	// Test PointsToNotSame success and failure cases
	t.Run("PointsToNotSame", func(t *testing.T) {
		t.Parallel()

		// Success case - pointers to different values
		value1 := "value1"
		value2 := "value2"
		ptr1 := &value1
		ptr2 := &value2
		PointsToNotSame(ptr1, ptr2, "should not panic")

		// Failure cases
		t.Run("same values", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected PointsToNotSame to panic on same values")
			}()

			value := "same"
			ptr1 := &value
			ptr2 := &value
			PointsToNotSame(ptr1, ptr2, "should panic")
		})

		t.Run("nil pointers", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected PointsToNotSame to panic on nil pointers")
			}()

			var nilPtr *string
			PointsToNotSame(ptr1, nilPtr, "should panic")
		})
	})
}

// TestSliceAssertions tests the SliceHas and SliceNotHas functions
func TestSliceAssertions(t *testing.T) {
	t.Parallel()

	// Test SliceHas success and failure cases
	t.Run("SliceHas", func(t *testing.T) {
		t.Parallel()

		// Success case
		SliceHas([]string{"foo", "bar", "baz"}, "bar", "should not panic")

		// Failure case
		t.Run("missing value", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected SliceHas to panic when value is missing")
			}()

			SliceHas([]string{"foo", "bar", "baz"}, "qux", "should panic")
		})
	})

	// Test SliceNotHas success and failure cases
	t.Run("SliceNotHas", func(t *testing.T) {
		t.Parallel()

		// Success case
		SliceNotHas([]string{"foo", "bar", "baz"}, "qux", "should not panic")

		// Failure case
		t.Run("contained value", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected SliceNotHas to panic when value is contained")
			}()

			SliceNotHas([]string{"foo", "bar", "baz"}, "bar", "should panic")
		})
	})
}

// TestMapAssertions tests the map-related assertion functions
func TestMapAssertions(t *testing.T) {
	t.Parallel()

	// Test MapHas success and failure cases
	t.Run("MapHas", func(t *testing.T) {
		t.Parallel()

		// Success case
		testMap := map[string]int{"foo": 1, "bar": 2, "baz": 3}
		MapHas(testMap, "bar", "should not panic")

		// Failure case
		t.Run("missing key", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected MapHas to panic when key is missing")
			}()

			MapHas(testMap, "qux", "should panic")
		})
	})

	// Test MapNotHas success and failure cases
	t.Run("MapNotHas", func(t *testing.T) {
		t.Parallel()

		// Success case
		testMap := map[string]int{"foo": 1, "bar": 2, "baz": 3}
		MapNotHas(testMap, "qux", "should not panic")

		// Failure case
		t.Run("existing key", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected MapNotHas to panic when key exists")
			}()

			MapNotHas(testMap, "bar", "should panic")
		})
	})

	// Test MapEmpty and MapNotEmpty success and failure cases
	t.Run("MapEmpty", func(t *testing.T) {
		t.Parallel()

		// Success case
		MapEmpty(map[string]int{}, "should not panic")

		// Failure case
		t.Run("non-empty map", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected MapEmpty to panic on non-empty map")
			}()

			MapEmpty(map[string]int{"foo": 1}, "should panic")
		})
	})

	t.Run("MapNotEmpty", func(t *testing.T) {
		t.Parallel()

		// Success case
		MapNotEmpty(map[string]int{"foo": 1}, "should not panic")

		// Failure case
		t.Run("empty map", func(t *testing.T) {
			defer func() {
				r := recover()
				assert.NotNil(t, r, "Expected MapNotEmpty to panic on empty map")
			}()

			MapNotEmpty(map[string]int{}, "should panic")
		})
	})
}

// TestIsEmpty tests the IsEmpty function
func TestIsEmpty(t *testing.T) {
	t.Parallel()

	// Success case
	IsEmpty([]int{}, "should not panic")

	// Failure case
	t.Run("non-empty slice", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotNil(t, r, "Expected IsEmpty to panic on non-empty slice")
		}()

		IsEmpty([]int{1, 2, 3}, "should panic")
	})
}
