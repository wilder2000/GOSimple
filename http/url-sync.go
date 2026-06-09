package http

import (
	"sync"

	"github.com/wilder2000/GOSimple/database"
	"github.com/wilder2000/GOSimple/dbmodel"
	"github.com/wilder2000/GOSimple/glog"
)

const (
	OPER_ID_ADMIN  = 10
	OPER_ID_VIEWER = 11
)

var (
	urlEntries = make(map[string]int32)
	urlMu      sync.RWMutex
)

func RegisterURL(url string, operatorID int32) {
	urlMu.Lock()
	defer urlMu.Unlock()
	if _, ok := urlEntries[url]; !ok {
		urlEntries[url] = operatorID
		glog.Logger.InfoF("RegisterURL: %s -> operator %d", url, operatorID)
	}
}

func CollectRegisteredURLs() []string {
	urlMu.RLock()
	defer urlMu.RUnlock()
	seen := make(map[string]bool)
	var result []string
	// admin URLs first
	for url, op := range urlEntries {
		if op == OPER_ID_ADMIN && !seen[url] {
			result = append(result, url)
			seen[url] = true
		}
	}
	// other operators
	for url, _ := range urlEntries {
		if !seen[url] {
			result = append(result, url)
			seen[url] = true
		}
	}
	return result
}

func UrlMappingsGrouped() map[int32][]string {
	urlMu.RLock()
	defer urlMu.RUnlock()
	grouped := make(map[int32][]string)
	for url, op := range urlEntries {
		grouped[op] = append(grouped[op], url)
	}
	return grouped
}

func SyncUrlMappings() {
	db := database.DBHander
	if db == nil {
		glog.Logger.InfoF("SyncUrlMappings: DB not initialized")
		return
	}

	grouped := UrlMappingsGrouped()
	totalInserted := 0

	for opID, urls := range grouped {
		for _, url := range urls {
			var count int64
			db.Model(&dbmodel.SUrlMapping{}).
				Where("operatorid = ? AND url = ?", opID, url).
				Count(&count)
			if count == 0 {
				if err := db.Create(&dbmodel.SUrlMapping{
					Operatorid: opID,
					Url:        url,
				}).Error; err != nil {
					glog.Logger.ErrorF("SyncUrlMappings: insert failed: url=%s op=%d err=%s", url, opID, err.Error())
				} else {
					totalInserted++
					glog.Logger.InfoF("SyncUrlMappings: inserted url=%s operator=%d", url, opID)
				}
			}
		}
	}

	glog.Logger.InfoF("SyncUrlMappings complete: %d new mappings inserted", totalInserted)
}
