package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	jbossHome       string = "/opt/server/wildfly"
	gobossArgs      string = "-b localhost"
	gobossBuildArgs string = "-T 1C,package"
	logLevel        string = "debug"
)

func TestMain(m *testing.M) {
	before()
	exitResult := m.Run()
	after()
	os.Exit(exitResult)
}

func before() {
	os.Setenv("JBOSS_HOME", jbossHome)
	os.Setenv("GOBOSS_ARGS", gobossArgs)
	os.Setenv("GOBOSS_BUILD_ARGS", gobossBuildArgs)
	os.Setenv("LOG_LEVEL", logLevel)

}

func after() {
	os.Unsetenv("JBOSS_HOME")
	os.Unsetenv("GOBOSS_ARGS")
	os.Unsetenv("GOBOSS_BUILD_ARGS")
	os.Unsetenv("LOG_LEVEL")
}

func TestConfigMustGet(t *testing.T) {
	var assert = assert.New(t)
	var cfg = MustGet()
	fmt.Printf("%s", cfg.BuildArgs)
	assert.Equal(jbossHome, cfg.JbossHome)
	assert.Equal(gobossArgs, cfg.Args)
	assert.ElementsMatch([...]string{"-T 1C", "package"}, cfg.BuildArgs)
	assert.Equal(logLevel, cfg.LogLevel)
}
