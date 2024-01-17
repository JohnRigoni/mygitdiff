package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)


var fileNameRgx = regexp.MustCompile(`diff --git a/(\S+)`)
var lineNumberRgx = regexp.MustCompile(`@@ -\d+,\d+ \+(\d+),\d+ @@`)
func main() {
	cmd := exec.Command("git", "diff")
	outputBytes, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	output := string(outputBytes)
	if output == "" {
		processAndPrintLine("No diff")
	}	
	var filename string
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		matches := fileNameRgx.FindStringSubmatch(line)
		if len(matches) > 1 {
			filename = matches[1]
			fmt.Println()
		}
		matches = lineNumberRgx.FindStringSubmatch(line)
		if len(matches) > 1 {
			lineNum := matches[1]
			out := fmt.Sprintf("./%s:%s", filename, lineNum)
			processAndPrintLine(out)
		}
		processAndPrintLine(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 50) //for :term in nvim
}


func processAndPrintLine(line string) {
	colorizedLine := processLineWithColor(line)
	fmt.Println(colorizedLine)
}

func processLineWithColor(line string) string {
	if strings.HasPrefix(line, "+") {
		return color.GreenString(line)
	}
	if strings.HasPrefix(line, "-") {
		return color.RedString(line)
	}
	if strings.HasPrefix(line, "diff --git") {
		return color.YellowString(line)
	}
	if strings.HasPrefix(line, "@@") {
		return color.BlueString(line)
	}
	if strings.HasPrefix(line, "./") {
		return color.HiMagentaString(line)
	}
	if strings.HasPrefix(line, "No diff") {
		return color.GreenString(line)
	}

	return line
}
