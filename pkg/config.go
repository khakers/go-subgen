package pkg

import (
	"path/filepath"

	"github.com/cristalhq/aconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	ModelType        Model         `default:"base_en" json:"model_type,omitempty"`
	TargetLang       string        `default:"en" json:"target_lang,omitempty"`
	IgnoreIfExisting bool          `default:"true" json:"ignore_if_existing,omitempty"`
	WhisperConf      WhisherConfig `json:"whisper_conf"`
	MaxConcurrency   uint          `json:"max_concurrency,omitempty" default:"1"`
	ModelDir         string        `json:"model_dir,omitempty"`
	LogLevel         log.Level     `json:"log_level" default:"debug"`
}

type WhisherConfig struct {
	Threads             uint    `json:"threads,omitempty"`
	WhisperSpeedup      bool    `json:"whisper_speedup,omitempty"`
	TokenThreshold      float32 `json:"token_threshold,omitempty"`
	TokenSumThreshold   float32 `json:"token_sum_threshold,omitempty"`
	MaxSegmentLength    uint    `json:"max_segment_length,omitempty"`
	MaxTokensPerSegment uint    `json:"max_tokens_per_segment,omitempty"`
}

//go:generate go-enum -type=Model -all=false -string=true -new=true -string=true -text=true -json=true -yaml=false

type Model uint8

const (
	tiny_en Model = iota
	tiny
	base_en
	base
	small_en
	small
	medium_en
	medium
	large_v1
	large
)

func LoadConfig() (Config, error) {
	loader := aconfig.LoaderFor(&Cfg, aconfig.Config{
		SkipFlags:        true,
		AllowUnknownEnvs: true,
	})
	err := loader.Load()
	if err != nil {
		return Config{}, err
	}

	return Cfg, nil
}

func GetModelLocation(config Config) string {
	return filepath.Join(config.ModelDir, "ggml-"+config.ModelType.String()+".bin")
}

var Cfg Config
