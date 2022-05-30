# Bashy
Bashy is a script collector for bash system. It allows to download and collect scripts, and facilitates argument parsing.

# Why I should use Baschy for my script
Bashy provides a easy way for resolving arguments and generating an useful Help.
You can simply have variable filled with values entered by the user by default, and the help infos build automatically from the script definition. This transofrms your bash script in a real console application with no pain

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
An external script can be included. External script is not exclusive: you can use it in conjunction with  `cmd, cmds` arguments, but loaded at the end. You can specify an local path (relative or absolute), or an URL. The remote script will be downloaded at the first usage then cached locally. Examples:
```yaml
script: http:/
```