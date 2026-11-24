package utils

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type EnvVariable int

const (
	envRoot EnvVariable = iota
	envConfigFile
	envReposDir
	envSSHKeyPath
	envSSHKeyPassword
)

var envMap = map[EnvVariable]string{
	envRoot:           "ENV_ROOT",
	envConfigFile:     "CONFIG_FILE",
	envReposDir:       "REPOS_DIR",
	envSSHKeyPath:     "SSH_KEY_PATH",
	envSSHKeyPassword: "SSH_KEY_PASSWORD",
}

func InitEnv() {
	envRoot := getEnv(envRoot)

	err := godotenv.Load(filepath.Join(envRoot, ".env.local"))
	if err != nil {
		LogWarning("fyi: couldn't open .env.local")
	}

	err = godotenv.Load(filepath.Join(envRoot, ".env"))
	CheckError(err, "error loading .env file")
}

func GetReposDir() string {
	return getEnv(envReposDir)
}

func GetConfigFile() string {
	return getEnv(envConfigFile)
}

func GetSSHKeyPath() string {
	return getEnv(envSSHKeyPath)
}

func GetSSHKeyPassword() string {
	return getEnv(envSSHKeyPassword)
}

func getEnv(env EnvVariable) string {
	return os.Getenv(envMap[env])
}
