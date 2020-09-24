package exqueue

import (
	"sync"
	"context"
	"github.com/lichao-mobanche/go-extractor-server/pkg/request"
	"github.com/lichao-mobanche/go-extractor-server/server/global"
	"github.com/lichao-mobanche/go-extractor-server/pkg/extract"

	"github.com/spf13/viper"
)

type ExQueue struct {

	Threads int
	*exList
	wake    chan struct{}
	mut     sync.Mutex // guards wake
	exitc	chan struct{}
}

type exList struct {
	MaxSize int
	lock *sync.RWMutex
	size    int
	first *exItem
	last *exItem
}

type exItem struct {
	Request *request.Request
	Next    *exItem
}

// New creates a new queue
func New(conf *viper.Viper) (*ExQueue, error){
	var workerNumber, maxQueueNum int
	if conf.IsSet("worker.number") {
		workerNumber = conf.GetInt("worker.number")
	}
	if conf.IsSet("queue.number") {
		maxQueueNum = conf.GetInt("queue.number")
	}
	return &ExQueue{
		workerNumber,
		&exList{MaxSize: maxQueueNum},
		nil,
		sync.Mutex{},
		make(chan struct{}),
	},nil
}

// AddRequest adds a new Request to the queue
func (q *ExQueue) AddRequest(r *request.Request) error {
	q.mut.Lock()
	waken := q.wake != nil
	q.mut.Unlock()
	if !waken {
		return global.QueueUnavailableError(r.URL)
	}
	err := q.Add(r)
	if err != nil {
		return err
	}
	q.wake <- struct{}{}
	return nil
}

func (q *ExQueue) run() error {
	q.mut.Lock()
	if q.wake != nil {
		q.mut.Unlock()
		panic("cannot call duplicate ExQueue.Run")
	}
	q.wake = make(chan struct{})
	q.mut.Unlock()

	requestc := make(chan *request.Request)
	complete, errc := make(chan struct{}), make(chan error, 1)
	for i := 0; i < q.Threads; i++ {
		go independentRunner(requestc, complete)
	}
	go q.loop(requestc, complete, errc)
	defer close(requestc)
	return <-errc
}

func (q *ExQueue) loop(requestc chan<- *request.Request, complete <-chan struct{}, errc chan<- error) {
	var active int
	
	for {
		sent := requestc
		req := q.Get()
		if req==nil{
			sent = nil
		}
	Sent:
		for {
			select {
			case sent <- req:
				active++
				break Sent
			case <-q.wake:
				if sent == nil {
					break Sent
				}
			case <-complete:
				active--
				if sent == nil && active == 0 {
					break Sent
				}
			case <-q.exitc:
				if active<=0{
					goto End
				}
				q.exitc <- struct{}{}
			}
		}
	}
	End:
		errc<-nil
}

// OnStart TODO
func (q *ExQueue) OnStart(ctx context.Context) (err error) {
	go q.run()
	return
}

// OnStop TODO
func (q *ExQueue) OnStop(ctx context.Context) (err error) {
	q.exitc <- struct{}{}
	return
}

func independentRunner(requestc <-chan *request.Request, complete chan<- struct{}) {
	for req := range requestc {
		extract.Extract(req)
		complete <- struct{}{}
	}
}

// Add add request to exList
func (l *exList) Add(r *request.Request) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	// Discard URLs if size limit exceeded
	if l.MaxSize > 0 && l.size >= l.MaxSize {
		return global.QueueFullError()
	}
	i := &exItem{Request: r}
	if l.first == nil {
		l.first = i
	} else {
		l.last.Next = i
	}
	l.last = i
	l.size++
	return nil
}

// Get get request from exList
func (l *exList) Get() *request.Request {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.size == 0 {
		return nil
	}
	r := l.first.Request
	l.first = l.first.Next
	l.size--
	return r
}

// Size get the size of exList
func (l *exList) Size() (int, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.size, nil
}
