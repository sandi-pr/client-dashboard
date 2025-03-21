<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { Search, Button as ButtonFlowbite } from "flowbite-svelte";

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
            cardDataArray = response.result.models.map((model, index) => {
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
                  // model.pointsOfInterest.map((point) => {
                  //   return [
                  //     {
                  //       icon: {
                  //         icon: "check2",
                  //       },
                  //       text: {
                  //         text: point,
                  //         color: "tertiary",
                  //         preset: "Text md/Regular",
                  //       },
                  //       justify: "justify-start",
                  //     },
                  //   ];
                  // }),
                  [
                    {
                      icon: {
                        icon: "check2",
                      },
                      text: {
                        text: "Finds hidden problems in data.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                      justify: "justify-start",
                    },
                  ],
                  [
                    {
                      icon: {
                        icon: "check2",
                      },
                      text: {
                        text: "Understands images and text.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                      justify: "justify-start",
                    },
                  ],
                  [
                    {
                      icon: {
                        icon: "check2",
                      },
                      text: {
                        text: "Predicts machine failures.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                      justify: "justify-start",
                    },
                  ],
                  [
                    {
                      icon: {
                        icon: "check2",
                      },
                      text: {
                        text: "Personalizes user experiences.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                      justify: "justify-start",
                    },
                  ],
                  [
                    {
                      icon: {
                        icon: "check2",
                      },
                      text: {
                        text: "Discovers new scientific insights.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                      justify: "justify-start",
                    },
                  ],
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
            console.log(cardDataArray);
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
        // masih dummy dibawah
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
                  text: "LLaMA 3.2 is an advanced large language model developed by Meta, building upon its predecessors in the LLaMA series. It is designed for high-performance natural language processing (NLP) tasks, including text generation, summarization, code assistance, and conversational AI. The model is optimized for efficiency, scalability, and improved contextual understanding, making it a strong competitor in the AI space.",
                  color: "secondary",
                  preset: "Text md/Regular",
                },
              },
              {
                button: {
                  class: "m-4",
                  color: "none",
                  onClick: () => {},
                  text: {
                    text: "Read More",
                    color: "brand-secondary",
                    preset: "Text md/Semibold",
                  },
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
                list: [
                  {
                    textGroup: [
                      {
                        text: "Improved Efficiency -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Optimized architecture for faster inference and lower computational costs",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "Enhanced Contextual Understanding -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Better ability to maintain context across longer conversations or documents.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "Multilingual Support -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Supports multiple languages, making it versatile for global applications.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "Open-Source Flexibility -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "If released with open access, it allows for greater customization and community contributions.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "Stronger Alignment & Safety -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Improved guardrails to reduce biases and prevent harmful outputs.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                ],
              },
              {
                button: {
                  class: "m-4",
                  color: "none",
                  onClick: () => {},
                  text: {
                    text: "Read More",
                    color: "brand-secondary",
                    preset: "Text md/Semibold",
                  },
                },
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
                list: [
                  {
                    textGroup: [
                      {
                        text: "Chatbots & Virtual Assistants -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Enhances user interactions with more human-like responses.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "Content Creation -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Generates high-quality text for articles, blogs, and creative writing.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "Code Assistance -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Helps developers with coding suggestions, debugging, and explanations.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "Summarization & Data Extraction -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Condenses long documents into concise summaries.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "Education & Tutoring -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Provides personalized learning support for students.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                ],
              },
              {
                button: {
                  class: "m-4",
                  color: "none",
                  onClick: () => {},
                  text: {
                    text: "Read More",
                    color: "brand-secondary",
                    preset: "Text md/Semibold",
                  },
                },
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
                list: [
                  {
                    textGroup: [
                      {
                        text: "Computational Demand -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "May require significant hardware resources for deployment and fine-tuning.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "Potential Ethical Risks -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Despite safety improvements, there may still be risks of misinformation or bias.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "Limited Real-Time Awareness -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Like other LLMs, it may not have real-time or up-to-date knowledge without external integrations.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "Fine-Tuning Complexity -",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Customizing the model for specific use cases can be challenging without expertise.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                ],
              },
              {
                button: {
                  class: "m-4",
                  color: "none",
                  onClick: () => {},
                  text: {
                    text: "Read More",
                    color: "brand-secondary",
                    preset: "Text md/Semibold",
                  },
                },
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
                text: {
                  text: "To run LLaMA 3.2 locally, your system should meet the following specifications:",
                  color: "secondary",
                  preset: "Text md/Regular",
                },
              },
              {
                text: {
                  text: "Hardware Requirements:",
                  color: "secondary",
                  preset: "Text md/Semibold",
                  propClass: "my-4",
                },
              },
              {
                list: [
                  {
                    textGroup: [
                      {
                        text: "GPU:",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "An NVIDIA GPU with CUDA support is recommended, with at least 16GB of VRAM. This is essential for efficient model training and inference.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "RAM:",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "A minimum of 32GB is required, with 64GB recommended for handling larger models and datasets.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "Storage:",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "At least 50GB of free disk space to accommodate the model files and dependencies.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                ],
              },
              {
                text: {
                  text: "Software Requirements:",
                  color: "secondary",
                  preset: "Text md/Semibold",
                  propClass: "my-4",
                },
              },
              {
                list: [
                  {
                    textGroup: [
                      {
                        text: "Operating System:",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Linux is preferred for optimal performance, but macOS and Windows are also supported.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "Python:",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Version 3.8 or higher is necessary to run the required scripts and libraries.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                  {
                    textGroup: [
                      {
                        text: "CUDA Toolkit:",
                        color: "tertiary",
                        preset: "Text md/Semibold",
                      },
                      {
                        text: "Version 11.6 or newer is required for GPU acceleration.",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    ],
                  },
                ],
              },
              {
                button: {
                  class: "m-4",
                  color: "none",
                  onClick: () => {},
                  text: {
                    text: "Read More",
                    color: "brand-secondary",
                    preset: "Text md/Semibold",
                  },
                },
              },
            ],
          },
        ];
        hiddenContent = new Array(detailDataArray.length).fill(true);
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

          {#if content?.list}
            <ul class="flex flex-col gap-2 list-disc mx-4">
              {#each content.list as list}
                <div class="flex items-center">
                  <li class="mr-2" />
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
