package whisper_cpp_generator

import (
	"bytes"
	"context"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"

	"github.com/kalafut/imohash"
	"github.com/nightlyone/lockfile"
	log "github.com/sirupsen/logrus"
	"go-subgen/internal"
	"go-subgen/internal/asr_job"
	"go-subgen/internal/configuration"
	"go-subgen/pkg"
	"go-subgen/pkg/model"
)

type QueuedSub struct {
	filePath string
	fileHash [16]byte
	AsrJobID uint
}

type SubtitleGenerator struct {
	jobChannel       chan QueuedSub
	conf             configuration.Config
	asrJobRepository asr_job.AsrJobRepository
}

func NewSubtitleGenerator(conf configuration.Config, asrJobRepository asr_job.AsrJobRepository) *SubtitleGenerator {
	return &SubtitleGenerator{
		// makes a channel of QueuedSubs with a capacity of 100.
		jobChannel:       make(chan QueuedSub, 100),
		conf:             conf,
		asrJobRepository: asrJobRepository,
	}
}

// todo handle cancellation

func (s SubtitleGenerator) EnqueueSub(job asr_job.FileAsrJob) {
	input, err := filepath.Abs(job.FilePath)
	log.Printf("Queueing file %v", input)

	// Check to make sure the file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.WithError(err).Error("File does not exist \"" + input + "\"")
		err := s.asrJobRepository.SetJobStatus(context.TODO(), job.ID, asr_job.Failed)
		if err != nil {
			return
		}
		return
	}

	filehash, err := GetFileHash(input)
	if err != nil {
		log.WithError(err).Errorln("failed to generate file hash")
		s.asrJobRepository.SetJobStatus(context.TODO(), job.ID, asr_job.Failed)
		return
	}

	s.jobChannel <- QueuedSub{
		filePath: input,
		fileHash: filehash,
		AsrJobID: job.ID,
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

func (s SubtitleGenerator) queueWorker(jobChan <-chan QueuedSub) {
	for job := range jobChan {
		s.process(job, s.conf)
	}
}

// todo do we really need 2 channels for this?
func (s SubtitleGenerator) newJobWorker(newJobChannel <-chan uint) {
	for jobID := range newJobChannel {
		job, err := s.asrJobRepository.GetJob(context.TODO(), jobID)
		if err != nil {
			log.WithError(err).Errorln("failed to retrieve job when attempting to enqueue it to the subtitle generator")
		}
		s.EnqueueSub(job)
	}
}

func (s SubtitleGenerator) StartWorkers() {
	// start the workers
	go s.newJobWorker(s.asrJobRepository.GetNewJobChannel())
	for i := uint(0); i < s.conf.MaxConcurrency; i++ {
		go s.queueWorker(s.jobChannel)
	}
}

func (s SubtitleGenerator) process(sub QueuedSub, conf configuration.Config) {

	log.Infof("Processing job for file %v, job id %v", sub.filePath, sub.AsrJobID)

	filehash, err := GetFileHash(sub.filePath)
	if err != nil {
		log.WithError(err).Println("failed to generate file hash")
		s.asrJobRepository.SetJobStatus(context.TODO(), sub.AsrJobID, asr_job.Failed)
		return
	}

	if filehash != sub.fileHash {
		log.Warnf("The hash for file \"%v\" has changed since it was queued", sub.filePath)
	}

	// We always want to use the most recent hash of the file
	hashString := hex.EncodeToString(filehash[:])

	// todo refactor locking
	lock, err := lockfile.New(filepath.Join(filepath.Dir(sub.filePath), hashString+".lock"))
	if err != nil {
		log.WithError(err).Errorln("failed to acquire file lock")
		s.asrJobRepository.SetJobStatus(context.TODO(), sub.AsrJobID, asr_job.Failed)
		return
	}
	log.Debugf("Locking file with hash %v at %v", hashString, filepath.Join(filepath.Dir(sub.filePath), hashString+".lock"))

	err = lock.TryLock()
	if err != nil {
		log.WithError(err).Errorln(err)
		s.asrJobRepository.SetJobStatus(context.TODO(), sub.AsrJobID, asr_job.Failed)
		return
	}
	defer lock.Unlock()

	s.asrJobRepository.SetJobStatus(context.TODO(), sub.AsrJobID, asr_job.InProgress)

	buffer := new(bytes.Buffer)

	start := time.Now()

	logger := log.New()
	logwriter := logger.WriterLevel(log.InfoLevel)

	err = pkg.StripAudioRaw(sub.filePath, buffer, logwriter)
	if err != nil {
		log.WithError(err).Errorln("Stripping audio failed")
		s.asrJobRepository.SetJobStatus(context.TODO(), sub.AsrJobID, asr_job.Failed)
		return
	}
	err = logwriter.Close()
	if err != nil {
		return
	}
	audioStripDuration := time.Since(start)

	log.Printf("completed audio stripping in %v seconds.", audioStripDuration.Seconds())

	err, subFileName := conf.GetSubtitleFileName(configuration.SubtitleTemplateData{
		FilePath:  sub.filePath,
		FileType:  "srt",
		FileName:  internal.GetFileName(sub.filePath),
		Lang:      conf.WhisperConf.TargetLang,
		FileHash:  hashString,
		ModelType: string(conf.ModelType),
	})
	if err != nil {
		log.WithError(err).Errorln("failed to template subtitle file name")
		s.asrJobRepository.SetJobStatus(context.TODO(), sub.AsrJobID, asr_job.Failed)
		return
	}

	subFilePath := filepath.Join(filepath.Dir(sub.filePath), subFileName)

	log.Printf("created srt file %v", subFilePath)

	subFile, err := os.Create(subFilePath)
	if err != nil {
		return
	}
	defer func(subFile *os.File) {
		err := subFile.Close()
		if err != nil {
			log.WithError(err).Errorln("failed to close file")
			s.asrJobRepository.SetJobStatus(context.TODO(), sub.AsrJobID, asr_job.Failed)
			return
		}
	}(subFile)

	if conf.FilePermissions.Gid != 0 || conf.FilePermissions.Uid != 0 {
		err = os.Chown(subFilePath, conf.FilePermissions.Uid, conf.FilePermissions.Gid)
		if err != nil {
			log.WithError(err).Errorln("failed to change file ownership")
		}
	}

	start = time.Now()

	progressChannel := make(chan float32, 2)

	go progressChannelWorker(progressChannel, s.asrJobRepository, sub.AsrJobID)

	err = Generate(model.GetModelPath(conf.ModelDir, conf.ModelType), buffer.Bytes(), subFile, progressChannel, context.TODO())
	if err != nil {
		log.WithError(err).Errorln("Generating subtitles failed")
		s.asrJobRepository.SetJobStatus(context.TODO(), sub.AsrJobID, asr_job.Failed)
		return
	}

	subDuration := time.Since(start)
	err = s.asrJobRepository.SetJobStatus(context.TODO(), sub.AsrJobID, asr_job.Complete)
	if err != nil {
		log.WithError(err).Errorln("failed to set job status on job ", sub.AsrJobID)
	}

	log.Infof("finished generating subtitles for \"%v\" in %v seconds. Sub file saved to \"%v\"", sub.filePath, subDuration.Seconds(), subFilePath)
}

// progressChannelWorker is a function that updates the progress of a job in the given ASR job repository.
// It listens to the provided channel for progress updates and calls the SetJobProgress function on the repository for each update.
// The function takes a channel of float32, an implementation of the AsrJobRepository interface, and a job ID as parameters.
// The channel is used to receive the progress updates, the repository is used to update the job progress, and the ID is used to identify the job.
// Each progress update is passed to the SetJobProgress function on the repository with the provided job ID.
// If an error occurs during the update, an error message is logged.
// The function continues to listen to the progress channel until it is closed.
// This function should be called as a goroutine, so that it can run concurrently with other tasks.
func progressChannelWorker(a chan float32, repository asr_job.AsrJobRepository, id uint) {
	for i := range a {
		err := repository.SetJobProgress(context.TODO(), id, i)
		if err != nil {
			log.WithError(err).Errorf("failed to update job progress on job %v", id)
		}
	}
}
