package configuration

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/cristalhq/aconfig"
	log "github.com/sirupsen/logrus"
	"go-subgen/internal"
)

type Config struct {
	ModelType            Model         `default:"base_en" json:"model_type,omitempty"`
	TargetLang           string        `default:"en" json:"target_lang,omitempty"`
	IgnoreIfExisting     bool          `default:"true" json:"ignore_if_existing,omitempty"`
	WhisperConf          WhisperConfig `json:"whisper_conf"`
	MaxConcurrency       uint          `json:"max_concurrency,omitempty" default:"1"`
	ModelDir             string        `json:"model_dir,omitempty" default:"/models/"`
	LogLevel             log.Level     `json:"log_level" default:"info"`
	VerifyModelHash      bool          `json:"verify_model_hash" default:"true"`
	Port                 uint16        `json:"port" default:"8095"`
	SubtitleNameTemplate string        `json:"subtitle_name_template" default:"{{.FileName}}.subgen.{{.Lang}}.{{.FileType}}"`
	FilePermissions      FileOwner     `json:"permissions"`
}

type WhisperConfig struct {
	Threads             uint    `json:"threads,omitempty"`
	WhisperSpeedup      bool    `json:"whisper_speedup,omitempty"`
	TokenThreshold      float32 `json:"token_threshold,omitempty"`
	TokenSumThreshold   float32 `json:"token_sum_threshold,omitempty"`
	MaxSegmentLength    uint    `json:"max_segment_length,omitempty"`
	MaxTokensPerSegment uint    `json:"max_tokens_per_segment,omitempty"`
}

type FileOwner struct {
	Uid int `json:"uid,omitempty" default:"0"`
	Gid int `json:"gid,omitempty" default:"0"`
}

//go:generate go-enum -type=Model -all=false -string=true -new=true -string=true -text=true -json=true -yaml=false

type Model uint8

const (
	Tiny_en Model = iota
	Tiny
	Base_en
	Base
	Small_en
	Small
	Medium_en
	Medium
	Large_v1
	Large
)

var Cfg Config

var SubFileTemplate template.Template

func LoadConfig() (Config, error) {
	loader := aconfig.LoaderFor(&Cfg, aconfig.Config{
		SkipFlags:        true,
		AllowUnknownEnvs: true,
	})
	err := loader.Load()
	if err != nil {
		return Config{}, err
	}

	tmpl, err := template.New("subfile").Parse(Cfg.SubtitleNameTemplate)

	SubFileTemplate = *tmpl

	return Cfg, nil
}

func (cfg Config) GetModelPathFromConfig() string {
	return GetModelPath(cfg, cfg.ModelType)
}

func GetModelPath(config Config, model Model) string {
	return filepath.Join(config.ModelDir, "ggml-"+model.String()+".bin")
}

func (cfg Config) GetSubtitleFileName(data internal.SubtitleTemplateData) (err error, str string) {
	buf := new(bytes.Buffer)
	err = SubFileTemplate.Execute(buf, data)
	return err, buf.String()
}
