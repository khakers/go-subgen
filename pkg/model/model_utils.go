package model

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

const modelSrc = "https://huggingface.co/ggerganov/whisper.cpp/"

// DownloadModel downloads the model to the model directory
// If verifyModelHash is true, the sha1 hash of the downloaded model will be verified against the expected hash for the model
func DownloadModel(model Model, modelDir string, verifyModelHash bool) error {
	modelPath := GetModelPath(modelDir, model)

	log.Printf("downloading model to '%v', this may take some time", modelPath)

	source := modelSrc + "resolve/main/" + "ggml-" + strings.ReplaceAll(model.String(), "_", ".") + ".bin"

	log.Debugf("downloading from %v", source)

	out, err := os.Create(modelPath)
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

	modelHash, err := GetModelSha1Hash(model)
	if err != nil {
		return err
	}

	if verifyModelHash {
		hash, err := verifyHash(modelPath, modelHash)
		if err != nil {
			return err
		}
		if hash == false {
			return errors.New("downloaded model failed hash verification")
		}
	}
	return err
}

// IsModelPresent checks if the model is present in the model directory without verifying the hash
func IsModelPresent(model Model, modelDir string) (bool, error) {
	stat, err := os.Stat(GetModelPath(modelDir, model))
	if err != nil {
		return false, err
	}
	if stat.IsDir() {
		return false, errors.New("model path was a directory")
	}
	return true, nil
}

// VerifyModelHash verifies the sha1 hash of the model file against the expected hash
func VerifyModelHash(model Model, modelDir string) (bool, error) {
	modelHash, err := GetModelSha1Hash(model)
	if err != nil {
		return false, err
	}
	return verifyHash(GetModelPath(modelDir, model), modelHash)
}

// verifyHash verifies the sha1 hash of a file against an expected hash
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
		log.Debugf("expected %v, got %v", expected, result)
		return false, nil
	}
}

// GetModelPath returns the path to the model file based on the given model directory and model type
// e.g. GetModelPath("/models/", Tiny) -> "/models/ggml-tiny.bin"
func GetModelPath(modelDir string, model Model) string {
	return filepath.Join(modelDir, "ggml-"+model.String()+".bin")
}
