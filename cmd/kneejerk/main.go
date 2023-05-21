package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// ASCII Banner
const banner = `
 _  __                _           _    
| |/ /               (_)         | |   
| ' / _ __   ___  ___ _  ___ _ __| | __
|  < | '_ \ / _ \/ _ | |/ _ | '__| |/ /
| . \| | | |  __|  __| |  __| |  |   < 
|_|\_|_| |_|\___|\___| |\___|_|  |_|\_\              
                    |__/                
                               v0.1.5
`

// Pattern for .js files
var jsFilePattern = regexp.MustCompile(`.*\.js`)

// Regex to find environment variables in both formats
var envVarPattern = regexp.MustCompile(`(\b(?:NODE|REACT|AWS)[A-Z_]*\b\s*:\s*".*?")|(process\.env\.[A-Z_][A-Z0-9_]*)`)

var foundVars = map[string]struct{}{}

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

var outputFileWriter *bufio.Writer = nil

func removeANSI(input string) string {
	return ansiEscape.ReplaceAllString(input, "")
}

func determineSeverity(envVar string) string {
	envVar = strings.ToUpper(envVar) // Ensure case-insensitive comparison
	if strings.Contains(envVar, "AWS") && (strings.Contains(envVar, "ACCESS") && (strings.Contains(envVar, "ID") || strings.Contains(envVar, "KEY"))) || strings.Contains(envVar, "SECRET") {
		return "high"
	} else if strings.Contains(envVar, "AWS") {
		return "medium"
	} else {
		return "info"
	}
}

func colorizeMessage(templateID string, outputType string, severity string, jsURL string, match string) (string, string) {
	templateIDColored := aurora.BrightGreen(templateID).String()
	outputTypeColored := aurora.BrightBlue(outputType).String()
	var severityColored string
	if severity == "high" {
		severityColored = aurora.Red(severity).String()
	} else if severity == "medium" {
		severityColored = aurora.Yellow(severity).String()
	} else {
		severityColored = aurora.Blue(severity).String()
	}
	coloredMessage := fmt.Sprintf("[%s] [%s] [%s] %s [%s]", templateIDColored, outputTypeColored, severityColored, jsURL, match)
	uncoloredMessage := fmt.Sprintf("[%s] [%s] [%s] %s [%s]", templateID, outputType, severity, jsURL, match)
	return coloredMessage, uncoloredMessage
}

func scrapeJSFiles(u string, debug bool) {
	// Remove ANSI escape sequences from the URL
	cleanUrl := removeANSI(u)

	res, err := http.Get(cleanUrl)
	if err != nil {
		fmt.Printf("Failed to get %s: %v\n", u, err)
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Printf("Failed to parse %s: %v\n", u, err)
		return
	}

	doc.Find("script, link").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if src == "" {
			src, _ = s.Attr("href")
		}
		if src != "" && strings.Contains(src, "/static/") && jsFilePattern.MatchString(src) {
			jsURL := urlJoin(u, src)
			jsRes, err := http.Get(jsURL)
			if err != nil {
				fmt.Printf("Failed to get %s: %v\n", jsURL, err)
				return
			}
			defer jsRes.Body.Close()

			jsContent, err := ioutil.ReadAll(jsRes.Body)
			if err != nil {
				fmt.Printf("Failed to read %s: %v\n", jsURL, err)
				return
			}

			// Remove ANSI escape sequences
			cleanJsContent := removeANSI(string(jsContent))

			matches := envVarPattern.FindAllString(cleanJsContent, -1)
			for _, match := range matches {
				if _, ok := foundVars[match]; !ok {
					foundVars[match] = struct{}{}
					severity := determineSeverity(match)
					coloredMessage, uncoloredMessage := colorizeMessage("kneejerk", "js", severity, jsURL, match)
					fmt.Println(coloredMessage)
					if outputFileWriter != nil {
						_, _ = outputFileWriter.WriteString(uncoloredMessage + "\n")
						_ = outputFileWriter.Flush()
					}
				}
			}
		}
	})
}

func urlJoin(baseURL string, relURL string) string {
	u, _ := url.Parse(baseURL)
	rel, _ := url.Parse(relURL)
	return u.ResolveReference(rel).String()
}

func main() {
	fmt.Println(banner)

	url := flag.String("u", "", "URL of the website to scan")
	list := flag.String("l", "", "Path to a file containing a list of URLs to scan")
	output := flag.String("o", "", "Path to output file")
	debug := flag.Bool("debug", false, "Print debugging statements")
	flag.Parse()

	if *output != "" {
		file, err := os.Create(*output)
		if err != nil {
			fmt.Printf("Failed to create %s: %v\n", *output, err)
			return
		}
		defer file.Close()

		outputFileWriter = bufio.NewWriter(file)
	}

	if *url != "" {
		scrapeJSFiles(*url, *debug)
	} else if *list != "" {
		file, err := os.Open(*list)
		if err != nil {
			fmt.Printf("Failed to open %s: %v\n", *list, err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Println(scanner.Text())                // print the input before processing
			cleanedInput := removeANSI(scanner.Text()) // Remove color codes
			if outputFileWriter != nil {
				_, _ = outputFileWriter.WriteString(cleanedInput + "\n")
				_ = outputFileWriter.Flush()
			}
			urlParts := strings.Split(cleanedInput, " ")
			if len(urlParts) > 3 {
				scrapeJSFiles(urlParts[3], *debug)
			} else {
				fmt.Println("Invalid input:", cleanedInput)
			}
		}
	} else if info, _ := os.Stdin.Stat(); info.Mode()&os.ModeCharDevice == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println(scanner.Text())                // print the input before processing
			cleanedInput := removeANSI(scanner.Text()) // Remove color codes
			if outputFileWriter != nil {
				_, _ = outputFileWriter.WriteString(cleanedInput + "\n")
				_ = outputFileWriter.Flush()
			}
			urlParts := strings.Split(cleanedInput, " ")
			if len(urlParts) > 3 {
				scrapeJSFiles(urlParts[3], *debug)
			} else {
				fmt.Println("Invalid input:", cleanedInput)
			}
		}
	}
}
