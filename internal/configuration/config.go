package configuration

import (
	"bytes"
	"text/template"

	"github.com/cristalhq/aconfig"
	log "github.com/sirupsen/logrus"
	"go-subgen/internal"
	"go-subgen/pkg"
)

type Config struct {
	ModelType            pkg.Model     `default:"base_en" json:"model_type,omitempty"`
	IgnoreIfExisting     bool          `default:"true" json:"ignore_if_existing,omitempty"`
	WhisperConf          WhisperConfig `json:"whisper_conf"`
	MaxConcurrency       uint          `json:"max_concurrency,omitempty" default:"1"`
	ModelDir             string        `json:"model_dir,omitempty" default:"/models/"`
	LogLevel             log.Level     `json:"log_level" default:"info"`
	VerifyModelHash      bool          `json:"verify_model_hash" default:"true"`
	ServerConfig         ServerConfig  `json:"server_config"`
	SubtitleNameTemplate string        `json:"subtitle_name_template" default:"{{.FileName}}.subgen.{{.Lang}}.{{.FileType}}"`
	FilePermissions      FileOwner     `json:"permissions"`
}

type ServerConfig struct {
	Port    uint16 `json:"port" default:"8095"`
	Address string `json:"address" default:"0.0.0.0"`
}

type WhisperConfig struct {
	Threads             uint    `json:"threads,omitempty"`
	WhisperSpeedup      bool    `json:"whisper_speedup,omitempty"`
	TokenThreshold      float32 `json:"token_threshold,omitempty"`
	TokenSumThreshold   float32 `json:"token_sum_threshold,omitempty"`
	MaxSegmentLength    uint    `json:"max_segment_length,omitempty"`
	MaxTokensPerSegment uint    `json:"max_tokens_per_segment,omitempty"`
	TargetLang          string  `default:"en" json:"target_lang,omitempty"`
}

type FileOwner struct {
	Uid int `json:"uid,omitempty" default:"0"`
	Gid int `json:"gid,omitempty" default:"0"`
}

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
	return GetModelPathFromConfig(cfg, cfg.ModelType)
}

func GetModelPathFromConfig(config Config, model pkg.Model) string {
	return pkg.GetModelPath(config.ModelDir, model)
}

func (cfg Config) GetSubtitleFileName(data internal.SubtitleTemplateData) (err error, str string) {
	buf := new(bytes.Buffer)
	err = SubFileTemplate.Execute(buf, data)
	return err, buf.String()
}
