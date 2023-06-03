package pkg

import (
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/kalafut/imohash"
	"github.com/nightlyone/lockfile"
	log "github.com/sirupsen/logrus"
	"go-subgen/internal"
	"go-subgen/pkg/configuration"
)

type QueuedSub struct {
	filepath string
	filehash [16]byte
}

func EnqueueSub(input string) {
	input, err := filepath.Abs(input)
	log.Printf("Queueing file %v", input)

	// Check to make sure the file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.WithError(err).Error("File does not exist \"" + input + "\"")
		return
	}

	filehash, err := GetFileHash(input)
	if err != nil {
		log.WithError(err).Errorln("failed to generate file hash")
		return
	}

	jobChannel <- QueuedSub{
		filepath: input,
		filehash: filehash,
	}
	log.Println("Job Queued")

	return
}

func GetFileHash(filePath string) (hash [16]byte, err error) {
	// Using imohash because I don't want to spend forever hashing a plausible worst case media file
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

var jobChannel chan QueuedSub

func StartWorkers(config configuration.Config) {
	// makes a channel of QueuedSubs with a capacity of 100.
	jobChannel = make(chan QueuedSub, 100)

	// start the worker
	for i := uint(0); i < config.MaxConcurrency; i++ {
		go queueWorker(jobChannel)
	}
}

func process(sub QueuedSub) {

	log.Infof("Processing job for file %v", sub.filepath)

	filehash, err := GetFileHash(sub.filepath)
	if err != nil {
		log.WithError(err).Println("failed to generate file hash")
		return
	}

	if filehash != sub.filehash {
		log.Warnf("The hash for file \"%v\" has changed since it was queued", sub.filepath)
	}

	// We always want to use the most recent hash of the file
	hashString := hex.EncodeToString(filehash[:])

	lock, err := lockfile.New(filepath.Join(filepath.Dir(sub.filepath), hashString+".lock"))
	if err != nil {
		log.WithError(err).Errorln("failed to acquire file lock")
		return
	}

	log.Debugf("Locking file with hash %v at %v", hashString, filepath.Join(filepath.Dir(sub.filepath), hashString+".lock"))

	err = lock.TryLock()
	if err != nil {
		log.WithError(err).Errorln(err)
		return
	}
	defer lock.Unlock()

	buffer := new(bytes.Buffer)

	start := time.Now()

	err = StripAudioRaw(sub.filepath, buffer, io.Discard)
	if err != nil {
		log.WithError(err).Errorln("Stripping audio failed")
		return
	}
	audioStripDuration := time.Since(start)

	log.Printf("completed audio stripping in %v seconds.", audioStripDuration.Seconds())

	err, subFileName := configuration.Cfg.GetSubtitleFileName(internal.SubtitleTemplateData{
		FilePath:  sub.filepath,
		FileType:  "srt",
		FileName:  internal.GetFileName(sub.filepath),
		Lang:      configuration.Cfg.TargetLang,
		FileHash:  hashString,
		ModelType: string(configuration.Cfg.ModelType),
	})
	if err != nil {
		log.WithError(err).Errorln("failed to template subtitle file name")
	}

	subFilePath := filepath.Join(filepath.Dir(sub.filepath), subFileName)

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

	if configuration.Cfg.FilePermissions.Gid != 0 || configuration.Cfg.FilePermissions.Uid != 0 {
		err = os.Chown(subFilePath, configuration.Cfg.FilePermissions.Uid, configuration.Cfg.FilePermissions.Gid)
		if err != nil {
			log.WithError(err).Errorln("failed to change file ownership")
		}
	}

	start = time.Now()

	err = Generate(configuration.Cfg.GetModelPathFromConfig(), buffer.Bytes(), subFile)
	if err != nil {
		log.WithError(err).Errorln("Generating subtitles failed")
		return
	}

	subDuration := time.Since(start)

	log.Infof("finished generating subtitles for \"%v\" in %v seconds. Sub file saved to \"%v\"", sub.filepath, subDuration.Seconds(), subFilePath)
}
