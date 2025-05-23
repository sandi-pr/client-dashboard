<script>
  import { browser } from "$app/environment";
  import { Document, Page, preferThisHeight } from "svelte-pdfjs";
  import { Range } from "flowbite-svelte";
  import { ArrowLeftOutline, ArrowRightOutline } from "flowbite-svelte-icons";

  import RefreshCwIcon from "$lib/icon/refresh-cw-01.svelte";
  import RefreshCcwIcon from "$lib/icon/refresh-ccw-01.svelte";
  import { onMount } from "svelte";

  export let pdfFile;
  export let pdfUrl;

  let scale = 1;
  let num = 1;
  let max_pages = 1;
  let renderTextLayer = false;
  let target_height = 500;
  let sizing = 1;
  let rotation = 0;
  let isMobile = false;

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

  function base64toLink(base64) {
    const blob = base64toBlob(base64);
    const fileUrl = URL.createObjectURL(blob);
    return fileUrl;
  }

  function base64toBlob(base64Data) {
    const byteCharacters = atob(base64Data);
    const byteArrays = [];

    for (let offset = 0; offset < byteCharacters.length; offset += 512) {
      const slice = byteCharacters.slice(offset, offset + 512);

      const byteNumbers = new Array(slice.length);
      for (let i = 0; i < slice.length; i++) {
        byteNumbers[i] = slice.charCodeAt(i);
      }

      const byteArray = new Uint8Array(byteNumbers);
      byteArrays.push(byteArray);
    }

    const blob = new Blob(byteArrays, { type: "application/pdf" });
    return blob;
  }

  const updateIsMobile = () => {
    isMobile = window.innerWidth <= 768; // Assume mobile screen is less than or equal to 768px
    if (isMobile) {
      scale = 0.5;
    }
  };

  onMount(() => {
    updateIsMobile();
    window.addEventListener("resize", updateIsMobile);

    return () => {
      window.removeEventListener("resize", updateIsMobile);
    };
  });
</script>

{#if browser}
  <div class="flex flex-col gap-6 justify-center items-center">
    <div class="pdf-container">
      <Document
        file={pdfFile ? base64toLink(pdfFile) : pdfUrl}
        on:loadsuccess={handleLoadSuccess}
        on:loaderror={handleLoadError}
      >
        <Page
          {scale}
          {num}
          {renderTextLayer}
          {rotation}
          getViewport={sizing === 1 ? undefined : preferThisHeight(target_height)}
        />
      </Document>
    </div>

    <div class="w-full md:max-w-max">
      <section class="flex items-center justify-between mb-2">
        <button
          on:click={prevPage}
          class="bg-fg-brand-primary-light dark:bg-fg-brand-primary-dark text-white p-2 rounded-lg shadow-md focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={num === 1}
        >
          <ArrowLeftOutline class="w-6 h-6" />
        </button>
        <span class="font-semibold text-primary-light dark:text-primary-dark fluid-base">
          Page {num} of {max_pages}
        </span>
        <button
          on:click={nextPage}
          class="bg-fg-brand-primary-light dark:bg-fg-brand-primary-dark text-white p-2 rounded-lg shadow-md focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={num === max_pages}
        >
          <ArrowRightOutline class="w-6 h-6" />
        </button>
      </section>

      {#if !isMobile}
        <section class="flex items-center justify-between gap-6">
          <div class="flex items-center gap-4">
            <label for="zoom" class="font-semibold text-primary-light dark:text-primary-dark"
              >Scale:</label
            >
            <Range id="range-steps" min="1" max="1.5" bind:value={scale} step="0.1" />
            <span class="font-mono text-primary-light dark:text-primary-dark"
              >{scale.toFixed(1)}x</span
            >
          </div>
          <div class="flex items-center gap-3">
            <button
              on:click={rotateCounterClockwise}
              class="bg-primary-light dark:bg-gray-800 text-primary-light dark:text-primary-dark p-2 rounded-lg shadow-md focus:outline-none"
            >
              <RefreshCcwIcon
                svgClass="w-6 h-6 text-fg-secondary-light dark:text-fg-secondary-dark"
              />
            </button>
            <button
              on:click={rotateClockwise}
              class="bg-primary-light dark:bg-gray-800 text-primary-light dark:text-primary-dark p-2 rounded-lg shadow-md focus:outline-none"
            >
              <RefreshCwIcon
                svgClass="w-6 h-6 text-fg-secondary-light dark:text-fg-secondary-dark"
              />
            </button>
          </div>
        </section>
      {/if}
    </div>
  </div>
{/if}

<style>
  .pdf-container {
    overflow: auto;
    border: 1px solid #ccc;
    position: relative;
    background-color: #f9f9f9;
  }
</style>
