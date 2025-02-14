<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { Search, Button as ButtonFlowbite } from "flowbite-svelte";

  import { Typography, CardFeature, Button, FeaturedIcon, Icon } from "$lib";
  import { getIcons } from "$lib";

  let idParams = "1";

  let hiddenContent = [true, true, true, true];

  let cardDataArray = [
    {
      iconData: {
        icon: "gemini",
      },
      contentCenter: [
        // Product Name
        [
          {
            text: {
              text: "Gemini",
              color: "brand-secondary",
              preset: "Text xl/Semibold",
            },
          },
        ],
        // Title
        [
          {
            text: {
              text: "DeepSeek R1",
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
            funcButton: () => (idParams = "1"),
          },
        },
      ],
    },
    {
      iconData: {
        icon: "gemini",
      },
      contentCenter: [
        // Product Name
        [
          {
            text: {
              text: "Gemini",
              color: "brand-secondary",
              preset: "Text xl/Semibold",
            },
          },
        ],
        // Title
        [
          {
            text: {
              text: "DeepSeek R1",
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
            funcButton: () => (idParams = "2"),
          },
        },
      ],
    },
    {
      iconData: {
        icon: "gemini",
      },
      contentCenter: [
        // Product Name
        [
          {
            text: {
              text: "Gemini",
              color: "brand-secondary",
              preset: "Text xl/Semibold",
            },
          },
        ],
        // Title
        [
          {
            text: {
              text: "DeepSeek R1",
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
            funcButton: () => (idParams = "3"),
          },
        },
      ],
    },
    {
      iconData: {
        icon: "gemini",
      },
      contentCenter: [
        // Product Name
        [
          {
            text: {
              text: "Gemini",
              color: "brand-secondary",
              preset: "Text xl/Semibold",
            },
          },
        ],
        // Title
        [
          {
            text: {
              text: "DeepSeek R1",
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
            funcButton: () => (idParams = "4"),
          },
        },
      ],
    },
  ];

  // let detailDataArray = [
  //   {
  //     headerText: "Overview",
  //     contentText:
  //       "LLaMA 3.2 is an advanced large language model developed by Meta, building upon its predecessors in the LLaMA series. It is designed for high-performance natural language processing (NLP) tasks, including text generation, summarization, code assistance, and conversational AI. The model is optimized for efficiency, scalability, and improved contextual understanding, making it a strong competitor in the AI space.",
  //   },
  // ];

  onMount(() => {
    const urlParams = new URLSearchParams(window.location.search);
    idParams = idParams ? idParams : urlParams.get("id");
    // then remove the query string
    window.history.replaceState({}, document.title, window.location.pathname);
  });
</script>

<div class="flex flex-col m-4 gap-4">
  {#if idParams}
    <div class="flex justify-between my-4">
      <Typography text="Deepseek R1" preset="Display sm/Semibold" />
      <Button
        text="Back"
        funcButton={() =>
          // window.location.href = "/llms/list"
          (idParams = "")}
      />
    </div>
    <div class="flex items-center gap-4">
      <FeaturedIcon classesSvg={getIcons("gemini")} />
      <div>
        <Typography text="Deepseek R1" preset="Text lg/Semibold" color="secondary" />
        <Typography text="Deepseek R1 by Deepseek" preset="Text sm/Regular" color="tertiary" />
      </div>
    </div>
    <!-- Overview -->
    <ButtonFlowbite
      class="mt-8 justify-between border border-secondary-light dark:border-secondary-dark focus:outline-none focus:ring-0"
      color="none"
      outline
      on:click={() => (hiddenContent[0] = !hiddenContent[0])}
    >
      <Typography text="Overview" preset="Text md/Semibold" color="secondary" />
      {#if hiddenContent[0]}
        <Icon name="chevronUp" />
      {:else}
        <Icon name="chevronDown" />
      {/if}
    </ButtonFlowbite>
    <div class={`mx-4 ${hiddenContent[0] ? "block" : "hidden"}`}>
      <Typography
        text="LLaMA 3.2 is an advanced large language model developed by Meta, building upon its predecessors in the LLaMA series. It is designed for high-performance natural language processing (NLP) tasks, including text generation, summarization, code assistance, and conversational AI. The model is optimized for efficiency, scalability, and improved contextual understanding, making it a strong competitor in the AI space."
        preset="Text md/Regular"
        color="secondary"
      />
      <ButtonFlowbite class="m-4" color="none" on:click={() => {}}>
        <Typography text="Read More" preset="Text md/Semibold" color="brand-secondary" />
      </ButtonFlowbite>
    </div>
    <!-- Advantages -->
    <ButtonFlowbite
      class="mt-8 justify-between border border-secondary-light dark:border-secondary-dark focus:outline-none focus:ring-0"
      color="none"
      outline
      on:click={() => (hiddenContent[0] = !hiddenContent[0])}
    >
      <Typography text="Advantages" preset="Text md/Semibold" color="secondary" />
      {#if hiddenContent[0]}
        <Icon name="chevronUp" />
      {:else}
        <Icon name="chevronDown" />
      {/if}
    </ButtonFlowbite>
    <div class={`mx-4 ${hiddenContent[0] ? "block" : "hidden"}`}>
      <ul class="flex flex-col gap-2 list-disc mx-4">
        <div class="flex items-center">
          <li class="mr-2" />
          <Typography text="Improved Efficiency -" preset="Text md/Semibold" color="tertiary" />
          &nbsp;
          <Typography
            text="Optimized architecture for faster inference and lower computational costs"
            preset="Text md/Regular"
            color="tertiary"
          />
        </div>
        <div class="flex items-center">
          <li class="mr-2" />
          <Typography
            text="Enhanced Contextual Understanding -"
            preset="Text md/Semibold"
            color="tertiary"
          />
          &nbsp;
          <Typography
            text="Better ability to maintain context across longer conversations or documents."
            preset="Text md/Regular"
            color="tertiary"
          />
        </div>
        <div class="flex items-center">
          <li class="mr-2" />
          <Typography text="Multilingual Support -" preset="Text md/Semibold" color="tertiary" />
          &nbsp;
          <Typography
            text="Supports multiple languages, making it versatile for global applications."
            preset="Text md/Regular"
            color="tertiary"
          />
        </div>
        <div class="flex items-center">
          <li class="mr-2" />
          <Typography text="Open-Source Flexibility -" preset="Text md/Semibold" color="tertiary" />
          &nbsp;
          <Typography
            text="If released with open access, it allows for greater customization and community contributions."
            preset="Text md/Regular"
            color="tertiary"
          />
        </div>
        <div class="flex items-center">
          <li class="mr-2" />
          <Typography
            text="Stronger Alignment & Safety -"
            preset="Text md/Semibold"
            color="tertiary"
          />
          &nbsp;
          <Typography
            text="Improved guardrails to reduce biases and prevent harmful outputs."
            preset="Text md/Regular"
            color="tertiary"
          />
        </div>
      </ul>
      <ButtonFlowbite class="m-4" color="none" on:click={() => {}}>
        <Typography text="Read More" preset="Text md/Semibold" color="brand-secondary" />
      </ButtonFlowbite>
    </div>
  {:else}
    <div class="flex justify-between my-4">
      <Typography text="Large Language Models Lists" preset="Display sm/Semibold" />
      <div>
        <Search class="bg-white dark:bg-hc-graydark-950" />
      </div>
    </div>
    <div class="grid grid-cols-3 gap-4">
      {#each cardDataArray as cardData}
        <CardFeature {cardData} />
      {/each}
    </div>
  {/if}
</div>
