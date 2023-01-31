package pkg

import (
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/kalafut/imohash"
	"github.com/nightlyone/lockfile"
	log "github.com/sirupsen/logrus"
)

type QueuedSub struct {
	filepath string
	filehash [16]byte
}

func EnqueueSub(input string) {
	input, err := filepath.Abs(input)
	log.Printf("Queueing file %v", input)

	filehash, err := getFileHash(input)
	if err != nil {
		log.WithError(err).Println("failed to generate file hash")
		return
	}

	jobChan <- QueuedSub{
		filepath: input,
		filehash: filehash,
	}
	log.Println("Job Queued")

	return
}

func getFileHash(filePath string) (hash [16]byte, err error) {
	hash, err = imohash.SumFile(filePath)
	if err != nil {
		return hash, err
	}
	return hash, nil
}

func queueWorker(jobChan <-chan QueuedSub) {
	for job := range jobChan {
		process(job)
	}
}

var jobChan chan QueuedSub

func StartWorkers(config Config) {
	// makes a channel of QueuedSubs with a capacity of 100.
	jobChan = make(chan QueuedSub, 100)

	// start the worker
	for i := uint(0); i < config.MaxConcurrency; i++ {
		go queueWorker(jobChan)
	}
}

func process(sub QueuedSub) {

	log.Infof("Processing job for file %v", sub.filepath)

	hashString := hex.EncodeToString(sub.filehash[:])

	lock, err := lockfile.New(filepath.Join(filepath.Dir(sub.filepath), hashString+".lock"))
	if err != nil {
		log.WithError(err).Errorln("failed to acquire file lock")
		return
	}

	log.Printf("Locking file with hash %v at %v", hashString, filepath.Join(filepath.Dir(sub.filepath), hashString+".lock"))

	err = lock.TryLock()
	if err != nil {
		log.WithError(err).Errorln(err)
		return
	}
	defer lock.Unlock()

	buffer := new(bytes.Buffer)

	err = StripAudioRaw(sub.filepath, buffer, io.Discard)
	if err != nil {
		log.WithError(err).Errorln("Stripping audio failed")
		return
	}

	log.Println("completed audio stripping")
	// todo remove ext from filename or use provided one
	subFilePath := filepath.Join(filepath.Dir(sub.filepath), filepath.Base(sub.filepath)+".subgen."+Cfg.TargetLang+".srt")
	log.Printf("created srt file %v", subFilePath)
	subFile, err := os.Create(subFilePath)
	if err != nil {
		return
	}
	defer func(subFile *os.File) {
		err := subFile.Close()
		if err != nil {
			log.WithError(err).Errorln("failed to close file")
			return
		}
	}(subFile)

	err = Generate(GetModelLocation(Cfg), buffer.Bytes(), subFile)
	if err != nil {
		log.WithError(err).Errorln("Generating subtitles failed")
		return
	}
	log.Infof("finished generated subtitles for \"%v\". Sub file is at \"%v\"", sub.filepath, subFilePath)
}
