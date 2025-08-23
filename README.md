# WildScript

## типы

WildScript - язык с динамической типизацией

функция `type` позволяет получить строковое представление типа

```wildscript
type("wild");   # string
type(1176);         # number
type(lambda() {});  # function
```

таблица типов и их аналогов в других языках

|  WildScript  |   Go    | JavaScript |  Python  |   Lua    |
| :----------: | :-----: | :--------: | :------: | :------: |
|   **nil**    |   nil   |    null    |   none   |   nil    |
|  **number**  | float64 |   number   |  float   |  number  |
|  **string**  | string  |   string   |   str    |  string  |
| **boolean**  |  bool   |  boolean   |   bool   | boolean  |
| **document** | struct  |   object   |  class   |  table   |
| **function** |  func   |  function  | function | function |

### nil

тип с единственным значением, обозначающим отсутствие значения

присваивается переменной при обьявлении без инициализации

```wildscript
let a;
let b = nil;
type(a);  # nil
type(b);  # nil
```

### number

представляет вещественные числа (с плавающей точкой двойной точности)

```wildscript
let a = 1;
let b = 10.76;
a = a + b;
a;        # 11.76
type(a);  # number
```

### string

представляет последовательность символов произвольной длинны

строки неизменяемы, поддерживают обьединение и срезы

```wildscript
let a = "hello";
let b = ", world!";
let hello = a + b;
hello;       # hello, world!
hello[1:5];  # ello
```

### boolean

представляет два традиционных логических значения `true` и `false`

только этот тип используется в условиях (без автоприведения)

```wildscript
let a = false;
let b = 1 < 2;
if a then {
    "hello"
} else {
    "world"
};        # world
b;        # true
type(b);  # boolean
```

для условного автоприведения рекомендуется создать отдельную функцию

```wildscript
function nonempty(object) {
    if type(object) == "number" then {
        if object == 0 then {
            return false
        } else {
            return true
        }
    } elif type(object) == "string" then {
        ...
    }
};

if nonempty(my_object) then {
    ...
}
```

### function

всегда возвращают одно значение (по умолчанию - nil)

подразделяются на 4 типа

1. function - обычная функция
2. lambda - анонимная функция (обьявляется внутри выражения)
3. method - автоматически принимает первым аргументом обьект, которому принадлежит
4. native - функция с нативной реализацией

```wildscript
function add(a, b) {
    return a + b
};

let sub = lambda(a, b) {
    return a - b
};

let object = {
    name = "WildScript",
    hello = method(self) {
        return "Hello! My name is " + self.name
    },
};

type(type);  # function
```

### document

тип для композиции и наследования

состоит из списка doc[index], словаря doc{key} и атрибутов doc.attr

```wildscipt
let doc = {
    "a", "b", "c",  # обьявление списка
    name = "Wild"  # обьявление атрибутов
    great = method(self) {
        println("My name is " + self.name + "!")
    },
    number = 1176,
    "string key": "string value",  # обьявление словаря
    123: "number value",
}

doc[1];
doc[2] = "d";
doc[3];  # panic -> index out of range

doc.great();  # My name is Wild!
doc.hello;  # panic -> attribute doesn't exists

doc{123};
let key = "string key";
doc{key};  # или doc{"string key"}
doc{"wrong key"};  # panic -> key doesn't exists
```

#### meta

функция `set_meta` позволяет задать документу метадоку (документ, атрибуты которого переопределяют поведение)

функция `get_meta` позволяет получит метадоку документа

```wildscipt
let md = {
    __call = method(self) {
        return self.name
    },
    __index = method(self, index) {
        return "haha, joke!"
    },
};

set_meta(doc, md);

doc();  # Wild
doc[1176];  # haha, joke!
```

## поток выполнения

### циклы

есть три вида циклов

1. перебор for
2. с предусловием while
3. с постусловием repeat

```wildscipt
for range(10) do {  # range возвращает итератор
    print("!")
};
println();

let list = {1, 2, 3};
for val in list[] do {  # [] возвращает итерируемый список
    print(val, " ");
};
println();

let i = 0
while i < 10 do {
    print(i, " ");
    i = i + 1
};
println();

repeat {
    print(i, " ");
    i = i - 1
} until i < 1;
println();
```

### ветвления

могут быть использованы внутри выражения, возвращая результат последней инсnрукции блока (без ; на конце)

```wildscipt
if 1 < 2 then {
    println(true)
} else {
    println(false)
};

let name = "Wild";
let hello = "Hello!" + if type(name) == "string" then {
    "My name is " + name
    } else { "" };
hello;  # Hello! My name is Wild
```

## panic ? Result

при взятии из документа атрибута/значения списка/значения словаря, которого не существует, произойдет паника

так же панику можно вызвать с помошью оператора panic

после panic передается строка - message, число - code или документ с этими атрибутами

```wildscipt
function unsafe() {
    panic "used unsafe function"
};
unsafe();  # panic -> {message = "used unsafe function"}

panic 404;  # panic -> {code = 404}

panic {
    message = "LOL",
    code = 1337,
    };  # panic -> {message = "LOL", code = 1337}

panic;  # panic -> {} не nil!
```

для безопасного взятия значения из документа или вызова функции используется оператор `?`

при его использовании результатом будет Result = `{value_, error_}`, где
value - возвращаемое значение или nil, если была остановлена паника,
error - ошибка или nil, если паники не было

Result содержит методы для работы с полученным значением

```wildscipt
doc.users?[36]?{"address"}?.or(nil);  # адрес 36-го пользователя или nil

let result = num("not a number")?;
result;  # {
    value_ = nil,
    error_ = {
        message = "could not convert string to number"
        },
    }
```

## классы

классы реализуются через цепочку присвоений атрибутов и метадоков

### конструктор

их можно реализовать вручную используя `set_meta` и `merge` функции
или воспользоваться конструктором

```wildscipt
define Human {
    __init = method(self, name) {
        self.name = name
    }
};

define Witch(Human) {
    __init = method(self, name) {
        super(name)
    },
    cast = method(self) {
        println(self.name + " is casting spell!")
    },
};

let zullie = Witch("Zullie");
zullie.cast();  # Zullie is casting spell!
```

### List Dict slice

для обработки данных документа можно взять эти данные в виде поддокумента вида

1. slice - копия участка исходного списка документа
2. list - полный список документа (ссылается на исходный документ)
3. dict - полный словарь исходного документа (ссылается на исходный документ)

все три подтипа обладают набором методов для работы с данными

чтобы вернуть обработанные данные в документ используется присваивание по

1. срезу `doc[start:end] =` - вставляет список внутрь среза
2. списку `doc[] =` - заменяет список
3. словарю `doc{} =` - заменяет словарь

```wildscipt
let doc = {1, 2, 3};
let slice = doc[:];
slice.append(4).reverse();
doc[] = slice;

doc[].append(0, -1).
    reverse();
```

на самом деле slice и list - это экземпляры одного класса List,
просто slice использует свой собственный список, а list использует `_ref` ссылку

slice.\_ref == nil; # true

то же самое справедливо и для Dict, но его вариант без `_ref` ссылки создается через метод `copy()`

## модули

для импорта модуля используется оператор `import`

после него через точку пишется путь к файлу (путь должен быть в формате синтаксиса переменных)

импортированные данные будут помещены в окружение под последним именем в пути

можно задать псевдоним импорта через `as`

```wildscipt
import mod;
import utils.mod as utils;

let constant = mod.CONSTANT;
utils.func(constant);  # 1176
```

в модулях необходимо определить, что они будут экспортировать через опертор `export`

он **завершает** выполнение модуля подобно `return` функции

после него указывается экспортируемый обьект

```wildscipt
# mod.sil
let constant = 1176;
export {CONSTANT = constant}

#utils/mod.sil
export {
    func = lambda(c) {
        println(c)
    },
}
```
