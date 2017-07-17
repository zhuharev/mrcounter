package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
)

// тут должны быть глобальные переменные, лол, но их не будет

func main() {

	var (
		// k - количество параллельных обработчиков
		k = 5

		// канал, по которому передаются все поступившие в stdin строки
		// канал имеет буфер k, которым мы будем ограничивать количество
		// параллельных обаботчиков
		jobChan = make(chan string, k)

		// total - общее количество всех вхождений подстроки
		total uint64

		// на всякий случай
		wg sync.WaitGroup
	)

	// обработчик, который будет дёргать url и считать количество вхождений
	work := func(uri string) (count uint64, err error) {
		resp, err := http.Get(uri)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		// используем bytes.Count вместо strings.Count для избежания лишних
		// аллокаций при приведении типов из []byte в string
		count = uint64(bytes.Count(bts, []byte("Go")))

		return
	}

	go func() {
		for {
			select {
			case uri := <-jobChan:
				count, err := work(uri)
				if err != nil {
					log.Fatalln(err)
				}
				atomic.AddUint64(&total, count)
				fmt.Printf("Count for %s: %d\n", uri, count)
				wg.Done()
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		wg.Add(1)
		jobChan <- scanner.Text()
	}

	wg.Wait()

	fmt.Printf("Total: %d\n", total)
}

// тут тоже не будет ничего лишнего, а так хотелось, так хотелось
