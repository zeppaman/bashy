package bashy

import (
	"bashy/src/logger"
	"bashy/src/utils"
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
	logger.BmagentaPrintln("Adding " + name)
	source := name
	cacheScript := filepath.Join(re.Tmp, utils.GetMD5Hash(name)+".yml")

	if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
		utils.DownloadFile(cacheScript, name)
		source = cacheScript
		logger.BmagentaPrintln(source)
	}

	scripts := re.loadConfigFromFiles(source)
	for _, script := range scripts {
		logger.BmagentaPrintln(" -  " + name)
		scriptName := filepath.Join(re.ScriptFolder, script.Name+".yml")
		binName := filepath.Join(re.BinFolder, script.Name)
		bin := utils.Transform("bin", script)
		ioutil.WriteFile(binName, []byte(bin), 0777)

		logger.ImagentaPrintln(scriptName, "imagenta")
		if !utils.Exists(scriptName) {
			if len(script.Script) == 0 {
				logger.GreenPrintln("Embedded script, nothing to include")
			} else if strings.HasPrefix(script.Script, "/") {
				//absolute local path, keep unchanged
				logger.GreenPrintln("Absoulte path used")
			} else if strings.HasPrefix(script.Script, "http://") || strings.HasPrefix(script.Script, "https://") {
				logger.GreenPrintln("Downloading and caching external script")
				//download it and chache
				cacheScript := filepath.Join(re.Home, "cache", script.Name+".sh")
				utils.Copy(script.Script, cacheScript)
				script.Script = cacheScript

			} else {
				logger.GreenPrintln("Moving relative path file to the cache")
				cacheScript := filepath.Join(re.Home, "cache", script.Name+".sh")
				utils.Copy(script.Script, cacheScript)
				script.Script = cacheScript
			}

			utils.SerializeToFile(scriptName, &script)

		}
	}

}

func (re *Bashy) listScript(name string) {
	logger.ImagentaPrintln("Loaded scripts form: " + re.ScriptFolder)
	for _, script := range re.scriptFiles(re.ScriptFolder) {
		logger.ImagentaPrintln(strings.Replace(script.Name(), filepath.Ext(script.Name()), "", 1) + " " + script.ModTime().Local().String())
	}

}

func (re *Bashy) dumpSettings() {
	logger.BgreenPrintln("Dump Settings: ")
	logger.IcyanPrintln("BinFolder: " + re.BinFolder)
	logger.IcyanPrintln("BinScriptTemplate: " + re.BinScriptTemplate)
	logger.IcyanPrintln("CacheFolder: " + re.CacheFolder)
	logger.IcyanPrintln("Home: " + re.Home)
	logger.IcyanPrintln("Instance: " + re.Instance)
	logger.IcyanPrintln("ScriptFolder: " + re.ScriptFolder)
	logger.IcyanPrintln("Tmp: " + re.Tmp)
}
