print("hello, world!\n");

3 {
    idx = @idx;
    print(idx, "\n");
};

print("GO, GO WILD!\n");

lin = new {
    name: "Lin",
    nick: "floverlin",
    age: 21,
    scream: fn() { print("AAA!\n"); },
    say: fn() {
        self = @self;
        print("Hello! My name is " + self.name + "!\n");
    }
};

lin.name = lin.name + "ovoe";
lin.say();