package bashy

import (
	"bashy/src/utils"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	//"bufio"
	"strings"
	//"net/http"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type Bashy struct {
	Home              string
	app               cli.App
	scripts           []*Script
	ScriptFolder      string
	CacheFolder       string
	BinFolder         string
	BinScriptTemplate string
	Tmp               string
	Instance          string
}

func (re *Bashy) Init() error {

	currentUser, _ := user.Current()
	re.Home = os.Getenv("BASHY_HOME")
	if re.Home == "" {
		re.Home = filepath.Join(currentUser.HomeDir, ".bashy")
	}

	re.Home, _ = filepath.Abs(re.Home)

	if utils.DirectoryExists(re.Home) == false {
		fmt.Println("creating missing folder" + re.Home)
		err := os.MkdirAll(re.Home, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	re.Instance = utils.GenerateToken(5)
	re.ScriptFolder = filepath.Join(re.Home, "scripts")
	re.CacheFolder = filepath.Join(re.Home, "cache")
	re.BinFolder = filepath.Join(re.Home, "bin")
	re.Tmp = filepath.Join(re.Home, "tmp", re.Instance)
	err := os.MkdirAll(re.ScriptFolder, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(re.CacheFolder, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(re.BinFolder, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(re.Tmp, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	cmds := re.loadCommands(re.ScriptFolder)
	cmds = append(cmds, re.LoadInternalCommands()...)
	if len(cmds) == 0 {
		return errors.New("no commands found")
	}

	re.app = cli.App{
		Name:        "Bashy",
		Commands:    cmds,
		Description: "hello! script loaded from " + re.Home,
	}
	re.app.Commands = cmds
	return nil
}

func (re *Bashy) ExecCommand(args []string, lines []string) {
	filename := utils.TempFileName(re.Tmp, ".sh")
	utils.WriteLinesToFile(filename, lines, 0777)

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

func (re *Bashy) scriptFiles(scriptPath string) []fs.FileInfo {
	result := []fs.FileInfo{}

	files, err := ioutil.ReadDir(scriptPath)
	if err != nil {
		fmt.Println(err)
	}
	manualfiles := os.Getenv("BASHY_FILES")
	if len(manualfiles) > 0 {
		for _, filestr := range strings.Split(manualfiles, ",") {
			file, _ := os.Stat(filestr)
			files = append(files, file)
		}
	}

	extrahome := os.Getenv("BASHY_EXTRA")
	if len(manualfiles) > 0 {
		files, err = ioutil.ReadDir(extrahome)
		//append script on main directory
		for _, file := range files {
			filename := file.Name()
			if filepath.Ext(filename) == ".yaml" || filepath.Ext(filename) == ".yml" {
				result = append(result, file)
			}
		}
	}

	for _, file := range files {
		filename := file.Name()
		if filepath.Ext(filename) == ".yaml" || filepath.Ext(filename) == ".yml" {
			result = append(result, file)
		}
	}

	//append script on main directory
	for _, file := range files {
		filename := file.Name()
		if filepath.Ext(filename) == ".yaml" || filepath.Ext(filename) == ".yml" {
			result = append(result, file)
		}
	}

	return result
}

func (re *Bashy) loadConfigs(scriptPath string) []*Script {
	files := re.scriptFiles(scriptPath)
	/* loop over command folder */
	configList := []*Script{}
	for _, file := range files {
		filename, _ := filepath.Abs(filepath.Join(scriptPath, file.Name()))
		if filepath.Ext(filename) == ".yaml" || filepath.Ext(filename) == ".yml" {
			configList = append(configList, re.loadConfigFromFiles(filename)...)
		}
	}
	return configList
}

func (re *Bashy) loadConfigFromFiles(filename string) []*Script {

	configList := []*Script{}

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	d := yaml.NewDecoder(f)
	for {
		// create new spec here
		config := new(Script)
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
		//config.Name = config.Name + "_" + strconv.Itoa(i)

		// if config.Script != "" {
		// 	if strings.HasPrefix(config.Script, "/") {
		// 		//absolute local path, unchanged
		// 	} else if strings.HasPrefix(config.Script, "http://") || strings.HasPrefix(config.Script, "https://") {

		// 		//download it and chache
		// 		cacheScript := filepath.Join(os.TempDir(), config.Name+".sh")

		// 		if !utils.Exists(cacheScript) {
		// 			fmt.Println("downloading..")
		// 			utils.DownloadFile(cacheScript, config.Script)
		// 		}

		// 		config.Script = cacheScript

		// 	} else {
		// 		//local relative path
		// 		parent := filepath.Dir(filename)
		// 		config.Script = filepath.Join(parent, config.Script)
		// 	}
		// }
		configList = append(configList, config)
	}
	return configList
}

func (re *Bashy) convertToCommands(configs []*Script) []*cli.Command {
	result := []*cli.Command{}
	for _, config := range configs {
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
			if config.Script != "" {
				script := utils.ReadFileLines(config.Script)
				lines = append(lines, script...)
			}
			re.ExecCommand(args, lines)
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

func (re *Bashy) loadCommands(home string) []*cli.Command {
	configs := re.loadConfigs(home)
	commands := re.convertToCommands(configs)
	return commands
}

func (re *Bashy) Run(args []string) {
	err := re.app.Run(args)
	if err != nil {
		log.Fatal(err)
	}
}

func (re *Bashy) Destroy() {
	os.RemoveAll(re.Tmp)
	os.Remove(re.Tmp)
}
