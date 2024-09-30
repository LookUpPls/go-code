package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	dbDriver      = "mysql"
	dbSource      = "root:Zknu728201@tcp(rm-uf6bwql5a390jn46nvo.mysql.rds.aliyuncs.com:3306)/word"
	insertQuery   = "INSERT INTO random_num(r) VALUES %s"
	maxGoroutines = 100  // 设置最大并发goroutines数量
	batchSize     = 1000 // 每个批次插入的行数
)

func main() {
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// 删除表里的所有数据
	_, _ = db.Exec("delete from random_num")
	fmt.Println("清空表完成")

	_, err = db.Exec("ALTER TABLE random_num AUTO_INCREMENT = 1")
	if err != nil {
		panic(err)
	}
	fmt.Println("自增设置为0")

	// 创建一个包含连续数字的切片
	numbers := make([]int, 200000)
	for i := range numbers {
		numbers[i] = i + 1
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxGoroutines)

	for i := 0; i < len(numbers); i += batchSize {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(startIndex int, db *sql.DB, wg *sync.WaitGroup) {
			defer wg.Done()
			insertBatch(numbers, startIndex, db)
			<-semaphore
		}(i, db, &wg)
	}

	wg.Wait()
	fmt.Println("Insertion complete.")
}

func insertBatch(numbers []int, startIndex int, db *sql.DB) {
	endIndex := startIndex + batchSize
	if endIndex > len(numbers) {
		endIndex = len(numbers)
	}

	values := make([]string, 0, endIndex-startIndex)
	args := make([]interface{}, 0, endIndex-startIndex)

	for _, num := range numbers[startIndex:endIndex] {
		values = append(values, "(?)")
		args = append(args, num)
	}

	stmt := fmt.Sprintf(insertQuery, strings.Join(values, ","))
	_, err := db.Exec(stmt, args...)
	if err != nil {
		fmt.Println("Batch insert error:", err)
	}
}
