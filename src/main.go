package main

import (
	"bashy/src/model"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io"
	"errors"
	"bufio"
	"strings"
	"net/http"
)

func exists(name string) bool {
    _, err := os.Stat(name)
    if err == nil {
        return true
    }
    if errors.Is(err, os.ErrNotExist) {
        return false
    }
    return false
}

func writeToFile(filename string, lines []string) string {
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

func execCommand(args []string, lines []string) {
	filename := writeToFile("", lines)

	cmd := exec.Command("sh", "-c", filename)
	cmd.Env = os.Environ()

	for _, arg := range args {
		cmd.Env = append(cmd.Env, arg)
	}

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}

func downloadFile(filepath string, url string) error {

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

func loadConfigs(scriptPath string)  []*model.Script {
	/* loop over command folder */
	configList:=[]*model.Script{}
	for i := 0; i < 10; i++ {
		filename, _ := filepath.Abs(scriptPath)

		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		d := yaml.NewDecoder(f)
		for {
			// create new spec here
			config := new(model.Script)
			// pass a reference to spec reference
			err := d.Decode(&config)
			// check it was parsed
			if config == nil {
				continue
			}
			// break the loop in case of EOF
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				panic(err)
			}
			config.Name = config.Name + "_" + strconv.Itoa(i)
			configList=append(configList,config)
			if(config.Script !=""){
				if(strings.HasPrefix(config.Script,"/")){ 
					//absolute local path, unchanged
				} else if strings.HasPrefix(config.Script,"http://") || strings.HasPrefix(config.Script,"https://") { 
				
					//download it and chache
					cacheScript:=filepath.Join(os.TempDir(),config.Name+".sh")
					if !exists(cacheScript){
						downloadFile(config.Script,cacheScript)
					}

					config.Script=cacheScript
					
				}else{
					//local relative path
					parent := filepath.Dir(filename)
					config.Script=filepath.Join(parent, config.Script)
				}
			}
			
		}
	}

	
	return configList
}
func readFileLines(filename string) []string {
	lines:=[]string{}
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

		lines=append(lines, strings.Replace(line,"\r","",-1))
    }
	return lines
}

func convertToCommands(configs []*model.Script) []*cli.Command{
	result:=[]*cli.Command{};
	for _,config := range configs{
	/* convert config to commands*/
		command := new(cli.Command)
		command.Name = config.Name
		command.ArgsUsage = config.ArgUsage
		command.Description = config.Description

		command.Action = func(c *cli.Context) error {

			args := []string{}
			for _, element := range c.FlagNames() {
				//fmt.Println(element)
				args = append(args, element+"="+c.String(element)+"")
			}

			lines := []string{}
			lines = append(lines, config.Cmds...)
			lines = append(lines, config.Cmd)
			if(config.Script != ""){				
				script:=readFileLines(config.Script);
				lines = append(lines, script...)
			}
			execCommand(args, lines)
			return nil
		}

		for _, element := range config.Params {
			param := new(cli.StringFlag)
			param.Name = element.Name
			param.Value = element.Default
			param.Usage = element.Desc
			param.Required = element.Required

			command.Flags = append(command.Flags, param)

		}
		result = append(result, command)
	}
	return result
}

func loadCommands() []*cli.Command {
	configs:= loadConfigs("samples/bash.yml")
	commands:= convertToCommands(configs)
	return commands
}
func main() {

	cmds := loadCommands()
	app := &cli.App{
		Name:     "Bashy",
		Commands: cmds,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
