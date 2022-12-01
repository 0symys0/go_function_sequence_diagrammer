package main

import (
//  "bytes"
//  "fmt"
  "log"
  "io/ioutil"

  "github.com/goccy/go-graphviz"
)

func main() {
	path := "./graphtest.gv"
	b, err := ioutil.ReadFile(path)
	if err != nil {
	  log.Fatal(err)
	}
	graph, err := graphviz.ParseBytes(b)
	if err != nil {
	  log.Fatal(err)
	}

	// create your graph
    g := graphviz.New()
	// 1. write encoded PNG data to buffer
	//var buf bytes.Buffer
	//if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
	//  log.Fatal(err)
	//}

	// 2. get as image.Image instance
	//image, err := g.RenderImage(graph)
	//if err != nil {
	//  log.Fatal(err)
	//}

	// 3. write to file directly
	if err := g.RenderFilename(graph, graphviz.PNG, "test.png"); err != nil {
	  log.Fatal(err)
	}
}
