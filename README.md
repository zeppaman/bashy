# Bashy
Bashy is a script collector for bash system. It allows to download and collect scripts, and facilitates argument parsing.

# Why I should use Bashy for my script
Bashy provides a easy way for resolving arguments and generating an useful Help.
You can simply have variable filled with values entered by the user by default, and the help infos build automatically from the script definition. This transofrms your bash script in a real console application with no pain

# How to install
There isnt any installer so far. You can install it by cloning the repo and running the install file.

# Linux or WSL
```bash
git clone <this repo url> bashy
cd bashy
bash ./install.sh
```
installation notes:
- this compile bashy form source, so you can delete the folder at the end
- the tool download and install go if not present. If present, nothing is touched but it have to be updated manually to the 1.19 version.

# Windows
```
git clone <this repo url> bashy
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; ./install.ps1
#rm -rf bashy optional
```
Despite that golang can produce multiplatform output and that this app could be improved to support multiple script engines (js,c#, etc..) and OS (windows, mac), at the moment the application is working only with linux and bash scripts.

# How is it possible?
The script hash a YAML file associated with it, so you can define all the infos. In the following example we are defining a command named `sample` with two parameters (`name, surname`).

```yaml
name: Name
description: the command description
argsusage: help text
params:
  - name: "name"
    desc: "enter your name"
  - name: "surname"
    desc: "enter your surname"
```
This will allow you to use directly the parameter names as named variables. So, your script will be able to do something similar to:

```bash
echo "$name"
echo "$surname"
echo "$name $surname"
```
Moreover, you will be able to list all the commands available, and for each commands the help usage.

# How can I share my script
Each yaml file can contains one or more command definition. Commands can be:
- embedded (the script is contained inside the yaml file definition)
- linked to an external files (local or remote)

You can work on-the-fly scripts on the current directory or move them to the internal report and make them available at OS level. Moreover you can publish scripts over HTTP and download it locally.

## Download a script and add it locally
```bash
bashy repo add https://xxx.yy/path.yml
```
This will download all related scripts and save them locally. The home of bashy is `.bashy` on the user's home.


## Run a script on current folder
You can run scripts on a foder witout installing them on bashy by specifing the `BASHY_EXTRA` (additional folder where the app looks for yaml) and `BASHY_FILES (a list of files path to add).
```bash
BASHY_EXTRA=samples/home go run ./src/main.go 
BASHY_FILES=myfile.yml,/var/scripts/myscript2.yml bashy
```

## Add a script to bashy
```bash
bashy repo add filepath.yml
```

# How to use it
Here some commands for usage.
```bash
bashy --help #list all commands with description

# NAME:
#    Bashy - A new cli application

# USAGE:
#    main [global options] command [command options] [arguments...]

# COMMANDS:
#    command1           
#    command2  

# GLOBAL OPTIONS:
#    --help, -h  show help (default: false)

bashy command1 --help # show infos about the command

# NAME:
#    command1 

# USAGE:
#    bashy command1 [command options] [arguments...]

# DESCRIPTION:
#    the command description

# OPTIONS:
#    --name value     enter your name
#    --surname value  enter your surname
#    --help, -h       show help (default: false)

bashy <command name> parameters args # execute the command

# the command output
```

# How to embed a script
Some samples about how to embed a script.
## Single command
You can write a command using 
```yaml
name: Name
description: the command description
argsusage: help text
params:
  - name: "name"
    desc: "enter your name"
  - name: "surname"
    desc: "enter your surname"
cmd: |
 echo "commmand"
 env
 echo "$name $surname"

```
or you can concatenate all commands in a single line if you want to keep things more hard to understand:
```yaml
name: Name
description: the command description
argsusage: help text
params:
  - name: "name"
    desc: "enter your name"
  - name: "surname"
    desc: "enter your surname"
cmd:  echo "commmand" &&  env  echo "$name $surname"
```

## Multiline script
In case you want to define multiple command you can define it by adding multiple `cmds` nodes. Each one can cantain multiple statement like in the Singleline script case. Cmd and Cmds can coexists: Cmd is executed AFTER the Cmds list.
```yaml
name: Name
description: the command description
argsusage: help text
cmd: echo "commmand" && env
params:
  - name: "name"
    desc: "enter your name"
  - name: "surname"
    desc: "enter your surname"
cmds:
    - echo "$name"
    - echo "$surname"
    - echo "$name $surname"
```

# How include an external script
An external script can be included. The external script is not exclusive: you can use it in conjunction with  `cmd, cmds` arguments, but loaded at the end. You can specify a local path (relative or absolute), or an URL. The remote script will be downloaded at the first usage then cached locally. Examples:
**Remote script**
```yaml
script: https://gist.githubusercontent.com/zeppaman/bbdfbac1304f88df1b905692e42f4d4e/raw/22ab3a2163f6ea481bca1b5d5570a964120a4f89/test-bashy.sh
```

**Local relative path**
The path is relative to the same location
```yaml
script: ./test-bashy.sh
```

**Local absolute path**
```yaml
script: /test-bashy.sh
```
Note: absolute paths are not changed during the import/execution process

# Extend Bashy 
The engine of bashy support any kind of script interpreter but it is shipped with sh for Linux and cmd for windows.
To add a new one just create a file like the following (this piece of code add the nodejs interpreter):
```yaml
name: node
params:
  - node
  - $filename
os: linux
installscript:
  - #!/bin/bash
  - echo "installing dependencies (node)"
  - '[[ "$(command -v apt)" ]] || { apt install nodejs -y; }'
  - exit 0
variabletemplate: var $name="$value";
```
Important fields are:
**name** the name of the interpreter. Into script YAML you will set the interpreter parameter to match this.
- **params** a list of parameters to run the script. The placeholder `$filename` will be replaced with the temporary path of the script that will be generated during the script run. The list in the example will be transformed into `node /`path/to/file`
- **os** the name of the os where the interpreter can be interpreted. If you have an interpreter that can be used in two different OS, you have to duplicate it
- **installscript** list of script steps for installing the interpreter. In the example above, it checks if apt is available, then installs node
- **variabletemplate** the template for adding variables on top of the script, `$name`` and `$value` are placeholders that will hold the param name and value as defined in the command YAML.

The installation of the interpreter can be done by using the `interpreter` command, like in the following example:

```bash
bash interpreter add samples/interpreter.yml
```

you can check the list of the loaded interpreters by typing:
Once the interpreter is installed, you can use it by adding the field `interpreter:

```yaml
name: Name
description: the command description
argsusage: help text
cmd: ...
interpreter: node
```

## Default interpreters

Bashy ships a set of default interpreters based on the OS. The list is the following.

| OS  | Name | Default |
| ------------- | ------------- | ------------- |
| LINUX | sh  | YES |
| WINDOWS  | cmd  | YES |

# Change default home

```bash
BASHY_HOME=samples bashy
```

# Debug
```bash
BASHY_HOME=tmp/home go run ./src/main.go 
BASHY_HOME=tmp/home  BASHY_EXTRA=samples/home go run ./src/main.go 
```
