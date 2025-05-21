package http

import (
	"github.com/wilder2000/GOSimple/comm"
	"github.com/wilder2000/GOSimple/config"
	"github.com/wilder2000/GOSimple/glog"
	"os"
	"path/filepath"
	"time"
)

type AttachManager struct {
	Home    string
	Expired time.Duration //seconds

}

const (
	AttachPrefix = "attach"
)

func (a *AttachManager) InitHome() {
	a.Home = filepath.Join(config.AConfig.StaticDir.AbsoluteFileDir, AttachPrefix)
	if e := os.MkdirAll(a.Home, os.ModePerm); e != nil {
		glog.Logger.ErrorF("init attach home failed. %s", e.Error())
	}
	a.Expired = 1 * time.Hour
	glog.Logger.InfoF("attach home:%s", a.Home)
	go a.checkAndDeleteOldFiles(a.Home, 10*time.Second)
}
func (a *AttachManager) RequestFile() (string, string) {
	fn := comm.LowerUUID() + ".xlsx"
	glog.Logger.InfoF("RequestFile found attach home:%s", a.Home)
	filename := filepath.Join(a.Home, fn)
	urlPath := filepath.Join(config.AConfig.StaticDir.RelativePath, AttachPrefix, fn)
	return filename, urlPath
}

// checkAndDeleteOldFiles 定时检查目录中的文件并删除超时的文件
func (a *AttachManager) checkAndDeleteOldFiles(directory string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 检查是否是文件（不是目录）
			if !info.IsDir() {
				// 获取文件的最后修改时间
				lastModified := info.ModTime()
				// 计算文件是否超时
				if time.Since(lastModified) > a.Expired {
					// 删除超时文件
					glog.Logger.InfoF("Attach Clean: Deleting old file: %s\n", path)
					err := os.Remove(path)
					if err != nil {
						glog.Logger.InfoF("Attach Clean: Failed to delete file: %s, error: %v\n", path, err)
					}
				}
			}
			return nil
		})

		if err != nil {
			glog.Logger.InfoF("Error walking the path %s: %v\n", directory, err)
		}
	}
}
