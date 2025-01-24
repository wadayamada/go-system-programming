package main

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	copy()
	createRandomByteFile()
	createZipFile()
	callCopyN()
	stream_Q3_6()
	server()
}

func copy() {
	old_file, err := os.Open("old.txt")
	if err != nil {
		panic(err)
	}
	defer old_file.Close()

	new_file, err := os.Create("new.txt")
	if err != nil {
		panic(err)
	}
	defer new_file.Close()

	io.Copy(new_file, old_file)
}

func createRandomByteFile() {
	buffer := make([]byte, 1024)
	rand.Read(buffer)

	file, err := os.Create("random_byte_file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(buffer)
}

func createZipFile() {
	zipFile, err := os.Create("file.zip")
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	inzip, err := zipWriter.Create("inzip.txt")
	if err != nil {
		panic(err)
	}
	inzip.Write([]byte("hi!"))
}

func callCopyN() {
	reader := strings.NewReader("Hello, World!")
	copyN(os.Stdout, reader, 5)
	os.Stdout.Write([]byte("\n"))
}

func copyN(dest io.Writer, source io.Reader, length int64) {
	buffer := make([]byte, length)
	io.ReadFull(source, buffer)
	io.Copy(dest, bytes.NewReader(buffer))
}

var (
	computer    = strings.NewReader("COMPUTER")
	system      = strings.NewReader("SYSTEM")
	programming = strings.NewReader("PROGRAMMING")
)

func stream_Q3_6() {
	var stream io.Reader

	// A
	sectionReaderA := io.NewSectionReader(programming, 5, 1)
	// S
	sectionReaderS := io.NewSectionReader(system, 0, 1)
	// C
	sectionReaderC := io.NewSectionReader(computer, 0, 1)
	// I1
	sectionReaderI1 := io.NewSectionReader(programming, 8, 1)
	// I2
	sectionReaderI2 := io.NewSectionReader(programming, 8, 1)

	stream = io.MultiReader(sectionReaderA, sectionReaderS, sectionReaderC, sectionReaderI1, sectionReaderI2)
	io.Copy(os.Stdout, stream)
	os.Stdout.Write([]byte("\n"))
}

func server() {
	http.HandleFunc("/sample_string", sampleStringHandler)
	http.HandleFunc("/zip", zipHandler)
	http.ListenAndServe(":8080", nil)
}

func sampleStringHandler(w http.ResponseWriter, r *http.Request) {
	str := "Hello, World!\n"
	w.Write([]byte(str))
}

func zipHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=sample.zip")
	file, err := os.Open("file.zip")
	if err != nil {
		panic(err)
	}
	buffer := make([]byte, 1024)
	file.Read(buffer)
	defer file.Close()
	w.Write(buffer)
}
