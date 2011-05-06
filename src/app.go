package main

import (
	"log"
	"flag"
	"./bf2go"
)

var source = flag.String("f","example.bf","brainfuck file to compile");
var dest = flag.String("o","bfsrc.go","output file name");
var debug = flag.Bool("d",false,"compile with debug information");

func main() {
	flag.Parse()
	log.Printf("source: %s", *source)
	log.Printf("dest: %s", *dest)
	log.Printf("debug: %t", *debug)
	bf2go.Translate(*source,*dest,*debug);
}
