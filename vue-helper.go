package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"errors"
)

func getCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0: i+1]), nil
}

func strFirstToUpper(str string) string {
	strArr := strings.Split(str, "-")
	var upperStr string
	for index, tempStr := range strArr {
		tempStrLen := len([]rune(tempStr))
		for i := 0; i < tempStrLen; i++ {
			if i == 0 {
				if index == 0 {
					upperStr += string([]rune(tempStr)[i])
				} else {
					upperStr += string([]rune(tempStr)[i] - 32)
				}
			} else {
				upperStr += string([]rune(tempStr)[i])
			}
		}
	}
	return upperStr
}

func readLines(lineNum int, content string, path string, source chan []string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lineTag = 1
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if lineTag == lineNum {
			lines = append(lines, "    "+content)
		}
		lines = append(lines, scanner.Text())
		lineTag++
	}
	source <- lines
	return lines, scanner.Err()
}

func createVueComponent(componentName string, insertContent string) {
	filename := "public/js/components/" + componentName + ".js"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	currentExecPath, _ := getCurrentPath()
	lines := make(chan []string)
	go readLines(7, insertContent, currentExecPath+"component.js", lines)
	bufferWriter := bufio.NewWriter(file)
	var lineTag = 1
	for _, line := range <-lines {
		if strings.Contains(line, "Component") {
			if lineTag == 3 {
				line = "Vue.component('" + componentName + "', $." + strFirstToUpper(componentName) + ");"
			} else {
				line = strings.Replace(line, "Component", strFirstToUpper(componentName), -1)
			}
			lineTag++
		}
		fmt.Fprintln(bufferWriter, line)
	}
	bufferWriter.Flush()
}

func main() {
	flag.Usage = func() {
		fmt.Printf("\nVue DevTools(0.0.1)\n")
		fmt.Printf("\nAvailable commands:\n")
		fmt.Println("component        (alias of: c)")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
	} else {
		switch flag.Args()[0] {
		case "component", "c":
			sourceDir := "app/views/" + flag.Args()[2] + "/template/"
			sourceFile := flag.Args()[1] + ".php"
			f, _ := os.Open(sourceDir + sourceFile)
			sourceContent, _ := ioutil.ReadAll(f)
			insertContent := "template:`" + string(sourceContent) + "`,"
			createVueComponent(flag.Args()[1], insertContent)
		}
	}
}
