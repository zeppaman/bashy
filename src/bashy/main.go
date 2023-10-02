package bashy

import (
	logger "bashy/src/logger"
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

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type Bashy struct {
	Home              string
	app               cli.App
	scripts           []*Script
	Interpreters      []*Interpreter
	ScriptFolder      string
	CacheFolder       string
	BinFolder         string
	BinScriptTemplate string
	Tmp               string
	Instance          string
	InterpreterFolder string
}

func (re *Bashy) Init() error {

	currentUser, _ := user.Current()
	re.Home = os.Getenv("BASHY_HOME")
	if re.Home == "" {
		re.Home = filepath.Join(currentUser.HomeDir, ".bashy")
	}

	re.Home, _ = filepath.Abs(re.Home)

	if !utils.DirectoryExists(re.Home) {
		logger.BgreenPrintln("creating missing folder" + re.Home)
		err := os.MkdirAll(re.Home, os.ModePerm)
		if err != nil {
			logger.Log("fatal", logger.JsonEncode(err))
		}
	}

	re.Instance = utils.GenerateToken(5)
	re.ScriptFolder = filepath.Join(re.Home, "scripts")
	re.CacheFolder = filepath.Join(re.Home, "cache")
	re.BinFolder = filepath.Join(re.Home, "bin")
	re.Tmp = filepath.Join(re.Home, "tmp", re.Instance)
	re.InterpreterFolder = filepath.Join(re.Home, "interpreters")
	err := os.MkdirAll(re.ScriptFolder, os.ModePerm)
	if err != nil {
		logger.Log("fatal", logger.JsonEncode(err))
	}
	err = os.MkdirAll(re.CacheFolder, os.ModePerm)
	if err != nil {
		logger.Log("fatal", logger.JsonEncode(err))
	}

	err = os.MkdirAll(re.BinFolder, os.ModePerm)
	if err != nil {
		logger.Log("fatal", logger.JsonEncode(err))
	}

	err = os.MkdirAll(re.Tmp, os.ModePerm)
	if err != nil {
		logger.Log("fatal", logger.JsonEncode(err))
	}

	err = os.MkdirAll(re.InterpreterFolder, os.ModePerm)
	if err != nil {
		logger.Log("fatal", logger.JsonEncode(err))
	}

	cmds := re.loadCommands(re.ScriptFolder)
	cmds = append(cmds, re.LoadInternalCommands()...)
	if len(cmds) == 0 {
		return errors.New("no commands found")
	}

	re.loadInterpreters()
	re.app = cli.App{
		Name:        "Bashy",
		Commands:    cmds,
		Description: "hello! script loaded from " + re.Home,
	}
	re.app.Commands = cmds
	return nil
}

func (re *Bashy) ExecCommand(interpreter *Interpreter, params map[string]string, lines []string) {
	filename := utils.TempFileName(re.Tmp, ".sh")
	linesToWrite := []string{}
	for name, value := range params {
		//fmt.Println("variable add" + interpreter.Variabletemplate)
		line := strings.ReplaceAll(interpreter.Variabletemplate, "$name", name)
		line = strings.ReplaceAll(line, "$value", value)
		linesToWrite = append(linesToWrite, line)
	}
	linesToWrite = append(linesToWrite, lines...)

	utils.WriteLinesToFile(filename, linesToWrite, 0777)
	//fmt.Printf("use " + interpreter.Name + " with file" + filename)

	commandArgs := []string{}
	for i, arg := range interpreter.Params {
		if i > 0 {
			arg = strings.ReplaceAll(arg, "$filename", filename)
			//fmt.Printf("arg:" + arg)
			if logger.IsDebug() {
				logger.Log("info", "arg:"+arg)
			}
			commandArgs = append(commandArgs, arg)
		}
	}
	cmd := exec.Command(interpreter.Params[0], commandArgs...)
	cmd.Env = os.Environ()

	// for _, arg := range args {
	// 	cmd.Env = append(cmd.Env, arg)
	// }

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		logger.Log("fatal", logger.JsonEncode(err))
	}
	fmt.Printf("%s\n", stdoutStderr)
}

func (re *Bashy) scriptFiles(scriptPath string) []fs.FileInfo {
	result := []fs.FileInfo{}

	files, err := ioutil.ReadDir(scriptPath)
	if err != nil {
		//fmt.Println(err)
		logger.Log("error", logger.JsonEncode(err))
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

	//TODO: Works only with --- ending. Fix it
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
			logger.Log("fatal", logger.JsonEncode(err))
			panic(err)
		}

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

			lines := []string{}
			lines = append(lines, config.Cmds...)
			lines = append(lines, config.Cmd)
			if config.Script != "" {
				script := utils.ReadFileLines(config.Script)
				lines = append(lines, script...)
			}
			//fmt.Println("Prepare flags ")
			params := make(map[string]string)
			names := c.FlagNames()
			for _, flag := range config.Params {
				logger.BcyanPrintln(" - "+flag.Name, "bcyan")
				if utils.Contains(names, flag.Name) {
					params[flag.Name] = c.String(flag.Name)
				} else {
					params[flag.Name] = flag.Default
				}
			}
			installInterpreter := re.GetInterpreterForCurrentOS(config.Interpreter)

			re.ExecCommand(installInterpreter, params, lines)
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
		logger.Log("fatal", logger.JsonEncode(err))
	}
}

func (re *Bashy) Destroy() {
	os.RemoveAll(re.Tmp)
	os.Remove(re.Tmp)
}

func (re *Bashy) LoadInternalCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "dumpsettings",
			Aliases: []string{"d"},
			Usage:   "dump settings for debug",
			Action: func(c *cli.Context) error {
				re.dumpSettings()
				return nil
			},
		},
		{
			Name:    "repo",
			Aliases: []string{"r"},
			Usage:   "manage local scripts",

			Subcommands: []*cli.Command{
				{
					Name:  "add",
					Usage: "add a new template",
					Action: func(c *cli.Context) error {
						re.addScript(c.Args().Get(0))
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) error {
						re.removeScript(c.Args().Get(0))
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "list installed scripts",
					Action: func(c *cli.Context) error {
						re.listScript(c.Args().Get(0))
						return nil
					},
				},
			},
		},
		{
			Name:    "interpreter",
			Aliases: []string{"i"},
			Usage:   "manage local interpreters",

			Subcommands: []*cli.Command{
				{
					Name:  "add",
					Usage: "add a new interpreter",
					Action: func(c *cli.Context) error {
						re.addInterpreter(c.Args().Get(0))
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing interpreter",
					Action: func(c *cli.Context) error {
						re.removeInterpreter(c.Args().Get(0))
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "list installed interpreter",
					Action: func(c *cli.Context) error {
						re.listInterpreters(c.Args().Get(0))
						return nil
					},
				},
			},
		},
	}
}
