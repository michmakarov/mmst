// performance
package main

import (
	"container/list"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type PerformRec struct {
	Id      int64
	AccName []byte
	R       *http.Request
	Start   time.Time
}

type PerfList struct {
	count    int64 //count of worked out requests
	notDone  *list.List
	totalDur time.Duration
	mtx      sync.Mutex
}

var perfList *PerfList

func init() {
	perfList = &PerfList{}
	perfList.notDone = list.New()
}

func (pl *PerfList) Reg(r *http.Request) {
	var newReq *PerformRec
	var recId int64
	pl.mtx.Lock()
	recId++
	pl.count = recId
	newReq = &PerformRec{
		Id:      recId,
		AccName: accountName(r),
		R:       r,
		Start:   time.Now(),
	}
	perfList.notDone.PushFront(newReq)
	pl.mtx.Unlock()
}

func (pl *PerfList) Done(r *http.Request) {
	var removedCount int
	pl.mtx.Lock()
	for e := pl.notDone.Front(); e != nil; e = e.Next() {
		if e.Value.(*PerformRec).R == r {
			pl.notDone.Remove(e)
			removedCount++

		}
	}
	if removedCount == 0 {
		pl.mtx.Unlock()
		panic(fmt.Sprintf("(pl *PerfList) Done: not a request my be done"))
	}
	if removedCount > 1 {
		pl.mtx.Unlock()
		panic(fmt.Sprintf("(pl *PerfList) Done: more than one request to done"))
	}

	pl.mtx.Unlock()
}

func (pl *PerfList) inPerforming(r *http.Request) bool {
	var performingCount int
	pl.mtx.Lock()
	for e := pl.notDone.Front(); e != nil; e = e.Next() {
		if e.Value.(*PerformRec).R == r {
			performingCount++

		}
	}
	if performingCount > 1 {
		pl.mtx.Unlock()
		panic(fmt.Sprintf("(pl *PerfList) inPerforming: more than one request registered"))
	}
	if performingCount == 1 {
		pl.mtx.Unlock()
		return true
	}

	pl.mtx.Unlock()
	return false
}

//220506 00:43
func (pl *PerfList) String(le string) (res string) {
	var pr *PerformRec
	var dur time.Duration
	var count int
	pl.mtx.Lock()
	res = fmt.Sprintf("not dane list  %s", le)
	for e := pl.notDone.Front(); e != nil; e = e.Next() {
		count++
		pr = e.Value.(*PerformRec)
		dur = time.Since(pr.Start)
		res = res + fmt.Sprintf("RA=%s;URL=%s;Start=%s; dur=%v%s", pr.R.RemoteAddr, pr.R.RequestURI, pr.Start.Format(timeFormat), dur, le)
	}
	res = res + fmt.Sprintf("-------count=%d%s", count, le)

	pl.mtx.Unlock()
	return
}
