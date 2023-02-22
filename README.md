# go-subgen

A Go adaptation of [McCloudS/subgen](https://github.com/McCloudS/subgen)

Runs a webserver that upon receiving a webhook with a file path will use whisper.cpp to generate subtitles in a .srt
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
* Can process multiple files at once if desired (files are locked so theoretically even multiple go-subgen instances
  should not conflict)
* Support for Sonarr & Radarr webhooks

## Todo/Future

* Finish this README
* Further software integrations
* Persistent Queue
* Translation and more advanced media checking (don't run if file already has subs, for example)
* A basic web ui to queue transcription and view queued/in progress tasks.

## Endpoints

Currently, go-subgen provides 2 webhook endpoints

### Tautulli

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

`/webhooks/radarr`

Accepts Radarr formatted webhooks

<details>
<summary>Example webhook configuration</summary>

![img_1.png](assets/radarr_webhook_config.png)

This example only sends the notification on series that have the 'whisper' tag, allowing you to only transcribe series
that need it.
</details>

### Sonarr

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