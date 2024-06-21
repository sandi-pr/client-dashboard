<script>
	import { browser } from '$app/environment';
  import pdfTemplate from '$lib/template.pdf';

	import { Document, Page, preferThisHeight } from 'svelte-pdfjs';

  export let pdfFile = pdfTemplate;

	let scale = 1;
	let num = 1;
	let max_pages = 1;
	let renderTextLayer = false;
	let target_height = 500;

	let sizing = 1;
</script>

{#if browser}
	<Document
		file={pdfFile}
		on:loadsuccess={(e) => {
			max_pages = e.detail.numPages;
			num = Math.min(num, max_pages);
		}}
		on:loaderror={(e) => alert(e.detail + '')}
	>
		<div class="border border-primary-light dark:border-primary-dark">
			<Page
				{scale}
				{num}
				{renderTextLayer}
				getViewport={sizing === 1 ? undefined : preferThisHeight(target_height)}
			/>
		</div>
	</Document>
{/if}

<section class="settings flex gap-2 items-center justify-end mt-4 mb-10">
	Page <input type="number" bind:value={num} step="1" min="1" max={max_pages} /> of {max_pages}
</section>