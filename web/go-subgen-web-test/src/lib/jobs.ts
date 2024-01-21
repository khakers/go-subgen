import { AsrJobStatus } from "$lib/asr.js";

export function decode(input: AsrJobStatus): string {
  switch (input) {
    case AsrJobStatus.Queued:
      return "Queued";
    case AsrJobStatus.AudioStripping:
      return "Audio Stripping";
    case AsrJobStatus.InProgress:
      return "In Progress";
    case AsrJobStatus.Canceled:
      return "Cancelled";
    case AsrJobStatus.Complete:
      return "Complete";
    case AsrJobStatus.Failed:
      return "Failed";
    default:
      return "Unknown";
  }
}

export function convertStatusToProgressSteps(input: AsrJobStatus): number {
  switch (input) {
    case AsrJobStatus.Queued:
      return 0;
    case AsrJobStatus.AudioStripping:
      return 1;
    case AsrJobStatus.InProgress:
      return 2;
    case AsrJobStatus.Canceled:
      return -1;
    case AsrJobStatus.Complete:
      return 3;
    case AsrJobStatus.Failed:
      return -1;
    default:
      return -1;
  }
}