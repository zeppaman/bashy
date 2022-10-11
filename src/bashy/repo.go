package bashy

import (
	"bashy/src/utils"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func (re *Bashy) removeScript(name string) {
	deletePattern := filepath.Join(re.ScriptFolder, name+".*")
	utils.RemoveAll(deletePattern)
	deletePattern = filepath.Join(re.CacheFolder, name+".*")
	utils.RemoveAll(deletePattern)

}
func (re *Bashy) addScript(name string) {

	source := name
	cacheScript := filepath.Join(re.Tmp, utils.GetMD5Hash(name)+".yml")

	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		utils.DownloadFile(cacheScript, name)
		source = cacheScript
		fmt.Println(source)
	}

	scripts := re.loadConfigFromFiles(source)
	for _, script := range scripts {
		scriptName := filepath.Join(re.ScriptFolder, script.Name+".yml")
		binName := filepath.Join(re.BinFolder, script.Name)
		bin := utils.Transform("bin", script)
		ioutil.WriteFile(binName, []byte(bin), 0777)

		if !utils.Exists(scriptName) {
			if strings.HasPrefix(script.Script, "/") {
				//absolute local path, keep unchanged
			} else if strings.HasPrefix(script.Script, "http://") || strings.HasPrefix(script.Script, "https://") {

				//download it and chache
				cacheScript := filepath.Join(re.Home, "cache", script.Name+".sh")
				utils.Copy(script.Script, cacheScript)
				script.Script = cacheScript

			} else {
				cacheScript := filepath.Join(re.Home, "cache", script.Name+".sh")
				utils.Copy(script.Script, cacheScript)
				script.Script = cacheScript
			}
			// Copy yaml file to script folder
			data, err := yaml.Marshal(&script)
			err = ioutil.WriteFile(scriptName, data, 0)

			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
func (re *Bashy) listScript(name string) {
	fmt.Println("Loaded scripts form: " + re.ScriptFolder)
	for _, script := range re.scriptFiles(re.ScriptFolder) {
		fmt.Println(strings.Replace(script.Name(), filepath.Ext(script.Name()), "", 1) + " " + script.ModTime().Local().String())
	}

}

func (re *Bashy) dumpSettings() {
	fmt.Println("Dump Settings ")
	fmt.Println("BinFolder: " + re.BinFolder)
	fmt.Println("BinScriptTemplate: " + re.BinScriptTemplate)
	fmt.Println("CacheFolder: " + re.CacheFolder)
	fmt.Println("Home: " + re.Home)
	fmt.Println("Instance: " + re.Instance)
	fmt.Println("ScriptFolder: " + re.ScriptFolder)
	fmt.Println("Tmp: " + re.Tmp)
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
		}}
}
