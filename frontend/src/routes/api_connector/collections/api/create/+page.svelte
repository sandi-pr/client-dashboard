<script>
  import FormApiConnector from "$lib/pages/FormApiConnector.svelte";
  import FormApiParams from "$lib/pages/FormApiParams.svelte";

  import { _ } from "svelte-i18n";

  let stateCreate = $state("formApi");

  let tabDatas = $state({
    selected: "API Information",
  });

  let paramsData = $state([
    {
      id: 1,
      key: "user_id",
      type: "INT",
      description: "Unique identifier for the user",
      actions: {
        objects: [
          {
            onclick: deleteParam(1),
            name: "Delete",
          icon: "trash",
          },
        ],
      },
    },
  ]);

  async function deleteParam(index) {
    // TODO: delete paramsData
  }

</script>

{#if stateCreate === "formApi"}
  <FormApiConnector bind:stateCreate {tabDatas} />
{:else if stateCreate === "formParams"}
  <FormApiParams bind:stateCreate bind:paramsData />
{/if}