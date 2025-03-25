<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { Search, Button as ButtonFlowbite } from "flowbite-svelte";
  import models from "$lib/utils/models.json";
  import ExMarkdown from "svelte-exmarkdown";
  import { gfmPlugin } from "svelte-exmarkdown/gfm";
  const plugins = [gfmPlugin()];

  import {
    Typography,
    CardFeature,
    Button,
    FeaturedIcon,
    Icon,
    call,
    triggerToast,
    StateCard,
  } from "$lib";
  import { getIcons } from "$lib";

  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";

  let idModel = "";

  let cardDataArray = [];

  let cardDataSelected = {};

  let detailDataArray = [];

  let hiddenContent = [];

  const findIconFromModelName = (modelName) => {
    let icon = "gemini";
    const lowerCaseModelName = modelName.toLowerCase();

    if (lowerCaseModelName.includes("llama")) {
      icon = "meta";
    } else if (lowerCaseModelName.includes("deepseek")) {
      icon = "gemini";
    } else if (lowerCaseModelName.includes("chatgpt")) {
      icon = "openai";
    }

    return icon;
  };

  const getModels = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        let response = await $call.call("hexa-ai-rag", "get-models", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        // console.log("models", response);
        if (response.status === 200) {
          if (response?.result && response.result.models.length > 0) {
            cardDataArray = models.map((model, index) => {
              return {
                id: index++,
                iconData: {
                  icon: findIconFromModelName(model.name),
                },
                contentCenter: [
                  // Product Name
                  [
                    {
                      text: {
                        text: "LLM",
                        color: "brand-secondary",
                        preset: "Text xl/Semibold",
                      },
                    },
                  ],
                  // Title
                  [
                    {
                      text: {
                        text: model.name,
                        color: "primary",
                        preset: "Display md/Semibold",
                      },
                    },
                  ],
                  // Subtitle
                  [
                    {
                      text: {
                        text: "Our most used models.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    },
                  ],
                  // Point of Interest
                  model.features.map((feature) => {
                    return {
                      icon: {
                        icon: "check2",
                      },
                      text: {
                        text: feature,
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                      justify: "justify-start",
                    };
                  }),
                ],
                actionFooter: [
                  {
                    button: {
                      text: "Try this model",
                      spaceY: "",
                      width: "w-full",
                      rightIcon: "true",
                      classesSvg: "arrowUpRight",
                      svgClassReplace: "w-5 h-5 text-white",
                      funcButton: () => (window.location.href = "/llms"),
                    },
                  },
                  {
                    button: {
                      text: "Learn More",
                      spaceY: "",
                      width: "w-full",
                      color: "alternative",
                      funcButton: async () => {
                        idModel = index;
                        await getDetailData();
                      },
                    },
                  },
                ],
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

  const getDetailData = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        cardDataSelected = cardDataArray[idModel - 1];
        detailDataArray = [
          {
            headerText: {
              text: "Overview",
              color: "secondary",
              preset: "Text md/Semibold",
            },
            content: [
              {
                text: {
                  text: cardDataSelected.raw.model.overview,
                  color: "secondary",
                  preset: "Text md/Regular",
                },
              },
            ],
          },
          {
            headerText: {
              text: "Advantages",
              color: "secondary",
              preset: "Text md/Semibold",
            },
            content: [
              {
                list: cardDataSelected.raw.model.advantages.map((advantage) => ({
                  textGroup: [
                    {
                      text: advantage.split(" – ")[0] + " –",
                      color: "tertiary",
                      preset: "Text md/Semibold",
                    },
                    {
                      text: advantage.split(" – ")[1],
                      color: "tertiary",
                      preset: "Text md/Regular",
                    },
                  ],
                })),
              },
            ],
          },
          {
            headerText: {
              text: "Best Use Case",
              color: "secondary",
              preset: "Text md/Semibold",
            },
            content: [
              {
                list: cardDataSelected.raw.model.base_use_cases.map((useCase) => ({
                  textGroup: [
                    {
                      text: useCase.split(" – ")[0] + " –",
                      color: "tertiary",
                      preset: "Text md/Semibold",
                    },
                    {
                      text: useCase.split(" – ")[1],
                      color: "tertiary",
                      preset: "Text md/Regular",
                    },
                  ],
                })),
              },
            ],
          },
          {
            headerText: {
              text: "Disadvantages",
              color: "secondary",
              preset: "Text md/Semibold",
            },
            content: [
              {
                list: cardDataSelected.raw.model.disadvantages.map((disadvantage) => ({
                  textGroup: [
                    {
                      text: disadvantage.split(" – ")[0] + " –",
                      color: "tertiary",
                      preset: "Text md/Semibold",
                    },
                    {
                      text: disadvantage.split(" – ")[1],
                      color: "tertiary",
                      preset: "Text md/Regular",
                    },
                  ],
                })),
              },
            ],
          },
          {
            headerText: {
              text: "Specifications",
              color: "secondary",
              preset: "Text md/Semibold",
            },
            content: [
              {
                textMD: {
                  text: cardDataSelected.raw.model.requirements,
                },
              },
            ],
          },
        ];
        hiddenContent = new Array(detailDataArray.length).fill(true);
        console.log("detailDataArray", detailDataArray);
      } catch (error) {
        console.log(error);
      }
    }
  };

  onMount(async () => {
    await getModels();

    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.has("id")) {
      idModel = idModel ? idModel : urlParams.get("id");
      await getDetailData();
      // then remove the query string
      window.history.replaceState({}, document.title, window.location.pathname);
    }
  });
</script>

<div class="flex flex-col m-4 gap-4 pb-4 md:pb-12">
  <div class="flex justify-between items-center my-4">
    <Typography
      text={cardDataSelected?.raw?.model?.name || "Large Language Models Lists"}
      preset="Display sm/Semibold"
    />
    {#if idModel}
      <Button
        text="Back"
        funcButton={() => {
          cardDataSelected = {};
          idModel = "";
        }}
      />
    {:else}
      <div>
        <Search class="bg-white dark:bg-hc-graydark-950" />
      </div>
    {/if}
  </div>
  {#if idModel}
    <div class="flex items-center gap-4">
      <FeaturedIcon classesSvg={getIcons(cardDataSelected?.iconData?.icon || "gemini")} />
      <div>
        <Typography
          text={cardDataSelected?.raw?.model?.name}
          preset="Text lg/Semibold"
          color="secondary"
        />
        <Typography
          text={cardDataSelected?.raw?.model?.description}
          preset="Text sm/Regular"
          color="tertiary"
        />
      </div>
    </div>
    {#each detailDataArray as detailData, index}
      <ButtonFlowbite
        class="mt-8 justify-between border border-secondary-light dark:border-secondary-dark focus:outline-none focus:ring-0"
        color="none"
        outline
        on:click={() => (hiddenContent[index] = !hiddenContent[index])}
      >
        {#if detailData?.headerText}
          <Typography {...detailData.headerText} />
        {:else}
          <div></div>
        {/if}
        {#if hiddenContent[index]}
          <Icon name="chevronUp" />
        {:else}
          <Icon name="chevronDown" />
        {/if}
      </ButtonFlowbite>
      <div class={`mx-4 ${hiddenContent[index] ? "block" : "hidden"}`}>
        {#each detailData.content as content}
          {#if content?.text}
            <Typography {...content.text} />
          {/if}

          {#if content?.textMD && content.textMD.text}
            <div class="prose dark:prose-dark">
              <ExMarkdown md={content.textMD.text} {plugins} />
            </div>
          {/if}

          {#if content?.list}
            <ul class="flex flex-col gap-2 mx-4">
              {#each content.list as list}
                <div class="flex items-center">
                  <li class="mr-2 before:content-['•'] before:text-black dark:before:text-white" />
                  {#each list.textGroup as textGroup, index}
                    <Typography {...textGroup} />
                    {#if list.textGroup.length > 1 && index < list.textGroup.length - 1}
                      <span>&nbsp;</span>
                    {/if}
                  {/each}
                </div>
              {/each}
            </ul>
          {/if}

          {#if content?.button}
            <ButtonFlowbite {...content.button}>
              <Typography {...content.button.text} />
            </ButtonFlowbite>
          {/if}
        {/each}
      </div>
    {/each}
  {:else if cardDataArray.length === 0}
    <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      {#each cardDataArray as cardData}
        <CardFeature {cardData} />
      {/each}
    </div>
  {/if}
</div>
