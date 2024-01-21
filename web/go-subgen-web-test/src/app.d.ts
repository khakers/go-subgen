// See https://kit.svelte.dev/docs/types#app

import type { AsrJob } from "$lib/asr";


// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		interface PageData {
			jobs: AsrJob[];
		}
		// interface Platform {}
	}
}

export {

};

