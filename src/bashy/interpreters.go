package bashy

import (
	"bashy/src/utils"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

func (re *Bashy) removeInterpreter(name string) {
	deletePattern := filepath.Join(re.InterpreterFolder, name+".*")
	utils.RemoveAll(deletePattern)
	deletePattern = filepath.Join(re.InterpreterFolder, name+".*")
	utils.RemoveAll(deletePattern)

}
func (re *Bashy) addInterpreter(path string) {
	fmt.Println("Add Interpreter " + path + " to " + re.InterpreterFolder)

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		cacheScript := filepath.Join(re.Tmp, utils.GetMD5Hash(path)+".yml")
		utils.DownloadFile(cacheScript, path)
		path = cacheScript
	}

	interpretersRead := re.loadInterpretersFromFile(path)
	interpretersToAdd := []*Interpreter{}

	for _, interpreter := range interpretersRead {
		if interpreter.Os == utils.CurrentOS() {
			interpretersToAdd = append(interpretersToAdd, interpreter)
			fmt.Println(interpreter.Name + " available for " + interpreter.Os)
		} else {
			fmt.Println("Skipped " + interpreter.Os + " of " + interpreter.Name)
		}
	}

	//remove interpreters with the same name for replacing and append
	tmp := []*Interpreter{}
	for _, interpreter := range re.Interpreters {
		hasInterpreter := false
		for _, element := range interpretersToAdd {
			hasInterpreter = hasInterpreter || element.Name == interpreter.Name
		}
		if !hasInterpreter {
			tmp = append(tmp, interpreter)
		}
	}
	re.Interpreters = tmp
	re.Interpreters = append(re.Interpreters, interpretersToAdd...)

	//execute installation scripts
	for _, interpreter := range interpretersToAdd {
		re.installInterpreter(*interpreter)
	}

	re.saveInterpreters()
	re.LoadInternalCommands()

}
func (re *Bashy) saveInterpreters() {
	deletePattern := filepath.Join(re.ScriptFolder, "*.yml")
	utils.RemoveAll(deletePattern)

	for index, interpreter := range re.Interpreters {
		path := filepath.Join(re.InterpreterFolder, interpreter.Name+"_"+strconv.Itoa(index)+".yml")
		utils.SerializeToFile(path, interpreter)
	}
}
func (re *Bashy) listInterpreters(name string) {
	fmt.Println("Loaded Interpreters form: " + re.InterpreterFolder)
	for _, Interpreter := range re.Interpreters {
		fmt.Println(Interpreter.Name)
	}

}
func (re *Bashy) loadInterpreters() {
	files := re.getInterpretersFileNames()
	for _, file := range files {
		fmt.Println("loading from " + file)
		interpreters := re.loadInterpretersFromFile(file)
		re.Interpreters = append(re.Interpreters, interpreters...)
	}
}

func (re *Bashy) getInterpretersFileNames() []string {
	result := []string{}

	files, err := ioutil.ReadDir(re.InterpreterFolder)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		filename := file.Name()

		fmt.Println(filename)
		if filepath.Ext(filename) == ".yaml" || filepath.Ext(filename) == ".yml" {
			result = append(result, filepath.Join(re.InterpreterFolder, filename))
		}
	}

	return result
}

func (re *Bashy) loadInterpretersFromFile(filename string) []*Interpreter {

	interpretersList := []*Interpreter{}

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	d := yaml.NewDecoder(f)
	for {
		// create new spec here
		config := new(Interpreter)
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

		interpretersList = append(interpretersList, config)
	}
	return interpretersList
}

func (re *Bashy) installInterpreter(interpreter Interpreter) {

	installInterpreter := re.GetInterpreterForCurrentOS(interpreter.Interpreter)
	fmt.Println("Installing " + interpreter.Name + " using " + installInterpreter.Name)

	re.ExecCommand(installInterpreter, make(map[string]string), interpreter.Installscript)
}

func (re *Bashy) GetInterpreterForCurrentOS(name string) *Interpreter {
	for _, interpreter := range re.Interpreters {
		if interpreter.Name == name && interpreter.Os == utils.CurrentOS() {
			return interpreter
		}
	}
	//if not found try to execute with the default one, first available for OS used
	for _, interpreter := range re.GetDefaultInterpreters() {
		if interpreter.Os == utils.CurrentOS() {
			return interpreter
		}
	}
	return nil
}

func (re *Bashy) GetDefaultInterpreters() []*Interpreter {
	interpreters := []*Interpreter{}
	//TODO add config file in assets, then add dinamically
	if "windows" == utils.CurrentOS() {
		interpreters = append(interpreters, &Interpreter{
			Name:   "bat",
			Os:     "windows",
			Params: []string{"RUN", "$filename"},
		})
	} else if "linux" == utils.CurrentOS() {

		interpreters = append(interpreters, &Interpreter{
			Name:   "sh",
			Os:     "linux",
			Params: []string{"sh", "-c", "$filename"},
		})

	}
	return interpreters

}
