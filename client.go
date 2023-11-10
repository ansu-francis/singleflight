package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	//Generate ransom user id
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomUserID := r1.Intn(100)

	totalRequests := 10

	fmt.Printf("Making %d concurrent requests without coalescing for user id %d\n", totalRequests, randomUserID)
	makeRequests(randomUserID, totalRequests,"v1")
	fmt.Printf("Making %d concurrent requests with coalescing for user id %d\n", totalRequests, randomUserID)
	makeRequests(randomUserID, totalRequests, "v2")
}

//Make request to endpoint 
func makeRequests(randomUserID, totalRequests int, version string) {
	now := time.Now()
	defer func() {
		fmt.Println("Execution time: ", time.Now().Sub(now))
	}()

	endpoint := fmt.Sprintf("http://localhost:8080/api/%s/user?id=%s", version, fmt.Sprint(randomUserID))

	//Make concurrent requests to the server
	
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			makeAPICall(endpoint)
		}(i)
	}
	wg.Wait()
}

func makeAPICall(endpoint string) {
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	result := string(body)
	fmt.Println("Response received: ", result)
}
