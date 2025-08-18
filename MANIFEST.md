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

__meta(my_object, meta_object)

println(my_object)  # >> "my object flux"

__meta(my_object)  # meta_object
```

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

# if else

```flux
if 1 < 2 {
    println(true)
} else {
    println(false)
}
```

# func

```flux
function hello(name) {
    println("Hello, " + name + "!")
    return true
}

let hello = lambda(name) {
    println("Hello, " + name + "!")
}  # default return nil
```
