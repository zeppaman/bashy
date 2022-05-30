package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	// "os/exec"
	// "path/filepath"
	// "strconv"
	// "github.com/urfave/cli/v2"
	// "gopkg.in/yaml.v2"
	"bufio"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

func DirectoryExists(name string) bool {
	return Exists(name)
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func WriteToFile(filename string, lines []string) string {
	file, err := os.CreateTemp(os.TempDir(), "bashy*.sh")
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range lines {
		file.WriteString(line + "\n")
	}
	file.Close()
	os.Chmod(file.Name(), 0777)
	return file.Name()
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

func ReadFileLines(filename string) []string {
	lines := []string{}
	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return lines
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("read file line error: %v", err)
			return lines
		}

		lines = append(lines, strings.Replace(line, "\r", "", -1))
	}
	return lines
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func RemoveAll(pattern string) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		fmt.Println("removing " + f)
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}
