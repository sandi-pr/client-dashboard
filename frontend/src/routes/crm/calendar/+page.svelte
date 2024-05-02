<script>
  import { onMount, onDestroy } from "svelte";
  import { format } from "date-fns";
  import { Button, Heading, Input } from "flowbite-svelte";
  import { ClockOutline } from "flowbite-svelte-icons";
  import InlineCalendar from "$lib/components/InlineCalendar/InlineCalendar.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { P } from "flowbite-svelte";

  let currentTheme = "";
  let selectedDate = new Date();
  let observer;

  const theme = {
    dark: {
      calendar: {
        width: "25rem",
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
        width: "25rem",
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

  let newSchedule = {
    date: format(new Date(), "yyyy-MM-dd"),
    time: format(new Date(), "kk:mm") + " - " + format(new Date(), "kk:mm"),
    title: "",
  };

  let dayAvailable = [true, false, true, true, true, false, false];
  let dayAvailableCopy = [...dayAvailable];
  let timeAvailable = ["08:00", "16:00"];
  let timeAvailableCopy = [...timeAvailable];

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
    console.log("date clicked backend logic here");
    selectedDate = date;
  };

  const handleCreateNewSchedule = () => {
    const isValidTimeString = (timeString) => {
      const regex = /^([01]\d|2[0-3]):([0-5]\d)$/;
      return regex.test(timeString);
    };

    if (
      isValidTimeString(newSchedule.time.split(" - ")[0]) === false ||
      isValidTimeString(newSchedule.time.split(" - ")[1]) === false
    ) {
      alert("Invalid time format");
    }

    const timeStart = newSchedule.time.split(" - ")[0];
    const timeEnd = newSchedule.time.split(" - ")[1];
    if (timeStart >= timeEnd) {
      alert("Invalid time range");
      return;
    }

    let timeStartObject = new Date();
    timeStartObject.setHours(timeStart.split(":")[0]);
    timeStartObject.setMinutes(timeStart.split(":")[1]);

    let timeEndObject = new Date();
    timeEndObject.setHours(timeEnd.split(":")[0]);
    timeEndObject.setMinutes(timeEnd.split(":")[1]);

    newSchedule.date = format(selectedDate, "yyyy-MM-dd");
    newSchedule.time = format(timeStartObject, "kk:mm") + " - " + format(timeEndObject, "kk:mm");

    schedule = [...schedule, newSchedule];

    console.log("new schedule backend logic here");

    newSchedule = {
      date: format(new Date(), "yyyy-MM-dd"),
      time: format(new Date(), "kk:mm") + " - " + format(new Date(), "kk:mm"),
      title: "",
    };
  };

  const checkConfigChange = () => {
    // for some reason this is not working
    if (
      JSON.stringify(dayAvailable) !== JSON.stringify(dayAvailableCopy) ||
      JSON.stringify(timeAvailable) !== JSON.stringify(timeAvailableCopy)
    ) {
      return true;
    }
    return false;
  };

  const handleSaveConfig = () => {
    console.log("save config backend logic here");
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
        <!-- <InlineCalendar theme={theme.light} {handleClickDate} {schedule} /> -->
      {:else if currentTheme === "dark"}
        <InlineCalendar theme={theme.dark} {handleClickDate} {schedule} />
        <!-- <InlineCalendar theme={theme.dark} {handleClickDate} {schedule} /> -->
      {/if}
      <div
        class="mt-4 p-4 w-[25rem] h-[12rem] border flex flex-col justify-center items-center rounded-lg border-hc-graylight-200 dark:border-hc-graydark-800"
      >
        <div class="flex">
          <Input
            type="text"
            placeholder={format(new Date(), "kk:mm")}
            on:keyup={(evt) => {
              newSchedule.time = evt.target.value + " - " + format(new Date(), "kk:mm");
            }}
          >
            <ClockOutline slot="right" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
          </Input>
          <P class="mx-4 text-center content-center">-</P>
          <Input
            type="text"
            placeholder={format(new Date(), "kk:mm")}
            on:keyup={(evt) => {
              newSchedule.time = format(new Date(), "kk:mm") + " - " + evt.target.value;
            }}
          >
            <ClockOutline slot="right" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
          </Input>
        </div>
        <Input
          type="text"
          placeholder="Title"
          defaultClass="my-4 w-full"
          bind:value={newSchedule.title}
        />
        <Button on:click={handleCreateNewSchedule} type="button" class="w-[40%]">Create New</Button>
      </div>
    </div>
  </div>
  <div class="my-8">
    <Divider />
  </div>
  <div class="flex">
    <div class="w-1/3"><P class="font-semibold">Availability</P></div>
    <div class="w-2/3">
      <!-- day availability selector -->
      <div class="flex">
        <P class="content-center mr-4">Select day</P>
        <div class="flex gap-3">
          {#each dayAvailable as available, index}
            <button
              type="button"
              class={`size-10 rounded-full text-center content-center font-semibold cursor-pointer ${!available ? "text-disabled-light dark:text-disabled-dark bg-quarterary-light dark:bg-quarterary-dark" : "text-white bg-fg-brand-primary-light dark:bg-fg-brand-primary-dark"}`}
              on:click={() => {
                dayAvailable[index] = !dayAvailable[index];
                checkConfigChange();
              }}
            >
              {format("2024-01-0" + (index + 1), "EEEEE")}
            </button>
          {/each}
        </div>
      </div>

      <!-- time availability selector -->
      <div class="flex mt-8 w-[30rem]">
        <Input
          type="text"
          placeholder={format(new Date(), "kk:mm")}
          bind:value={timeAvailable[0]}
          on:keyup={checkConfigChange}
        >
          <ClockOutline slot="right" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
        </Input>
        <P class="mx-4 text-center content-center">-</P>
        <Input
          type="text"
          placeholder={format(new Date(), "kk:mm")}
          bind:value={timeAvailable[1]}
          on:keyup={checkConfigChange}
        >
          <ClockOutline slot="right" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
        </Input>
      </div>

      <!-- save config -->
      <div class="flex justify-end">
        <!-- <div>Switch input</div>
        &nbsp;
        <P>Set individual hours</P> -->
        {#if checkConfigChange() === true}
          <Button on:click={handleSaveConfig}>Save</Button>
        {:else}
          <Button disabled>Save</Button>
        {/if}
      </div>
    </div>
  </div>
</div>
