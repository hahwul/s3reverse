package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

const (
	s3path = "^(s3-|s3\\.)?(s3.*)\\.amazonaws\\.com"
	s3url  = "^s3://*"
	s3vh   = "(s3.*)\\.amazonaws\\.com$"
)

func banner() {
	VERSION := "v1.0.0"
	fmt.Printf(`
8""""8 eeee       8"""8  8"""" 88   8 8"""" 8"""8  8""""8 8"""" 
8         8       8   8  8     88   8 8     8   8  8      8     
8eeeee    8       8eee8e 8eeee 88  e8 8eeee 8eee8e 8eeeee 8eeee 
    88 eee8  eeee 88   8 88    "8  8  88    88   8     88 88    
e   88    88      88   8 88     8  8  88    88   8 e   88 88    
8eee88 eee88      88   8 88eee  8ee8  88eee 88   8 8eee88 88eee  
`)
	fmt.Println("by @hahwul | "+VERSION)
	fmt.Println("")
}


func main() {
	// input options
	iL := flag.String("iL", "", "input List")
	// to options
	tN := flag.Bool("tN", false, "to name")
	tS := flag.Bool("tS", false, "to s3 url")
	tP := flag.Bool("tP", false, "to path-style")
	tV := flag.Bool("tV", false, "to virtual-hosted-style")
	verify := flag.Bool("verify", false, "testing bucket(acl,takeover)")
	// output options
	oN := flag.String("oN", "", "Write output in Normal format (optional)")
	oA := flag.String("oA", "", "Write output in Array format (optional)")
	var s3Buckets []string
	flag.Parse()
	if flag.NFlag() == 0 {
		banner()
		flag.Usage()
		return
	}

	// accept domains on stdin
	if *iL == "" {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			target := strings.ToLower(sc.Text())
			var s3 = identifys3(target)
			if s3 != "" {
				s3Buckets = append(s3Buckets, s3)
			}
		}
	} else {
		target, err := readLinesOrLiteral(*iL)
		if err != nil {
			fmt.Println(err)
		}
		for _, s := range target {
			var s3 = identifys3(s)
			if s3 != "" {
				s3Buckets = append(s3Buckets, s3)
			}
		}

	}
	// Remove Deplicated value
	s3Buckets = unique(s3Buckets)
	// Printing
	if *verify {
		var wg sync.WaitGroup
		for _, s := range s3Buckets {
			wg.Add(1)
			go func(s string) {
				defer wg.Done()
				var DefaultTransport http.RoundTripper = &http.Transport{}
				var resp *http.Response
				req, err := http.NewRequest("GET", "https://s3.amazonaws.com/"+s, nil)
				if err != nil {
					fmt.Printf("%s", err)
				} else {
					resp, err = DefaultTransport.RoundTrip(req)
					if err != nil {
						fmt.Printf("%s", err)
						return
					}
				}
				defer resp.Body.Close()
				contents, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("%s", err)
				}
				if strings.Contains(string(contents), "<Code>NoSuchBucket</Code>") {
					fmt.Println("[NoSuchBucket] " + s)
				} else if strings.Contains(string(contents), "<Code>AccessDenied</Code>") {
					fmt.Println("[PublicAccessDenied] " + s)
				} else {
					fmt.Println("[PublicAccessGranted] " + s)
				}
			}(s)
		}
		wg.Wait()
	} else {
		for _, s := range s3Buckets {
			if *tN {
				fmt.Println(s)
			}
			if *tS {
				fmt.Println("s3://" + s)
			}
			if *tP {
				fmt.Println("https://s3.amazonaws.com/" + s)
			}
			if *tV {
				fmt.Println(s + ".s3.amazonaws.com")
			}
		}
	}
	_ = oN
	_ = oA
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func identifys3(t string) string {
	// images.skypicker.com-dev.s3-website-eu => images.skypicker.com-dev1
	target := strings.Replace(t, "http://", "", 1)
	target = strings.Replace(target, "https://", "", 1)
	target = strings.Replace(target, "s3://", "s3:////", 1)
	target = strings.Replace(target, "//", "", 1)

	path, _ := regexp.MatchString(s3path, target)
	vh, _ := regexp.MatchString(s3vh, target)
	url, _ := regexp.MatchString(s3url, target)

	if path {
		target = strings.Replace(target, "s3.amazonaws.com/", "", 1)
		target = strings.Split(target, "/")[0]
	} else if vh {
		target = strings.Replace(target, ".s3.amazonaws.com", "", 1)
		target = strings.Split(target, "/")[0]
	} else if url {
		target = strings.Replace(target, "s3://", "", 1)
		target = strings.Split(target, "/")[0]
	}
	return target
}

// readLines reads all of the lines from a text file in to
// a slice of strings, returning the slice and any error
func readLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	lines := make([]string, 0)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	return lines, sc.Err()
}

// readLinesOrLiteral tries to read lines from a file, returning
// the arg in a string slice if the file doesn't exist, unless
// the arg matches its default value
func readLinesOrLiteral(arg string) ([]string, error) {
	if isFile(arg) {
		return readLines(arg)
	}

	// if the argument isn't a file, but it is the default, don't
	// treat it as a literal value

	return []string{arg}, nil
}

// isFile returns true if its argument is a regular file
func isFile(path string) bool {
	f, err := os.Stat(path)
	return err == nil && f.Mode().IsRegular()
}
