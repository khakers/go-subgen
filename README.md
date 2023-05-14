# go-subgen

A Go adaptation of [McCloudS/subgen](https://github.com/McCloudS/subgen)

Runs a webserver that upon receiving a webhook with a file path will use whisper.cpp to generate subtitles for your media and output them as an .srt
file.

Whisper.cpp is relatively CPU and RAM intensive depending on the model you use, and go-subgen also stores the stripped
audio in memory instead of saving it to the filesystem. If you have large media files

## Differences from Subgen

* Written in Go (if you care)
* Subtitle filename templating (see [subfile-name-templating](#subfile-name-templating))
* Docker uses static ffmpeg build from mwader
* no temp audio files (extracted audio is kept entirely in memory)
* Queues files
* Doesn't have direct plex webhook integration
* Can process multiple files at once if desired
* Support for Sonarr & Radarr webhooks

## Todo/Future

* Further software integrations
* Persistent Queue 
* Translation and more advanced media checking (don't run if file already has subs, for example)
* A basic web ui to queue transcription and view queued/in progress tasks.
* Support for different models of whisper? (faster-whisper, openai whisper)

## Web API Endpoints

### Healthcheck

GET

`/healthcheck`

Returns 200 if the server is running

### Tautulli

POST

`/webhooks/tautulli`

Which is designed around being compatible with subgens Tautulli webhook and accepts the same json payload

However, go-subgen currently only uses 'file' and ignores the rest of the json data.

```json
{
  "event": "",
  "file": "{file}",
  "filename": "{filename}",
  "mediatype": "{media_type}"
}
```

### Generic

POST

`/webhooks/generic`

A very basic post endpoint that accepts a json array of file paths

```json
{
  "files": [
    "/path/to/file.mp4",
    "/path/to/other/file.mkv"
  ]
}
```

### Radarr

POST

`/webhooks/radarr`

Accepts Radarr formatted webhooks

<details>
<summary>Example webhook configuration</summary>

![img_1.png](assets/radarr_webhook_config.png)

This example only sends the notification on series that have the 'whisper' tag, allowing you to only transcribe series
that need it.
</details>

### Sonarr

POST

`/webhooks/sonarr`

Accepts Sonarr formatted webhooks

<details>
<summary>Example webhook configuration</summary>

![Sonarr webhook connection example image](assets/sonarr_webhook_config.png)

This example only sends the notification on series that have the 'whisper' tag, allowing you to only transcribe series
that need it.

</details>

## Configuration

### Subfile name templating

Go-Subgen allows you to configure how subtitle files are name using Go templates.

The default filename templates is as follows:
`{{.FileName}}.subgen.{{.Lang}}.{{.FileType}}`

You can set your own template by setting the environment variable `SUBTITLE_NAME_TEMPLATE`. The template is created
using the struct below. You can use any of the variables provided and any features of the Go templating system, but keep
in mind no escaping is applied to the result. Additionally, FileHash ***is not*** a SHA hash. It is a hash generated
using imohash
which hashes only portions of the file using murmur3

```go
type SubtitleTemplateData struct {
  FilePath  string
  FileName  string
  Lang      string
  FileHash  string
  FileType  string
  ModelType string
}
```

## Options

| Environment Variable                | Type      | Default                                        | Description                                                                                                  |
|-------------------------------------|-----------|------------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| MODEL_TYPE                          | Model     | base_en                                        | Whisper.cpp Model version `Tiny_en, Tiny Base_en, Base, Small_en, Small, Medium_en, Medium, Large_v1, Large` |
| TARGET_LANG                         | string    | en                                             |                                                                                                              |
| IGNORE_IF_EXISTING                  | bool      | true                                           |                                                                                                              |
| MAX_CONCURRENCY                     | uint      | 1                                              |                                                                                                              |
| MODEL_DIR                           | string    | /models/                                       |                                                                                                              |
| LOG_LEVEL                           | log.Level | info                                           |                                                                                                              |
| VERIFY_MODEL_HASH                   | bool      | true                                           | Verify that the downloaded model mashes the expected hash                                                    |
| PORT                                | uint8     | 8095                                           | Web server port                                                                                              |
| SUBTITLE_NAME_TEMPLATE              | string    | "{{.FileName}}.subgen.{{.Lang}}.{{.FileType}}" | [See Subfile-name-templating](#Subfile-name-templating)                                                      |
| WHISPER_CONF_THREADS                | uint      |                                                | Number of threads to run Whisper.cpp on                                                                      |
| WHISPER_CONF_WHISPER_SPEEDUP        | bool      |                                                |                                                                                                              |
| WHISPER_CONF_TOKEN_THRESHOLD        | float32   |                                                |                                                                                                              |
| WHISPER_CONF_TOKEN_SUM_THRESHOLD    | float32   |                                                |                                                                                                              |
| WHISPER_CONF_MAX_SEGMENT_LENGTH     | uint      |                                                |                                                                                                              |
| WHISPER_CONF_MAX_TOKENS_PER_SEGMENT | uint      |                                                |                                                                                                              |
