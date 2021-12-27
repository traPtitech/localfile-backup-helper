package main

import "sync"

type Result struct {
	errs      []error
	objectNum int
	sync.Mutex
}

func (r *Result) appendError(err error) {
	r.Lock()
	defer r.Unlock()
	r.errs = append(r.errs, err)
}

func (r *Result) addObjectNum() {
	r.Lock()
	defer r.Unlock()
	r.objectNum++
}
