package pkg

import "fmt"

func Multiply(str string, num int) string {
	var out string
	for range num {
		out += str
	}
	return out
}

func Cover(str, title string, cov string, num int) string {
	var out string
	out += fmt.Sprint(Multiply(cov, num) + title + Multiply(cov, num) + "\n")
	out += fmt.Sprint(str + "\n")
	out += fmt.Sprint(Multiply(cov, num*2+len(title)) + "\n")
	return out
}
