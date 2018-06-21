## Goboss

- Goboss will build your java project, copy to jboss and run-it.
Tested with Jboss EAP 6.4/Wildfly10 on Ubuntu 16.04

## How To

1) Modify `projects.yml` adding your project information.

- dir: must be your project folder, for instance: /home/dude/git/myproject
- target-dir: must be your project target folder, where your main artifact will be generated. for instance: /home/dude/git/myproject/target (necessary when building multi-module projects)
- artifact: must be the main artifact generated file, for instance: myproject.ear

2) Set in your environment the necessary environment variables:

- JBOSS_HOME: the folder where your jboss server is located
- GOBOSS_ARGS: arguments that will be passed to the jboss to initialize, for instance: -b localhost --server-config=standalone.xml
- GOBOSS_BUILD_ARGS: arguments that will be passed to the maven to build the artifact, for instance: clean install -skipSomePlugin

## Running

First of all, setup your GO environment, GOPATH, etc.

Ensure that you have [dep](https://github.com/golang/dep) installed as well.
```
$ go get
```

```
$ dep ensure & go build main.go & ./main

or

$ dep ensure & go run main.go
```
