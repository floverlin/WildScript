print("hello, world!");

3 {
    idx = @idx;
    print(idx);
};

print("GO, GO WILD!");

lin = {
    name: "Lin",
    nick: "floverlin",
    age: 21,
    scream: fn() { print("AAA!"); },
    say: fn() {
        self = @self;
        print("Hello! My name is " + self.name + "!");
    }
};

lin.name = lin.name + "ovoe";
lin.say();