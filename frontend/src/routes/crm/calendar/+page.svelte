<script>
  import { onMount, onDestroy } from "svelte";
  import { format } from "date-fns";
  import { Button, Heading, Input, Checkbox } from "flowbite-svelte";
  import { ClockOutline } from "flowbite-svelte-icons";
  import InlineCalendar from "$lib/components/InlineCalendar/InlineCalendar.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { P } from "flowbite-svelte";

  let currentTheme = "";
  let selectedDateFrom = new Date();
  let selectedDateTo = new Date();
  let workHourChecked = false;
  let weekendChecked = false;
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
    {
      title: "test 1",
      dateEvent: [
        { date: "2024-05-01", timeStart: "10:00", timeEnd: "12:00" },
        { date: "2024-05-02", timeStart: "10:00", timeEnd: "12:00" },
      ],
    },
    {
      title: "test 2",
      dateEvent: [
        { date: "2024-05-01", timeStart: "10:00", timeEnd: "12:00" },
        { date: "2024-05-10", timeStart: "10:00", timeEnd: "12:00" },
      ],
    },
  ];

  let newSchedule = {
    title: "",
    dateEvent: [
      {
        date: format(new Date(), "yyyy-MM-dd"),
        timeStart: format(new Date(), "kk:mm"),
        timeEnd: format(new Date(), "kk:mm"),
      },
      {
        date: format(new Date(), "yyyy-MM-dd"),
        timeStart: format(new Date(), "kk:mm"),
        timeEnd: format(new Date(), "kk:mm"),
      },
    ],
  };

  let dayAvailable = [true, false, true, true, true, false, false];
  let dayAvailableCopy = [...dayAvailable];
  let timeAvailable = ["08:00", "16:00"];
  let timeAvailableCopy = [...timeAvailable];

  const handleClickDate = (fromToDate, date) => {
    console.log("date clicked backend logic here");
    if (fromToDate === "from") selectedDateFrom = date;
    else if (fromToDate === "to") selectedDateTo = date;
  };

  const handleCreateNewSchedule = () => {
    // check if the time range is valid
    newSchedule.dateEvent.forEach((date) => {
      if (date.timeStart > date.timeEnd) {
        alert("Invalid time range");
        return;
      }
    });

    // console.log(selectedDateFrom, selectedDateTo);

    // check if the selectedDateFrom is before or equal to selectedDateTo
    if (selectedDateFrom > selectedDateTo) {
      alert("From date must be before or equal to To date");
      return;
    }

    // apply workdate and workhour
    newSchedule.dateEvent.forEach((date, index) => {
      index === 0
        ? (date.date = format(selectedDateFrom, "yyyy-MM-dd"))
        : (date.date = format(selectedDateTo, "yyyy-MM-dd"));
      date.timeStart = newSchedule.dateEvent[0].timeStart;
      date.timeEnd = newSchedule.dateEvent[0].timeEnd;
    });

    // generate date from selectedDateFrom to selectedDateTo
    const dateFrom = new Date(selectedDateFrom);
    const dateTo = new Date(selectedDateTo);
    const dateDiff = Math.abs(dateTo - dateFrom);
    const dateDiffDays = Math.ceil(dateDiff / (1000 * 60 * 60 * 24));

    for (let i = 1; i < dateDiffDays; i++) {
      const newDate = new Date(dateFrom);
      newDate.setDate(dateFrom.getDate() + i);
      newSchedule.dateEvent.push({
        date: format(newDate, "yyyy-MM-dd"),
        timeStart: newSchedule.dateEvent[0].timeStart,
        timeEnd: newSchedule.dateEvent[0].timeEnd,
      });
    }

    // sort the dateEvent array
    newSchedule.dateEvent.sort((a, b) => new Date(a.date) - new Date(b.date));
    console.log(newSchedule);

    // remove weekend from dateEvent if weekendChecked is false
    newSchedule.dateEvent = newSchedule.dateEvent.filter((date) => {
      const day = new Date(date.date).getDay();
      return weekendChecked || (day !== 0 && day !== 6);
    });
    console.log(newSchedule);

    schedule = [...schedule, newSchedule];

    console.log("new schedule backend logic here");

    newSchedule = {
      title: "",
      dateEvent: [
        {
          date: format(new Date(), "yyyy-MM-dd"),
          timeStart: format(new Date(), "kk:mm"),
          timeEnd: format(new Date(), "kk:mm"),
        },
        {
          date: format(new Date(), "yyyy-MM-dd"),
          timeStart: format(new Date(), "kk:mm"),
          timeEnd: format(new Date(), "kk:mm"),
        },
      ],
    };
  };

  const checkTimeFormat = (evt, time) => {
    if (/[^0-9:]/.test(evt.target.value)) {
      // Remove any invalid characters
      evt.target.value = evt.target.value.replace(/[^0-9:]/g, "");
    }
    // check if target value contains a colon
    if (evt.target.value.includes(":")) {
      const targetValueSplit = evt.target.value.split(":");

      // Check the hour part
      if (
        !targetValueSplit[0] ||
        targetValueSplit[0].length > 2 ||
        parseInt(targetValueSplit[0]) > 23
      ) {
        targetValueSplit[0] = targetValueSplit[0] ? targetValueSplit[0].slice(-2) : "01";
      }
      if (parseInt(targetValueSplit[0]) > 23) {
        targetValueSplit[0] = "01";
      }

      // Check the minute part
      if (
        !targetValueSplit[1] ||
        targetValueSplit[1].length > 2 ||
        parseInt(targetValueSplit[1]) > 59
      ) {
        targetValueSplit[1] = targetValueSplit[1] ? targetValueSplit[1].slice(0, 2) : "01";
      }
      if (parseInt(targetValueSplit[1]) > 59) {
        targetValueSplit[1] = "01";
      }

      evt.target.value = targetValueSplit.join(":");

      const isValidTimeFormat =
        targetValueSplit[0].length === 2 && targetValueSplit[1].length === 2;

      const timeVariables = {
        newTimeStart: newSchedule.dateEvent[0].timeStart,
        newTimeEnd: newSchedule.dateEvent[0].timeEnd,
        availTimeStart: timeAvailable[0],
        availTimeEnd: timeAvailable[1],
      };

      if (isValidTimeFormat) {
        switch (time) {
          case "newTimeStart":
            newSchedule.dateEvent[0].timeStart = evt.target.value;
            break;
          case "newTimeEnd":
            newSchedule.dateEvent[0].timeEnd = evt.target.value;
            break;
          case "availTimeStart":
            timeAvailable[0] = evt.target.value;
            break;
          case "availTimeEnd":
            timeAvailable[1] = evt.target.value;
            break;
        }
        console.log(newSchedule, timeAvailable);
      } else {
        timeVariables[time] = timeVariables[time];
      }
    } else {
      const timeVariables = {
        newTimeStart: newSchedule.dateEvent[0].timeStart,
        newTimeEnd: newSchedule.dateEvent[0].timeEnd,
        availTimeStart: timeAvailable[0],
        availTimeEnd: timeAvailable[1],
      };

      timeVariables[time] = timeVariables[time];
    }
  };

  const checkTimeFinal = (evt, time) => {
    // Split the input into hours and minutes
    let [hours, minutes] = evt.target.value.split(":");

    // Add leading zeros if necessary
    hours = hours.padStart(2, "0");
    minutes = minutes.padStart(2, "0");

    // Combine the hours and minutes back into a time string
    evt.target.value = `${hours}:${minutes}`;

    const isValidTimeFormat = hours.length === 2 && minutes.length === 2;

    const timeVariables = {
      newTimeStart: newSchedule.dateEvent[0].timeStart,
      newTimeEnd: newSchedule.dateEvent[0].timeEnd,
      availTimeStart: timeAvailable[0],
      availTimeEnd: timeAvailable[1],
    };

    if (isValidTimeFormat) {
      switch (time) {
        case "newTimeStart":
          newSchedule.dateEvent[0].timeStart = evt.target.value;
          break;
        case "newTimeEnd":
          newSchedule.dateEvent[0].timeEnd = evt.target.value;
          break;
        case "availTimeStart":
          timeAvailable[0] = evt.target.value;
          break;
        case "availTimeEnd":
          timeAvailable[1] = evt.target.value;
          break;
      }
    } else {
      timeVariables[time] = timeVariables[time];
    }
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

  // onMount(() => {
  //   const intervalId = setInterval(() => {
  //     const isDark = document.documentElement.classList.value.includes("dark");
  //     currentTheme = isDark ? "dark" : "light";
  //   }, 100);

  //   return () => {
  //     clearInterval(intervalId);
  //   };
  // });

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
    <div class="w-2/3 flex flex-col">
      <div class="flex gap-4">
        {#if currentTheme === "light"}
          <div class="flex flex-col">
            <P class="w-[25rem] mb-4 text-center text-xl font-semibold">From Date</P>
            <InlineCalendar theme={theme.light} {handleClickDate} {schedule} fromToDate="from" />
          </div>
          <div class="flex flex-col">
            <P class="w-[25rem] mb-4 text-center text-xl font-semibold">To Date</P>
            <InlineCalendar theme={theme.light} {handleClickDate} {schedule} fromToDate="to" />
          </div>
        {:else if currentTheme === "dark"}
          <div class="flex flex-col">
            <P class="w-[25rem] mb-4 text-center text-xl font-semibold">From Date</P>
            <InlineCalendar theme={theme.dark} {handleClickDate} {schedule} fromToDate="from" />
          </div>
          <div class="flex flex-col">
            <P class="w-[25rem] mb-4 text-center text-xl font-semibold">To Date</P>
            <InlineCalendar theme={theme.dark} {handleClickDate} {schedule} fromToDate="to" />
          </div>
        {/if}
      </div>
      <div
        class="mt-4 p-4 w-[51rem] h-[14rem] border flex flex-col justify-center items-center rounded-lg border-hc-graylight-200 dark:border-hc-graydark-800"
      >
        <div class="flex gap-4">
          <div class="flex">
            <Input
              type="text"
              placeholder={format(new Date(), "kk:mm")}
              value={newSchedule.dateEvent[0].timeStart}
              on:keyup={(evt) => {
                checkTimeFormat(evt, "newTimeStart");
              }}
              on:change={(evt) => {
                checkTimeFinal(evt, "availTimeEnd");
              }}
            >
              <ClockOutline slot="right" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
            </Input>
            <P class="mx-4 text-center content-center">-</P>
            <Input
              type="text"
              placeholder={format(new Date(), "kk:mm")}
              value={newSchedule.dateEvent[0].timeEnd}
              on:keyup={(evt) => {
                checkTimeFormat(evt, "newTimeEnd");
              }}
              on:change={(evt) => {
                checkTimeFinal(evt, "availTimeEnd");
              }}
            >
              <ClockOutline slot="right" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
            </Input>
          </div>
          <!-- <div class="flex">
            <Input
              type="text"
              placeholder={format(new Date(), "kk:mm")}
              value={newSchedule.timeStart}
              on:keyup={(evt) => {
                checkTimeFormat(evt, "newTimeStart");
              }}
              on:change={(evt) => {
                checkTimeFinal(evt, "availTimeEnd");
              }}
            >
              <ClockOutline slot="right" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
            </Input>
            <P class="mx-4 text-center content-center">-</P>
            <Input
              type="text"
              placeholder={format(new Date(), "kk:mm")}
              value={newSchedule.timeEnd}
              on:keyup={(evt) => {
                checkTimeFormat(evt, "newTimeEnd");
              }}
              on:change={(evt) => {
                checkTimeFinal(evt, "availTimeEnd");
              }}
            >
              <ClockOutline slot="right" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
            </Input>
          </div> -->
        </div>
        <div class="flex gap-4">
          <Checkbox
            bind:checked={workHourChecked}
            on:click={() => {
              newSchedule.dateEvent[0].timeStart = timeAvailable[0];
              newSchedule.dateEvent[0].timeEnd = timeAvailable[1];
            }}
            class="mt-2">Use Workhour?</Checkbox
          ><Checkbox bind:checked={weekendChecked} class="mt-2">Include Weekend?</Checkbox>
        </div>
        <Input
          type="text"
          placeholder="Title"
          defaultClass="my-4 w-1/2"
          bind:value={newSchedule.title}
        />
        <Button on:click={handleCreateNewSchedule} type="button" class="w-[25%]">Create New</Button>
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
          value={timeAvailable[0]}
          on:keyup={(evt) => {
            checkTimeFormat(evt, "availTimeStart");
            checkConfigChange();
          }}
          on:change={(evt) => {
            checkTimeFinal(evt, "availTimeEnd");
          }}
        >
          <ClockOutline slot="right" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
        </Input>
        <P class="mx-4 text-center content-center">-</P>
        <Input
          type="text"
          placeholder={format(new Date(), "kk:mm")}
          value={timeAvailable[1]}
          on:keyup={(evt) => {
            checkTimeFormat(evt, "availTimeEnd");
            checkConfigChange();
          }}
          on:change={(evt) => {
            checkTimeFinal(evt, "availTimeEnd");
          }}
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
