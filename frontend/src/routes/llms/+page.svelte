<script>
  import { Typography, Divider, Card, ActionCard, CardFeature, Button } from "$lib";

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
            text: "View All",
            spaceY: "",
            width: "w-full",
            funcButton: () => (window.location.href = "/llms/list"),
          },
        },
        {
          button: {
            text: "Learn More",
            spaceY: "",
            width: "w-full",
            color: "alternative",
            funcButton: () => (window.location.href = "/llms"),
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
            text: "View All",
            spaceY: "",
            width: "w-full",
            funcButton: () => (window.location.href = "/llms/list"),
          },
        },
        {
          button: {
            text: "Learn More",
            spaceY: "",
            width: "w-full",
            color: "alternative",
            funcButton: () => (window.location.href = "/llms"),
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
            text: "View All",
            spaceY: "",
            width: "w-full",
            funcButton: () => (window.location.href = "/llms/list"),
          },
        },
        {
          button: {
            text: "Learn More",
            spaceY: "",
            width: "w-full",
            color: "alternative",
            funcButton: () => (window.location.href = "/llms"),
          },
        },
      ],
    },
  ];
</script>

<div class="flex flex-col m-4 gap-4">
  <Typography text="Large Language Models Overview" preset="Display sm/Semibold" />
  <Typography text="Quick Actions" preset="Text lg/Semibold" propClass="my-4" />
  <Divider />
  <ActionCard {actionItems} />
  <div class="flex justify-between mt-4">
    <Typography text="Large Language Models Lists" preset="Text lg/Semibold" propClass="my-4" />
    <Button
      text="View All"
      funcButton={() => (window.location.href = "/llms/list")}
      color="alternative"
    />
  </div>
  <Divider />
  <div class="grid grid-cols-3 gap-4">
    {#each cardDataArray as cardData}
      <CardFeature {cardData} />
    {/each}
  </div>
  <!-- <CardFeature {cardData} /> -->
  <div class="flex justify-between mt-4">
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
  </div>
</div>
