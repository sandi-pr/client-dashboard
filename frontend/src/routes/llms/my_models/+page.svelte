<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { Search, Button as ButtonFlowbite } from "flowbite-svelte";
  import { goto } from "$app/navigation";

  import { Typography, Icon, call, triggerToast, StateCard, BotCard } from "$lib";

  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";

  const isFrontend = process.env.NODE_ENV === "development";

  let cardDataArray = $state([]);

  let cardDataSelected = $state({});

  let botsList = [
    { botsID: 6969, botsName: "LLaMA 3.2", botDescription: "lorem ipsum dolor sit amet" },
    { botsID: 6970, botsName: "DeepSeek", botDescription: "lorem ipsum dolor sit amet" },
    { botsID: 6971, botsName: "ChatGPT", botDescription: "lorem ipsum dolor sit amet" },
  ];

  let dropdownOpen = $state(new Array(botsList.length).fill(false));

  let token = $state(null);

  let searchTimeout;

  const getModels = async (search) => {
    // console.log("search", search);
    if (token) {
      try {
        let response = await $call.call("hexa-ai-rag", "get-models", {
          URLQuery: {},
          URLParams: {},
          Params: {
            modelName: search || "",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        // console.log("models", response);

        if (response.status === 200) {
          if (response?.result && response.result.models.length > 0) {
            // BRI ADJUSTMENT
            response.result.models = response.result.models.filter(
              (model) => !model.name.includes("Qurani"),
            );

            cardDataArray = response.result.models.map((model) => {
              return {
                botsID: model.id,
                botsName: model.name,
                botDescription: model.description,
                raw: {
                  model,
                },
              };
            });
            // console.log("cardDataArray", cardDataArray);
          }
        } else {
          triggerToast("error");
          console.log(response);
        }
      } catch (error) {
        console.log(error);
      }
    }
  };

  const delayedSearchModel = (value) => {
    // Clear the previous timeout if it exists
    if (searchTimeout) {
      clearTimeout(searchTimeout);
    }

    // Set a new timeout to call getModels after 500ms
    searchTimeout = setTimeout(() => {
      getModels(value);
    }, 500);
  };

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    for (let i = 0; i < dropdownOpen.length; i++) {
      dropdownOpen[i] = false;
    }
  };

  onMount(async () => {
    token = localStorage.getItem("token");

    await getModels();

    // const urlParams = new URLSearchParams(window.location.search);
  });
</script>

<div class="flex flex-col m-4 gap-4 pb-4 md:pb-12">
  <div class="flex flex-col md:flex-row justify-between items-center my-4 gap-4">
    <Typography
      text={cardDataSelected?.raw?.model?.name || "Large Language Models Lists"}
      preset="Display sm/Semibold"
    />
    <div class="flex gap-4">
      <Search
        class="bg-white dark:bg-hc-graydark-950"
        on:input={(evt) => {
          delayedSearchModel(evt.target.value);
        }}
      />
      <ButtonFlowbite
        color="primary"
        outline
        class="flex items-center gap-2 whitespace-nowrap [&>svg]:hover:!text-white"
        onclick={async () => {
          if (isFrontend) {
            goto("/llms/my_models/create");
          } else {
            goto(await $call.getPage("client-dashboard", "llms-my-models-create"));
          }
        }}
      >
        <Icon name="uploadCloud01" color="fg-brand-primary" />
        {$_("deployNewModel")}
      </ButtonFlowbite>
    </div>
  </div>
  {#if cardDataArray.length === 0}
    <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      {#each cardDataArray as bot, index}
        <BotCard
          {bot}
          {index}
          {dropdownOpen}
          config={{
            delete: true,
            detail: true,
            detailView: true,
            detailTry: true,
            isLlm: true,
            isFrontend,
          }}
          {handleDropdownFocusLoss}
        />
      {/each}
    </div>
  {/if}
</div>
