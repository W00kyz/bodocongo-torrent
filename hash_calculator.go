package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return nil, err
	}
	return data, nil
}

func sum(filePath string) (int, error) {
	data, err := readFile(filePath)
	if err != nil {
		return 0, err
	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}

	return _sum, nil
}

func worker(filePath string, barrier *sync.WaitGroup, sumMapChan chan map[int][]string) {
	defer barrier.Done()

	_sum, err := sum(filePath)
	if err != nil {
		return
	}

	sumMap := make(map[int][]string)
	sumMap[_sum] = append(sumMap[_sum], filePath)
	sumMapChan <- sumMap
}

func ProcessDirectory(dirPath string) (map[int][]string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var barrier sync.WaitGroup
	sumMapChan := make(chan map[int][]string, len(files))
	resultMap := make(map[int][]string)

	for _, file := range files {
		if !file.IsDir() {
			filePath := dirPath + string(os.PathSeparator) + file.Name()
			barrier.Add(1)
			go worker(filePath, &barrier, sumMapChan)
		}
	}

	go func() {
		barrier.Wait()
		close(sumMapChan)
	}()

	for sumMap := range sumMapChan {
		for hash, paths := range sumMap {
			resultMap[hash] = append(resultMap[hash], paths...)
		}
	}

	return resultMap, nil
}


