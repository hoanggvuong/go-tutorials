| Category / Type         | Description                                                                                   | Example                        |
|------------------------|-----------------------------------------------------------------------------------------------|--------------------------------|
| Basic Types            | Fundamental types for representing simple values.                                             | `var x int = 5`                |
| Boolean                | `bool` (stores true or false).                                                               | `var b bool = true`            |
| Signed Integers        | `int8`, `int16`, `int32`, `int64`, `int` (platform-dependent, typically 32-bit or 64-bit).   | `var i int32 = -42`            |
| Unsigned Integers      | `uint8`, `uint16`, `uint32`, `uint64`, `uint` (platform-dependent).                          | `var u uint16 = 42`            |
| Aliases                | `byte` (alias for `uint8`), `rune` (alias for `int32`, representing a Unicode code point).   | `var ch rune = 'ą'`            |
| Floating-Point Numbers | `float32`, `float64`.                                                                        | `var f float64 = 3.14`         |
| Complex Numbers        | `complex64`, `complex128`.                                                                   | `var c complex64 = 2 + 3i`     |
| String                 | `string` (immutable sequence of bytes, typically UTF-8 encoded).                             | `var s string = "hello"`       |
| Composite Types        | Types composed of other types.                                                               |                                |
| Arrays                 | Fixed-size sequences of elements of the same type.                                           | `var arr [3]int = [3]int{1,2,3}`|
| Structs                | Collections of named fields, each with a potentially different type.                         | `type Point struct {X, Y int}` |
| Reference Types        | Types that refer to other values or enable advanced features.                                |                                |
| Slices                 | Dynamically-sized, flexible views into arrays.                                               | `var sl []int = []int{1,2,3}`  |
| Maps                   | Unordered collections of key-value pairs.                                                    | `var m map[string]int`         |
| Pointers               | Store the memory address of another variable.                                                | `var p *int = &x`              |
| Channels               | Provide a way for goroutines to communicate safely.                                          | `ch := make(chan int)`         |
| Functions              | First-class citizens, can be assigned to variables and passed as arguments.                  | `func add(a, b int) int { ... }`|
| Interfaces             | Define a set of method signatures that a type must implement.                                | `type Reader interface { Read([]byte) (int, error)`