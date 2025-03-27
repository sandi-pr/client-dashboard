<script>
  import { page } from "$app/stores";
  import PageFormMyModels from "$lib/components/PageFormMyModels.svelte";

  let botId = $page.url.searchParams.get("botId");
  let edit = $page.url.searchParams.get("edit");
</script>

<PageFormMyModels {botId} editable={edit} />
