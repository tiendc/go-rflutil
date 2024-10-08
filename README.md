[![Go Version][gover-img]][gover] [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![GoReport][rpt-img]][rpt]

# Reflection utility functions for Go

## Installation

```shell
go get github.com/tiendc/go-rflutil
```

## Usage

- [Slice functions](#slice-functions)
- [Map functions](#map-functions)
- [Struct functions](#struct-functions)
- [Common functions](#common-functions)

### Slice functions

#### SliceLen

```go
slice := []int{1, 2, 3}
v, err := SliceLen(reflect.ValueOf(slice)) // v == 3
```

#### SliceGet

```go
slice := []int{1, 2, 3}
v, err := SliceGet[int](reflect.ValueOf(slice), 1)     // v == 2
v, err := SliceGet[int](reflect.ValueOf(slice), 3)     // err is ErrIndexOutOfRange
v, err := SliceGet[string](reflect.ValueOf(slice), 1)  // err is ErrTypeUnmatched
```

#### SliceSet

```go
slice := []int{1, 2, 3}
err := SliceSet(reflect.ValueOf(slice), 1, 22)    // slice[1] == 22
err := SliceSet(reflect.ValueOf(slice), 3, 44)    // err is ErrIndexOutOfRange
err := SliceSet(reflect.ValueOf(slice), 1, "22")  // err is ErrTypeUnmatched
```

#### SliceAppend

```go
slice := []int{1, 2, 3}
slice2, err := SliceAppend(reflect.ValueOf(slice), 4)   // slice2 == []int{1, 2, 3, 4}
slice2, err := SliceAppend(reflect.ValueOf(slice), "4") // err is ErrTypeUnmatched
```

#### SliceGetAll

```go
slice := []int{1, 2, 3}
s, err := SliceGetAll(reflect.ValueOf(slice)) // returns []reflect.Value
```

#### SliceAs

```go
slice := []any{1, 2, 3}
s, err := SliceAs[int64](reflect.ValueOf(slice)) // s == []int64{1,2,3}
```

### Map functions

#### MapLen

```go
aMap := map[int]string{1: "11", 2: "22", 3: "33"}
v, err := MapLen(reflect.ValueOf(aMap), 1) // v == 3
```

#### MapGet

```go
aMap := map[int]string{1: "11", 2: "22", 3: "33"}
v, err := MapGet[string](reflect.ValueOf(aMap), 1) // v == "11"
v, err := MapGet[int](reflect.ValueOf(aMap), 3)    // err is ErrTypeUnmatched
```

#### MapSet

```go
aMap := map[int]string{1: "11", 2: "22", 3: "33"}
err := MapSet(reflect.ValueOf(aMap), 1, "111") // success
err := MapSet(reflect.ValueOf(aMap), 4, "444") // success
err := MapSet(reflect.ValueOf(aMap), 5, 555)   // err is ErrTypeUnmatched
```

#### MapDelete

```go
aMap := map[int]string{1: "11", 2: "22", 3: "33"}
err := MapDelete(reflect.ValueOf(aMap), 1) // success
err := MapDelete(reflect.ValueOf(aMap), 4) // success
```

#### MapKeys

```go
aMap := map[int]string{1: "11", 2: "22", 3: "33"}
keys, err := MapKeys(reflect.ValueOf(aMap)) // returns a slice []reflect.Value of keys
```

#### MapEntries

```go
aMap := map[int]string{1: "11", 2: "22", 3: "33"}
entries, err := MapEntries(reflect.ValueOf(aMap)) // returns a slice []MapEntry of entries
```

### Struct functions

#### StructGetField

```go
type Struct struct {
    I int
    S string
}
s := Struct{I: 1, S: "11"}
v, err := StructGetField[int](reflect.ValueOf(s), "I", true)     // v == 1
v, err := StructGetField[string](reflect.ValueOf(s), "S", true)  // v == "11"
v, err := StructGetField[string](reflect.ValueOf(s), "s", false) // v == "11"
v, err := StructGetField[string](reflect.ValueOf(s), "s", true)  // err is ErrNotFound
```

#### StructSetField

```go
type Struct struct {
    I int
    S string
}
s := Struct{I: 1, S: "11"}
err := StructSetField[int](reflect.ValueOf(&s), "I", 11, true)        // success
err := StructSetField[string](reflect.ValueOf(&s), "S", "111", true)  // success
err := StructSetField[string](reflect.ValueOf(&s), "s", "111", false) // success
err := StructSetField[string](reflect.ValueOf(&s), "s", "111", true)  // err is ErrNotFound
```

#### StructListFields

```go
type Base struct {
    I  int
    S2 string
}
type Struct struct {
    Base
    I int
    S string
}
s := Struct{}
fields, err := StructListFields(reflect.ValueOf(&s), false) // returns []string{"Base", "I", "S"}
fields, err := StructListFields(reflect.ValueOf(&s), true)  // returns []string{"S2", "I", "S"}
```

#### StructToMap

```go
type Base struct {
    I  int     `json:"i"`
    S2 string  `json:"s2"`
}
type Struct struct {
    Base
    I int    `json:"i"`
    S string `json:"s"`
}

s := Struct{I: 1, S: "S", Base: Base{I: 2, S2: "S2"}}

// Converts without flattening the embedded struct
m, err := StructToMap(reflect.ValueOf(&s), "", false)    // m == map[string]any{"Base": Base{I: 2, S2: "S2"}, "I": 1, "S": "S"}

// Converts with flattening the embedded struct
m, err := StructToMap(reflect.ValueOf(&s), "", true)     // m == map[string]any{"S2": "S2", "I": 1, "S": "S"}

// Converts with parsing json tag
m, err := StructToMap(reflect.ValueOf(&s), "json", true) // m == map[string]any{"s2": "S2", "i": 1, "s": "S"}
```

#### ParseTag / ParseTagOf / ParseTagsOf

```go
type S struct {
    I int    `mytag:"i,optional,k=v"`
    S string `mytag:"s,optional,k1=v1,k2=v2,omitempty"`
    U uint   `mytag:"-,optional"`
}

s := S{I: 1, S: "11", U: 10}
sVal := reflect.ValueOf(s)
iField, _ := sVal.Type().FieldByName("I")
tag, err := ParseTag(&iField, "mytag", ",") // tag.Name == i
                                            // tag.Attrs == map[string]string{"optional": "", "k", "v"}
```

### Common functions

#### ValueAs

```go
v, err := ValueAs[float32](reflect.ValueOf(97)) // v == float32(97)
v, err := ValueAs[string](reflect.ValueOf(97))  // v == "a"
```

## Contributing

- You are welcome to make pull requests for new functions and bug fixes.

## License

- [MIT License](LICENSE)

[doc-img]: https://pkg.go.dev/badge/github.com/tiendc/go-rflutil
[doc]: https://pkg.go.dev/github.com/tiendc/go-rflutil
[gover-img]: https://img.shields.io/badge/Go-%3E%3D%201.18-blue
[gover]: https://img.shields.io/badge/Go-%3E%3D%201.18-blue
[ci-img]: https://github.com/tiendc/go-rflutil/actions/workflows/go.yml/badge.svg
[ci]: https://github.com/tiendc/go-rflutil/actions/workflows/go.yml
[cov-img]: https://codecov.io/gh/tiendc/go-rflutil/branch/main/graph/badge.svg
[cov]: https://codecov.io/gh/tiendc/go-rflutil
[rpt-img]: https://goreportcard.com/badge/github.com/tiendc/go-rflutil
[rpt]: https://goreportcard.com/report/github.com/tiendc/go-rflutil
