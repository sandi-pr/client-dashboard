<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { Search } from "flowbite-svelte";

  import { Typography, CardFeature, Button } from "$lib";

  let idParams = "";

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

  onMount(() => {
    const urlParams = new URLSearchParams(window.location.search);
    idParams = urlParams.get("id");
    // then remove the query string
    window.history.replaceState({}, document.title, window.location.pathname);
  });
</script>

<div class="flex flex-col m-4 gap-4">
  {#if idParams}
    <div class="flex justify-between my-4">
      <Typography text="Large Language Models Lists" preset="Display sm/Semibold" />
      <Button
        text="Back"
        funcButton={() =>
          // window.location.href = "/llms/list"
          (idParams = "")}
      />
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
