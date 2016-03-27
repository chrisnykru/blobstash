/*

f is a tiny command line too to interact with BlobStash filetree extension.

Support BlobFS natively.

*/
package main

import (
	"flag"
	"fmt"
	"os"
	_ "syscall"

	"github.com/kuba--/xattr"
	_ "github.com/tsileo/blobstash/client/blobstore"
)

// TODO(tsileo): zsh autocomplete support

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: f [options] [<path>]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func getRef(path string) (string, bool) {
	ref, _ := xattr.Getxattr(path, "blobfs.ref")
	if len(ref) > 0 {
		return string(ref), true
	}
	return "", false
}

func main() {
	// var ns string
	// var comment string
	// var save bool
	// flag.StringVar(&ns, "ns", "", "Optional namespace (for upload)")
	// flag.StringVar(&comment, "comment", "", "Optional comment (for upload)")
	// flag.BoolVar(&save, "save", false, "Save the hash in the log (for upload)")
	// TODO(tsileo):
	// Store the latest uploaded blob as feature? log the uploaded? with an optional comment (import vkv)
	// sharing feature (with a blob app to register?)
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
	}
	// TODO(tsileo): an "extra" sub command
	path := flag.Arg(0)
	ref, _ := getRef(path)
	fmt.Printf("%s", ref)
	if ref == "" {
		fmt.Printf("Not inside a BlobFS\n")
		os.Exit(1)
	}
	// switch flag.Arg(0) {
	// case "extra":
	// 	if flag.NArg() == 1 {
	// 		fmt.Printf("get")
	// 	}
	// 	switch flag.Arg(1) {
	// 	case "get":
	// 		fmt.Printf("extra get")
	// 	}
	// }
	// opts := blobstore.DefaultOpts()
	// opts.SetHost(os.Getenv("BLOBSTASH_API_HOST"), os.Getenv("BLOBSTASH_API_KEY"))
	// blobstore := blobstore.New(opts)
	// if flag.Arg(0) == "-" {
	// blob, err := blobstore.Get(flag.Arg(0))
	os.Exit(0)
}
