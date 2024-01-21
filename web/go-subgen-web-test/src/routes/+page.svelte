<script lang="ts">
  import { decode } from "$lib/jobs";
  import {
    Column,
    DataTable,
    Grid,
    ProgressBar,
    Row,
  } from "carbon-components-svelte";
  import { _progressBarStatus } from "./+page.js";

  import { AsrJobChangeType, type AsrJob, type AsrJobEvent, type AsrProgressEvent } from "$lib/asr";

  import '@formatjs/intl-durationformat/polyfill';

  export let data;

  // const url = new URL("/events?stream=progress", {base});
  const eventSource = new EventSource(
    "http://localhost:8095/api/events?stream=progress"
  );

  eventSource.addEventListener("progressUpdate", (event) => {
    const eventData: AsrProgressEvent = JSON.parse(event.data);
    console.log("new progress event");
    console.log(eventData);
    // update the data of the job with the given ID

    console.log(data);

    for (var i in data.jobs) {
      if (data.jobs[i].ID === eventData.JobId) {
        data.jobs[i].Progress = eventData.Progress;
      }
    }

    // update the data
    data.jobs = data.jobs;
    console.log(data.jobs);
  });

  eventSource.addEventListener("jobUpdate", (event) => {
    const eventData: AsrJobEvent = JSON.parse(event.data);
    console.log("new jobUpdate event");
    console.log(eventData);

    switch (eventData.ChangeType) {
      // Create
      case AsrJobChangeType.New:
        console.log("new job");
        data.jobs.push(eventData.Job);
        break;
      // Update
      case AsrJobChangeType.Update:
        console.log("update job");
        data.jobs = data.jobs.filter(
          (job: AsrJob) => job.ID !== eventData.JobId
        );
        break;
      // Delete
      case AsrJobChangeType.Delete:
        console.log("delete job");
        for (var i in data.jobs) {
          if (data.jobs[i].ID === eventData.JobId) {
            data.jobs[i] = eventData.Job;
          }
        }
        data.jobs = data.jobs.filter(
          (job: AsrJob) => job.ID !== eventData.JobId
        );
        break;
      default:
        console.error("unknown jobUpdate event type");
        break;
    }

    // update the data
    console.log(data.jobs);
  });

  eventSource.addEventListener("statusChange", (event) => {
    const eventData = JSON.parse(event.data);
    console.log("new status event");
    console.log(eventData);

    // update the data of the job with the given ID
    for (var i in data.jobs) {
      if (data.jobs[i].ID === eventData.JobId) {
        data.jobs[i].Status = eventData.Status;
      }
    }

    // update the data
    data.jobs = data.jobs;
    console.log(data.jobs);
  });

  eventSource.onmessage = (event) => {
    const eventData = JSON.parse(event.data);
    console.log(eventData);
    // update the data of the job with the given ID

    console.log(data);
    console.log(data.jobs);
  };
  // @ts-ignore
  let durationFormatter  = new Intl.DurationFormat("en", { style: "digital" });
</script>

<DataTable
  title="ASR Jobs"
  description="A list of all ASR jobs"
  rows={Array.from(data.jobs, (_, i) => ({
    id: data.jobs[i].ID,
    ID: data.jobs[i].ID,
    FilePath: data.jobs[i].FilePath,
    Lang: data.jobs[i].Lang,
    CreatedAt: data.jobs[i].CreatedAt,
    UpdatedAt: data.jobs[i].UpdatedAt,
    Status: decode(data.jobs[i].Status),
    Progress: data.jobs[i].Progress,
    Duration: durationFormatter.format({milliseconds: Math.round(data.jobs[i].DurationMS)}),
  }))}
  headers={[
    {
      key: "ID",
      value: "ID",
    },
    {
      key: "FilePath",
      value: "File Path",
    },
    {
      key: "Lang",
      value: "Language",
    },
    {
      key: "CreatedAt",
      value: "Created At",
    },
    {
      key: "UpdatedAt",
      value: "Updated At",
    },
    {
      key: "Status",
      value: "Status",
    },
    {
      key: "Progress",
      value: "Progress",
    },
    {
      key: "Duration",
      value: "Duration",
    },
  ]}
/>

<hr />

<Grid>
  <Row padding>
    <Column>ID</Column>
    <Column>Status</Column>
    <Column>File Path</Column>
    <Column>Progress</Column>
  </Row>
  {#each data.jobs as job}
    <Row>
      <Column>{job.ID}</Column>
      <Column>
        {decode(job.Status)}
        <!-- <ProgressIndicator preventChangeOnClick currentIndex={convertStatusToProgressSteps(job.Status)}>
            <ProgressStep
              label="Queued"
            />
            <ProgressStep
              label="Audio Stripping"
            />
            <ProgressStep
              label="Generating Subtitles"
            />
          </ProgressIndicator> -->
      </Column>
      <Column>{job.FilePath}</Column>
      <Column>
        <ProgressBar
          labelText="ASR progress"
          hideLabel={false}
          value={job.Status != 0 ? job.Progress : undefined}
          status={_progressBarStatus(job.Status)}
          helperText="{job.Progress}% complete"
        />
      </Column>
    </Row>
  {/each}
</Grid>
