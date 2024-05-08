<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Taqie Fadlillah <taqief@hexacode.org>
  * Dean Ramhan <dean@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<script>
  import { onMount, onDestroy } from "svelte";
  import { format, startOfWeek, endOfWeek } from "date-fns";
  import { P, Button, Heading, Input, Modal, Label } from "flowbite-svelte";
  import {
    ClockOutline,
    ChevronLeftOutline,
    ChevronRightOutline,
    PlusOutline,
  } from "flowbite-svelte-icons";
  import DatePicker from "$lib/components/DatePicker.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import Switch from "$lib/components/Switch.svelte";
  import { TabBar } from "$lib";
  import { CR3Call } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);

  let currentTheme = "";
  let currentDate = new Date(); // Current date
  let currentDateEvents = []; // Events for the current date
  let currentCalendarMode = "month"; // Current calendar mode
  let modal = false;
  let modalMode = "form";
  let observer;
  let dataReady = false;

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

  let schedule = [];

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
  let timeAvailableIndividualMode = "off";
  let timeAvailableIndividual = [
    { start: "08:00", end: "12:00" },
    { start: "13:00", end: "16:00" },
    { start: "08:00", end: "12:00" },
    { start: "08:00", end: "12:00" },
    { start: "08:00", end: "12:00" },
    { start: "13:00", end: "16:00" },
    { start: "13:00", end: "16:00" },
  ];
  let timeAvailableIndividualCopy = [...timeAvailableIndividual];
  let daysInMonth = [];
  let daysInWeek = [];
  let hoursOfDay = [];

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
            date: dateEvent.date,
            timeStart: dateEvent.timeStart,
            timeEnd: dateEvent.timeEnd,
          })),
      );

      daysOfWeek.push({
        dayName: dayOfWeek,
        dateNumber: dateOfWeek,
        dateObject: new Date(day),
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

  function getEventsForDate(date, hour = -1) {
    let formattedDate = date.toISOString().split("T")[0]; // Get YYYY-MM-DD format
    let events = [];

    schedule.forEach((item) => {
      item.dateEvent.forEach((event) => {
        if (event.date === formattedDate) {
          // If hour is specified and is within the event time range, add the event to the list
          if (hour !== -1) {
            let eventStartHour = parseInt(event.timeStart.split(":")[0]);
            let eventEndHour = parseInt(event.timeEnd.split(":")[0]);
            if (hour >= eventStartHour && hour < eventEndHour) {
              events.push({
                title: item.title,
                timeStart: event.timeStart,
                timeEnd: event.timeEnd,
              });
            }
          } else {
            // If hour is not specified, add all events for the date
            events.push({
              title: item.title,
              timeStart: event.timeStart,
              timeEnd: event.timeEnd,
            });
          }
        }
      });
    });

    return events;
  }

  const handleClickDate = (date) => {
    // console.log(getEventsForDate(date.dateObject));
    const events = getEventsForDate(date.dateObject);
    modalMode = "event";
    modal = true;
    currentDateEvents = events;
    console.log("date clicked backend logic here");
  };

  const handleClickHour = (date, hour, mode) => {
    let currentDayOrDate;
    if (mode === "day") {
      currentDayOrDate = getDaysOfWeek(date).find((days) => {
        return days.dateNumber === date.getDate();
      });
    } else {
      currentDayOrDate = date;
    }

    currentDateEvents = []; // Initialize as empty array

    if (currentDayOrDate.events.length > 0) {
      for (let i = 0; i < currentDayOrDate.events.length; i++) {
        const event = currentDayOrDate.events[i];
        if (
          hour >= parseInt(event.timeStart.split(":")[0]) &&
          hour < parseInt(event.timeEnd.split(":")[0])
        ) {
          const events = getEventsForDate(currentDayOrDate.dateObject, hour);
          currentDateEvents = events;
          break; // Exit the loop if the event is found
        }
      }
    }

    modalMode = "event";
    modal = true;

    console.log("date clicked backend logic here");
  };

  const handleChangeDate = (direction) => {
    let currentHour = new Date().getHours();
    let currentMinute = new Date().getMinutes();
    let currentSecond = new Date().getSeconds();

    if (currentCalendarMode === "month") {
      if (direction === "previous") {
        currentDate = new Date(
          currentDate.getFullYear(),
          currentDate.getMonth() - 1,
          currentDate.getDate(),
          currentHour,
          currentMinute,
          currentSecond,
        );
      } else if (direction === "next") {
        currentDate = new Date(
          currentDate.getFullYear(),
          currentDate.getMonth() + 1,
          currentDate.getDate(),
          currentHour,
          currentMinute,
          currentSecond,
        );
      }

      // Update the daysInMonth array based on the new currentDate
      daysInMonth = getDaysInMonth(currentDate.getMonth(), currentDate.getFullYear());
    } else if (currentCalendarMode === "week") {
      if (direction === "previous") {
        currentDate = new Date(
          currentDate.getFullYear(),
          currentDate.getMonth(),
          currentDate.getDate() - 7,
          currentHour,
          currentMinute,
          currentSecond,
        );
      } else if (direction === "next") {
        currentDate = new Date(
          currentDate.getFullYear(),
          currentDate.getMonth(),
          currentDate.getDate() + 7,
          currentHour,
          currentMinute,
          currentSecond,
        );
      }

      // Update the daysInWeek array based on the new currentDate
      daysInWeek = getDaysOfWeek(new Date(currentDate));
    } else if (currentCalendarMode === "day") {
      if (direction === "previous") {
        currentDate = new Date(
          currentDate.getFullYear(),
          currentDate.getMonth(),
          currentDate.getDate() - 1,
          currentHour,
          currentMinute,
          currentSecond,
        );
      } else if (direction === "next") {
        currentDate = new Date(
          currentDate.getFullYear(),
          currentDate.getMonth(),
          currentDate.getDate() + 1,
          currentHour,
          currentMinute,
          currentSecond,
        );
      }
    }
  };

  const handleCreateNewSchedule = async () => {
    newSchedule.dateEvent[1].timeStart = newSchedule.dateEvent[0].timeStart;
    newSchedule.dateEvent[1].timeEnd = newSchedule.dateEvent[0].timeEnd;
    // console.log(newSchedule);

    const postProcessedDate = newSchedule.dateEvent.map((date, index) => {
      return {
        title: newSchedule.title,
        datetime: new Date(
          index === 0
            ? newSchedule.dateEvent[0].date + "T" + date.timeStart + ":00"
            : newSchedule.dateEvent[1].date + "T" + date.timeEnd + ":00",
        ),
      };
    });

    // disini mas dean
    let date_1;
    let date_2;
    let dateSame = false;
    for (let i = 0; i < postProcessedDate.length; i++) {
      const element = postProcessedDate[i];

      if (i === 0) {
        date_1 = element.datetime.toISOString().split("T")[0];
      } else if (i === 1) {
        date_2 = element.datetime.toISOString().split("T")[0];
        if (date_1 === date_2) {
          dateSame = true;
        }
      }
    }

    if (dateSame) {
      let token = sessionStorage.getItem("token");
      try {
        let response = await call.call("client-dashboard", "create-schedule", {
          URLQuery: {},
          URLParams: {},
          Params: {
            title: postProcessedDate[0].title,
            date: postProcessedDate[0].datetime.toISOString().split("T")[0],
            start_time:
              postProcessedDate[0].datetime.getHours() +
              ":" +
              postProcessedDate[0].datetime.getMinutes() +
              ":00",
            end_time:
              postProcessedDate[1].datetime.getHours() +
              ":" +
              postProcessedDate[1].datetime.getMinutes() +
              ":00",
            is_active: "Yes",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status === 200) {
          console.log(response);
          modal = false;
        }
      } catch (error) {
        console.log(error);
      }
    } else {
      for (let i = 0; i < postProcessedDate.length; i++) {
        let token = sessionStorage.getItem("token");
        try {
          let response = await call.call("client-dashboard", "create-schedule", {
            URLQuery: {},
            URLParams: {},
            Params: {
              title: postProcessedDate[i].title,
              date: postProcessedDate[i].datetime.toISOString().split("T")[0],
              start_time:
                postProcessedDate[0].datetime.getHours() +
                ":" +
                postProcessedDate[0].datetime.getMinutes() +
                ":00",
              end_time:
                postProcessedDate[1].datetime.getHours() +
                ":" +
                postProcessedDate[1].datetime.getMinutes() +
                ":00",
              is_active: true,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });
          if (response.status === 200) {
            console.log(response);
            modal = false;
          }
        } catch (error) {
          console.log(error);
        }
      }

      modal = false;
    }
  };

  const checkTimeFormat = (evt, time, index = -1) => {
    const checkAndCorrectTime = (time) => {
      if (!time || time.length > 2 || parseInt(time) > 23) {
        return time ? time.slice(-2) : "01";
      }
      return time;
    };

    const handleTimeChange = (time, value) => {
      switch (time) {
        case "newTimeStart":
          if (index === -1) {
            newSchedule.dateEvent[0].timeStart = value;
          } else {
            newSchedule.dateEvent[index].timeStart = value;
          }
          break;
        case "newTimeEnd":
          if (index === -1) {
            newSchedule.dateEvent[0].timeEnd = value;
          } else {
            newSchedule.dateEvent[index].timeEnd = value;
          }
          break;
        case "availTimeStart":
          if (index === -1) {
            timeAvailable[0] = value;
          } else {
            timeAvailableIndividual[index].start = value;
          }
          break;
        case "availTimeEnd":
          if (index === -1) {
            timeAvailable[1] = value;
          } else {
            timeAvailableIndividual[index].end = value;
          }
          break;
      }
    };

    if (/[^0-9:]/.test(evt.target.value)) {
      evt.target.value = evt.target.value.replace(/[^0-9:]/g, "");
    }

    if (evt.target.value.includes(":")) {
      const targetValueSplit = evt.target.value.split(":");
      targetValueSplit[0] = checkAndCorrectTime(targetValueSplit[0]);
      targetValueSplit[1] = checkAndCorrectTime(targetValueSplit[1]);
      evt.target.value = targetValueSplit.join(":");

      const isValidTimeFormat =
        targetValueSplit[0].length === 2 && targetValueSplit[1].length === 2;

      if (isValidTimeFormat) {
        handleTimeChange(time, evt.target.value);
        console.log(newSchedule);
        console.log(timeAvailable);
      }
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
      JSON.stringify(timeAvailable) !== JSON.stringify(timeAvailableCopy) ||
      JSON.stringify(timeAvailableIndividual) !== JSON.stringify(timeAvailableIndividualCopy)
    ) {
      return true;
    }
    return false;
  };

  const handleSaveConfig = () => {
    console.log("save config backend logic here");
  };

  async function getSchedule() {
    let today = new Date();

    let month = today.getMonth();

    let firstDayOfMonth = new Date(today.getFullYear(), month, 1);

    let firstDayOfWeek = firstDayOfMonth.getDay();

    if (firstDayOfWeek >= 1 && firstDayOfWeek <= 6) {
      firstDayOfMonth.setDate(firstDayOfMonth.getDate() - firstDayOfWeek);
    }
    let lastDayOfMonth = new Date(today.getFullYear(), month + 1, 0);

    let lastDayOfWeek = lastDayOfMonth.getDay();

    if (lastDayOfWeek >= 1 && lastDayOfWeek <= 6) {
      lastDayOfMonth.setDate(lastDayOfMonth.getDate() + (6 - lastDayOfWeek));
    }

    let token = sessionStorage.getItem("token");

    try {
      let response = await call.call("client-dashboard", "get-schedule", {
        URLQuery: {},
        URLParams: {},
        Params: {
          start_date: firstDayOfMonth.toISOString().split("T")[0],
          end_date: lastDayOfMonth.toISOString().split("T")[0],
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });
      console.log("ayam", response);
      if (response.status === 200) {
        let promise = response.result.data.map((item) => {
          let date = new Date(item.attributes.date);
          date.setDate(date.getDate() - 1);
          return {
            title: item.attributes.title,
            dateEvent: [
              {
                date: date.toISOString().split("T")[0],
                timeStart: item.attributes.start_time.substring(0, 5),
                timeEnd: item.attributes.end_time.substring(0, 5),
              },
            ],
          };
        });
        schedule = await Promise.all(promise);
        daysInMonth = getDaysInMonth(currentDate.getMonth(), currentDate.getFullYear());
        daysInWeek = getDaysOfWeek(new Date(currentDate));
        hoursOfDay = getHoursOfDay();
        dataReady = true;
        console.log(schedule);
      }
    } catch (error) {
      console.log(error);
    }
  }

  onMount(() => {
    getSchedule();
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
  <div class="flex justify-between items-center">
    <Heading
      tag="h3"
      class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
    >
      Calendar
    </Heading>
    <Button
      class="text-nowrap"
      on:click={() => {
        modalMode = "form";
        modal = true;
      }}><PlusOutline class="size-6" /> Add new booking</Button
    >
  </div>
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
          {#if dataReady}
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
                      "dd MMMM yyyy",
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
                <!-- Scrollable Grid Container -->
                <div class="overflow-y-auto max-h-[60vh]">
                  <!-- Sticky Header -->
                  <div
                    class="grid sticky top-0 z-10 grid-cols-8 grid-rows-1 bg-primary-light dark:bg-primary-dark"
                  >
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
                          on:click={() => handleClickHour(_day, hour, "week")}
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
              {:else if currentCalendarMode === "day"}
                <!-- Scrollable Grid Container -->
                <div class="overflow-y-auto max-h-[60vh]">
                  <!-- Sticky Header -->
                  <div
                    class="grid sticky top-0 z-10 grid-cols-2 grid-rows-1 bg-primary-light dark:bg-primary-dark"
                    style="grid-template-columns: 1fr 4fr;"
                  >
                    <!-- Placeholder for the top-left corner -->
                    <div class="border border-secondary-light dark:border-secondary-dark"></div>
                    <P
                      class="grid-item day-header text-center p-2 min-h-24 text-sm border border-secondary-light dark:border-secondary-dark"
                    >
                      {currentDate.toLocaleDateString("en-US", { weekday: "long" })}<br />
                      <span class="text-3xl font-normal">{currentDate.getDate()}</span>
                    </P>
                  </div>

                  <!-- Calendar Grid -->
                  <div
                    class="grid grid-cols-2 grid-rows-25"
                    style="grid-template-columns: 1fr 4fr;"
                  >
                    {#each hoursOfDay as hour}
                      <P
                        class="grid-item hour-header text-center min-h-40 border border-secondary-light dark:border-secondary-dark pt-4"
                      >
                        {hour}:00
                      </P>
                      <button
                        on:click={() => handleClickHour(currentDate, hour, "day")}
                        class="grid-item hour-header border border-secondary-light dark:border-secondary-dark p-2 px-4 flex flex-col"
                      >
                        {#each daysInWeek as _day}
                          {#if _day.events.length > 0}
                            {#each _day.events as event}
                              {#if currentDate.getDate() === new Date(event.date).getDate()}
                                {#if hour >= parseInt(event.timeStart.split(":")[0]) && hour < parseInt(event.timeEnd.split(":")[0])}
                                  <div class="rounded-md w-full bg-brand-solid-hover-light mx-auto">
                                    <P class="my-1 mx-2 text-primary-dark">{event.title}</P>
                                  </div>
                                {/if}
                              {/if}
                            {/each}
                          {/if}
                        {/each}
                      </button>
                    {/each}
                  </div>
                </div>
              {/if}
            </div>
          {/if}
        </div>
      </div>
    </div>
  {:else if tabDatas.selected === "availability"}
    <div class="flex">
      <div class="w-1/3"><P class="font-semibold">Availability</P></div>
      <div class="w-2/3">
        {#if timeAvailableIndividualMode === "off"}
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
                checkTimeFinal(evt, "availTimeStart");
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
        {:else}
          <!-- day availability selector -->
          <div class="flex flex-col gap-3">
            {#each dayAvailable as available, index}
              <div class="flex">
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
                <div class="flex mx-4 w-full">
                  <Input
                    type="text"
                    placeholder={format(new Date(), "kk:mm")}
                    value={timeAvailableIndividual[index].start}
                    class="bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark focus:ring-0 border-gray-300 focus:border-gray-300 dark:border-gray-500 focus:dark:border-gray-500 flex-grow"
                    on:keyup={(evt) => {
                      checkTimeFormat(evt, "availTimeStart", index);
                      checkConfigChange();
                    }}
                    on:change={(evt) => {
                      checkTimeFinal(evt, "availTimeStart", index);
                    }}
                  >
                    <ClockOutline slot="right" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
                  </Input>
                  <P class="mx-4 text-center content-center">-</P>
                  <Input
                    type="text"
                    placeholder={format(new Date(), "kk:mm")}
                    value={timeAvailableIndividual[index].end}
                    class="bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark focus:ring-0 border-gray-300 focus:border-gray-300 dark:border-gray-500 focus:dark:border-gray-500 flex-grow"
                    on:keyup={(evt) => {
                      checkTimeFormat(evt, "availTimeEnd", index);
                      checkConfigChange();
                    }}
                    on:change={(evt) => {
                      checkTimeFinal(evt, "availTimeEnd", index);
                    }}
                  >
                    <ClockOutline slot="right" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
                  </Input>
                </div>
              </div>
            {/each}
          </div>
        {/if}

        <!-- individual time toggle -->
        <Switch
          bind:value={timeAvailableIndividualMode}
          label="Set individual hour"
          fontSize={16}
          design="slider"
        />

        <!-- save config -->
        <div class="flex justify-end">
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

<Modal
  title={modalMode === "form" ? "Create new schedule" : "Event detail"}
  bind:open={modal}
  outsideclose
>
  {#if modalMode === "form"}
    <!-- create new schedule modal -->
    <div>
      <Label class="ml-2">Title</Label>
      <Input
        bind:value={newSchedule.title}
        class="bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark focus:ring-0 border-gray-300 focus:border-gray-300 dark:border-gray-500 focus:dark:border-gray-500 cursor-default"
      />
    </div>

    <div class="flex gap-4">
      <div class="flex-grow">
        <Label class="ml-2">Date start</Label>
        <DatePicker
          on:datepicked={(evt) => {
            // console.log(evt.detail.datepicked);
            newSchedule.dateEvent[0].date = evt.detail.datepicked;
          }}
          customclass="block w-full disabled:cursor-not-allowed disabled:opacity-50 rtl:text-right p-2.5 bg-primary-light dark:bg-primary-dark text-gray-900 dark:text-white dark:placeholder-gray-400 border-gray-300 dark:border-gray-500 text-sm rounded-lg"
        />
      </div>

      <div class="flex-grow">
        <Label class="ml-2">Date end</Label>
        <DatePicker
          bind:value={newSchedule.dateEvent[1].date}
          on:datepicked={(evt) => {
            // console.log(evt.detail.datepicked);
            newSchedule.dateEvent[1].date = evt.detail.datepicked;
          }}
          customclass="block w-full disabled:cursor-not-allowed disabled:opacity-50 rtl:text-right p-2.5 bg-primary-light dark:bg-primary-dark text-gray-900 dark:text-white dark:placeholder-gray-400 border-gray-300 dark:border-gray-500 text-sm rounded-lg"
        />
      </div>
    </div>

    <div class="flex gap-4">
      <div class="flex-grow">
        <Label class="ml-2">Time start</Label>
        <Input
          type="text"
          placeholder={format(new Date(), "kk:mm")}
          value={newSchedule.dateEvent[0].timeStart}
          class="bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark focus:ring-0 border-gray-300 focus:border-gray-300 dark:border-gray-500 focus:dark:border-gray-500 flex-grow"
          on:keyup={(evt) => {
            checkTimeFormat(evt, "newTimeStart");
            checkConfigChange();
          }}
          on:change={(evt) => {
            checkTimeFinal(evt, "newTimeStart");
          }}
        >
          <ClockOutline slot="right" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
        </Input>
      </div>

      <div class="flex-grow">
        <Label class="ml-2">Time end</Label>
        <Input
          type="text"
          placeholder={format(new Date(), "kk:mm")}
          value={newSchedule.dateEvent[0].timeEnd}
          class="bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark focus:ring-0 border-gray-300 focus:border-gray-300 dark:border-gray-500 focus:dark:border-gray-500 flex-grow"
          on:keyup={(evt) => {
            checkTimeFormat(evt, "newTimeEnd");
            checkConfigChange();
          }}
          on:change={(evt) => {
            checkTimeFinal(evt, "newTimeEnd");
          }}
        >
          <ClockOutline slot="right" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
        </Input>
      </div>
    </div>
  {:else if modalMode === "event"}
    <!-- event detail event modal -->
    {#if currentDateEvents.length > 0}
      <!-- content here -->
      {#each currentDateEvents as event}
        <div class="p-4 border border-secondary-light dark:border-secondary-dark rounded-lg">
          <Label class="ml-2 text-lg">Title</Label>
          <Input
            bind:value={event.title}
            readonly
            class="bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark focus:ring-0 border-gray-300 focus:border-gray-300 dark:border-gray-500 focus:dark:border-gray-500 cursor-default"
          />
          <div class="mt-4 flex items-center w-full gap-4">
            <div class="flex-grow">
              <Label class="ml-2 text-lg">Time Start</Label>
              <Input
                bind:value={event.timeStart}
                readonly
                class="bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark focus:ring-0 border-gray-300 focus:border-gray-300 dark:border-gray-500 focus:dark:border-gray-500 cursor-default"
              />
            </div>
            <div class="flex-grow">
              <Label class="ml-2 text-lg">Time End</Label>
              <Input
                bind:value={event.timeEnd}
                readonly
                class="bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark focus:ring-0 border-gray-300 focus:border-gray-300 dark:border-gray-500 focus:dark:border-gray-500 cursor-default"
              />
            </div>
          </div>
        </div>
      {/each}
    {:else}
      <P>No event found for the selected {currentCalendarMode === "month" ? "date" : "hour"}.</P>
    {/if}
  {:else}
    <!-- placeholder -->
  {/if}

  <svelte:fragment slot="footer">
    {#if modalMode === "form"}
      <!-- create new schedule modal -->
      <Button
        on:click={() => {
          handleCreateNewSchedule();
          modal = false;
          newSchedule = {
            title: "",
            dateEvent: [
              {
                date: new Date(),
                timeStart: format(new Date(), "kk:mm"),
                timeEnd: format(new Date(), "kk:mm"),
              },
              {
                date: new Date(),
                timeStart: format(new Date(), "kk:mm"),
                timeEnd: format(new Date(), "kk:mm"),
              },
            ],
          };
        }}>Create</Button
      >
    {:else if modalMode === "event"}
      <!-- view detail event modal -->
      <Button
        class="text-nowrap"
        on:click={() => {
          modalMode = "form";
          modal = true;
        }}><PlusOutline class="size-6" /> Add new booking</Button
      >
    {:else}
      <!-- place holder -->
    {/if}
    <Button
      color={modalMode === "form" ? "alternative" : "red"}
      on:click={() => {
        modal = false;
        newSchedule = {
          title: "",
          dateEvent: [
            {
              date: new Date(),
              timeStart: format(new Date(), "kk:mm"),
              timeEnd: format(new Date(), "kk:mm"),
            },
            {
              date: new Date(),
              timeStart: format(new Date(), "kk:mm"),
              timeEnd: format(new Date(), "kk:mm"),
            },
          ],
        };
      }}
    >
      {modalMode === "form" ? "Cancel" : "Close"}
    </Button>
  </svelte:fragment>
</Modal>
