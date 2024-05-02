<script>
  import { onMount, onDestroy } from "svelte";
  import { Heading } from "flowbite-svelte";
  import InlineCalendar from "$lib/components/InlineCalendar/InlineCalendar.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { P } from "flowbite-svelte";

  let currentTheme = "";
  let observer;

  const theme = {
    dark: {
      calendar: {
        width: "20rem",
        maxWidth: "100vw",
        legend: {
          height: "45px",
        },
        shadow: "0px 10px 26px rgba(0, 0, 0, 0.25)",
        colors: {
          text: {
            primary: "#eee",
            highlight: "#fff",
          },
          background: {
            primary: "#333",
            highlight: "#7F56D9",
            hover: "#222",
          },
          border: "#222",
        },
        font: {
          regular: "1.5em",
          large: "37em",
        },
        grid: {
          disabledOpacity: ".35",
          outsiderOpacity: ".6",
        },
      },
    },
    light: {
      calendar: {
        width: "20rem",
        maxWidth: "100vw",
        legend: {
          height: "45px",
        },
        shadow: "0px 10px 26px rgba(0, 0, 0, 0.25)",
        colors: {
          text: {
            primary: "#333",
            highlight: "#fff",
          },
          background: {
            primary: "#fff",
            highlight: "#7F56D9",
            hover: "#eee",
          },
          border: "#eee",
        },
        font: {
          regular: "1.5em",
          large: "37em",
        },
        grid: {
          disabledOpacity: ".35",
          outsiderOpacity: ".6",
        },
      },
    },
  };

  let schedule = [
    { date: "2024-05-01", time: "10:00 - 12:00", title: "test 1" },
    { date: "2024-05-10", time: "10:00 - 12:00", title: "test 2" },
  ];

  // onMount(() => {
  //   const intervalId = setInterval(() => {
  //     const isDark = document.documentElement.classList.value.includes("dark");
  //     currentTheme = isDark ? "dark" : "light";
  //   }, 100);

  //   return () => {
  //     clearInterval(intervalId);
  //   };
  // });

  const handleClickDate = (date) => {
    console.log(date);
  };

  onMount(() => {
    observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.attributeName === "class") {
          const isDark = document.documentElement.classList.value.includes("dark");
          currentTheme = isDark ? "dark" : "light";
        }
      });
    });

    observer.observe(document.documentElement, { attributes: true });

    currentTheme = document.documentElement.classList.value.includes("dark") ? "dark" : "light";
  });

  onDestroy(() => {
    observer.disconnect();
  });
</script>

<div class="m-4">
  <Heading tag="h3" class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50">
    Calendar
  </Heading>
  <div class="my-8">
    <Divider />
  </div>
  <div class="flex">
    <div class="w-1/3"><P class="font-semibold">Calendar</P></div>
    <div class="w-2/3">
      {#if currentTheme === "light"}
        <InlineCalendar theme={theme.light} {handleClickDate} {schedule} />
        <InlineCalendar theme={theme.light} {handleClickDate} {schedule} />
      {:else if currentTheme === "dark"}
        <InlineCalendar theme={theme.dark} {handleClickDate} {schedule} />
        <InlineCalendar theme={theme.dark} {handleClickDate} {schedule} />
      {/if}
    </div>
  </div>
  <div class="my-8">
    <Divider />
  </div>
  <div class="flex">
    <div class="w-1/3"><P class="font-semibold">Availability</P></div>
    <div class="w-2/3">
      <div class="flex">
        <P>Select day</P>
        &nbsp;
        <div class="flex">Day here</div>
      </div>
      <div class="flex justify-between">
        <div>
          <P>Start time</P>
          <div>Time input here</div>
        </div>
        <div>-</div>
        <div>
          <P>End time</P>
          <div>Time input here</div>
        </div>
      </div>
      <div class="flex">
        <div>Switch input</div>
        &nbsp;
        <P>Set individual hours</P>
      </div>
    </div>
  </div>
</div>
