package magejiTx

import "fmt"

type Transactor interface {
	Commit() error
	Rollback() error
}

type transactorStack struct {
	errChan chan error
	stack   []Transactor
}

type TransactorStack struct {
	s transactorStack
}

func (r transactorStack) Error() <-chan error {
	return r.errChan
}

func (t transactorStack) Commit() {
	t.rangez(func(ts Transactor) {
		if err := ts.Commit(); err != nil {
			t.errChan <- fmt.Errorf("commit: %q", err)
		}
	})
}
func (t transactorStack) Rollback() {
	t.rangez(func(ts Transactor) {
		if err := ts.Rollback(); err != nil {
			t.errChan <- fmt.Errorf("rollback: %q", err)
		}
	})
}
func (t transactorStack) Add(ts Transactor) {
	t.push(ts)
}
func (t transactorStack) rangez(fn func(ts Transactor)) {
	for i := len(t.stack) - 1; i > 0; i-- {
		fn(t.stack[i])
	}
}
func (t transactorStack) push(ts Transactor) {
	t.stack = append(t.stack, ts)
}
func (t transactorStack) pop() Transactor {
	ts := t.stack[len(t.stack)-1]
	newstack := make([]Transactor, 0, len(t.stack)-1)
	copy(newstack, t.stack[:len(t.stack)-1])
	t.stack = newstack
	return ts
}
func (t transactorStack) len() int {
	return len(t.stack)
}

func New(ts ...Transactor) TransactorStack {
	return TransactorStack{
		transactorStack{errChan: make(chan error), stack: ts},
	}
}
