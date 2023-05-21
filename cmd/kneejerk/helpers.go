package main

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"net/url"
	"regexp"
	"strings"
)

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Helper Functions
func debugLog(debug bool, format string, v ...interface{}) {
	if debug {
		fmt.Printf(format, v...)
	}
}

func removeANSI(input string) string {
	return ansiEscape.ReplaceAllString(input, "")
}

func determineSeverity(envVar string) string {
	envVar = strings.ToUpper(envVar) // Ensure case-insensitive comparison
	if strings.Contains(envVar, "AWS") && (strings.Contains(envVar, "ACCESS") && (strings.Contains(envVar, "ID") || strings.Contains(envVar, "KEY"))) || strings.Contains(envVar, "SECRET") {
		return "high"
	} else if strings.Contains(envVar, "AWS") {
		return "medium"
	} else if strings.Contains(envVar, "API") && (strings.Contains(envVar, "URL") || strings.Contains(envVar, "HOST") || strings.Contains(envVar, "ROOT")) {
		return "low"
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
	} else if severity == "low" {
		severityColored = aurora.Green(severity).String()
	} else {
		severityColored = aurora.Blue(severity).String()
	}
	coloredMessage := fmt.Sprintf("[%s] [%s] [%s] %s [%s]", templateIDColored, outputTypeColored, severityColored, jsURL, match)
	uncoloredMessage := fmt.Sprintf("[%s] [%s] [%s] %s [%s]", templateID, outputType, severity, jsURL, match)
	return coloredMessage, uncoloredMessage
}

func urlJoin(baseURL string, relURL string) string {
	u, _ := url.Parse(baseURL)
	rel, _ := url.Parse(relURL)
	return u.ResolveReference(rel).String()
}
