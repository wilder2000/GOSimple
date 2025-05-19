package pool

import (
	"GOSimple/glog"
	"errors"
	"time"
)

var (
	AppendTimeout time.Duration
)

// PoolEngine 队列处理引擎
type PoolEngine[T interface{}] struct {
	TaskList chan *T       //任务例表
	timeout  time.Duration //超时时间
}

// Append 添加一个任务到队列中，如果队列满了，会有一个超时时间，过了超时时间返回失败的错误
func (de *PoolEngine[T]) Append(taskInfo *T) (bool, error) {
	timeout := time.After(time.Second * de.timeout)
	for {
		select {
		case de.TaskList <- taskInfo:
			return true, nil
		case <-timeout:
			return false, errors.New("try to append task to pool timeout")
		}
	}

}

// New 初始化创建一个文件处理引擎
func New[T interface{}](pooSize int, currTh int, proc ExeProcess[T]) *PoolEngine[T] {
	AppendTimeout = 10
	dEngine := &PoolEngine[T]{}
	dEngine.TaskList = make(chan *T, pooSize)
	dEngine.timeout = AppendTimeout
	for i := 0; i < currTh; i++ {
		go run[T](proc, dEngine)
	}
	return dEngine
}

// 启动引警
func run[T interface{}](proc ExeProcess[T], engine *PoolEngine[T]) {
	for {
		select {
		case task := <-engine.TaskList:
			proc(task)

		case <-time.After(2 * time.Second):
			glog.Logger.DebugF("wait 2 second. array length:%d", len(engine.TaskList))
		}

	}

}

type ExeProcess[T interface{}] func(task *T)
