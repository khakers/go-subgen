# go-subgen
A Go adaptation of [McCloudS/subgen](https://github.com/McCloudS/subgen)

Runs a webserver that upon receiving a webhook with a file path will use whisper.cpp to generate subtitles in a .srt file.

Whisper.cpp is relatively CPU and RAM intensive depending on the model you use, and go-subgen also stores the stripped audio in memory instead of saving it to the filesystem. If you have large media files

## Differences from Subgen
* Written in Go (if you care)
* Docker uses static ffmpeg build from mwader
* no temp audio files (kept entirely in memory)
* Queues files
* Doesn't have direct plex webhook integration

## Todo/Future
* Finish this README
* Further integrations/webhooks (at least webhooks from the *arrs)
* Subtitle file templating (you provide a template for what you want your file names to look like based off its variables)
* Persistent Queue
* Translation and more advanced media checking (don't run if file already has subs, for example)

## Endpoints
Currently, go-subgen provides 2 webhook endpoints

### Tautulli
Which is designed around being compatible with subgens Tautullo webhook and accepts the same json payload

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

A very basic post endpoint that accepts a json array of file paths

## Configuration

