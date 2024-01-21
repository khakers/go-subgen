export enum AsrJobStatus {
  Queued,
  AudioStripping,
  InProgress,
  Canceled,
  Complete,
  Failed,
}

export type AsrJob = {
  ID: number;
  FilePath: string;
  Lang: string;
  CreatedAt: string;
  UpdatedAt: string;
  Status: AsrJobStatus;
  Progress: number;
  DurationMS: number;
};

export type AsrProgressEvent = {
  JobId: number;
  Progress: number;
  Time?: string;
};

export enum AsrJobChangeType {
  New,
  Update,
  Delete,
}

export type AsrJobEvent = {
  JobId: number;
  ChangeType: AsrJobChangeType;
  Job: AsrJob;
};

export type AsrStatusEvent = {
  JobId: number;
  Status: number;
  Time: string;
};
