<script>
  import { onMount, onDestroy } from "svelte";
  import { format, startOfWeek, endOfWeek } from "date-fns";
  import { P, Button, Heading, Input, Checkbox } from "flowbite-svelte";
  import { ClockOutline, ChevronLeftOutline, ChevronRightOutline } from "flowbite-svelte-icons";
  import { TabBar } from "$lib";
  import InlineCalendar from "$lib/components/InlineCalendar/InlineCalendar.svelte";
  import Divider from "$lib/components/Divider.svelte";

  let currentTheme = "";
  let currentDate = new Date(); // Current date
  let currentCalendarMode = "week"; // Current calendar mode
  let selectedDateFrom = new Date();
  let selectedDateTo = new Date();
  let workHourChecked = false;
  let weekendChecked = false;
  let observer;

  let tabDatas = {
    selected: "calendar",
    data: [
      {
        name: "calendar",
      },
      {
        name: "availability",
      },
    ],
  };

  let schedule = [
    {
      title: "test 1",
      dateEvent: [
        { date: "2024-05-02", timeStart: "10:00", timeEnd: "12:00" },
        { date: "2024-05-10", timeStart: "10:00", timeEnd: "13:00" },
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
  let daysInMonth = getDaysInMonth(currentDate.getMonth(), currentDate.getFullYear());
  let daysInWeek = getDaysOfWeek(new Date(currentDate));

  const changeTabbar = (name) => {
    tabDatas.selected = name;
  };

  function getDaysInMonth(month, year) {
    // Create a new date object for the first day of the specified month and year
    let firstDay = new Date(year, month, 1);
    let days = [];

    // Add days from the previous month
    for (let i = firstDay.getDay(); i > 0; i--) {
      let prevDate = new Date(year, month, -i + 1);
      days.push({
        date: prevDate.getDate(),
        dayName: prevDate.toLocaleDateString("en-US", { weekday: "short" }),
        dateObject: prevDate,
        isCurrentMonth: false,
        dateEvent: getEventsForDate(prevDate),
      });
    }

    // Add days from the current month
    for (let i = 1; i <= new Date(year, month + 1, 0).getDate(); i++) {
      let currentDate = new Date(year, month, i);
      days.push({
        date: i,
        dayName: currentDate.toLocaleDateString("en-US", { weekday: "short" }),
        dateObject: currentDate,
        isCurrentMonth: true,
        dateEvent: getEventsForDate(currentDate),
      });
    }

    // Add days from the next month to fill the grid
    let nextMonthDays = 7 * 5 - days.length;
    for (let i = 1; i <= nextMonthDays; i++) {
      let nextDate = new Date(year, month + 1, i);
      days.push({
        date: nextDate.getDate(),
        dayName: nextDate.toLocaleDateString("en-US", { weekday: "short" }),
        dateObject: nextDate,
        isCurrentMonth: false,
        dateEvent: getEventsForDate(nextDate),
      });
    }

    return days;
  }

  function getDaysOfWeek(startDay) {
    let daysOfWeek = [];
    let day = new Date(startDay);

    // Set the day to the nearest Sunday before the specified startDay
    day.setDate(day.getDate() - day.getDay());

    for (let i = 0; i < 7; i++) {
      let dayOfWeek = day.toLocaleDateString("en-US", { weekday: "short" });
      let dateOfWeek = day.getDate();
      let isCurrentDate = day.toDateString() === currentDate.toDateString(); // Check if it's the current date

      let eventsForDay = schedule.flatMap((event) =>
        event.dateEvent
          .filter((dateEvent) => new Date(dateEvent.date).getDate() === dateOfWeek)
          .map((dateEvent) => ({
            title: event.title,
            timeStart: dateEvent.timeStart,
            timeEnd: dateEvent.timeEnd,
          }))
      );

      daysOfWeek.push({
        dayName: dayOfWeek,
        dateNumber: dateOfWeek,
        isCurrentDate: isCurrentDate,
        events: eventsForDay,
      });

      day.setDate(day.getDate() + 1);
    }

    return daysOfWeek;
  }

  function getHoursOfDay() {
    let hoursOfDay = [];

    for (let i = 1; i <= 24; i++) {
      hoursOfDay.push(i);
    }

    return hoursOfDay;
  }

  function handleClickHour(date, hour) {
    console.log("Clicked date:", date);
    console.log("Clicked hour:", hour);
    // Add your logic here to handle the click event
  }

  function getEventsForDate(date) {
    let formattedDate = date.toISOString().split("T")[0]; // Get YYYY-MM-DD format
    let events = [];

    schedule.forEach((item) => {
      item.dateEvent.forEach((event) => {
        if (event.date === formattedDate) {
          events.push({
            title: item.title,
            timeStart: event.timeStart,
            timeEnd: event.timeEnd,
          });
        }
      });
    });

    return events;
  }

  const handleClickDate = (date) => {
    console.log(date);
    console.log("date clicked backend logic here");
  };

  const handleChangeDate = (direction) => {
    console.log(currentCalendarMode);
    console.log(direction);
    if (currentCalendarMode === "month") {
      if (direction === "previous") {
        currentDate = new Date(currentDate.getFullYear(), currentDate.getMonth() - 1);
      } else if (direction === "next") {
        currentDate = new Date(currentDate.getFullYear(), currentDate.getMonth() + 1);
      }

      // Update the daysInMonth array based on the new currentDate
      daysInMonth = getDaysInMonth(currentDate.getMonth(), currentDate.getFullYear());
    } else if (currentCalendarMode === "week") {
      if (direction === "previous") {
        currentDate = new Date(
          currentDate.getFullYear(),
          currentDate.getMonth(),
          currentDate.getDate() - 7
        );
      } else if (direction === "next") {
        currentDate = new Date(
          currentDate.getFullYear(),
          currentDate.getMonth(),
          currentDate.getDate() + 7
        );
      }

      // Update the daysInWeek array based on the new currentDate
      daysInWeek = getDaysOfWeek(new Date(currentDate));
    }
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
  <TabBar {tabDatas} {changeTabbar} />
  <div class="my-8">
    <Divider />
  </div>
  {#if tabDatas.selected === "calendar"}
    <div class="flex flex-col">
      <div class="flex gap-4">
        <div class="flex flex-col w-full">
          <!-- calendar mode selector -->
          <div class="flex mb-4 justify-end gap-2">
            <Button
              color="light"
              class={currentCalendarMode === "day"
                ? "bg-brand-solid-hover-light dark:bg-brand-solid-hover-light text-primary-dark hover:bg-brand-primary-dark dark:hover:bg-brand-primary-dark"
                : ""}
              on:click={() => (currentCalendarMode = "day")}>Day</Button
            >
            <Button
              color="light"
              class={currentCalendarMode === "week"
                ? "bg-brand-solid-hover-light dark:bg-brand-solid-hover-light text-primary-dark hover:bg-brand-primary-dark dark:hover:bg-brand-primary-dark"
                : ""}
              on:click={() => (currentCalendarMode = "week")}>Week</Button
            >
            <Button
              color="light"
              class={currentCalendarMode === "month"
                ? "bg-brand-solid-hover-light dark:bg-brand-solid-hover-light text-primary-dark hover:bg-brand-primary-dark dark:hover:bg-brand-primary-dark"
                : ""}
              on:click={() => (currentCalendarMode = "month")}>Month</Button
            >
          </div>

          <!-- create a calendar -->
          <div class="w-full border border-secondary-light dark:border-secondary-dark rounded-lg">
            <!-- calendar header -->
            <div class="px-8 py-4 flex justify-between items-center">
              <Button class="m-0 p-0" color="light" on:click={() => handleChangeDate("previous")}
                ><ChevronLeftOutline
                  class="size-8 text-fg-secondary-light dark:text-fg-secondary-dark"
                /></Button
              >
              {#if currentCalendarMode === "month"}
                <P class="font-semibold"
                  >{daysInMonth[15].dateObject.toLocaleDateString("en-US", {
                    month: "long",
                    year: "numeric",
                  })}</P
                >
              {:else if currentCalendarMode === "week"}
                <P class="font-semibold">
                  {format(startOfWeek(new Date(currentDate)), "dd MMMM")} - {format(
                    endOfWeek(new Date(currentDate)),
                    "dd MMMM yyyy"
                  )}
                </P>
              {:else if currentCalendarMode === "day"}
                <!-- else content here -->
              {/if}
              <Button class="m-0 p-0" color="light" on:click={() => handleChangeDate("next")}
                ><ChevronRightOutline
                  class="size-8 text-fg-secondary-light dark:text-fg-secondary-dark"
                /></Button
              >
            </div>

            <!-- calendar content -->
            {#if currentCalendarMode === "month"}
              <div class="grid grid-cols-7">
                {#each daysInMonth as day}
                  <button
                    type="button"
                    on:click={() => handleClickDate(day)}
                    class={`flex flex-col p-2 min-h-40 justify-start items-center border border-secondary-light dark:border-secondary-dark ${day.isCurrentMonth ? "bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark" : "bg-secondary-light dark:bg-secondary-dark hover:bg-primary-light hover:dark:bg-primary-dark"}`}
                  >
                    <P>{day.dayName}</P>
                    <P class="mb-4">{day.date}</P>
                    <div class="rounded-md bg-brand-solid-hover-light">
                      {#each day.dateEvent as event}
                        <P class="m-3 text-primary-dark">{event.title}</P>
                      {/each}
                    </div>
                  </button>
                {/each}
              </div>
            {:else if currentCalendarMode === "week"}
              <div class="container mx-auto">
                <!-- Scrollable Grid Container -->
                <div class="overflow-y-auto max-h-[60vh]">
                  <!-- Sticky Header -->
                  <div class="grid sticky top-0 z-10 grid-cols-8 grid-rows-1">
                    <div class="border border-secondary-light dark:border-secondary-dark"></div>
                    <!-- Placeholder for the top-left corner -->
                    {#each daysInWeek as day}
                      <P
                        class="grid-item day-header text-center p-2 min-h-24 text-sm border border-secondary-light dark:border-secondary-dark"
                      >
                        {day.dayName}<br />
                        <span class="text-3xl font-normal">{day.dateNumber}</span>
                      </P>
                    {/each}
                  </div>

                  <!-- Calendar Grid -->
                  <div class="grid grid-cols-8 grid-rows-25">
                    {#each getHoursOfDay() as hour}
                      <P
                        class="grid-item hour-header text-center pt-4 border border-secondary-light dark:border-secondary-dark"
                        >{hour}:00</P
                      >
                      {#each daysInWeek as _day}
                        <button
                          on:click={() => handleClickHour(_day, hour)}
                          class="grid-item min-h-40 flex flex-col items-center border border-secondary-light dark:border-secondary-dark"
                        >
                          {#if _day.events.length > 0}
                            {#each _day.events as event}
                              {#if hour >= parseInt(event.timeStart.split(":")[0]) && hour < parseInt(event.timeEnd.split(":")[0])}
                                <div class="rounded-md bg-brand-solid-hover-light mt-2 mx-auto">
                                  <P class="my-2 mx-3 text-primary-dark">{event.title}</P>
                                </div>
                              {/if}
                            {/each}
                          {/if}
                        </button>
                      {/each}
                    {/each}
                  </div>
                </div>
              </div>
            {:else if currentCalendarMode === "day"}
              <P>not yet</P>
            {/if}
          </div>
        </div>
      </div>
    </div>
  {:else if tabDatas.selected === "availability"}
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
  {:else}
    <P>Loading...</P>
  {/if}
</div>
