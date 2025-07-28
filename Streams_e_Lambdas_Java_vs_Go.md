# 📖 Streams and Lambdas — Java vs. Go

Summary of Java's `Stream` API operations and their idiomatic equivalents in Go.

| Java `Stream` API                         | What it does                                                           | Possible idiomatic equivalent in Go                                           |
|:-----------------------------------------|:-----------------------------------------------------------------------|:-----------------------------------------------------------------------------|
| `filter(Predicate<T>)`                   | Filters elements that satisfy the predicate                            | `for` + `if` to build a new slice                                            |
| `map(Function<T,R>)`                     | Transforms elements into another type or value                         | `for` + `append(f(item))` or custom `Map([]T, func(T) R) []R` function       |
| `forEach(Consumer<T>)`                   | Executes an operation for each element                                  | `for` with inline function call                                               |
| `collect(Collectors.toList())`           | Collects elements into a list                                          | Slice created via `append`                                                   |
| `findFirst()`                            | Returns the first element, if any                                       | `if len(slice) > 0 { return slice[0] }` or `First([]T) Optional[T]` method    |
| `findAny()`                              | Returns any element (in sequential stream it's the first)              | Same as `findFirst()` in Go, since Go doesn't have native parallel streams    |
| `anyMatch(Predicate<T>)`                 | Returns `true` if any element satisfies the predicate                   | `for` + `if`                                                                 |
| `allMatch(Predicate<T>)`                 | Returns `true` if all elements satisfy the predicate                    | `for` + `if !cond { return false }`                                           |
| `noneMatch(Predicate<T>)`                | Returns `true` if no element satisfies the predicate                    | `for` + `if cond { return false }`                                            |
| `sorted()`                               | Sorts the elements                                                     | `sort.Slice(slice, func(i, j int) bool { ... })`                             |
| `distinct()`                             | Removes duplicate elements                                             | `map` or `set` to filter duplicates                                          |
| `reduce(BinaryOperator<T>)`              | Reduces elements to a single value                                      | `for` + `accumulator variable`                                               |
| `limit(long maxSize)`                    | Limits the number of elements                                          | Slice with cut: `slice[:maxSize]`                                             |
| `skip(long n)`                           | Skips the first `n` elements                                           | Slice with offset: `slice[n:]`                                                |
| `peek(Consumer<T>)`                      | Performs an intermediate operation without modifying the stream         | `for` with side effect (debug, log)                                          |
| `toArray()`                              | Converts the stream to an array                                        | Direct `slice`                                                               |
| `flatMap(Function<T,Stream<R>>) `         | Maps each element to a stream and flattens it                         | `for` + `append` with slice unification operations                           |
| `count()`                                | Counts the number of elements                                          | `len(slice)`                                                                 |
| `min(Comparator<T>)`                     | Returns the smallest element according to the comparator               | `sort.Slice(slice, func(i, j int) bool { return i < j })`                    |
| `max(Comparator<T>)`                     | Returns the largest element according to the comparator                | `sort.Slice(slice, func(i, j int) bool { return i > j })`                    |

## 💡 Considerations
- The Go language doesn't have a native Stream API like Java, but it's possible to simulate many operations efficiently with functions, `for` loops, and libraries like `sort`.
- Operations like `map`, `filter`, `reduce` can be implemented with simplicity, but without the fluency of the Java API.
- Go doesn't have an equivalent concept to Java's parallel `Stream`, but parallelism can be achieved using goroutines and channels if needed.
