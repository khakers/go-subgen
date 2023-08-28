# go-subgen

Automatic subtitle generation for your media using [whisper.cpp](https://github.com/ggerganov/whisper.cpp).

Runs a webserver that upon receiving a webhook with a file path will use whisper.cpp to generate subtitles for your media and output them as srt
files. Supports webhooks from Radarr & Sonarr.

Whisper.cpp is resource intensive, and go-subgen stores the stripped
audio in memory instead of saving it to the filesystem. If you have large media files and/or use large models you will need lots of ram.

## Features

* Sonarr & Radarr webhook support
* basic Tautulli webhook support
* Subtitle filename templating (see [subfile-name-templating](#subfile-name-templating))
* File queueing

## Todo/Future

Create an issue or Discussion if you have other feature requests or comments.

* file permissions and user handling
* Further software integrations
* Persistent Queue
* Translation and more advanced media checking (don't run if file already has subs, for example)
* A basic web ui to queue transcription and view queued/in progress tasks.
* Filesystem watching support
* Integration into Bazarr as a provider (Not sure how to go about this, please open an issue if you have any ideas)
* CLBLAST & CUBLAS support for GPU acceleration (CUBLAS support is available in the CUBLAS branch)

## Docker

### Docker Compose

#### Pre-built image

If you encounter illegal instruction errors when the model attempts to run, you may need to build the image locally. See [Locally built image](#locally-built-image)
The image currently requires your CPU support avx2, avx, and sse3.

```yaml
version: "3.8"
services:
  subgen:
    image: ghcr.io/khakers/go-subgen:latest
    restart: unless-stopped
    ports:
      - "8095:8095"
    volumes:
      # Path must be identical to the path in your other services
      - /path/to/your/media:/media
      # Models will be downloaded to this directory
      - models:/models
    environment:
      - MODEL_TYPE=base_en
volumes:
  models:
```

#### Locally built image

In many cases you may want or even *need* to build the image locally. Building locally can give you better optimizations
for your specific hardware.
Unfortunately, this means you will need to manually trigger rebuilds (`docker compose build`) and change versions
yourself instead of letting a service like watchtower pull newer images for you.

```yaml
version: "3.8"
services:
  subgen:
    build:
      context: https://github.com/khakers/go-subgen.git#v0.1.0
    restart: unless-stopped
    ports:
      - "8095:8095"
    volumes:
      # Path must be identical to the path in your other services
      - /path/to/your/media:/media
      # Models will be downloaded to this directory
      - models:/models
    environment:
      - MODEL_TYPE=base_en
volumes:
  models:
```

## Web API Endpoints

### Healthcheck

GET

`/healthcheck`

Returns 200 if the server is running

### Tautulli

POST

`/webhooks/tautulli`

Designed around being compatible with subgens Tautulli webhook and accepts the same json payload

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
| MODEL_TYPE                          | Model     | base_en                                        | Whisper.cpp Model version `tiny_en, tiny base_en, base, small_en, small, medium_en, medium, large_v1, large` |
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
| FILE_PERMISSIONS_UID                | uint      | 0                                              | UID of the saved subtitle file                                                                               |
| FILE_PERMISSIONS_GID                | uint      | 0                                              | GID of the saved subtitle file                                                                               |

## Acknowledgements

Inspired by [McCloudS/subgen](https://github.com/McCloudS/subgen) but written in Go and for more service agnostic (not designed around plex) usage.
