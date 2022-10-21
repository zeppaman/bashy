package bashy

import (
	"bashy/src/utils"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func (re *Bashy) removeScript(name string) {
	deletePattern := filepath.Join(re.ScriptFolder, name+".*")
	utils.RemoveAll(deletePattern)
	deletePattern = filepath.Join(re.CacheFolder, name+".*")
	utils.RemoveAll(deletePattern)

}
func (re *Bashy) addScript(name string) {
	fmt.Println("Adding " + name)
	source := name
	cacheScript := filepath.Join(re.Tmp, utils.GetMD5Hash(name)+".yml")

	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		utils.DownloadFile(cacheScript, name)
		source = cacheScript
		fmt.Println(source)
	}

	scripts := re.loadConfigFromFiles(source)
	for _, script := range scripts {
		fmt.Println(" -  " + name)
		scriptName := filepath.Join(re.ScriptFolder, script.Name+".yml")
		binName := filepath.Join(re.BinFolder, script.Name)
		bin := utils.Transform("bin", script)
		ioutil.WriteFile(binName, []byte(bin), 0777)

		fmt.Println(scriptName)
		if !utils.Exists(scriptName) {
			if len(script.Script) == 0 {
				fmt.Println("Embedded script, nothing to include")
			} else if strings.HasPrefix(script.Script, "/") {
				//absolute local path, keep unchanged
				fmt.Println("Absoulte path used")
			} else if strings.HasPrefix(script.Script, "http://") || strings.HasPrefix(script.Script, "https://") {
				fmt.Println("Downloading and caching external script")
				//download it and chache
				cacheScript := filepath.Join(re.Home, "cache", script.Name+".sh")
				utils.Copy(script.Script, cacheScript)
				script.Script = cacheScript

			} else {
				fmt.Println("Moving relative path file to the cache")
				cacheScript := filepath.Join(re.Home, "cache", script.Name+".sh")
				utils.Copy(script.Script, cacheScript)
				script.Script = cacheScript
			}

			utils.SerializeToFile(scriptName, &script)

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
