# Generic Collections

Package generic_collections provides generic versions of common collection types. You will find each collection in its own sub-package. The code is based on ideas from other go packages and .NET generic collections.

Collections are optionally thread-safe (see [Thread Safety](#thread-safety)), and optionally use concurrency in some methods (see [Concurrency](#concurrency)).

Elements of any type may be stored in these collections. Most basic types are supported out the box, whilst custom types will require the developer to provide functions to support the implementation.

## Supported Types

The following types are supported directly by this library without a requirement to provide custom functions for [hashing](#hashfunc) (required for HashSet) and [comparer](#comparerfunc) (required for all collections).

* Numeric types - all classes of integer and float.
* Pointers - Out of the box, the pointers themselves (i.e. the memory address) are compared by value. If you want to compare/hash on what is being pointed to, then you must provide implementations for these operations.
* `string`
* `bool`
* `time.Time`
* Any type directly castable to one of the above, e.g. `time.Duration` which is in effect `int64`.

Anything else e.g. structs require implementation of these functions, or the collection will panic.

## Collection Hierarchy

- Collection
  - Lists
    - SList - A singly linked list
    - DList - A doubly linked list.
  - Stacks
    - Stack - A slice-backed LIFO stack.
  - Queues
    - Queue - A slice-backed FIFO queue.
    - RingBuffer - A slice-backed circular buffer
  - Sets
    - HashSet - An unordered collection of unique items. Implemented as a hash table.
    - OrderedSet - An ordered collection of unique items. Implemented as a red-black tree.

## Thread Safety

You have the option of whether or not a collection supports thread safety by way of a constructor option. Thread safety is not enabled by default. Enabling thread safety engages a [sync.RWMutex](https://pkg.go.dev/sync#RWMutex) taking the appropriate lock for the type of operation. Some overhead is incurred if thread safety is enabled. . Use the `WithThreadSafe()` constructor option to enable. See [Benchmarks](#benchmarks).

```go
stk := stack.New[int](WithThreadSafe[int]())
```

## Concurrency

In a few places within the sub-packages, concurrency may be enabled to improve performance of some operations. Concurrency is not enabled by default. This is currently limited in scope and may be expanded in future versions. Use the `WithConcurrent()` constructor option to enable. See [benchmarks](#benchamrks) to see where this applies.

Concurrency features only kick in on large collections (currently > 64K elements). For small collections, using concurrency would slow things down due to the overhead involved in managing goroutines.

```go
stk := stack.New[int](WithConcurrent[int]())
```

## Error Handling

Contrary to the more common pattern of returning an error interface as a second argument, I took the decision to panic in case of errors. Common errors include reading from an empty collection, and modifying an underlying collection while an iteration is in progress. If user code is well behaved, then you should be able to avoid these. All collections can be tested for being empty, and many have "Try" versions of methods that return an additional `bool` on some operations that would panic.

## Iteration

All collections are iterable via a common Iterator interface that yields `Element[T]` interface permitting interaction with the values stored in the collections. Collections may be iterated forwards (start to end), reverse (end to start), or forwards with a filter (`TakeWhile()`) It has the following interface:

```go
type Iterator[T any] interface {
	Start() Element[T]
	Next() Element[T]
}
```

An iteration can be performed as follows

```go
ll := slist.New[int]()
// add values, then...
iter := ll.Iterator()

for e := iter.Start() ; e != nil; e = iter.Next() {
    // do something with e.Value() or e.ValuePtr()
}
```

Collections must not be modified during iteration. Modification of the collection will cause iterators to panic on the next call to `Start()` or `Next()`.

### Element

Iteration yields `Element[T]` permitting access to the value stored in the collection at that point. It has the following methods:

```go
type Element[T any] interface {
	Value() T
	ValuePtr() *T
}
```

Note that attempting to modify an item in a collection that implements `Set[T]` via `ValuePtr()` will panic as changing a value breaks the implementation of a set.

## Functions

Some function signatures are provided for you to create your own logic to support various collection operations

### ComparerFunc

This function allows the developer to supply a custom compare function to a collection for a type that is not one of the [supported types](#supported-types).

A compare function should return `< 0` if `a < b`, `0` if `a == b` else `> 0` and looks like this, e.g. for a collection typed on `int`

```go
func myComparer(int a, int b) int {
    return a - b
}
```

### PredicateFunc

For many of the Enumerable methods, a predicate function must be given as an argument. A value is selected when the predicate function returns `true`. For instance, to filter all even numbers from a collection of `int` it might look like this

```go
func filterEvens(int value) bool {
    return value % 2 == 0
}
```

### HashFunc

For collections that use a hash table to store values (currently only `HashSet`), a function to compute a hash value for types not [supported by default](#supported-types) (as per [ComparerFunc](#comparerfunc)) must be provided.

The hash algorithms for the supported types are exported as function variables by the `hashset` sub-package so can be used to construct hashes for struct types.

### DeepCopyFunc

The default action if an instance of this function is not passed to the collection constructor is that when making copies of collection elements, they will be copied by value. If the element type is a pointer, or a struct containing pointers this may not be what you want.

Should you need to deep-copy elements, supply an implementation of this function to the collection's constructor.

The function should return a new instance of the type which is a deep copy of the instance passed as an argument.

## Enumerable

Enumerable defines a set of methods for enumerating a collection in various ways. All collections are enumerable.

```go
type Enumerable[T any] interface {
	Any(predicate PredicateFunc[T]) bool
	All(predicate PredicateFunc[T]) bool
    Find(predicate functions.PredicateFunc[T]) Element[T]
    FindAll(predicate functions.PredicateFunc[T]) []Element[T]
	ForEach(func(Element[T]))
    Min() T
    Max() T
	Map(func(T) T) Collection[T]
	Select(PredicateFunc[T]) Collection[T]
    SelectDeep(functions.PredicateFunc[T]) Collection[T]
}

```

## Benchmarks

In the following tables, the data in the columns have the following meanings

* Elements - Number of elements stored in collection.
* Operation
    * Add/Push/Enqueue - Starting with an empty collection, add the given number of elements one at a time to what is logically the end of the collection. `ns/op` is the time to add all the elements.
    * Remove/Pop/Dequeue - Starting with a full collection, items are all removed one at a time. `ns/op` is the time to remove all the elements.
    * Sort - For collections that support sorting, take a collection collection containing the given number of elements. `ns/op` is the time to sort the collection.
    * Contains - Collection is filled with given number of elements. Source element slice is shuffled and the values it contains are looked up in the collection `ns/op` is the average time to lookup any single element in the collection.
    * Min, Max  - Collection is filled with given number of elements. Min and Max are then taken `b.N` times.
* TS (ThreadSafe)
    * :heavy_check_mark: - Thread safety enabled. Lock is acquired/released on entry/exit of collection method under test.
    * :x: - Thread safety not enabled.
    * Empty - Locking is negligible compared to total run time of the method so not benchmarked.
* PS (Pre-sized)
    * :heavy_check_mark: - Backing store initialized to given number of elements
    * :x: - Backing store not pre-initialized meaning that it has to grow periodically.
    * Empty - Collection does not support a backing store that can be pre-sized.
* CC (Concurrent)
    * :heavy_check_mark: - Concurrent algorithm enabled.
    * :x: - Concurrent algorithm not enabled.
    * Empty - No concurrency option on this algorithm.


Benchmarks are run on collections of `int`

<details>
<summary>Intel(R) Core(TM) i7-7800X</summary>

[CPU Specification](https://www.intel.co.uk/content/www/uk/en/products/sku/123589/intel-core-i77800x-xseries-processor-8-25m-cache-up-to-4-00-ghz/specifications.html)

* Lists

    <details>
    <summary>Expand</summary>

    | Collection | Operation | Elements | TS | PS | CC | ns/op | B/op | allocs/op |
    |------------|-----------|---------:|:--:|:--:|:--:|------:|-----:|----------:|
    | DList | Add | 100 | :x: |  |  | `     5,162    ` | 3,200 | 100 |
    | DList | Add | 1,000 | :x: |  |  | `    40,489    ` | 32,000 | 1,000 |
    | DList | Add | 10,000 | :x: |  |  | `   377,561    ` | 320,000 | 10,000 |
    | DList | Add | 100,000 | :x: |  |  | ` 5,236,563    ` | 3,200,022 | 100,000 |
    | DList | Add | 100 | :heavy_check_mark: |  |  | `     6,947    ` | 3,200 | 100 |
    | DList | Add | 1,000 | :heavy_check_mark: |  |  | `    72,904    ` | 32,000 | 1,000 |
    | DList | Add | 10,000 | :heavy_check_mark: |  |  | `   613,653    ` | 320,000 | 10,000 |
    | DList | Add | 100,000 | :heavy_check_mark: |  |  | ` 7,706,782    ` | 3,200,038 | 100,000 |
    | DList | Remove | 100 | :x: |  |  | `     1,456    ` | 0 | 0 |
    | DList | Remove | 1,000 | :x: |  |  | `    12,370    ` | 0 | 0 |
    | DList | Remove | 10,000 | :x: |  |  | `   127,025    ` | 0 | 0 |
    | DList | Remove | 100,000 | :x: |  |  | ` 1,736,280    ` | 1 | 0 |
    | DList | Remove | 100 | :heavy_check_mark: |  |  | `     4,059    ` | 0 | 0 |
    | DList | Remove | 1,000 | :heavy_check_mark: |  |  | `    37,468    ` | 0 | 0 |
    | DList | Remove | 10,000 | :heavy_check_mark: |  |  | `   363,614    ` | 0 | 0 |
    | DList | Remove | 100,000 | :heavy_check_mark: |  |  | ` 3,959,233    ` | 1 | 0 |
    | DList | Sort | 100 |  |  |  | `     5,612    ` | 0 | 0 |
    | DList | Sort | 1,000 |  |  |  | `   101,728    ` | 0 | 0 |
    | DList | Sort | 10,000 |  |  |  | ` 1,612,152    ` | 0 | 0 |
    | DList | Sort | 100,000 |  |  |  | `30,913,500    ` | 0 | 0 |
    | DList | Min | 100 |  |  |  | `       354.1  ` | 0 | 0 |
    | DList | Min | 1,000 |  |  |  | `     3,336    ` | 0 | 0 |
    | DList | Min | 10,000 |  |  |  | `    32,851    ` | 0 | 0 |
    | DList | Min | 100,000 |  |  |  | `   330,851    ` | 0 | 0 |
    | DList | Max | 100 |  |  |  | `       353.1  ` | 0 | 0 |
    | DList | Max | 1,000 |  |  |  | `     3,395    ` | 0 | 0 |
    | DList | Max | 10,000 |  |  |  | `    33,297    ` | 0 | 0 |
    | DList | Max | 100,000 |  |  |  | `   332,389    ` | 0 | 0 |
    | DList | Contains | 100 |  |  |  | `       160.3  ` | 0 | 0 |
    | DList | Contains | 1,000 |  |  |  | `     1,421    ` | 0 | 0 |
    | DList | Contains | 10,000 |  |  |  | `    14,289    ` | 0 | 0 |
    | DList | Contains | 100,000 |  |  |  | `   142,327    ` | 0 | 0 |
    | SList | Add | 100 | :x: |  |  | `     4,252    ` | 2,400 | 100 |
    | SList | Add | 1,000 | :x: |  |  | `    39,132    ` | 24,000 | 1,000 |
    | SList | Add | 10,000 | :x: |  |  | `   314,291    ` | 240,000 | 10,000 |
    | SList | Add | 100,000 | :x: |  |  | ` 3,662,899    ` | 2,400,011 | 100,000 |
    | SList | Add | 100 | :heavy_check_mark: |  |  | `     8,217    ` | 2,400 | 100 |
    | SList | Add | 1,000 | :heavy_check_mark: |  |  | `    54,659    ` | 24,000 | 1,000 |
    | SList | Add | 10,000 | :heavy_check_mark: |  |  | `   552,099    ` | 240,000 | 10,000 |
    | SList | Add | 100,000 | :heavy_check_mark: |  |  | ` 5,989,855    ` | 2,400,000 | 100,000 |
    | SList | Remove | 100 | :x: |  |  | `     1,185    ` | 0 | 0 |
    | SList | Remove | 1,000 | :x: |  |  | `    12,092    ` | 0 | 0 |
    | SList | Remove | 10,000 | :x: |  |  | `   104,223    ` | 0 | 0 |
    | SList | Remove | 100,000 | :x: |  |  | ` 1,033,820    ` | 0 | 0 |
    | SList | Remove | 100 | :heavy_check_mark: |  |  | `     3,688    ` | 0 | 0 |
    | SList | Remove | 1,000 | :heavy_check_mark: |  |  | `    34,738    ` | 0 | 0 |
    | SList | Remove | 10,000 | :heavy_check_mark: |  |  | `   350,273    ` | 0 | 0 |
    | SList | Remove | 100,000 | :heavy_check_mark: |  |  | ` 3,309,956    ` | 0 | 0 |
    | SList | Sort | 100 |  |  |  | `     5,121    ` | 0 | 0 |
    | SList | Sort | 1,000 |  |  |  | `    93,342    ` | 0 | 0 |
    | SList | Sort | 10,000 |  |  |  | ` 1,461,912    ` | 0 | 0 |
    | SList | Sort | 100,000 |  |  |  | `25,846,042    ` | 0 | 0 |
    | SList | Min | 100 |  |  |  | `       346.0  ` | 0 | 0 |
    | SList | Min | 1,000 |  |  |  | `     3,280    ` | 0 | 0 |
    | SList | Min | 10,000 |  |  |  | `    32,769    ` | 0 | 0 |
    | SList | Min | 100,000 |  |  |  | `   328,320    ` | 0 | 0 |
    | SList | Max | 100 |  |  |  | `       348.2  ` | 0 | 0 |
    | SList | Max | 1,000 |  |  |  | `     3,341    ` | 0 | 0 |
    | SList | Max | 10,000 |  |  |  | `    33,216    ` | 0 | 0 |
    | SList | Max | 100,000 |  |  |  | `   331,937    ` | 0 | 0 |
    | SList | Contains | 100 |  |  |  | `       168.1  ` | 0 | 0 |
    | SList | Contains | 1,000 |  |  |  | `     1,525    ` | 0 | 0 |
    | SList | Contains | 10,000 |  |  |  | `    15,196    ` | 0 | 0 |
    | SList | Contains | 100,000 |  |  |  | `   154,434    ` | 0 | 0 |

    </details>

* Queues

    <details>
    <summary>Expand</summary>

    | Collection | Operation | Elements | TS | PS | CC | ns/op | B/op | allocs/op |
    |------------|-----------|---------:|:--:|:--:|:--:|------:|-----:|----------:|
    | Queue | Enqueue | 100 | :x: | :x: |  | `     2,227    ` | 1,792 | 3 |
    | Queue | Enqueue | 1,000 | :x: | :x: |  | `    15,661    ` | 16,128 | 6 |
    | Queue | Enqueue | 10,000 | :x: | :x: |  | `   168,893    ` | 261,888 | 10 |
    | Queue | Enqueue | 100,000 | :x: | :x: |  | ` 1,630,451    ` | 2,096,901 | 13 |
    | Queue | Enqueue | 100 | :x: | :heavy_check_mark: |  | `       661.9  ` | 0 | 0 |
    | Queue | Enqueue | 1,000 | :x: | :heavy_check_mark: |  | `     8,476    ` | 0 | 0 |
    | Queue | Enqueue | 10,000 | :x: | :heavy_check_mark: |  | `   121,440    ` | 0 | 0 |
    | Queue | Enqueue | 100,000 | :x: | :heavy_check_mark: |  | ` 1,248,483    ` | 2 | 0 |
    | Queue | Enqueue | 100 | :heavy_check_mark: | :x: |  | `     4,003    ` | 1,792 | 3 |
    | Queue | Enqueue | 1,000 | :heavy_check_mark: | :x: |  | `    45,827    ` | 16,128 | 6 |
    | Queue | Enqueue | 10,000 | :heavy_check_mark: | :x: |  | `   367,966    ` | 261,888 | 10 |
    | Queue | Enqueue | 100,000 | :heavy_check_mark: | :x: |  | ` 3,742,788    ` | 2,096,899 | 13 |
    | Queue | Enqueue | 100 | :heavy_check_mark: | :heavy_check_mark: |  | `     2,564    ` | 0 | 0 |
    | Queue | Enqueue | 1,000 | :heavy_check_mark: | :heavy_check_mark: |  | `    36,120    ` | 0 | 0 |
    | Queue | Enqueue | 10,000 | :heavy_check_mark: | :heavy_check_mark: |  | `   321,078    ` | 0 | 0 |
    | Queue | Enqueue | 100,000 | :heavy_check_mark: | :heavy_check_mark: |  | ` 3,303,376    ` | 0 | 0 |
    | Queue | Dequeue | 100 | :x: |  |  | `       638.3  ` | 0 | 0 |
    | Queue | Dequeue | 1,000 | :x: |  |  | `    13,667    ` | 0 | 0 |
    | Queue | Dequeue | 10,000 | :x: |  |  | `   129,131    ` | 0 | 0 |
    | Queue | Dequeue | 100,000 | :x: |  |  | ` 1,258,166    ` | 1 | 0 |
    | Queue | Dequeue | 100 | :heavy_check_mark: |  |  | `     2,435    ` | 0 | 0 |
    | Queue | Dequeue | 1,000 | :heavy_check_mark: |  |  | `    33,463    ` | 0 | 0 |
    | Queue | Dequeue | 10,000 | :heavy_check_mark: |  |  | `   324,595    ` | 0 | 0 |
    | Queue | Dequeue | 100,000 | :heavy_check_mark: |  |  | ` 3,283,561    ` | 0 | 0 |
    | Queue | Sort | 100 |  |  |  | `     4,867    ` | 928 | 2 |
    | Queue | Sort | 1,000 |  |  |  | `    88,174    ` | 8,224 | 2 |
    | Queue | Sort | 10,000 |  |  |  | ` 1,333,856    ` | 81,952 | 2 |
    | Queue | Sort | 100,000 |  |  |  | `16,674,301    ` | 802,849 | 2 |
    | Queue | Min | 100 |  |  | :x: | `       960.2  ` | 0 | 0 |
    | Queue | Min | 1,000 |  |  | :x: | `     3,300    ` | 0 | 0 |
    | Queue | Min | 10,000 |  |  | :x: | `    33,681    ` | 0 | 0 |
    | Queue | Min | 100,000 |  |  | :x: | `   337,106    ` | 0 | 0 |
    | Queue | Min | 100,000 |  |  | :heavy_check_mark: | `   107,563    ` | 1,893 | 35 |
    | Queue | Max | 100 |  |  | :x: | `       954.1  ` | 0 | 0 |
    | Queue | Max | 1,000 |  |  | :x: | `     3,376    ` | 0 | 0 |
    | Queue | Max | 10,000 |  |  | :x: | `    32,882    ` | 0 | 0 |
    | Queue | Max | 100,000 |  |  | :x: | `   329,761    ` | 0 | 0 |
    | Queue | Max | 100,000 |  |  | :heavy_check_mark: | `   107,083    ` | 1,881 | 35 |
    | Queue | Contains | 100 |  |  |  | `       595.5  ` | 0 | 0 |
    | Queue | Contains | 1,000 |  |  |  | `     5,750    ` | 0 | 0 |
    | Queue | Contains | 10,000 |  |  |  | `    58,685    ` | 0 | 0 |
    | Queue | Contains | 100,000 |  |  |  | `   585,437    ` | 0 | 0 |
    | RingBuffer | Enqueue | 100 | :x: |  |  | `       790.1  ` | 0 | 0 |
    | RingBuffer | Enqueue | 1,000 | :x: |  |  | `    18,859    ` | 0 | 0 |
    | RingBuffer | Enqueue | 10,000 | :x: |  |  | `   138,050    ` | 0 | 0 |
    | RingBuffer | Enqueue | 100,000 | :x: |  |  | ` 1,368,396    ` | 0 | 0 |
    | RingBuffer | Enqueue | 100 | :heavy_check_mark: |  |  | `     3,292    ` | 0 | 0 |
    | RingBuffer | Enqueue | 1,000 | :heavy_check_mark: |  |  | `    36,771    ` | 0 | 0 |
    | RingBuffer | Enqueue | 10,000 | :heavy_check_mark: |  |  | `   341,061    ` | 0 | 0 |
    | RingBuffer | Enqueue | 100,000 | :heavy_check_mark: |  |  | ` 3,539,813    ` | 2 | 0 |
    | RingBuffer | Dequeue | 100 | :x: |  |  | `       609.1  ` | 0 | 0 |
    | RingBuffer | Dequeue | 1,000 | :x: |  |  | `     3,954    ` | 0 | 0 |
    | RingBuffer | Dequeue | 10,000 | :x: |  |  | `    36,054    ` | 0 | 0 |
    | RingBuffer | Dequeue | 100,000 | :x: |  |  | `   366,657    ` | 1 | 0 |
    | RingBuffer | Dequeue | 100 | :heavy_check_mark: |  |  | `     2,786    ` | 0 | 0 |
    | RingBuffer | Dequeue | 1,000 | :heavy_check_mark: |  |  | `    28,678    ` | 0 | 0 |
    | RingBuffer | Dequeue | 10,000 | :heavy_check_mark: |  |  | `   290,307    ` | 0 | 0 |
    | RingBuffer | Dequeue | 100,000 | :heavy_check_mark: |  |  | ` 2,791,561    ` | 1 | 0 |
    | RingBuffer | Sort | 100 |  |  |  | `    10,418    ` | 928 | 2 |
    | RingBuffer | Sort | 1,000 |  |  |  | `   113,191    ` | 8,224 | 2 |
    | RingBuffer | Sort | 10,000 |  |  |  | ` 1,444,259    ` | 81,952 | 2 |
    | RingBuffer | Sort | 100,000 |  |  |  | `17,667,739    ` | 802,849 | 2 |
    | RingBuffer | Min | 100 |  |  |  | `       943.2  ` | 0 | 0 |
    | RingBuffer | Min | 1,000 |  |  |  | `     9,339    ` | 0 | 0 |
    | RingBuffer | Min | 10,000 |  |  |  | `    95,596    ` | 0 | 0 |
    | RingBuffer | Min | 100,000 |  |  |  | `   953,840    ` | 0 | 0 |
    | RingBuffer | Max | 100 |  |  |  | `       938.7  ` | 0 | 0 |
    | RingBuffer | Max | 1,000 |  |  |  | `     9,404    ` | 0 | 0 |
    | RingBuffer | Max | 10,000 |  |  |  | `    94,722    ` | 0 | 0 |
    | RingBuffer | Max | 100,000 |  |  |  | `   945,439    ` | 0 | 0 |
    | RingBuffer | Contains | 100 |  |  |  | `       598.5  ` | 0 | 0 |
    | RingBuffer | Contains | 1,000 |  |  |  | `     5,884    ` | 0 | 0 |
    | RingBuffer | Contains | 10,000 |  |  |  | `    58,109    ` | 0 | 0 |
    | RingBuffer | Contains | 100,000 |  |  |  | `   589,616    ` | 0 | 0 |

    </details>

* Sets

    <details>
    <summary>Expand</summary>

    | Collection | Operation | Elements | TS | PS | CC | ns/op | B/op | allocs/op |
    |------------|-----------|---------:|:--:|:--:|:--:|------:|-----:|----------:|
    | HashSet | Add | 100 | :x: | :x: |  | `     8,841    ` | 9,964 | 107 |
    | HashSet | Add | 1,000 | :x: | :x: |  | `   153,902    ` | 180,734 | 1,028 |
    | HashSet | Add | 10,000 | :x: | :x: |  | ` 1,518,866    ` | 1,432,842 | 10,209 |
    | HashSet | Add | 100,000 | :x: | :x: |  | `17,702,102    ` | 12,233,593 | 103,929 |
    | HashSet | Add | 100 | :x: | :heavy_check_mark: |  | `     6,414    ` | 2,113 | 101 |
    | HashSet | Add | 1,000 | :x: | :heavy_check_mark: |  | `    80,157    ` | 16,000 | 1,000 |
    | HashSet | Add | 10,000 | :x: | :heavy_check_mark: |  | `   842,747    ` | 160,004 | 10,000 |
    | HashSet | Add | 100,000 | :x: | :heavy_check_mark: |  | `11,545,505    ` | 2,076,248 | 101,653 |
    | HashSet | Add | 100 | :heavy_check_mark: | :x: |  | `    19,281    ` | 9,965 | 107 |
    | HashSet | Add | 1,000 | :heavy_check_mark: | :x: |  | `   175,249    ` | 180,701 | 1,028 |
    | HashSet | Add | 10,000 | :heavy_check_mark: | :x: |  | ` 1,678,002    ` | 1,432,757 | 10,209 |
    | HashSet | Add | 100,000 | :heavy_check_mark: | :x: |  | `19,997,582    ` | 12,236,830 | 103,940 |
    | HashSet | Add | 100 | :heavy_check_mark: | :heavy_check_mark: |  | `    12,774    ` | 2,113 | 101 |
    | HashSet | Add | 1,000 | :heavy_check_mark: | :heavy_check_mark: |  | `    98,353    ` | 16,000 | 1,000 |
    | HashSet | Add | 10,000 | :heavy_check_mark: | :heavy_check_mark: |  | ` 1,100,168    ` | 160,000 | 10,000 |
    | HashSet | Add | 100,000 | :heavy_check_mark: | :heavy_check_mark: |  | `13,549,729    ` | 2,077,883 | 101,659 |
    | HashSet | Remove | 100 | :x: |  |  | `     6,100    ` | 0 | 0 |
    | HashSet | Remove | 1,000 | :x: |  |  | `    52,680    ` | 0 | 0 |
    | HashSet | Remove | 10,000 | :x: |  |  | `   525,006    ` | 0 | 0 |
    | HashSet | Remove | 100,000 | :x: |  |  | ` 7,590,020    ` | 0 | 0 |
    | HashSet | Remove | 100 | :heavy_check_mark: |  |  | `     7,676    ` | 0 | 0 |
    | HashSet | Remove | 1,000 | :heavy_check_mark: |  |  | `    77,053    ` | 0 | 0 |
    | HashSet | Remove | 10,000 | :heavy_check_mark: |  |  | `   832,967    ` | 0 | 0 |
    | HashSet | Remove | 100,000 | :heavy_check_mark: |  |  | ` 9,916,081    ` | 0 | 0 |
    | HashSet | Min | 100 |  |  |  | `     1,378    ` | 0 | 0 |
    | HashSet | Min | 1,000 |  |  |  | `    14,594    ` | 0 | 0 |
    | HashSet | Min | 10,000 |  |  |  | `   136,113    ` | 0 | 0 |
    | HashSet | Min | 100,000 |  |  |  | ` 1,925,612    ` | 0 | 0 |
    | HashSet | Max | 100 |  |  |  | `     1,392    ` | 0 | 0 |
    | HashSet | Max | 1,000 |  |  |  | `    14,448    ` | 0 | 0 |
    | HashSet | Max | 10,000 |  |  |  | `   135,533    ` | 0 | 0 |
    | HashSet | Max | 100,000 |  |  |  | ` 1,915,966    ` | 0 | 0 |
    | HashSet | Contains | 100 |  |  |  | `        27.39 ` | 0 | 0 |
    | HashSet | Contains | 1,000 |  |  |  | `        37.85 ` | 0 | 0 |
    | HashSet | Contains | 10,000 |  |  |  | `        44.72 ` | 0 | 0 |
    | HashSet | Contains | 100,000 |  |  |  | `        78.73 ` | 0 | 0 |
    | OrderedSet | Add | 100 | :x: |  |  | `     5,782    ` | 4,800 | 100 |
    | OrderedSet | Add | 1,000 | :x: |  |  | `   139,177    ` | 48,000 | 1,000 |
    | OrderedSet | Add | 10,000 | :x: |  |  | ` 1,777,309    ` | 480,001 | 10,000 |
    | OrderedSet | Add | 100,000 | :x: |  |  | `27,412,678    ` | 4,800,011 | 100,000 |
    | OrderedSet | Add | 100 | :heavy_check_mark: |  |  | `    17,031    ` | 4,800 | 100 |
    | OrderedSet | Add | 1,000 | :heavy_check_mark: |  |  | `   132,474    ` | 48,000 | 1,000 |
    | OrderedSet | Add | 10,000 | :heavy_check_mark: |  |  | ` 1,897,430    ` | 480,001 | 10,000 |
    | OrderedSet | Add | 100,000 | :heavy_check_mark: |  |  | `29,175,012    ` | 4,800,014 | 100,000 |
    | OrderedSet | Remove | 100 | :x: |  |  | `     6,572    ` | 0 | 0 |
    | OrderedSet | Remove | 1,000 | :x: |  |  | `    98,252    ` | 0 | 0 |
    | OrderedSet | Remove | 10,000 | :x: |  |  | ` 1,371,899    ` | 0 | 0 |
    | OrderedSet | Remove | 100,000 | :x: |  |  | `22,635,769    ` | 0 | 0 |
    | OrderedSet | Remove | 100 | :heavy_check_mark: |  |  | `     5,266    ` | 0 | 0 |
    | OrderedSet | Remove | 1,000 | :heavy_check_mark: |  |  | `   101,616    ` | 0 | 0 |
    | OrderedSet | Remove | 10,000 | :heavy_check_mark: |  |  | ` 1,364,635    ` | 0 | 0 |
    | OrderedSet | Remove | 100,000 | :heavy_check_mark: |  |  | `23,862,490    ` | 0 | 0 |
    | OrderedSet | Min | 100 |  |  |  | `         6.216` | 0 | 0 |
    | OrderedSet | Min | 1,000 |  |  |  | `         6.660` | 0 | 0 |
    | OrderedSet | Min | 10,000 |  |  |  | `         9.111` | 0 | 0 |
    | OrderedSet | Min | 100,000 |  |  |  | `        11.13 ` | 0 | 0 |
    | OrderedSet | Max | 100 |  |  |  | `         4.560` | 0 | 0 |
    | OrderedSet | Max | 1,000 |  |  |  | `         5.843` | 0 | 0 |
    | OrderedSet | Max | 10,000 |  |  |  | `         6.875` | 0 | 0 |
    | OrderedSet | Max | 100,000 |  |  |  | `         8.876` | 0 | 0 |
    | OrderedSet | Contains | 100 |  |  |  | `        33.14 ` | 0 | 0 |
    | OrderedSet | Contains | 1,000 |  |  |  | `        75.97 ` | 0 | 0 |
    | OrderedSet | Contains | 10,000 |  |  |  | `       111.6  ` | 0 | 0 |
    | OrderedSet | Contains | 100,000 |  |  |  | `       197.8  ` | 0 | 0 |

    </details>

* Stacks

    <details>
    <summary>Expand</summary>

    | Collection | Operation | Elements | TS | PS | CC | ns/op | B/op | allocs/op |
    |------------|-----------|---------:|:--:|:--:|:--:|------:|-----:|----------:|
    | Stack | Push | 100 | :x: | :x: |  | `     1,122    ` | 2,016 | 3 |
    | Stack | Push | 1,000 | :x: | :x: |  | `     6,525    ` | 18,656 | 6 |
    | Stack | Push | 10,000 | :x: | :x: |  | `    50,456    ` | 299,239 | 10 |
    | Stack | Push | 100,000 | :x: | :x: |  | `   751,467    ` | 2,363,622 | 13 |
    | Stack | Push | 100 | :x: | :heavy_check_mark: |  | `       586.5  ` | 0 | 0 |
    | Stack | Push | 1,000 | :x: | :heavy_check_mark: |  | `     3,763    ` | 0 | 0 |
    | Stack | Push | 10,000 | :x: | :heavy_check_mark: |  | `    43,139    ` | 0 | 0 |
    | Stack | Push | 100,000 | :x: | :heavy_check_mark: |  | `   388,153    ` | 3 | 0 |
    | Stack | Push | 100 | :heavy_check_mark: | :x: |  | `     3,068    ` | 2,016 | 3 |
    | Stack | Push | 1,000 | :heavy_check_mark: | :x: |  | `    19,165    ` | 18,656 | 6 |
    | Stack | Push | 10,000 | :heavy_check_mark: | :x: |  | `   323,252    ` | 299,233 | 10 |
    | Stack | Push | 100,000 | :heavy_check_mark: | :x: |  | ` 3,271,576    ` | 2,363,618 | 13 |
    | Stack | Push | 100 | :heavy_check_mark: | :heavy_check_mark: |  | `     3,276    ` | 0 | 0 |
    | Stack | Push | 1,000 | :heavy_check_mark: | :heavy_check_mark: |  | `    34,538    ` | 0 | 0 |
    | Stack | Push | 10,000 | :heavy_check_mark: | :heavy_check_mark: |  | `   274,057    ` | 0 | 0 |
    | Stack | Push | 100,000 | :heavy_check_mark: | :heavy_check_mark: |  | ` 2,857,594    ` | 0 | 0 |
    | Stack | Pop | 100 | :x: |  |  | `       720.5  ` | 0 | 0 |
    | Stack | Pop | 1,000 | :x: |  |  | `     3,349    ` | 0 | 0 |
    | Stack | Pop | 10,000 | :x: |  |  | `    36,320    ` | 1 | 0 |
    | Stack | Pop | 100,000 | :x: |  |  | `   342,374    ` | 2 | 0 |
    | Stack | Pop | 100 | :heavy_check_mark: |  |  | `     3,014    ` | 0 | 0 |
    | Stack | Pop | 1,000 | :heavy_check_mark: |  |  | `    22,547    ` | 0 | 0 |
    | Stack | Pop | 10,000 | :heavy_check_mark: |  |  | `   276,337    ` | 0 | 0 |
    | Stack | Pop | 100,000 | :heavy_check_mark: |  |  | ` 2,669,541    ` | 1 | 0 |
    | Stack | Sort | 100 |  |  |  | `     4,634    ` | 32 | 1 |
    | Stack | Sort | 1,000 |  |  |  | `   107,059    ` | 32 | 1 |
    | Stack | Sort | 10,000 |  |  |  | ` 1,295,618    ` | 32 | 1 |
    | Stack | Sort | 100,000 |  |  |  | `16,371,947    ` | 33 | 1 |
    | Stack | Min | 100 |  |  | :x: | `       371.1  ` | 0 | 0 |
    | Stack | Min | 1,000 |  |  | :x: | `     3,618    ` | 0 | 0 |
    | Stack | Min | 10,000 |  |  | :x: | `    36,231    ` | 0 | 0 |
    | Stack | Min | 100,000 |  |  | :x: | `   357,545    ` | 0 | 0 |
    | Stack | Min | 100,000 |  |  | :heavy_check_mark: | `   108,811    ` | 1,882 | 35 |
    | Stack | Max | 100 |  |  | :x: | `       379.1  ` | 0 | 0 |
    | Stack | Max | 1,000 |  |  | :x: | `     3,639    ` | 0 | 0 |
    | Stack | Max | 10,000 |  |  | :x: | `    35,200    ` | 0 | 0 |
    | Stack | Max | 100,000 |  |  | :x: | `   354,018    ` | 0 | 0 |
    | Stack | Max | 100,000 |  |  | :heavy_check_mark: | `   109,385    ` | 1,882 | 35 |
    | Stack | Contains | 100 |  |  | :x: | `       156.6  ` | 0 | 0 |
    | Stack | Contains | 1,000 |  |  | :x: | `     1,363    ` | 0 | 0 |
    | Stack | Contains | 10,000 |  |  | :x: | `    14,037    ` | 0 | 0 |
    | Stack | Contains | 100,000 |  |  | :x: | `   139,171    ` | 0 | 0 |
    | Stack | Contains | 100,000 |  |  | :heavy_check_mark: | `   100,202    ` | 1,600 | 29 |

    </details>


</details>

<details>
<summary>Intel(R) Core(TM) i7-12700H</summary>

[CPU Specification](https://ark.intel.com/content/www/us/en/ark/products/132228/intel-core-i712700h-processor-24m-cache-up-to-4-70-ghz.html)

* Lists

    <details>
    <summary>Expand</summary>

    | Collection | Operation | Elements | TS | PS | CC | ns/op | B/op | allocs/op |
    |------------|-----------|---------:|:--:|:--:|:--:|------:|-----:|----------:|
    | DList | Add | 100 | :x: |  |  | `     2,988    ` | 3,200 | 100 |
    | DList | Add | 1,000 | :x: |  |  | `    26,507    ` | 32,000 | 1,000 |
    | DList | Add | 10,000 | :x: |  |  | `   243,373    ` | 320,000 | 10,000 |
    | DList | Add | 100,000 | :x: |  |  | ` 3,558,804    ` | 3,200,034 | 100,000 |
    | DList | Add | 100 | :heavy_check_mark: |  |  | `     6,702    ` | 3,200 | 100 |
    | DList | Add | 1,000 | :heavy_check_mark: |  |  | `    51,964    ` | 32,000 | 1,000 |
    | DList | Add | 10,000 | :heavy_check_mark: |  |  | `   540,395    ` | 320,000 | 10,000 |
    | DList | Add | 100,000 | :heavy_check_mark: |  |  | ` 6,676,934    ` | 3,200,027 | 100,000 |
    | DList | Remove | 100 | :x: |  |  | `       866.4  ` | 0 | 0 |
    | DList | Remove | 1,000 | :x: |  |  | `    11,226    ` | 0 | 0 |
    | DList | Remove | 10,000 | :x: |  |  | `   107,712    ` | 0 | 0 |
    | DList | Remove | 100,000 | :x: |  |  | `   988,447    ` | 1 | 0 |
    | DList | Remove | 100 | :heavy_check_mark: |  |  | `     3,527    ` | 0 | 0 |
    | DList | Remove | 1,000 | :heavy_check_mark: |  |  | `    30,529    ` | 0 | 0 |
    | DList | Remove | 10,000 | :heavy_check_mark: |  |  | `   322,074    ` | 1 | 0 |
    | DList | Remove | 100,000 | :heavy_check_mark: |  |  | ` 2,764,714    ` | 0 | 0 |
    | DList | Sort | 100 |  |  |  | `     4,165    ` | 0 | 0 |
    | DList | Sort | 1,000 |  |  |  | `    71,783    ` | 0 | 0 |
    | DList | Sort | 10,000 |  |  |  | ` 1,325,536    ` | 0 | 0 |
    | DList | Sort | 100,000 |  |  |  | `21,481,198    ` | 0 | 0 |
    | DList | Min | 100 |  |  |  | `       142.2  ` | 0 | 0 |
    | DList | Min | 1,000 |  |  |  | `     1,224    ` | 0 | 0 |
    | DList | Min | 10,000 |  |  |  | `    14,411    ` | 0 | 0 |
    | DList | Min | 100,000 |  |  |  | `   128,253    ` | 0 | 0 |
    | DList | Max | 100 |  |  |  | `       141.4  ` | 0 | 0 |
    | DList | Max | 1,000 |  |  |  | `     1,288    ` | 0 | 0 |
    | DList | Max | 10,000 |  |  |  | `    13,221    ` | 0 | 0 |
    | DList | Max | 100,000 |  |  |  | `   130,500    ` | 0 | 0 |
    | DList | Contains | 100 |  |  |  | `        74.73 ` | 0 | 0 |
    | DList | Contains | 1,000 |  |  |  | `       682.6  ` | 0 | 0 |
    | DList | Contains | 10,000 |  |  |  | `     7,257    ` | 0 | 0 |
    | DList | Contains | 100,000 |  |  |  | `    75,901    ` | 0 | 0 |
    | SList | Add | 100 | :x: |  |  | `     3,204    ` | 2,400 | 100 |
    | SList | Add | 1,000 | :x: |  |  | `    29,243    ` | 24,000 | 1,000 |
    | SList | Add | 10,000 | :x: |  |  | `   245,878    ` | 240,000 | 10,000 |
    | SList | Add | 100,000 | :x: |  |  | ` 2,632,169    ` | 2,400,007 | 100,000 |
    | SList | Add | 100 | :heavy_check_mark: |  |  | `     6,467    ` | 2,400 | 100 |
    | SList | Add | 1,000 | :heavy_check_mark: |  |  | `    42,146    ` | 24,000 | 1,000 |
    | SList | Add | 10,000 | :heavy_check_mark: |  |  | `   410,458    ` | 240,000 | 10,000 |
    | SList | Add | 100,000 | :heavy_check_mark: |  |  | ` 4,451,332    ` | 2,400,007 | 100,000 |
    | SList | Remove | 100 | :x: |  |  | `     1,003    ` | 0 | 0 |
    | SList | Remove | 1,000 | :x: |  |  | `     9,262    ` | 0 | 0 |
    | SList | Remove | 10,000 | :x: |  |  | `    82,960    ` | 0 | 0 |
    | SList | Remove | 100,000 | :x: |  |  | `   444,766    ` | 0 | 0 |
    | SList | Remove | 100 | :heavy_check_mark: |  |  | `     3,390    ` | 0 | 0 |
    | SList | Remove | 1,000 | :heavy_check_mark: |  |  | `    29,534    ` | 0 | 0 |
    | SList | Remove | 10,000 | :heavy_check_mark: |  |  | `   258,203    ` | 0 | 0 |
    | SList | Remove | 100,000 | :heavy_check_mark: |  |  | ` 2,518,566    ` | 0 | 0 |
    | SList | Sort | 100 |  |  |  | `     3,844    ` | 0 | 0 |
    | SList | Sort | 1,000 |  |  |  | `    60,946    ` | 0 | 0 |
    | SList | Sort | 10,000 |  |  |  | ` 1,115,609    ` | 0 | 0 |
    | SList | Sort | 100,000 |  |  |  | `19,240,269    ` | 0 | 0 |
    | SList | Min | 100 |  |  |  | `       145.1  ` | 0 | 0 |
    | SList | Min | 1,000 |  |  |  | `     1,359    ` | 0 | 0 |
    | SList | Min | 10,000 |  |  |  | `    13,594    ` | 0 | 0 |
    | SList | Min | 100,000 |  |  |  | `   134,319    ` | 0 | 0 |
    | SList | Max | 100 |  |  |  | `       151.6  ` | 0 | 0 |
    | SList | Max | 1,000 |  |  |  | `     1,394    ` | 0 | 0 |
    | SList | Max | 10,000 |  |  |  | `    13,703    ` | 0 | 0 |
    | SList | Max | 100,000 |  |  |  | `   131,296    ` | 0 | 0 |
    | SList | Contains | 100 |  |  |  | `        75.89 ` | 0 | 0 |
    | SList | Contains | 1,000 |  |  |  | `       676.4  ` | 0 | 0 |
    | SList | Contains | 10,000 |  |  |  | `     6,859    ` | 0 | 0 |
    | SList | Contains | 100,000 |  |  |  | `    69,234    ` | 0 | 0 |

    </details>

* Queues

    <details>
    <summary>Expand</summary>

    | Collection | Operation | Elements | TS | PS | CC | ns/op | B/op | allocs/op |
    |------------|-----------|---------:|:--:|:--:|:--:|------:|-----:|----------:|
    | Queue | Enqueue | 100 | :x: | :x: |  | `     1,447    ` | 1,792 | 3 |
    | Queue | Enqueue | 1,000 | :x: | :x: |  | `     9,382    ` | 16,128 | 6 |
    | Queue | Enqueue | 10,000 | :x: | :x: |  | `    86,595    ` | 261,889 | 10 |
    | Queue | Enqueue | 100,000 | :x: | :x: |  | `   805,183    ` | 2,096,904 | 13 |
    | Queue | Enqueue | 100 | :x: | :heavy_check_mark: |  | `       687.1  ` | 0 | 0 |
    | Queue | Enqueue | 1,000 | :x: | :heavy_check_mark: |  | `     5,328    ` | 0 | 0 |
    | Queue | Enqueue | 10,000 | :x: | :heavy_check_mark: |  | `    56,027    ` | 0 | 0 |
    | Queue | Enqueue | 100,000 | :x: | :heavy_check_mark: |  | `   556,338    ` | 1 | 0 |
    | Queue | Enqueue | 100 | :heavy_check_mark: | :x: |  | `     3,599    ` | 1,792 | 3 |
    | Queue | Enqueue | 1,000 | :heavy_check_mark: | :x: |  | `    31,312    ` | 16,128 | 6 |
    | Queue | Enqueue | 10,000 | :heavy_check_mark: | :x: |  | `   283,523    ` | 261,888 | 10 |
    | Queue | Enqueue | 100,000 | :heavy_check_mark: | :x: |  | ` 2,672,735    ` | 2,096,897 | 13 |
    | Queue | Enqueue | 100 | :heavy_check_mark: | :heavy_check_mark: |  | `     2,985    ` | 0 | 0 |
    | Queue | Enqueue | 1,000 | :heavy_check_mark: | :heavy_check_mark: |  | `    28,162    ` | 0 | 0 |
    | Queue | Enqueue | 10,000 | :heavy_check_mark: | :heavy_check_mark: |  | `   233,762    ` | 0 | 0 |
    | Queue | Enqueue | 100,000 | :heavy_check_mark: | :heavy_check_mark: |  | ` 2,338,408    ` | 0 | 0 |
    | Queue | Dequeue | 100 | :x: |  |  | `       688.3  ` | 0 | 0 |
    | Queue | Dequeue | 1,000 | :x: |  |  | `     5,216    ` | 0 | 0 |
    | Queue | Dequeue | 10,000 | :x: |  |  | `    57,024    ` | 0 | 0 |
    | Queue | Dequeue | 100,000 | :x: |  |  | `   472,703    ` | 0 | 0 |
    | Queue | Dequeue | 100 | :heavy_check_mark: |  |  | `     2,709    ` | 0 | 0 |
    | Queue | Dequeue | 1,000 | :heavy_check_mark: |  |  | `    26,728    ` | 0 | 0 |
    | Queue | Dequeue | 10,000 | :heavy_check_mark: |  |  | `   237,120    ` | 0 | 0 |
    | Queue | Dequeue | 100,000 | :heavy_check_mark: |  |  | ` 2,201,107    ` | 0 | 0 |
    | Queue | Sort | 100 |  |  |  | `     5,050    ` | 928 | 2 |
    | Queue | Sort | 1,000 |  |  |  | `    73,676    ` | 8,224 | 2 |
    | Queue | Sort | 10,000 |  |  |  | `   948,724    ` | 81,952 | 2 |
    | Queue | Sort | 100,000 |  |  |  | `11,564,077    ` | 802,848 | 2 |
    | Queue | Min | 100 |  |  | :x: | `       222.1  ` | 0 | 0 |
    | Queue | Min | 1,000 |  |  | :x: | `     1,160    ` | 0 | 0 |
    | Queue | Min | 10,000 |  |  | :x: | `    12,438    ` | 0 | 0 |
    | Queue | Min | 100,000 |  |  | :x: | `   128,021    ` | 0 | 0 |
    | Queue | Min | 100,000 |  |  | :heavy_check_mark: | `    64,953    ` | 3,036 | 52 |
    | Queue | Max | 100 |  |  | :x: | `       223.1  ` | 0 | 0 |
    | Queue | Max | 1,000 |  |  | :x: | `     1,250    ` | 0 | 0 |
    | Queue | Max | 10,000 |  |  | :x: | `    11,444    ` | 0 | 0 |
    | Queue | Max | 100,000 |  |  | :x: | `   113,885    ` | 0 | 0 |
    | Queue | Max | 100,000 |  |  | :heavy_check_mark: | `    60,476    ` | 3,034 | 52 |
    | Queue | Contains | 100 |  |  |  | `       238.6  ` | 0 | 0 |
    | Queue | Contains | 1,000 |  |  |  | `     2,336    ` | 0 | 0 |
    | Queue | Contains | 10,000 |  |  |  | `    23,404    ` | 0 | 0 |
    | Queue | Contains | 100,000 |  |  |  | `   236,049    ` | 0 | 0 |
    | RingBuffer | Enqueue | 100 | :x: |  |  | `       800.2  ` | 0 | 0 |
    | RingBuffer | Enqueue | 1,000 | :x: |  |  | `     6,197    ` | 0 | 0 |
    | RingBuffer | Enqueue | 10,000 | :x: |  |  | `    57,253    ` | 0 | 0 |
    | RingBuffer | Enqueue | 100,000 | :x: |  |  | `   554,793    ` | 1 | 0 |
    | RingBuffer | Enqueue | 100 | :heavy_check_mark: |  |  | `     3,140    ` | 0 | 0 |
    | RingBuffer | Enqueue | 1,000 | :heavy_check_mark: |  |  | `    25,130    ` | 0 | 0 |
    | RingBuffer | Enqueue | 10,000 | :heavy_check_mark: |  |  | `   239,471    ` | 0 | 0 |
    | RingBuffer | Enqueue | 100,000 | :heavy_check_mark: |  |  | ` 2,340,666    ` | 0 | 0 |
    | RingBuffer | Dequeue | 100 | :x: |  |  | `       472.7  ` | 0 | 0 |
    | RingBuffer | Dequeue | 1,000 | :x: |  |  | `     3,147    ` | 0 | 0 |
    | RingBuffer | Dequeue | 10,000 | :x: |  |  | `    32,156    ` | 0 | 0 |
    | RingBuffer | Dequeue | 100,000 | :x: |  |  | `   276,973    ` | 3 | 0 |
    | RingBuffer | Dequeue | 100 | :heavy_check_mark: |  |  | `     2,860    ` | 0 | 0 |
    | RingBuffer | Dequeue | 1,000 | :heavy_check_mark: |  |  | `    25,374    ` | 0 | 0 |
    | RingBuffer | Dequeue | 10,000 | :heavy_check_mark: |  |  | `   221,808    ` | 0 | 0 |
    | RingBuffer | Dequeue | 100,000 | :heavy_check_mark: |  |  | ` 2,314,409    ` | 1 | 0 |
    | RingBuffer | Sort | 100 |  |  |  | `     5,245    ` | 928 | 2 |
    | RingBuffer | Sort | 1,000 |  |  |  | `    67,790    ` | 8,224 | 2 |
    | RingBuffer | Sort | 10,000 |  |  |  | `   921,211    ` | 81,952 | 2 |
    | RingBuffer | Sort | 100,000 |  |  |  | `12,147,900    ` | 802,848 | 2 |
    | RingBuffer | Min | 100 |  |  |  | `       227.5  ` | 0 | 0 |
    | RingBuffer | Min | 1,000 |  |  |  | `     2,249    ` | 0 | 0 |
    | RingBuffer | Min | 10,000 |  |  |  | `    22,563    ` | 0 | 0 |
    | RingBuffer | Min | 100,000 |  |  |  | `   224,869    ` | 0 | 0 |
    | RingBuffer | Max | 100 |  |  |  | `       229.2  ` | 0 | 0 |
    | RingBuffer | Max | 1,000 |  |  |  | `     2,262    ` | 0 | 0 |
    | RingBuffer | Max | 10,000 |  |  |  | `    22,853    ` | 0 | 0 |
    | RingBuffer | Max | 100,000 |  |  |  | `   226,723    ` | 0 | 0 |
    | RingBuffer | Contains | 100 |  |  |  | `       249.9  ` | 0 | 0 |
    | RingBuffer | Contains | 1,000 |  |  |  | `     2,392    ` | 0 | 0 |
    | RingBuffer | Contains | 10,000 |  |  |  | `    23,815    ` | 0 | 0 |
    | RingBuffer | Contains | 100,000 |  |  |  | `   243,623    ` | 0 | 0 |

    </details>

* Sets

    <details>
    <summary>Expand</summary>

    | Collection | Operation | Elements | TS | PS | CC | ns/op | B/op | allocs/op |
    |------------|-----------|---------:|:--:|:--:|:--:|------:|-----:|----------:|
    | HashSet | Add | 100 | :x: | :x: |  | `    10,294    ` | 9,963 | 107 |
    | HashSet | Add | 1,000 | :x: | :x: |  | `   107,654    ` | 180,694 | 1,028 |
    | HashSet | Add | 10,000 | :x: | :x: |  | ` 1,170,594    ` | 1,432,765 | 10,209 |
    | HashSet | Add | 100,000 | :x: | :x: |  | `11,240,331    ` | 12,234,374 | 103,932 |
    | HashSet | Add | 100 | :x: | :heavy_check_mark: |  | `     9,209    ` | 2,113 | 101 |
    | HashSet | Add | 1,000 | :x: | :heavy_check_mark: |  | `    66,145    ` | 16,000 | 1,000 |
    | HashSet | Add | 10,000 | :x: | :heavy_check_mark: |  | `   688,392    ` | 160,003 | 10,000 |
    | HashSet | Add | 100,000 | :x: | :heavy_check_mark: |  | ` 7,238,287    ` | 2,078,747 | 101,662 |
    | HashSet | Add | 100 | :heavy_check_mark: | :x: |  | `    14,157    ` | 9,965 | 107 |
    | HashSet | Add | 1,000 | :heavy_check_mark: | :x: |  | `   133,932    ` | 180,693 | 1,028 |
    | HashSet | Add | 10,000 | :heavy_check_mark: | :x: |  | ` 1,225,502    ` | 1,433,053 | 10,210 |
    | HashSet | Add | 100,000 | :heavy_check_mark: | :x: |  | `11,446,165    ` | 12,236,271 | 103,938 |
    | HashSet | Add | 100 | :heavy_check_mark: | :heavy_check_mark: |  | `    11,180    ` | 2,114 | 101 |
    | HashSet | Add | 1,000 | :heavy_check_mark: | :heavy_check_mark: |  | `    71,554    ` | 16,000 | 1,000 |
    | HashSet | Add | 10,000 | :heavy_check_mark: | :heavy_check_mark: |  | `   755,077    ` | 160,001 | 10,000 |
    | HashSet | Add | 100,000 | :heavy_check_mark: | :heavy_check_mark: |  | ` 8,826,382    ` | 2,078,576 | 101,661 |
    | HashSet | Remove | 100 | :x: |  |  | `     6,245    ` | 0 | 0 |
    | HashSet | Remove | 1,000 | :x: |  |  | `    47,646    ` | 0 | 0 |
    | HashSet | Remove | 10,000 | :x: |  |  | `   410,425    ` | 0 | 0 |
    | HashSet | Remove | 100,000 | :x: |  |  | ` 4,901,910    ` | 0 | 0 |
    | HashSet | Remove | 100 | :heavy_check_mark: |  |  | `     7,600    ` | 0 | 0 |
    | HashSet | Remove | 1,000 | :heavy_check_mark: |  |  | `    58,179    ` | 0 | 0 |
    | HashSet | Remove | 10,000 | :heavy_check_mark: |  |  | `   468,389    ` | 0 | 0 |
    | HashSet | Remove | 100,000 | :heavy_check_mark: |  |  | ` 6,180,152    ` | 0 | 0 |
    | HashSet | Min | 100 |  |  |  | `       738.1  ` | 0 | 0 |
    | HashSet | Min | 1,000 |  |  |  | `     9,053    ` | 0 | 0 |
    | HashSet | Min | 10,000 |  |  |  | `    85,572    ` | 0 | 0 |
    | HashSet | Min | 100,000 |  |  |  | `   821,231    ` | 0 | 0 |
    | HashSet | Max | 100 |  |  |  | `       741.4  ` | 0 | 0 |
    | HashSet | Max | 1,000 |  |  |  | `     9,181    ` | 0 | 0 |
    | HashSet | Max | 10,000 |  |  |  | `    86,746    ` | 0 | 0 |
    | HashSet | Max | 100,000 |  |  |  | `   830,899    ` | 0 | 0 |
    | HashSet | Contains | 100 |  |  |  | `         9.696` | 0 | 0 |
    | HashSet | Contains | 1,000 |  |  |  | `        15.71 ` | 0 | 0 |
    | HashSet | Contains | 10,000 |  |  |  | `        31.41 ` | 0 | 0 |
    | HashSet | Contains | 100,000 |  |  |  | `        37.66 ` | 0 | 0 |
    | OrderedSet | Add | 100 | :x: |  |  | `     7,095    ` | 4,800 | 100 |
    | OrderedSet | Add | 1,000 | :x: |  |  | `   104,035    ` | 48,000 | 1,000 |
    | OrderedSet | Add | 10,000 | :x: |  |  | ` 1,293,838    ` | 480,004 | 10,000 |
    | OrderedSet | Add | 100,000 | :x: |  |  | `19,638,457    ` | 4,800,034 | 100,000 |
    | OrderedSet | Add | 100 | :heavy_check_mark: |  |  | `    10,053    ` | 4,800 | 100 |
    | OrderedSet | Add | 1,000 | :heavy_check_mark: |  |  | `   108,203    ` | 48,000 | 1,000 |
    | OrderedSet | Add | 10,000 | :heavy_check_mark: |  |  | ` 1,414,259    ` | 480,003 | 10,000 |
    | OrderedSet | Add | 100,000 | :heavy_check_mark: |  |  | `21,027,762    ` | 4,800,031 | 100,000 |
    | OrderedSet | Remove | 100 | :x: |  |  | `     3,186    ` | 0 | 0 |
    | OrderedSet | Remove | 1,000 | :x: |  |  | `    71,386    ` | 0 | 0 |
    | OrderedSet | Remove | 10,000 | :x: |  |  | ` 1,015,825    ` | 4 | 0 |
    | OrderedSet | Remove | 100,000 | :x: |  |  | `16,077,759    ` | 0 | 0 |
    | OrderedSet | Remove | 100 | :heavy_check_mark: |  |  | `     5,102    ` | 0 | 0 |
    | OrderedSet | Remove | 1,000 | :heavy_check_mark: |  |  | `    81,914    ` | 0 | 0 |
    | OrderedSet | Remove | 10,000 | :heavy_check_mark: |  |  | ` 1,126,826    ` | 8 | 0 |
    | OrderedSet | Remove | 100,000 | :heavy_check_mark: |  |  | `17,350,911    ` | 0 | 0 |
    | OrderedSet | Min | 100 |  |  |  | `         2.071` | 0 | 0 |
    | OrderedSet | Min | 1,000 |  |  |  | `         2.428` | 0 | 0 |
    | OrderedSet | Min | 10,000 |  |  |  | `         3.484` | 0 | 0 |
    | OrderedSet | Min | 100,000 |  |  |  | `         4.856` | 0 | 0 |
    | OrderedSet | Max | 100 |  |  |  | `         2.239` | 0 | 0 |
    | OrderedSet | Max | 1,000 |  |  |  | `         3.639` | 0 | 0 |
    | OrderedSet | Max | 10,000 |  |  |  | `         5.302` | 0 | 0 |
    | OrderedSet | Max | 100,000 |  |  |  | `         4.986` | 0 | 0 |
    | OrderedSet | Contains | 100 |  |  |  | `        11.33 ` | 0 | 0 |
    | OrderedSet | Contains | 1,000 |  |  |  | `        46.42 ` | 0 | 0 |
    | OrderedSet | Contains | 10,000 |  |  |  | `        81.52 ` | 0 | 0 |
    | OrderedSet | Contains | 100,000 |  |  |  | `       180.9  ` | 0 | 0 |

    </details>

* Stacks

    <details>
    <summary>Expand</summary>

    | Collection | Operation | Elements | TS | PS | CC | ns/op | B/op | allocs/op |
    |------------|-----------|---------:|:--:|:--:|:--:|------:|-----:|----------:|
    | Stack | Push | 100 | :x: | :x: |  | `     1,281    ` | 2,016 | 3 |
    | Stack | Push | 1,000 | :x: | :x: |  | `     6,305    ` | 18,656 | 6 |
    | Stack | Push | 10,000 | :x: | :x: |  | `    76,407    ` | 299,232 | 10 |
    | Stack | Push | 100,000 | :x: | :x: |  | `   582,318    ` | 2,363,629 | 13 |
    | Stack | Push | 100 | :x: | :heavy_check_mark: |  | `       494.4  ` | 0 | 0 |
    | Stack | Push | 1,000 | :x: | :heavy_check_mark: |  | `     3,266    ` | 0 | 0 |
    | Stack | Push | 10,000 | :x: | :heavy_check_mark: |  | `    29,567    ` | 0 | 0 |
    | Stack | Push | 100,000 | :x: | :heavy_check_mark: |  | `   277,348    ` | 2 | 0 |
    | Stack | Push | 100 | :heavy_check_mark: | :x: |  | `     3,637    ` | 2,016 | 3 |
    | Stack | Push | 1,000 | :heavy_check_mark: | :x: |  | `    26,152    ` | 18,656 | 6 |
    | Stack | Push | 10,000 | :heavy_check_mark: | :x: |  | `   264,761    ` | 299,232 | 10 |
    | Stack | Push | 100,000 | :heavy_check_mark: | :x: |  | ` 2,559,970    ` | 2,363,617 | 13 |
    | Stack | Push | 100 | :heavy_check_mark: | :heavy_check_mark: |  | `     2,672    ` | 0 | 0 |
    | Stack | Push | 1,000 | :heavy_check_mark: | :heavy_check_mark: |  | `    25,598    ` | 0 | 0 |
    | Stack | Push | 10,000 | :heavy_check_mark: | :heavy_check_mark: |  | `   230,666    ` | 0 | 0 |
    | Stack | Push | 100,000 | :heavy_check_mark: | :heavy_check_mark: |  | ` 2,219,721    ` | 0 | 0 |
    | Stack | Pop | 100 | :x: |  |  | `       408.0  ` | 0 | 0 |
    | Stack | Pop | 1,000 | :x: |  |  | `     2,888    ` | 0 | 0 |
    | Stack | Pop | 10,000 | :x: |  |  | `    29,767    ` | 0 | 0 |
    | Stack | Pop | 100,000 | :x: |  |  | `   265,711    ` | 6 | 0 |
    | Stack | Pop | 100 | :heavy_check_mark: |  |  | `     2,639    ` | 0 | 0 |
    | Stack | Pop | 1,000 | :heavy_check_mark: |  |  | `    24,835    ` | 0 | 0 |
    | Stack | Pop | 10,000 | :heavy_check_mark: |  |  | `   224,791    ` | 0 | 0 |
    | Stack | Pop | 100,000 | :heavy_check_mark: |  |  | ` 2,188,625    ` | 3 | 0 |
    | Stack | Sort | 100 |  |  |  | `     4,914    ` | 32 | 1 |
    | Stack | Sort | 1,000 |  |  |  | `    66,146    ` | 32 | 1 |
    | Stack | Sort | 10,000 |  |  |  | `   910,777    ` | 32 | 1 |
    | Stack | Sort | 100,000 |  |  |  | `11,592,799    ` | 32 | 1 |
    | Stack | Min | 100 |  |  | :x: | `       145.6  ` | 0 | 0 |
    | Stack | Min | 1,000 |  |  | :x: | `     1,273    ` | 0 | 0 |
    | Stack | Min | 10,000 |  |  | :x: | `    11,584    ` | 0 | 0 |
    | Stack | Min | 100,000 |  |  | :x: | `   114,872    ` | 0 | 0 |
    | Stack | Min | 100,000 |  |  | :heavy_check_mark: | `    61,672    ` | 3,036 | 52 |
    | Stack | Max | 100 |  |  | :x: | `       149.2  ` | 0 | 0 |
    | Stack | Max | 1,000 |  |  | :x: | `     1,269    ` | 0 | 0 |
    | Stack | Max | 10,000 |  |  | :x: | `    11,410    ` | 0 | 0 |
    | Stack | Max | 100,000 |  |  | :x: | `   114,919    ` | 0 | 0 |
    | Stack | Max | 100,000 |  |  | :heavy_check_mark: | `    59,504    ` | 3,034 | 52 |
    | Stack | Contains | 100 |  |  | :x: | `        57.69 ` | 0 | 0 |
    | Stack | Contains | 1,000 |  |  | :x: | `       468.5  ` | 0 | 0 |
    | Stack | Contains | 10,000 |  |  | :x: | `     4,514    ` | 0 | 0 |
    | Stack | Contains | 100,000 |  |  | :x: | `    45,569    ` | 0 | 0 |
    | Stack | Contains | 100,000 |  |  | :heavy_check_mark: | `    55,627    ` | 2,496 | 45 |

    </details>


</details>