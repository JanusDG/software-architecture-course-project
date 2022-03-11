package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hazelcast/hazelcast-go-client"
)

func main() {
	args := os.Args[1:]
	m := map[string]interface{}{
		"1":  task1,
		"2":  task2,
		"3":  task3,
		"-h": help,
	}
	if len(args) > 2 {
		return
	}
	m[args[0]].(func(string))(args[1])

}

// 1.
// Продемонструйте роботу Distributed Map
// http://docs.hazelcast.org/docs/latest/manual/html-single/index.html#map
// використовуючи API створіть Distributed Map
// запишіть в неї 1000 значень з ключем від 0 до 1к
// за допомогою Management Center (https://hazelcast.org/imdg/get-started/) подивиться на розподіл значень по нодах
// подивитись як зміниться розподіл якщо відключити одну ноду/дві ноди. Чи буде втрата даних?

func task1(a string) {
	ctx := context.TODO()
	// create the client and connect to the cluster on localhost
	client, err := hazelcast.StartNewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	dmp, err := client.GetMap(ctx, "default")

	for i := 1; i < 1001; i++ {
		// dmp.Remove(ctx, i)
		dmp.Put(ctx, i, i)
		time.Sleep(time.Second / 1000)
	}

	age, err := dmp.Get(ctx, 1000)
	fmt.Printf("%d\n", age)
	client.Shutdown(ctx)
}

// 2.
// Продемонструйте роботу Distributed Map with locks
// використовуючи 3 підключення (чи підключення з 3х клієнтів)
// одночансо запустіть приклади за посиланням з документації з підрахунку значень в циклі:
//  a) без блокування; б) з песимістичним; в) з оптимістичним блокуванням
// http://docs.hazelcast.org/docs/latest/manual/html-single/index.html#locking-maps
// порівняйте результати кожного з запусків

func task2_no_lock() {

	ctx := context.TODO()
	// create the client and connect to the cluster on localhost
	client, err := hazelcast.StartNewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	dmp, err := client.GetMap(ctx, "default")
	key := "1"
	dmp.Put(ctx, key, 1)
	for i := 0; i < 1000; i++ {
		value, _ := dmp.Get(ctx, key)
		intValue := value.(int64) + 1
		dmp.Put(ctx, key, intValue)
		time.Sleep(time.Second / 1000)
	}

	value, err := dmp.Get(ctx, key)
	fmt.Println("Finished! Result =", value)
	client.Shutdown(ctx)

}

func task2_pes_lock() {
	ctx := context.TODO()
	hz, _ := hazelcast.StartNewClient(ctx)
	myMap, _ := hz.GetMap(ctx, "default")
	key := "1"
	myMap.Put(ctx, key, 1)
	lockCtx := myMap.NewLockContext(ctx)
	for k := 0; k < 1000; k++ {
		myMap.Lock(lockCtx, key)
		value, _ := myMap.Get(lockCtx, key)
		intValue := value.(int64) + 1
		myMap.Put(lockCtx, key, intValue)
		myMap.Unlock(lockCtx, key)
	}
	value, _ := myMap.Get(lockCtx, key)
	fmt.Println("Finished! Result =", value)
	hz.Shutdown(ctx)
}

func task2_opt_lock() {
	ctx := context.TODO()
	// create the client and connect to the cluster on localhost
	client, err := hazelcast.StartNewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	dmp, err := client.GetMap(ctx, "default")
	var key = "1"
	dmp.Put(ctx, key, 1)

	// var value
	for k := 0; k < 1000; k++ {
		for {
			var oldValue, _ = dmp.Get(ctx, key)
			var newValue = oldValue
			time.Sleep(time.Second / 1000)
			newValue = newValue.(int64) + 1
			l, err := dmp.Put(ctx, key, newValue)
			if l.(int64) > -1 || err != nil {
				break
			}

		}
	}
	value, _ := dmp.Get(ctx, key)
	fmt.Println("Finished! Result =", value)
	client.Shutdown(ctx)

}

func task2(id string) {
	m := map[string]interface{}{
		"nl": task2_no_lock,
		"pl": task2_pes_lock,
		"ol": task2_opt_lock,
	}
	var k = m[id]

	// task1("")
	now := time.Now()
	k.(func())()
	fmt.Println((time.Now()).Sub(now))
}

// 3.
// Налаштуйте Bounded queue
// на основі Distributed Queue (http://docs.hazelcast.org/docs/latest/manual/html-single/index.html#queue)
//  налаштуйте Bounded queue (http://docs.hazelcast.org/docs/latest/manual/html-single/index.html#setting-a-bounded-queue)
// з однієї ноди (клієнта) йде запис, а на двох інших читання
// перевірте яка буде поведінка на запис якщо відсутнє читання, і черга заповнена
// як будуть вичитуватись значення з черги якщо є декілька читачів

func t3_read_from_queue(id int, wg *sync.WaitGroup) {
	ctx := context.TODO()
	// create the client and connect to the cluster on localhost
	client, err := hazelcast.StartNewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	dmp, err := client.GetQueue(ctx, "queue")
	// defer wg.Done()
	for {
		item, _ := dmp.Take(ctx)
		time.Sleep(time.Second / 100)
		fmt.Printf("reader %d took item %v\n", id, item)
	}

}

func t3_write_in_queue(wg *sync.WaitGroup) {
	ctx := context.TODO()
	client, err := hazelcast.StartNewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	dmp, err := client.GetQueue(ctx, "queue")
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		dmp.Put(ctx, "value"+fmt.Sprint(i))
		time.Sleep(time.Second / 100)
		fmt.Printf("writer put item %d\n", i)
	}
}

func task3(k string) {
	var wg sync.WaitGroup
	for i := 1; i <= 1; i++ {
		wg.Add(1)
	}
	go t3_write_in_queue(&wg)
	go t3_read_from_queue(1, &wg)
	go t3_read_from_queue(2, &wg)
	wg.Wait()
}
