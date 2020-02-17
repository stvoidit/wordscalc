package wordscalc

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"
)

// Ваш HR не смог открыть код из 7zip архива...
// Попросила выслать в ворде или pdf... Вау

// RuneGraph - хэш мапа для подсчета символов
type RuneGraph struct {
	Graph map[rune]uint
	lock  *sync.RWMutex
	err   error
}

// CountInFile - подсчет символов в одном файле
func CountInFile(filename string) (map[rune]uint, error) {
	var overallRuneMap = NewMap()
	ReadParse(filename, &overallRuneMap, nil)
	return overallRuneMap.Graph, overallRuneMap.err
}

// CountInFolderFiles - подсчитать сиволы в файлах из папки
func CountInFolderFiles(foldername string) (map[rune]uint, error) {
	files, err := FilesInFolder("files")
	if err != nil {
		return nil, err
	}
	var overallRuneMap = NewMap()
	var wg = new(sync.WaitGroup)
	wg.Add(len(files))
	for _, filename := range files {
		go ReadParse(filename, &overallRuneMap, wg)
	}
	wg.Wait()
	return overallRuneMap.Graph, overallRuneMap.err
}

// Add - безопасная инкрементация в мапу
// * читать из мапы без мьютексов - безопасно
// * писать - ошибка конкрурентности всегда
func (rg *RuneGraph) Add(r rune) {
	rg.lock.Lock()
	rg.Graph[r]++
	rg.lock.Unlock()
}

// SetErr - вставка ошибки
func (rg *RuneGraph) SetErr(err error) {
	rg.lock.Lock()
	rg.err = err
	rg.lock.Unlock()
}

// NewMap - новый объект для графа символов
func NewMap() RuneGraph {
	return RuneGraph{
		Graph: make(map[rune]uint, 0),
		lock:  new(sync.RWMutex)}
}

// FilesInFolder - сбор списка файлов в папке с путями и именами файлов
func FilesInFolder(folder string) ([]string, error) {
	var files = make([]string, 0)
	var walker = func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return err
	}
	if err := filepath.Walk(folder, walker); err != nil {
		return files, err
	}
	return files, nil
}

// ReadParse - чтение файла и подсчет символов
func ReadParse(f string, om *RuneGraph, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	// если уже кто-то записал ошибку - не тратим время на файл
	if om.err != nil {
		return
	}
	file, err := os.Open(f)
	if err != nil {
		om.SetErr(err)
		return
	}
	defer file.Close()
	var buff bytes.Buffer
	if n, err := buff.ReadFrom(file); err != nil {
		om.err = err
		return
	} else if n == 0 {
		return
	}
	defer buff.Reset()
	for _, word := range buff.String() {
		om.Add(word)
	}
}
