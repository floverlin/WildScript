# FLUX

## types

|   Flux   |   Go    |   JavaScript   |  Python   |   Lua   |
| :------: | :-----: | :------------: | :-------: | :-----: |
| **nil**  |   nil   | null/undefined |   none    |   nil   |
| **num**  | float64 |     number     |   float   | number  |
| **str**  | string  |     string     |    str    | string  |
| **bool** |  bool   |    boolean     |   bool    | boolean |
| **doc**  |   map   |     object     | list/dict |  table  |

## doc

```flux
let my_object = {
    "a",
    "b",
    "c",
    name = "Lin"
    say = method(self) {
        println("Hello, World!")
        println("My name is " + self.name + "!")
    },
    number = 1176,
    "string key": "string value",
    123: "number value",
}

my_object[1]  # "b"
my_object.say()  # >> "Hello, World!"
my_object{123}  # "number value"
let key = "string key"
my_object{key}  # "string value"
```

## meta

```flux
let meta_object = {
    __str = method(self) {
        return "my object " + self.name
    },
}

my_object = {
    name = "flux",
}

set_meta(my_object, meta_object)

println(my_object)  # >> "my object flux"

get_meta(my_object)  # meta_object
```

- \_\_add
- \_\_sub
- \_\_mul
- \_\_div
- \_\_floor_div
- \_\_mod
- \_\_pow

- \_\_unm
- \_\_not

- \_\_eq
- \_\_lt
- \_\_le

- \_\_str
- \_\_num
- \_\_bool

- \_\_call
- \_\_len

- \_\_slice

- \_\_index
- \_\_safe_index
- \_\_set_index
- \_\_set_list

- \_\_key
- \_\_safe_key
- \_\_set_key
- \_\_set_dict

- \_\_attribute
- \_\_safe_attribute
- \_\_set_attribute

## for while repeat

```flux
for idx, val in my_list {
    println(idx, val)
}

let i = 0
while i < 10 {
    println(i)
    i = i + 1
}

repeat {
    println(i)
    i = i - 1
} until i < 1
```

## if else

```flux
if 1 < 2 {
    println(true)
} else {
    println(false)
}
```

## func

```flux
function hello(name) {
    println("Hello, " + name + "!")
    return true
}

let hello = lambda(name) {
    println("Hello, " + name + "!")
}  # default return nil
```

## slice

```flux
let slice = doc[:]
slice.append("another one lin")
doc[] = slice;
doc[] = doc[2:-2].
    append("and another one lin").
    reverse()
```

## safe

```flux
let res = my_doc?.users?[16]?{"address"}
res  # { value = "rolotushkina", ok = true }
res = my_doc?.users?[16]?{"adres"}
res  # { value = nil, ok = false}
res.or("no address")  # "no address"
```
