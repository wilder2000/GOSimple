package config

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

const (
	FSep = "/"
	Logo = "  \n" +
		"         #        #      #      #\n" +
		"        #   #    #      #      #\n" +
		"       #  # #   #      ########\n" +
		"      # #   #  #      #      #\n" +
		"     ##     # #      #      #\n" +
		"    #       ##      #      #\n"
	LogoTitle               = "Power by GOSimple, A golang base framework."
	GOGO_HOME_ENV_KEY       = "GOGO_HOME"
	APP_CONFIG_FILE_ENV_KEY = "GOGO_CONFIG_FILE"
)

func AppConfigFile() string {

	appConfig := os.Getenv(APP_CONFIG_FILE_ENV_KEY)
	if appConfig == "" {
		return "Application.yaml"
	} else {
		return appConfig
	}

}

// AppDir 返回当前应用目录
func AppDir() string {
	path, _ := os.Getwd()

	gogohome := os.Getenv(GOGO_HOME_ENV_KEY)
	if gogohome == "" {
		//fmt.Printf("app dir:%s\n", path)
		return path
	} else {
		//fmt.Printf("app dir:%s\n", gogohome)
		return gogohome
	}

}

// ConfDir 返回当前配置目录
func ConfDir() string {

	return AppDir() + FSep + "conf" + FSep

}
func TempDir() string {
	tempDir := filepath.Join(AppDir(), "temp")
	_, err := os.Stat(tempDir)
	if err != nil && os.IsNotExist(err) {
		err2 := os.MkdirAll(tempDir, os.ModePerm)
		if err2 != nil {
			panic(err2)
		}
	}
	return tempDir
}
func HandleSignals() {
	// Wait for SIGINT, SIGQUIT, or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Println(fmt.Sprintf("shutting down in response to received signal,signal:%s", sig))
	//log.Logger.DebugF("shutting down in response to received signal,signal:%s", sig)
}
