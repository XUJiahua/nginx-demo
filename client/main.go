package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
)

const NumOfWorkers = 10

var throttle = make(chan int, NumOfWorkers)

var count = 0
var errCount = 0

func main() {
	// go run main.go http://192.168.1.196:10080
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s URL\n", os.Args[0])
		os.Exit(1)
	}
	for {
		count ++
		throttle <- 1
		go func() {
			response, err := http.Get(os.Args[1])
			if err != nil {
				println("get count:", count)
				log.Println(err)
			} else {
				defer response.Body.Close()
				b, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Println("unexpected error")
					println("get count:", count)
					log.Fatal(err)
				}
				str := string(b)
				if str != "10000" && str != "20000" {
					errCount ++
					log.Println(str)
					println("get count:", count)

					if errCount > 10 {
						os.Exit(0)
					}
				}
			}
			<-throttle
		}()
	}
}
