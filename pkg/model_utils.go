package pkg

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"go-subgen/pkg/configuration"
)

const modelSrc = "https://huggingface.co/ggerganov/whisper.cpp/"

func DownloadModel(model configuration.Model) error {
	log.Printf("downloading model to %v", configuration.GetModelPath(configuration.Cfg, model))

	source := modelSrc + "resolve/main/" + "ggml-" + strings.ReplaceAll(model.String(), "_", ".") + ".bin"

	log.Debugf("downloading from %v", source)

	out, err := os.Create(configuration.GetModelPath(configuration.Cfg, model))
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(source)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	log.Println("model downloaded")

	if configuration.Cfg.VerifyModelHash {
		hash, err := verifyHash(configuration.GetModelPath(configuration.Cfg, model), GetModelShaHash(model))
		if err != nil {
			return err
		}
		if hash == false {
			return errors.New("hash mismatch")
		}
	}
	return err
}

func IsModelDownloaded(model configuration.Model) (bool, error) {
	stat, err := os.Stat(configuration.GetModelPath(configuration.Cfg, model))
	if err != nil {
		return false, err
	}
	if stat.IsDir() {
		return false, errors.New("model path was a directory")
	}

	if configuration.Cfg.VerifyModelHash {
		hash, err := verifyHash(configuration.GetModelPath(configuration.Cfg, model), GetModelShaHash(configuration.Cfg.ModelType))
		if err != nil {
			return false, err

		}
		if hash {
			return true, nil
		} else {
			return false, errors.New("model hash verification failed")

		}
	}
	return true, nil
}

func verifyHash(path string, expected string) (matches bool, err error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, err
	}

	result := hex.EncodeToString(h.Sum(nil))
	if expected == result {
		return true, nil
	} else {
		log.Debugf("expected %v, got %v", GetModelShaHash(configuration.Cfg.ModelType), result)
		return false, nil
	}
}
