<script>
  import { onMount } from "svelte";
  import {
    Typography,
    Divider,
    Card,
    ActionCard,
    CardFeature,
    Button,
    call,
    triggerToast,
  } from "$lib";
  import { goto } from "$app/navigation";
  import { _ } from "svelte-i18n";

  const isFrontend = process.env.NODE_ENV === "development";

  let actionItems = [];

  let recentItems = [
    {
      type: "ledger",
      href: "/llms",
      icon: "openAi",
      title: "My Foundation Model 3",
    },
    {
      type: "ledger",
      href: "/llms",
      icon: "gemini",
      title: "My Foundation Model 2",
    },
    {
      type: "ledger",
      href: "/llms",
      icon: "meta",
      title: "My Foundation Model 1",
    },
  ];

  let cardDataArray = [
    // {
    //   iconData: {
    //     icon: "gemini",
    //   },
    //   contentCenter: [
    //     // Product Name
    //     [
    //       {
    //         text: {
    //           text: "Gemini",
    //           color: "brand-secondary",
    //           preset: "Text xl/Semibold",
    //         },
    //       },
    //     ],
    //     // Title
    //     [
    //       {
    //         text: {
    //           text: "DeepSeek R1",
    //           color: "primary",
    //           preset: "Display md/Semibold",
    //         },
    //       },
    //     ],
    //     // Subtitle
    //     [
    //       {
    //         text: {
    //           text: "Our most used models.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //       },
    //     ],
    //     // Point of Interest
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Finds hidden problems in data.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Understands images and text.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Predicts machine failures.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Personalizes user experiences.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Discovers new scientific insights.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //   ],
    //   actionFooter: [
    //     {
    //       button: {
    //         text: "Try this model",
    //         spaceY: "",
    //         width: "w-full",
    //         rightIcon: "true",
    //         classesSvg: "arrowUpRight",
    //         svgClassReplace: "w-5 h-5 text-white",
    //         funcButton: () => (window.location.href = "/llms"),
    //       },
    //     },
    //     {
    //       button: {
    //         text: "Learn More",
    //         spaceY: "",
    //         width: "w-full",
    //         color: "alternative",
    //         funcButton: () => (window.location.href = "/llms/list?id=1"),
    //       },
    //     },
    //   ],
    // },
    // {
    //   iconData: {
    //     icon: "gemini",
    //   },
    //   contentCenter: [
    //     // Product Name
    //     [
    //       {
    //         text: {
    //           text: "Gemini",
    //           color: "brand-secondary",
    //           preset: "Text xl/Semibold",
    //         },
    //       },
    //     ],
    //     // Title
    //     [
    //       {
    //         text: {
    //           text: "DeepSeek R1",
    //           color: "primary",
    //           preset: "Display md/Semibold",
    //         },
    //       },
    //     ],
    //     // Subtitle
    //     [
    //       {
    //         text: {
    //           text: "Our most used models.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //       },
    //     ],
    //     // Point of Interest
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Finds hidden problems in data.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Understands images and text.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Predicts machine failures.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Personalizes user experiences.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Discovers new scientific insights.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //   ],
    //   actionFooter: [
    //     {
    //       button: {
    //         text: "Try this model",
    //         spaceY: "",
    //         width: "w-full",
    //         rightIcon: "true",
    //         classesSvg: "arrowUpRight",
    //         svgClassReplace: "w-5 h-5 text-white",
    //         funcButton: () => (window.location.href = "/llms"),
    //       },
    //     },
    //     {
    //       button: {
    //         text: "Learn More",
    //         spaceY: "",
    //         width: "w-full",
    //         color: "alternative",
    //         funcButton: () => (window.location.href = "/llms/list?id=2"),
    //       },
    //     },
    //   ],
    // },
    // {
    //   iconData: {
    //     icon: "gemini",
    //   },
    //   contentCenter: [
    //     // Product Name
    //     [
    //       {
    //         text: {
    //           text: "Gemini",
    //           color: "brand-secondary",
    //           preset: "Text xl/Semibold",
    //         },
    //       },
    //     ],
    //     // Title
    //     [
    //       {
    //         text: {
    //           text: "DeepSeek R1",
    //           color: "primary",
    //           preset: "Display md/Semibold",
    //         },
    //       },
    //     ],
    //     // Subtitle
    //     [
    //       {
    //         text: {
    //           text: "Our most used models.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //       },
    //     ],
    //     // Point of Interest
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Finds hidden problems in data.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Understands images and text.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Predicts machine failures.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Personalizes user experiences.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //     [
    //       {
    //         icon: {
    //           icon: "check2",
    //         },
    //         text: {
    //           text: "Discovers new scientific insights.",
    //           color: "tertiary",
    //           preset: "Text md/Regular",
    //         },
    //         justify: "justify-start",
    //       },
    //     ],
    //   ],
    //   actionFooter: [
    //     {
    //       button: {
    //         text: "Try this model",
    //         spaceY: "",
    //         width: "w-full",
    //         rightIcon: "true",
    //         classesSvg: "arrowUpRight",
    //         svgClassReplace: "w-5 h-5 text-white",
    //         funcButton: () => (window.location.href = "/llms"),
    //       },
    //     },
    //     {
    //       button: {
    //         text: "Learn More",
    //         spaceY: "",
    //         width: "w-full",
    //         color: "alternative",
    //         funcButton: () => (window.location.href = "/llms/list?id=3"),
    //       },
    //     },
    //   ],
    // },
  ];

  const findIconFromModelName = (modelName) => {
    let icon = "gemini";
    const lowerCaseModelName = modelName.toLowerCase();

    if (lowerCaseModelName.includes("llama")) {
      icon = "meta";
    } else if (lowerCaseModelName.includes("deepseek")) {
      icon = "deepseek";
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
            
            response.result.models = response.result.models.filter(
              (model) => !model.name.includes("Qurani"),
            );

            let models = response.result.models;

            for (let i = 0; i < models.length; i++) {
              models[i].points = models[i].features.map((point) => {
                  return [
                    {
                      icon: {
                        icon: "check2",
                      },
                      text: {
                        text: point,
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                      justify: "justify-start",
                    },
                  ];
                });
            }

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
                        text: model.name.includes("Qurani") ? "Finetuned Model" : "LLM",
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
                        text: "ourMostUsedModels",
                        color: "tertiary",
                        preset: "Text md/Regular",
                      },
                    },
                  ],

                  //tambahkan isi dari points disini
                  model.points ? model.points.flat() : [],
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
                      funcButton: async () => {
                        let resultPage = await $call.getPage("hexabot", "home-page");
                        window.open(
                          resultPage +
                            "?model=" +
                            model.model +
                            "&id=" +
                            localStorage.getItem("botLlmId"),
                          "_blank",
                        );
                        // goto(resultPage + "?model=" + model.model + "&id=" + localStorage.getItem("botLlmId"));
                      },
                    },
                  },
                  // {
                  //   button: {
                  //     text: "Learn More",
                  //     spaceY: "",
                  //     width: "w-full",
                  //     color: "alternative",
                  //     funcButton: async () => (window.location.href = "/llms/list?id=" + index++),
                  //   },
                  // },
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

  onMount(async () => {
    await getModels();
    actionItems = [
      {
        href: isFrontend
          ? "/llms/my_models/create"
          : await $call.getPage("client-dashboard", "llms-my-models-create"),
        icon: "uploadCloud02",
        actionName: "deploy your model",
        desc: "desc deploy your model",
      },
      {
        href: isFrontend
          ? "/llms/my_models"
          : await $call.getPage("client-dashboard", "llms-my-models"),
        icon: "dataflow04",
        actionName: "view your model",
        desc: "desc view your model",
      },
    ];
  });
</script>

<div class="flex flex-col m-4 gap-4 pb-4 md:pb-12">
  <Typography text={$_("largeLanguageModelsOverview")} preset="Display sm/Semibold" />
  <!-- <Typography text={$_("quickActions")} preset="Text lg/Semibold" propClass="mt-4" />
  <Divider />
  {#if actionItems.length > 0}
    <ActionCard {actionItems} />
  {/if} -->
  <div class="flex justify-between items-center mt-4">
    <!-- <Typography
      text={$_("largeLanguageModelsLists")}
      preset="Text lg/Semibold"
      propClass="md:mt-4"
    /> -->
    <!-- <Button
      spaceY="m-0 max-h-min"
      text="View All"
      funcButton={async () => {
        let link = isFrontend ? "/llms/list" : await $call.getPage("client-dashboard", "llms-list");
        goto(link);
      }}
      color="alternative"
    /> -->
  </div>
  <!-- <Divider /> -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    {#each cardDataArray as cardData}
      <CardFeature {cardData} />
    {/each}
  </div>
  <!-- <div class="flex justify-between mt-4">
    <Typography text="Recently Deployed Models" preset="Text lg/Semibold" propClass="my-4" />
    <Button text="View All" funcButton={() => goto("/llms/my_models")} color="alternative" />
  </div> -->
  <!-- <Divider />
  <div class="grid grid-cols-3 gap-4">
    {#each recentItems as item}
      <Card type={item.type} href={item.href} icon={item.icon} title={item.title} />
    {/each}
  </div> -->
</div>
