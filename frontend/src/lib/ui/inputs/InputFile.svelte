<script lang="ts">
  import uploadIcon from '$lib/icon/upload-icon.svg';
  import pdfIcon from '$lib/icon/pdf-icon.svg';

  let { value, placeholder, maxSize } : {
    value: File | undefined;
    placeholder: string | undefined;
    maxSize?: string;
  } = $props();
</script>

<div class="flex items-center justify-center w-full">
  <label for="dropzone-file" class="flex flex-col items-center justify-center w-full h-fit border-2 border-hc-graylight-300 border-dashed rounded-lg cursor-pointer bg-white dark:hover:bg-bray-800 dark:bg-gray-700 hover:bg-gray-100 dark:border-gray-600 dark:hover:border-gray-500 dark:hover:bg-gray-600">
    <div class="flex flex-col items-center justify-center pt-5 pb-6">
      <img src={uploadIcon} class="w-10 h-10 ml-3" alt="uploadIcon" />
      <p class="mb-2 text-sm text-hc-graylight-600 dark:text-gray-400">
        <span class="font-semibold text-hc-brand-700">Click to upload</span>
        or drag and drop
      </p>
      <p class="text-xs text-gray-500 dark:text-gray-400">{placeholder} File ({maxSize})</p>
    </div>
    <input
      id="dropzone-file"
      type="file"
      class="hidden"
      on:change={(event) => {
        value = event.target.files[0];
      }} />
  </label>
</div>

{#if value !== undefined}
  <div class="flex flex-col gap-2">
    <div class="flex justify-between items-center bg-white px-3 py-3 rounded-lg border-2 border-hc-graylight-200">
      <div class="flex justify-start items-center gap-3">
        <img src={pdfIcon} class="w-10 h-10" alt="PDF Icon" />
        <div class="flex flex-col justify-start items-start truncate">
          <p class="text-hc-graylight-600 text-sm font-semibold dark:text-gray-400">
            {value.name}
          </p>
          <p class="text-hc-graylight-600 dark:text-gray-400 text-sm">{value.size}</p>
        </div>
      </div>
      <div class="flex justify-center items-center">
        <!-- <CheckCircleOutline class="w-5 h-5 ml-3 text-hc-brand-700" />
        <TrashBinOutline
          class="w-5 h-5 ml-3 text-hc-error-700 cursor-pointer"
          on:click={() => {
            value = null;
          }} /> -->
      </div>
    </div>
  </div>
{/if}
