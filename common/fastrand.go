package common

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var (
	randCh     = make(chan *rand.Rand, runtime.NumCPU())
	randChOnce sync.Once
)

const (
	chars          = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsLen       = len(chars)
	addrCharSet    = "abcdef0123456789"
	addrCharSetLen = len(addrCharSet)
)

func initRandCh() {
	for i := 0; i < runtime.NumCPU(); i++ {
		randCh <- rand.New(rand.NewSource(time.Now().Add(time.Second * time.Duration(i)).UnixNano()))
	}
}

func fastRandomString(l uint) string {
	randChOnce.Do(initRandCh)

	r := <-randCh
	s := make([]byte, l)
	for i := 0; i < int(l); i++ {
		s[i] = chars[r.Intn(charsLen)]
	}
	randCh <- r
	return string(s)
}

func fastRandomInt(min, max int) int {
	randChOnce.Do(initRandCh)
	r := <-randCh
	i := r.Intn(max-min) + min
	randCh <- r
	return i
}

func fastRandomInt63() int64 {
	randChOnce.Do(initRandCh)
	r := <-randCh
	i := r.Int63()
	randCh <- r
	return i
}

func fastRandomIntn(num int) int {
	randChOnce.Do(initRandCh)
	r := <-randCh
	i := r.Intn(num)
	randCh <- r
	return i
}

func fastRandomAddr() string {
	randChOnce.Do(initRandCh)

	r := <-randCh
	s := make([]byte, 40)
	for i := 0; i < 40; i++ {
		s[i] = addrCharSet[r.Intn(addrCharSetLen)]
	}
	randCh <- r
	return string(s)
}
