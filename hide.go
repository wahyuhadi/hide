package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"

	"github.com/auyer/steganography"
)

var (
	img      = flag.String("image", "", "image location")
	message  = flag.String("message", "This is default message", "Message to  hide")
	function = flag.String("func", "encode", "encode or decode")
)

func main() {
	flag.Parse()
	if *function == "encode" {
		Hide(*img, *message)
	} else {
		Uhide("new-image.png", "tes")
	}
}

func Hide(pictureInputFile, message string) {

	inFile, err := os.Open(pictureInputFile) // Opens input file provided in the flags
	if err != nil {
		log.Fatalf("Error opening file %s: %v", pictureInputFile, err)
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile) // Reads binary data from picture file
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatalf("Error opening file %v", err)
	}
	encodedImg := new(bytes.Buffer)
	err = steganography.Encode(encodedImg, img, []byte(message)) // Calls library and Encodes the message into a new buffer
	if err != nil {
		log.Fatalf("Error encoding message into file  %v", err)
	}
	outFile, err := os.Create("new-image.png") // Creates file to write the message into
	if err != nil {
		log.Fatalf("Error creating file %s: %v", "new-image.png", err)
	}
	bufio.NewWriter(outFile).Write(encodedImg.Bytes()) // writes file to disk
}

func Uhide(pictureInputFile, messageOutputFile string) {
	inFile, err := os.Open(pictureInputFile) // Opens input file provided in the flags
	if err != nil {
		log.Fatalf("Error opening file %s: %v", pictureInputFile, err)
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal("error decoding file", img)
	}

	sizeOfMessage := steganography.GetMessageSizeFromImage(img) // Uses the library to check the message size

	msg := steganography.Decode(sizeOfMessage, img) // Read the message from the picture file
	fmt.Println(msg)
	// if the user specifies a location to write the message to...
	if messageOutputFile != "" {
		err := ioutil.WriteFile(messageOutputFile, msg, 0644) // write the message to the given output file

		if err != nil {
			fmt.Println("There was an error writing to file: ", messageOutputFile)
		}
	} else { // otherwise, print the message to STDOUT
		for i := range msg {
			fmt.Printf("%c", msg[i])
		}
	}
}
