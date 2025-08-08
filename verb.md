| Verb | Description                                   | Example Usage                | Output Example                |
|------|-----------------------------------------------|------------------------------|-------------------------------|
| %v   | Default format for value                      | `fmt.Printf("%v", 123)`      | `123`                         |
| %w   | Error wrapping (used with fmt.Errorf)         | `fmt.Errorf("fail: %w", err)`| Wraps error for unwrapping    |
| %s   | String or []byte                              | `fmt.Printf("%s", "hi")`     | `hi`                          |
| %d   | Decimal integer                               | `fmt.Printf("%d", 42)`       | `42`                          |
| %f   | Floating-point number                         | `fmt.Printf("%f", 3.14)`     | `3.140000`                    |
| %t   | Boolean                                       | `fmt.Printf("%t", true)`     | `true`                        |
| %q   | Quoted string                                 | `fmt.Printf("%q", "hi")`     | `"hi"`                        |
| %x   | Hexadecimal (int, byte slice, string)         | `fmt.Printf("%x", 255)`      | `ff`                          |
| %T   | Type of value                                 | `fmt.Printf("%T", 123)`      | `int`                         |
| %#v  | Go-syntax representation                      | `fmt.Printf("%#v", 123)`     | `123`                         |
| %p   | Pointer address                               | `fmt.Printf("%p", &x)`       | `0xc000012080`                |
| %b   | Binary representation (int, byte slice)       | `fmt.Printf("%b", 5)`        | `101`                         |
| %o   | Octal representation (int, byte slice)        | `fmt.Printf("%o", 8)`        | `10`                          |
| %c   | Character (rune)                              | `fmt.Printf("%c", 'A')`      | `A`                           |
| %U   | Unicode code point (rune)                     | `fmt.Printf("%U", 'A')`     | `U+0041 'A'`                  |
| %e   | Scientific notation (float)                   | `fmt.Printf("%e", 123456.789)` | `1.234568e+05`               |
| %g   | Compact representation (float)                | `fmt.Printf("%g", 123456.789)` | `123456.789`                  |
| %v   | Default format for value (with type)          | `fmt.Printf("%v", myStruct)` | `{{Field1:Value1 Field2:Value2}}` |
| %s   | String representation of a struct field       | `fmt.Printf("%s", myStruct.Field1)` | `Value1`                     |
| %d   | Decimal representation of a struct field      | `fmt.Printf("%d", myStruct.Field2)` | `42`                          |
| %p   | Pointer address of a struct field             | `fmt.Printf("%p",  &myStruct.Field1)` | `0xc000012080`                |
| %v   | Default format for a slice or array           | `fmt.Printf("%v", mySlice)` | `[1 2 3]`                     |
| %s   | String representation of a slice or array     | `fmt.Printf("%s", mySlice)` | `"[1 2 3]"`                   |
| %d   | Decimal representation of a slice or array    | `fmt.Printf("%d    ", len(mySlice))` | `3`                           |
| %p   | Pointer address of a slice or array           | `fmt.Printf("%p    ", mySlice)` | `0xc000012080`                |
| %v   | Default format for a map                       | `fmt.Printf("%v", myMap)` | `map[key1:value1 key2:value2]` |
| %s   | String representation of a map                | `fmt.Printf("%s", myMap)` | `"{key1:value1 key2:value2}"` |
| %d   | Decimal representation of a map (key count)   | `fmt.Printf("%d", len(myMap))` | `2`                           |
| %p   | Pointer address of a map                      | `fmt.Printf("%p", myMap)` | `0xc000012080`                |
| %v   | Default format for a channel                   | `fmt.Printf("%v", myChannel)` | `<-chan int`                  |
| %s   | String representation of a channel            | `fmt.Printf("%s", myChannel)` | `"<-chan int"`                |
| %d   | Decimal representation of a channel (buffer size) | `fmt.Printf("%d", cap(myChannel))` | `0`                           |
| %p   | Pointer address of a channel                  | `fmt.Printf("%p", myChannel)` | `0xc000012080`                |
| %v   | Default format for a function                  | `fmt.Printf("%v", myFunction)` | `func() { ... }`              |
| %s   | String representation of a function           | `fmt.Printf("%s", myFunction)` | `"func() { ... }"`            |
| %d   | Decimal representation of a function (arity)  | `fmt.Printf("%d", runtime.NumIn(myFunction))` | `0`                           |
| %p   | Pointer address of a function                 | `fmt.Printf("%p", myFunction)` | `0xc000012080`                |
| %v   | Default format for an interface                | `fmt.Printf("%v", myInterface)` | `interface{}`                |
| %s   | String representation of an interface         | `fmt.Printf("%s", myInterface)` | `"interface{}"`               |
| %d   | Decimal representation of an interface (method count) | `fmt.Printf("%d", reflect.TypeOf(myInterface).NumMethod())` | `0`                           |
| %p   | Pointer address of an interface               | `fmt.Printf("%p", myInterface)` | `0xc000012080`               |