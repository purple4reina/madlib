package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var client = new(http.Client)

const (
	wordURL  = "https://reminiscent-steady-albertosaurus.glitch.me/"
	sentence = `It was a %s day. I went downstairs to see if I could %s dinner. I asked, "Does the stew need fresh %s?"`
)

type words struct {
	noun      string
	verb      string
	adjective string
}

func getWord(wordType string, respChan chan<- string, errChan chan<- error) {
	defer close(respChan)
	resp, err := client.Get(wordURL + wordType)
	if err != nil {
		errChan <- err
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errChan <- err
		return
	}
	respChan <- strings.Trim(string(body), `"`)
}

func getWords() (*words, error) {
	w := new(words)
	nounChan := make(chan string)
	verbChan := make(chan string)
	adjChan := make(chan string)
	errChan := make(chan error)
	defer close(errChan)

	go getWord("noun", nounChan, errChan)
	go getWord("verb", verbChan, errChan)
	go getWord("adjective", adjChan, errChan)

	for n := range nounChan {
		w.noun = n
	}
	for v := range verbChan {
		w.verb = v
	}
	for a := range adjChan {
		w.adjective = a
	}

	select {
	case e := <-errChan:
		return nil, e
	default:
		return w, nil
	}
}

func madlibEndpoint(c *gin.Context) {
	words, err := getWords()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"madlib": fmt.Sprintf(sentence, words.adjective, words.verb, words.noun),
	})
}

func main() {
	r := gin.Default()
	r.GET("/madlib", madlibEndpoint)
	go r.Run() // listen and serve on 0.0.0.0:8080
	getMadlib()
}

func getMadlib() {
	resp, err := client.Get("http://localhost:8080/madlib")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(body))
}
