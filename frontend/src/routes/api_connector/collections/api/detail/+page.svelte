<script>
  import FormApiConnector from "$lib/pages/FormApiConnector.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import FormApiParams from "$lib/pages/FormApiParams.svelte";

  import { _ } from "svelte-i18n";

  let stateCreate = $state("formApi");

  let paramsData = $state([
    {
      id: 1,
      key: "user_id",
      type: "INT",
      description: "Unique identifier for the user",
    },
  ]);

</script>

{#if stateCreate === "formApi"}
  <FormApiConnector bind:stateCreate bind:paramsData />
{:else if stateCreate === "formParams"}
  <FormApiParams bind:stateCreate bind:paramsData />
{/if}