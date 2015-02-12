package main

import (
	"archive/zip"
	"strings"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"github.com/yani-/wpress"
)

func usage() {
	// Create our usage string
	usage := fmt.Sprintf(
		"Usage: %s COMMAND file\n" +
			"\n" +
			"WPRESS utility by ServMask Inc.\n" +
			"\n" +
			"Commands:\n" +
			"    extract    Extract a wpress file\n" +
			"    compress   Compress a folder or file to wpress archive\n" +
			"    convert    Convert a zip archive to wpress archive\n", filepath.Base(os.Args[0]))

	// Display the usage string
	fmt.Println(usage)

	// Exit normally
	os.Exit(0)
}

func main() {
	// Do we have the exact number of arguments?
	if len(os.Args) != 3 {
		usage()
	}

	task := os.Args[1]
	file := os.Args[2]

	switch task {
		case "extract":
			fmt.Println(file)
			//extract(file)
		case "compress":
			fmt.Println(file)
			//compress(file)
		case "convert":
			// Set destination folder
			dest := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

			// Make destination folder
			os.MkdirAll(dest, 0777)

			// Extract zip archive
			unzip(file, dest)

			// Compress into wpress format
			compress(dest)

			// Remove destination folder
			// os.RemoveAll(dest)

			// Print OK
			fmt.Println("OK")
		default:
			usage()
	}
}

func compress(src string) {
	archiver, _ := wpress.NewWriter(fmt.Sprintf("%s.wpress", src))

	// Go to destination folder
	os.Chdir(src)

	// Add directory recursively
	archiver.AddDirectory(".")
	archiver.Close()
}

func unzip(zipfile string, dest string) {
	reader, err := zip.OpenReader(zipfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer reader.Close()

	for _, f := range reader.Reader.File {
		zipped, err := f.Open()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer zipped.Close()

		// Get the individual file name and extract the current directory
		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			fmt.Println("Creating directory", path)
		} else {
			writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, f.Mode())

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			defer writer.Close()

			if _, err = io.Copy(writer, zipped); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("Decompressing : ", path)
		}
	}
}