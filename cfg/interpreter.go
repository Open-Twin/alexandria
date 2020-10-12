package cfg

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var NumberConf map[string]int
var StringConf map[string]string

func ReadConf() []string {
	absPath, _ := filepath.Abs("alexandria/cfg/configuration.alex")
	file, err := os.Open(absPath)

	if err != nil {
		log.Fatalf("failed to open config")
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	var confText []string
	for scanner.Scan() {
		confText = append(confText, scanner.Text())
	}

	file.Close()

	processConf(confText)

	return confText
}

func processConf(filecontent []string) {
	NumberConf = make(map[string]int)
	StringConf = make(map[string]string)

	for i, s := range filecontent {
		// Zeile am = aufteilen auf 2 Substrings
		field := strings.SplitN(s, "=", 2)
		if len(field) == 2 {
			key := field[0]
			value := field[1]

			if len(value) > 3 && value[0] == '"' && value[len(value)-1] == '"' {
				StringConf[key] = value[1 : len(field[1])-1]
			} else {
				numberVal, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					fmt.Printf("Config-Error: Wrong formatted number on %d", i)
				} else {
					NumberConf[key] = int(numberVal)
				}
			}
		} else {
			fmt.Printf("Config-Error on line %d", i)
		}
	}
}
