<script>
  import { P } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import Markdown from "svelte-exmarkdown";

  export let datas;

  // Process the data to replace each newline with two spaces and a newline
  $: processedDatas = datas.map((data) => ({
    ...data,
    result: data.result.replace(/(\n+)/g, "  \n"),
  }));
</script>

<div class="max-h-[46rem] overflow-auto hide-scrollbar">
  <div class={`p-10 bg-tertiary-light rounded-lg`}>
    <div class="bg-white w-full h-full rounded-sm">
      <div class="grid grid-cols-12 p-12">
        {#each processedDatas as data, index}
          <div class={`col-span-1 ${!index ? "pb-4" : "py-4"}`}>{index + 1}.</div>
          <P class={`col-span-3 ${!index ? "pb-4" : "py-4"}`}>{data.name}</P>
          <div class={`col-span-8 ${!index ? "pb-4" : "py-4"}`}>
            <Markdown md={data.result} />
          </div>
          <div class="col-span-12 border-b"></div>
        {/each}
      </div>
    </div>
  </div>
</div>

<style>
  .hide-scrollbar::-webkit-scrollbar {
    display: none;
  }

  .hide-scrollbar {
    -ms-overflow-style: none; /* IE and Edge */
    scrollbar-width: none; /* Firefox */
  }
</style>
