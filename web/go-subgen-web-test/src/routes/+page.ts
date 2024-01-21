import type { ProgressBarProps } from "carbon-components-svelte/types/ProgressBar/ProgressBar.svelte";

import { AsrJobStatus, type AsrJob } from "$lib/asr";

export async function load({ fetch, params }) {
  const url = "http://localhost:8095/api/v1/jobs";

  const res = await fetch(url);

  const jobs: AsrJob[] = await res.json();
  console.log(jobs);

  return { jobs: jobs };
}

// export async function loadJobs(fe:type) {

// } loadJobs({ fetch, params }) {

export function _progressBarStatus(status: AsrJobStatus): ProgressBarProps["status"] {
  switch (status) {
    case AsrJobStatus.Queued:
      return "active";
    case AsrJobStatus.AudioStripping:
      return "active";
    case AsrJobStatus.InProgress:
      return "active";
    case AsrJobStatus.Canceled:
      return "error";
    case AsrJobStatus.Complete:
      return "finished";
    case AsrJobStatus.Failed:
      return "error";
    default:
      return "error";
  }
}
