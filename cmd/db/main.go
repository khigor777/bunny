package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/khigor777/bunny"
)

const ttl = 30 //ttl in seconds

type Cache interface {
	Get(key string) []byte
	Set(key string, value []byte)
	Delete(key string)
}

func main() {

	fmt.Println("db starts ...")
	db := bunny.Open()
	setCh := make(chan *bunny.Set, 100)
	getCh := make(chan string, 100)
	deleteCh := make(chan *bunny.Set, 100)

	go readFromStdin(getCh, setCh)
	go deleteKeys(db, deleteCh)

	process(db, setCh, getCh, deleteCh)

}

//handle get and set
func process(db Cache, setCh chan *bunny.Set, getCh chan string, deleteCh chan *bunny.Set) {
	for {
		select {
		case setVal := <-setCh:
			db.Set(setVal.Key, setVal.Value)
			deleteCh <- setVal
			fmt.Println("ok")
		case val := <-getCh:
			res := db.Get(val)
			db.Delete(val)
			if _, err := io.Copy(os.Stdout, bytes.NewReader(res)); err != nil {
				fmt.Print(err.Error())
			}
		default:
		}
	}
}

//ttl delete keys by ttl
func deleteKeys(c Cache, deleteCh chan *bunny.Set) {
	for {
		select {
		case val := <-deleteCh:
			go func(val *bunny.Set) {
				select {
				case <-val.Ctx.Done():
					if c.Get(val.Key) != nil {
						c.Delete(val.Key)
					}
				}
			}(val)
		default:
		}
	}
}

//read keys get and set from stdin
func readFromStdin(getCh chan string, setCh chan *bunny.Set) {
	reader := bufio.NewReader(os.Stdin)
	for {
		b, _ := reader.ReadBytes('\n')
		res := bytes.Split(b, []byte(" "))
		key := string(bytes.Trim(res[1], "\n"))

		switch string(res[0]) {
		case "get":
			go func(key string) {
				getCh <- key
			}(key)
		case "set":
			go func(key string, value []byte) {
				ctx, _ := context.WithTimeout(context.Background(), time.Second*ttl)
				setCh <- &bunny.Set{Key: key, Value: value, Ctx: ctx}
			}(key, res[2])
		}
	}
}
