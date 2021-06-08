package mimetype

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var typeMap = make(map[string]string)
func init() {
	lines := ReadFile()
	ParseFile(lines)
}

func GetMimeTypeMap() map[string]string{
	return typeMap
}

func ReadFile()[]string{

	f, err := os.Open("mimetype/mime.types")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	lines := make([]string,0)
	for {
		line, err := r.ReadString('\n')
		lines = append(lines, line)
		if err != nil {
			break
		}
	}
	return lines
}

func ParseFile(lines []string) {
	for i:= 0; i < len(lines) ; i++ {
		lines[i] = strings.TrimSpace(lines[i])
		if lines[i] == "" {
			continue
		}
		if strings.Contains(lines[i],"{"){
			continue
		}
		if strings.Contains(lines[i],"}"){
			continue
		}
		var mimeType []byte
		var suffix []byte
		var flag = 0
		for j := 0; j < len(lines[i]); j++ {
			if lines[i][j] != ' ' && flag ==0 {
				for{
					if j == len(lines[i]) || lines[i][j] == ' '{
						break
					}
					mimeType = append(mimeType,lines[i][j])
					j++
				}
				flag ++
			}

			if lines[i][j] != ' ' && flag ==1 {
				for{
					if j == len(lines[i]) || lines[i][j] == ' '{
						break
					}
					suffix = append(suffix,lines[i][j])
					j++
				}
				flag ++
			}
			if flag == 2 {
				break
			}

		}
		typeMap[string(mimeType)] = string(suffix)
	}
}

