package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Hello, LightCache")

	ln, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer ln.Close()

	conn, err := ln.Accept()

	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()

		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expects array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expects array length > 0")
			continue
		}
		fmt.Println(value)

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)

		writer.Write(result)
	}
}
