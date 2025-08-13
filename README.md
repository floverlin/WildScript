# WildScript - минималистичный язык программирования

- высокого уровня
- динамической сильной типизации

## Типы

| WildScript |   Go    |   JavaScript   | Python |
| :--------: | :-----: | :------------: | :----: |
|  **num**   | float64 |     number     | float  |
|  **str**   | string  |     string     |  str   |
|  **nil**   |   nil   | null/undefined |  none  |
|  **bool**  |  bool   |    boolean     |  bool  |
|  **func**  |  func   |    function    |  def   |
|  **obj**   |   map   |     object     |  dict  |
|  **list**  |  slice  |     array      |  list  |

## Переменные

```wildscript
age = 21;
name = "lin";
job = nil; jobs;
married = false;
scream = () { print("AAA!"); };
inventory = {
    book: "Holy Bible",
};
marks = [name, inventory, nil];

a = true;
a = "another";
a  # "another"

# b; -> panic: undefined variable b
```

## Блоки

- возвращают значение последней инструкции
- инструкции разделяются ;
- последняя инструкция не содержит после себя ;
- результат пустой инструкции = nil
- основная программа - тоже блок кода, но без {}
- обращение к внешним переменным через &

```
{ 1 + 2 };  # 3
{ 1 + 2; };  # nil

a = 1;
b = 1;
{
    # a = a + 1; -> panic: undefined variable a
    a = &a + 1;
    &b = a;
};
a;  # 1
b;  # 2
```

## Функции

- всегда возвращают одно значение
- состоит из списка параметров и блока

```wildscript
f = (a, b) { a * b };
f(2, 2);  # 4

(a, b) { a * b }(2, 2)  # 4
```

## Объекты

```wildscript
obj = {
    a: 1,
    b: "b",
    c: () { print("hello, world!"); },
};

obj.a;  # 1
obj.b;  # "b"
obj.c();  # выведет "hello, world!"

obj.d = "d";  # or obj.set("d", "d")

# obj.d; -> panic: undefined obj field d
obj.get("d");  # nil

print(obj);  # выведет {a: 1, b: "b", c: func, d: "d"}
```

## Списки

```wildscript
list = [10, "10", [nil, true]];

l[2][1];  # true

list.[1] = 20;

# list[4]; -> panic: index out of range

print(list);  # выведет [10, 20, [nil, true]]
```

## Преобразование в логический тип

- !true = false
- false = true

- !!true = true

- !!0 = false
- !!(num != 0) = true

- !!(other object) = !!len(object)

## Ветвления

- состоит из условия, операторов и блоков
- условие только bool типа

```wildscript
a > b
    ? { print("a > b"); }
    : { print("a <= b"); };

a > b ? {
    print("a > b");
} : a > c ? {
    print("a <= b && a > c");
} : {
    print("a <= b && a <= c");
};

5 + (true ? { 5 })  # 10
```

## Циклы

- состоит из условия и блока
- условие любого типа

|   type    | iterations |
| :-------: | :--------: |
| **bool**  | until true |
| **other** | len(other) |

```wildscript
true {
    print("eternity");
};

i = 1;
i < 10 {
    print(i);
    i = i + 1;
};

i = 0
5 + (5 { &i = &i + 1; &i })  # 10
5 + ("wild" { &i = &i + 1; &i })  # 14
```

## Математика

- только с num типом

- \+
- \-
- \*
- /
- //
- %
- ^

```wildscript
2.2 + 4.4;  # 6.6
4.4 - 2;  # 2.4

2.3 * 2;  # 4.6
10.5 / 5;  # 2.1
10.5 // 5;  # 2
10.5 % 5;  # 0.5

2^4;  # 16

# 1 / 0 -> panic: division by zero
# 1 // 0 -> panic: division by zero
# 1 % 0 -> panic: modulo by zero
```

## Логика

- &&
- ||

- ==
- !=

- <
- \>

- <=
- \>=


## Стандартные функции

- работают с обьектами всех типов

### print

Выводит текстовые представления обьектов в консоль

- return nil

```wildscript
print(1, "2", true);  выведет 1 "2" true
```

### len

Возвращает числовое представление обьекта

- return num

|   type   |   value   |       operation       | result |
| :------: | :-------: | :-------------------: | :----: |
| **num**  |   12.34   |      math.floor       |   12   |
| **str**  | "flower"  |        length         |   6    |
| **nil**  |    nil    |           0           |   0    |
| **bool** |   true    | true -> 1; false -> 0 |   1    |
| **func** | (a, b) {} |    length of args     |   2    |
| **obj**  |  {a: 1}   |    number of keys     |   1    |
| **list** | [1, 2, 3] |  number of elements   |   3    |

```wildscript
len(3.14);  # 3
```

### type

Возвращает текстовое представление типа обьекта

- return str

|   type   |   value   | result |
| :------: | :-------: | :----: |
| **num**  |   12.34   | "num"  |
| **str**  | "flower"  | "str"  |
| **nil**  |    nil    | "nil"  |
| **bool** |   true    | "bool" |
| **func** | (a, b) {} | "func" |
| **obj**  |  {a: 1}   | "obj"  |
| **list** | [1, 2, 3] | "list" |

```wildscript
type(3.14);  # "num"
```
