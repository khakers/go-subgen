<!-- src/routes/about.svelte -->
<script lang="ts">
  import type { PageServerData } from "../$types";

    // import type { PageServerData } from './$types';
    export let data: PageServerData;
</script>

<h1>About</h1>
<p>{data.message}</p>