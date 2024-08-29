<script>
	import { browser } from "$app/environment";
	import { Document, Page, preferThisHeight } from "svelte-pdfjs";
	import {
		ArrowLeftOutline,
		ArrowRightOutline,
		ZoomOutSolid,
		ZoomInSolid,
		ArrowLeftToBracketOutline,
		ArrowRightToBracketOutline,
	} from "flowbite-svelte-icons";

	export let pdfFile;

	let scale = 1;
	let num = 1;
	let max_pages = 1;
	let renderTextLayer = false;
	let target_height = 500;
	let sizing = 1;
	let rotation = 0;

	const handleLoadSuccess = (e) => {
		max_pages = e.detail.numPages;
		num = Math.min(num, max_pages);
	};

	const handleLoadError = (e) => {
		alert(`Error loading PDF: ${e.detail}`);
	};

	const nextPage = () => {
		if (num < max_pages) {
			num += 1;
		}
	};

	const prevPage = () => {
		if (num > 1) {
			num -= 1;
		}
	};

	const rotateClockwise = () => {
		rotation = (rotation + 90) % 360;
	};

	const rotateCounterClockwise = () => {
		rotation = (rotation - 90 + 360) % 360;
	};
</script>

{#if browser}
	<section class="flex items-center justify-between mb-6">
		<button
			on:click={prevPage}
			class="bg-blue-600 text-white p-2 rounded-lg shadow-md focus:outline-none hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
			disabled={num === 1}
		>
			<ArrowLeftOutline class="w-6 h-6" />
		</button>
		<span class="font-semibold text-xl">
			Page {num} of {max_pages}
		</span>
		<button
			on:click={nextPage}
			class="bg-blue-600 text-white p-2 rounded-lg shadow-md focus:outline-none hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
			disabled={num === max_pages}
		>
			<ArrowRightOutline class="w-6 h-6" />
		</button>
	</section>

	<!-- <section class="flex items-center justify-between mb-6">
		<div class="flex items-center gap-4">
			<label for="zoom" class="font-semibold">Zoom:</label>
			<input
				id="zoom"
				type="range"
				min="0.2"
				max="2"
				step="0.1"
				bind:value={scale}
				class="w-60 cursor-pointer"
			/>
			<span class="font-mono text-lg">{scale.toFixed(1)}x</span>
		</div>
		<div class="flex items-center gap-4">
			<button
				on:click={rotateCounterClockwise}
				class="bg-gray-600 text-white p-2 rounded-lg shadow-md focus:outline-none hover:bg-gray-700"
			>
				<ArrowLeftToBracketOutline class="w-6 h-6" />
			</button>
			<button
				on:click={rotateClockwise}
				class="bg-gray-600 text-white p-2 rounded-lg shadow-md focus:outline-none hover:bg-gray-700"
			>
				<ArrowRightToBracketOutline class="w-6 h-6" />
			</button>
		</div>
	</section> -->

	<Document file={pdfFile} on:loadsuccess={handleLoadSuccess} on:loaderror={handleLoadError}>
		<div class="border border-gray-300 dark:border-gray-700 rounded-lg shadow-lg overflow-hidden">
			<Page
				{scale}
				{num}
				{renderTextLayer}
				{rotation}
				getViewport={sizing === 1 ? undefined : preferThisHeight(target_height)}
			/>
		</div>
	</Document>
{/if}
