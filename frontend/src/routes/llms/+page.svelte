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

  let actionItems = [
    {
      // href: await $call.getPage("client-dashboard", "bots-create"),
      href: "/llms",
      icon: "uploadCloud02",
      actionName: "deploy your model",
      desc: "desc deploy your model",
    },
  ];

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
            let llama31 = ["Expanded Model Sizes", "Multilingual Processing", "Extended Context Length", "Open-Source Accessibility", "Enhanced Performance"];
            let llama32 = ["Multimodal Capabilities", "Cost Efficiency", "High Performance", "Edge Deployment Suitability", "Broad Flexibility"];
            let gemma = ["Real-Time Multilingual Support", "High Contextual Understanding", "Efficient Processing", "Scalability", "User-Friendly Integration"];
            let deepseek = ["Modular Design", "Cost-Effective Reasoning", "Efficient Resource Utilization", "Open-Source Accessibility", "Competitive Performance"];
            let qurani = ["Finds hidden problems in data.", "Understands images and text.", "Predicts machine failures.", "Personalizes user experiences.", "Discovers new scientific insights."];
           
            let models = response.result.models;

            for (let i = 0; i < models.length; i++) {
              if (models[i].name === "Llama3.1") {
                models[i].points = llama31.map((point) => {
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
                  ]
                });
              } else if (models[i].name === "Llama3.2") {
                models[i].points = llama32.map((point) => {
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
                  ]
                });
              } else if (models[i].name === "DeepSeek R1") {
                models[i].points = deepseek.map((point) => {
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
                  ]
                });
              } else if (models[i].name === "Gemma 2") {
                models[i].points = gemma.map((point) => {
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
                  ]
                });
              } else if (models[i].name.includes("Qurani")) {
                models[i].points = qurani.map((point) => {
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
                  ]
                });
              }
            }
            console.log("masukkkkkkkkk", models)

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
                        text: "Our most used models.",
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
                        window.open(resultPage + "?model=" + model.model + "&id=" + localStorage.getItem("botLlmId"), "_blank");
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
  });
</script>

<div class="flex flex-col m-4 gap-4 pb-4 md:pb-12">
  <!-- <Typography text="Large Language Models Overview" preset="Display sm/Semibold" />
  <Typography text="Quick Actions" preset="Text lg/Semibold" propClass="my-4" />
  <Divider />
  <ActionCard {actionItems} /> -->
  <div class="flex justify-between items-center mt-4">
    <Typography text="Large Language Models Lists" preset="Text lg/Semibold" propClass="md:my-4" />
    <Button
      spaceY="m-0 max-h-min"
      text="View All"
      funcButton={ async() => {
        let link = await $call.getPage("client-dashboard", "llms-list");
        goto(link);
      }}
      color="alternative"
    />
  </div>
  <!-- <Divider /> -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    {#each cardDataArray as cardData}
      <CardFeature {cardData} />
    {/each}
  </div>
  <!-- <CardFeature {cardData} /> -->
  <!-- <div class="flex justify-between mt-4">
    <Typography text="Recently Deployed Models" preset="Text lg/Semibold" propClass="my-4" />
    <Button
      text="View All"
      funcButton={() => (window.location.href = "/llms/list")}
      color="alternative"
    />
  </div>
  <Divider />
  <div class="grid grid-cols-3 gap-4">
    {#each recentItems as item}
      <Card type={item.type} href={item.href} icon={item.icon} title={item.title} />
    {/each}
  </div> -->
</div>
