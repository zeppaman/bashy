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

func (re Bashy) removeScript(name string) {
	deletePattern := filepath.Join(re.ScriptFolder, name+".*")
	utils.RemoveAll(deletePattern)
	deletePattern = filepath.Join(re.CacheFolder, name+".*")
	utils.RemoveAll(deletePattern)

}
func (re Bashy) addScript(name string) {

	source := name
	cacheScript := filepath.Join(re.Home, "tmp", utils.GetMD5Hash(name)+".sh")

	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		utils.DownloadFile(cacheScript, name)
		source = cacheScript
	}

	scripts := re.loadConfigFromFiles(source)
	for _, script := range scripts {
		scriptName := filepath.Join(re.ScriptFolder, script.Name+".yml")
		fmt.Println(scriptName)
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
			data, err := yaml.Marshal(&script)
			err = ioutil.WriteFile(scriptName, data, 0)

			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
func (re Bashy) listScript(name string) {
	fmt.Println("Loaded scripts form: " + re.Home)
	for _, script := range re.scriptFiles(re.Home) {
		fmt.Println(strings.Replace(script.Name(), filepath.Ext(script.Name()), "", 1) + " " + script.ModTime().Local().String())
	}

}
func (re Bashy) LoadInternalCommands() []*cli.Command {
	return []*cli.Command{
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
