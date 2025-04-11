<script>
  import { page } from "$app/stores";
  import PageFormMyModels from "$lib/components/PageFormMyModels.svelte";

  let botId = $page.url.searchParams.get("botId");
  let modelId = $page.url.searchParams.get("modelId");
  let edit = $page.url.searchParams.get("edit");
</script>

<PageFormMyModels {botId} {modelId} editable={edit} />
