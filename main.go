package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"pipeline/buffer"
	"pipeline/pipeline"
	"strconv"
	"strings"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sourceChan := make(chan int64)
	go dataSource(
		ctx,
		cancel,
		sourceChan,
	)

	consumer(
		ctx,
		cancel,
		pipeline.New(
			negativeFilterStage,
			mult3FilterStage,
			bufferStage,
		).Run(ctx, sourceChan),
	)

	os.Exit(0)
}

// генерация данных
func dataSource(ctx context.Context, ctxCancel context.CancelFunc, ch chan<- int64) {
	scanner := bufio.NewScanner(os.Stdin)
	var data string

	for {
		select {
		case <-ctx.Done():
			return
		default:
			scanner.Scan()
			data = scanner.Text()
			if strings.EqualFold(data, "quit") {
				fmt.Println("Программа завершила работу!")
				ctxCancel()
				return
			}

			i, err := strconv.ParseInt(data, 10, 64)
			if err != nil {
				fmt.Println("Программа обрабатывает только целые числа!")
				continue
			}

			ch <- i
		}
	}
}

// стадия, фильтрующая отрицательные числа
func negativeFilterStage(ctx context.Context, inp <-chan int64) (out chan int64) {
	out = make(chan int64)

	go func() {
		for {
			select {
			case data := <-inp:
				if data >= 0 {
					select {
					case out <- data:
					case <-ctx.Done():
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// стадия, фильтрующая числа, не кратные 3
func mult3FilterStage(ctx context.Context, inp <-chan int64) (out chan int64) {
	out = make(chan int64)
	go func() {
		for {
			select {
			case data := <-inp:
				if data != 0 && data%3 == 0 {
					select {
					case out <- data:
					case <-ctx.Done():
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// стадия буферизации
func bufferStage(ctx context.Context, inp <-chan int64) (out chan int64) {
	out = make(chan int64)
	buf := buffer.New(100)

	go func() {
		for {
			select {
			case data := <-inp:
				buf.Push(data)
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(30 * time.Second):
				bufferData := buf.Get()
				if bufferData != nil {
					for _, data := range bufferData {
						select {
						case out <- data:
						case <-ctx.Done():
							return
						}
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// Потребитель данных от пайплайна
func consumer(ctx context.Context, ctxCancel context.CancelFunc, inp <-chan int64) {
	for {
		select {
		case data := <-inp:
			fmt.Printf("Обработаны данные: %d\n", data)
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Hour): // пусть консумер тоже имеет возможность завершать пайплайн если не получает данные долгое время
			ctxCancel()
			return
		}
	}
}
