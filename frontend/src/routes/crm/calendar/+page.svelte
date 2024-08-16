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
  import { _, locale } from "svelte-i18n";
  import { DatePicker } from "@svelte-plugins/datepicker";
  import { format, startOfWeek, endOfWeek, addHours, addMinutes } from "date-fns";
  import { id as localeId, enUS as localeEn } from "date-fns/locale";
  import {
    P,
    Button,
    Heading,
    Input,
    Modal,
    Label,
    Textarea,
    MultiSelect,
    Helper,
  } from "flowbite-svelte";
  import {
    PlusOutline,
    ClockOutline,
    TrashBinOutline,
    ChevronLeftOutline,
    ChevronRightOutline,
  } from "flowbite-svelte-icons";
  import AutoComplete from "simple-svelte-autocomplete";
  // import DatePicker from "$lib/components/DatePicker.svelte";
  import { ToastContainer, triggerToast } from "$lib";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import Switch from "$lib/components/Switch.svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import CalendarPlus02 from "$lib/icon/calendar-plus-02.svelte";

  import { call } from "$lib";

  let currentTheme = "";
  let currentDate = new Date(); // Current date
  let currentDateEvents = []; // Events for the current date
  let currentCalendarMode = "month"; // Current calendar mode
  let modal = false;
  let modalMode = "form";
  let observer;
  let dataReady = false;
  let frontEndOnly = false;
  let notAvailableBookings = false;
  let selectedDate = new Date();
  let selectedTime = null;
  let startDate = new Date();
  let endDate = null;
  let dateFormat = "d MMM yyyy";
  let isOpen = false;

  let tabDatas = {
    selected: "Calendar",
    data: [
      {
        name: "Calendar",
      },
      {
        name: "Availability Slots",
      },
    ],
  };

  let schedule = [];
  let personClient = [];

  let newSchedule = {
    title: "",
    description: "",
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
    person: [],
  };

  let idCalendar;
  let dayAvailable = [];
  let dayAvailableCopy = [];
  let timeAvailable = [];
  let timeAvailableCopy = [];
  let timeAvailableIndividualMode = "off";
  let timeAvailableIndividual = [];
  let timeAvailableIndividualCopy = [];
  let days = [];
  let daysInMonth = [];
  let daysInWeek = [];
  let hoursOfDay = [];

  let buttonDisable = true;

  const toggleDatePicker = () => {
    isOpen = !isOpen;
  };
  const formatDate = (dateString) => (dateString && format(new Date(dateString), dateFormat)) || "";

  $: formattedStartDate = formatDate(startDate);
  $: formattedEndDate = formatDate(endDate);

  $: startDate = new Date(selectedDate);
  $: endDate = new Date(selectedDate);

  $: newSchedule.dateEvent[0].date = format(new Date(startDate), "yyyy-MM-dd");
  $: newSchedule.dateEvent[1].date = format(new Date(endDate), "yyyy-MM-dd");

  $: newSchedule.dateEvent[0].timeStart = selectedTime ? selectedTime : format(new Date(), "kk:mm");
  $: newSchedule.dateEvent[1].timeStart = selectedTime ? selectedTime : format(new Date(), "kk:mm");
  $: newSchedule.dateEvent[0].timeEnd = selectedTime
    ? addOneHour(selectedTime)
    : format(addHours(new Date(), 1), "kk:mm");
  $: newSchedule.dateEvent[1].timeEnd = selectedTime
    ? addOneHour(selectedTime)
    : format(addHours(new Date(), 1), "kk:mm");

  $: buttonDisable = newSchedule.title.length < 1;

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
        dayName:
          $locale === "id"
            ? prevDate.toLocaleDateString("id", { weekday: "short" })
            : prevDate.toLocaleDateString("en-US", { weekday: "short" }),
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
        dayName:
          $locale === "id"
            ? currentDate.toLocaleDateString("id", { weekday: "short" })
            : currentDate.toLocaleDateString("en-US", { weekday: "short" }),
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
        dayName:
          $locale === "id"
            ? nextDate.toLocaleDateString("id", { weekday: "short" })
            : nextDate.toLocaleDateString("en-US", { weekday: "short" }),
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
      let dayOfWeek =
        $locale === "id"
          ? day.toLocaleDateString("id", { weekday: "short" })
          : day.toLocaleDateString("en-US", { weekday: "short" });
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
      let hour = i.toString().padStart(2, "0");
      hoursOfDay.push(hour);
    }

    return hoursOfDay;
  }

  function addOneHour(time) {
    const [hours, minutes] = time.split(":").map(Number);
    const date = new Date();
    date.setHours(hours);
    date.setMinutes(minutes);
    const newDate = addMinutes(date, 60);
    return format(newDate, "kk:mm");
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
                description: item.description,
                person: item.person.join(", "),
              });
            }
          } else {
            // If hour is not specified, add all events for the date
            events.push({
              title: item.title,
              timeStart: event.timeStart,
              timeEnd: event.timeEnd,
              description: item.description,
              person: item.person,
            });
          }
        }
      });
    });

    return events;
  }

  const handleClickDate = (date) => {
    const today = new Date();

    today.setHours(0, 0, 0, 0);

    if (date.dateObject <= today) {
      notAvailableBookings = true;
    } else {
      notAvailableBookings = false;
    }
    const events = getEventsForDate(date.dateObject);
    selectedDate = date.dateObject;
    selectedTime = null;
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

    selectedDate = mode === "day" ? currentDayOrDate.dateObject : date.dateObject;
    selectedTime = hour + ":00";

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

  const resetSchedule = () => {
    newSchedule = {
      title: "",
      description: "",
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
      person: [],
    };
  };

  const handleCreateNewSchedule = async () => {
    newSchedule.dateEvent[1].timeStart = newSchedule.dateEvent[0].timeStart;
    newSchedule.dateEvent[1].timeEnd = newSchedule.dateEvent[0].timeEnd;
    // console.log(newSchedule);

    const postProcessedDate = newSchedule.dateEvent.map((date, index) => {
      return {
        title: newSchedule.title,
        description: newSchedule.description,
        datetime: new Date(
          index === 0
            ? newSchedule.dateEvent[0].date + "T" + date.timeStart + ":00"
            : newSchedule.dateEvent[1].date + "T" + date.timeEnd + ":00",
        ),
      };
    });

    console.log("itu", postProcessedDate);

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

    let isAvailable = checkAvailability(postProcessedDate);
    if (!isAvailable) {
      return;
    }

    if (dateSame) {
      let token = sessionStorage.getItem("token");
      try {
        let response = await $call.call("client-dashboard", "create-schedule", {
          URLQuery: {},
          URLParams: {},
          Params: {
            title: postProcessedDate[0].title,
            date: postProcessedDate[0].datetime.toISOString().split("T")[0],
            start_time:
              ("0" + postProcessedDate[0].datetime.getHours()).slice(-2) +
              ":" +
              ("0" + postProcessedDate[0].datetime.getMinutes()).slice(-2) +
              ":" +
              ("0" + postProcessedDate[0].datetime.getSeconds()).slice(-2) +
              "." +
              ("00" + postProcessedDate[0].datetime.getMilliseconds()).slice(-3),
            end_time:
              ("0" + postProcessedDate[1].datetime.getHours()).slice(-2) +
              ":" +
              ("0" + postProcessedDate[1].datetime.getMinutes()).slice(-2) +
              ":" +
              ("0" + postProcessedDate[1].datetime.getSeconds()).slice(-2) +
              "." +
              ("00" + postProcessedDate[1].datetime.getMilliseconds()).slice(-3),
            description: postProcessedDate[0].description,
            is_active: "Yes",
            person: JSON.stringify(newSchedule.person),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status === 200) {
          console.log(response);
          triggerToast("created");
          await getSchedule();
          resetSchedule();
          modal = false;
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    } else {
      const startDate = postProcessedDate[0].datetime;
      const endDate = postProcessedDate[1].datetime;

      const timeDifference = Math.abs(endDate.getTime() - startDate.getTime());

      const daysDifference = Math.ceil(timeDifference / (1000 * 3600 * 24));

      console.log("Jarak tanggal:", daysDifference);

      for (let i = 0; i < daysDifference; i++) {
        const currentDate = new Date(startDate);
        currentDate.setDate(startDate.getDate() + i);
        let token = sessionStorage.getItem("token");
        try {
          let response = await $call.call("client-dashboard", "create-schedule", {
            URLQuery: {},
            URLParams: {},
            Params: {
              title: postProcessedDate[0].title,
              date: currentDate.toISOString().split("T")[0],
              start_time:
                ("0" + postProcessedDate[0].datetime.getHours()).slice(-2) +
                ":" +
                ("0" + postProcessedDate[0].datetime.getMinutes()).slice(-2) +
                ":" +
                ("0" + postProcessedDate[0].datetime.getSeconds()).slice(-2) +
                "." +
                ("00" + postProcessedDate[0].datetime.getMilliseconds()).slice(-3),
              end_time:
                ("0" + postProcessedDate[1].datetime.getHours()).slice(-2) +
                ":" +
                ("0" + postProcessedDate[1].datetime.getMinutes()).slice(-2) +
                ":" +
                ("0" + postProcessedDate[1].datetime.getSeconds()).slice(-2) +
                "." +
                ("00" + postProcessedDate[1].datetime.getMilliseconds()).slice(-3),
              description: postProcessedDate[0].description,
              is_active: "Yes",
              person: JSON.stringify(newSchedule.person),
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });
          if (response.status === 200) {
            console.log("test", i);
            console.log(currentDate.toISOString().split("T")[0]);
            triggerToast("created");
            await getSchedule();
            modal = false;
          }
        } catch (error) {
          triggerToast("error");
          console.log(error);
        }
      }

      modal = false;
    }
  };

  const checkAvailability = (postProcessedDate) => {
    for (let i = 0; i < postProcessedDate.length; i++) {
      const dateObj = postProcessedDate[i];
      const date = new Date(dateObj.datetime);
      let dayIndex = date.getDay() - 1;

      console.log("day index", dayIndex);

      if (dayIndex === -1) {
        dayIndex = 6;
      }

      if (!dayAvailableCopy[dayIndex]) {
        console.log(`Hari yang Anda pilih tidak tersedia untuk booking.`);
        triggerToast("dayNotAvailable");
        return false;
      }

      if (
        !isTimeAvailable(
          `${postProcessedDate[0].datetime.getHours()}:${postProcessedDate[0].datetime.getMinutes()}`,
          `${postProcessedDate[1].datetime.getHours()}:${postProcessedDate[1].datetime.getMinutes()}`,
          dayIndex,
        )
      ) {
        console.log(`Waktu yang Anda pilih tidak tersedia untuk booking.`);
        triggerToast("timeNotAvailable");
        return false;
      }
    }

    return true;
  };

  const isTimeAvailable = (timeStart, timeEnd, dayIndex) => {
    console.log(timeStart, timeEnd);
    const availableTime = timeAvailableIndividualCopy[dayIndex];
    const start = availableTime.start;
    const end = availableTime.end;

    const [availableStartHour, availableStartMinute] = start.split(":").map(Number);
    const [availableEndHour, availableEndMinute] = end.split(":").map(Number);

    const [bookStartHour, bookStartMinute] = timeStart.split(":").map(Number);
    const [bookEndHour, bookEndMinute] = timeEnd.split(":").map(Number);

    const availableStartTime = availableStartHour * 60 + availableStartMinute;
    const availableEndTime = availableEndHour * 60 + availableEndMinute;
    const bookStartTime = bookStartHour * 60 + bookStartMinute;
    const bookEndTime = bookEndHour * 60 + bookEndMinute;

    console.log(bookStartTime, ">=", availableStartTime, "&&", bookEndTime, "<=", availableEndTime);

    return bookStartTime >= availableStartTime && bookEndTime <= availableEndTime;
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

  const handleSaveConfig = async () => {
    // console.log(dayAvailable);
    // console.log(timeAvailable);
    // console.log(timeAvailableIndividual);
    // console.log(timeAvailableIndividualMode);
    // console.log("save config backend logic here");

    let timeAvailability = [];

    if (timeAvailableIndividualMode === "on") {
      for (let i = 0; i < days.length; i++) {
        timeAvailability.push({
          name: days[i],
          start: timeAvailableIndividual[i].start,
          end: timeAvailableIndividual[i].end,
          is_active: dayAvailable[i],
        });
      }

      console.log(timeAvailability);
    } else {
      for (let i = 0; i < days.length; i++) {
        timeAvailability.push({
          name: days[i],
          start: timeAvailable[0],
          end: timeAvailable[1],
          is_active: dayAvailable[i],
        });
      }

      console.log(timeAvailability);
    }

    if (timeAvailability.length === 0) {
      alert("Please choose at least one day and time");
      return;
    }

    try {
      if (idCalendar === undefined) {
        let token = sessionStorage.getItem("token");
        let response = await $call.call("client-dashboard", "create-availability", {
          URLQuery: {},
          URLParams: {},
          Params: {
            availability: JSON.stringify(timeAvailability),
            individualMode: timeAvailableIndividualMode === "on" ? "Yes" : "No",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          triggerToast("created");
          console.log(response);
        } else {
          triggerToast("error");
        }
      } else {
        let token = sessionStorage.getItem("token");
        let response = await $call.call("client-dashboard", "update-availability", {
          URLQuery: {},
          URLParams: {
            id: idCalendar.toString(),
          },
          Params: {
            availability: JSON.stringify(timeAvailability),
            individualMode: timeAvailableIndividualMode === "on" ? "Yes" : "No",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("ayam", response);

        if (response.status === 200) {
          triggerToast("created");
          console.log(response);
        } else {
          triggerToast("error");
        }
      }

      await getSchedule();
    } catch (error) {
      triggerToast("error");
      console.log(error);
    }
  };

  async function getSchedule() {
    let today = new Date();

    let month = today.getMonth();

    let firstDayOfMonth = new Date(today.getFullYear(), month, 1);

    let firstDayOfWeek = firstDayOfMonth.getDay();

    if (firstDayOfWeek >= 1 && firstDayOfWeek <= 6) {
      firstDayOfMonth.setDate(firstDayOfMonth.getDate() - firstDayOfWeek);
    }
    let lastDayOfMonth = new Date(today.getFullYear(), month + 1, 1);

    let lastDayOfWeek = lastDayOfMonth.getDay();

    if (lastDayOfWeek >= 1 && lastDayOfWeek <= 6) {
      lastDayOfMonth.setDate(lastDayOfMonth.getDate() + (6 - lastDayOfWeek));
    }

    let token = sessionStorage.getItem("token");

    try {
      if (!frontEndOnly) {
        var response = await $call.call("client-dashboard", "get-schedule", {
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
        if (response.status === 200) {
          let promise = response.result.booking.data.map((item) => {
            let date = new Date(item.attributes.date);
            date.setDate(date.getDate() - 1);
            return {
              title: item.attributes.title,
              description: item.attributes.description,
              dateEvent: [
                {
                  date: date.toISOString().split("T")[0],
                  timeStart: item.attributes.start_time.substring(0, 5),
                  timeEnd: item.attributes.end_time.substring(0, 5),
                },
              ],
              person: item.attributes.person.data.map((item) => {
                return item.attributes.fullname;
              }),
            };
          });
          schedule = await Promise.all(promise);

          console.log("itu", schedule);

          daysInMonth = getDaysInMonth(currentDate.getMonth(), currentDate.getFullYear());
          daysInWeek = getDaysOfWeek(new Date(currentDate));
          hoursOfDay = getHoursOfDay();

          // console.log("cek", response.result.calendar);

          if (response.result.calendar === null) {
            days = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];
            dayAvailable = [false, false, false, false, false, false, false];
            dayAvailableCopy = [...dayAvailable];
            timeAvailable = ["08:00", "17:00"];
            timeAvailableCopy = [...timeAvailable];
            timeAvailableIndividual = [
              {
                start: "08:00",
                end: "17:00",
              },
              {
                start: "08:00",
                end: "17:00",
              },
              {
                start: "08:00",
                end: "17:00",
              },
              {
                start: "08:00",
                end: "17:00",
              },
              {
                start: "08:00",
                end: "17:00",
              },
              {
                start: "08:00",
                end: "17:00",
              },
              {
                start: "08:00",
                end: "17:00",
              },
            ];
          } else {
            idCalendar = response.result.calendar.id;
            console.log(idCalendar);
            let availability = response.result.calendar.attributes.availability;
            days = availability.map((day) => {
              return day.name;
            });
            dayAvailable = availability.map((days) => {
              return days.is_active;
            });
            dayAvailableCopy = [...dayAvailable];
            timeAvailable = [availability[0].start, availability[0].end];
            console.log("ini", timeAvailable);
            timeAvailableCopy = [...timeAvailable];
            timeAvailableIndividual = availability.map((time) => {
              return {
                start: time.start,
                end: time.end,
              };
            });
            timeAvailableIndividualCopy = [...timeAvailableIndividual];
            timeAvailableIndividualMode =
              response.result.calendar.attributes.individualMode === "Yes" ? "on" : "off";
          }

          dataReady = true;
        } else {
          triggerToast("error");
        }
      } else {
        daysInMonth = getDaysInMonth(currentDate.getMonth(), currentDate.getFullYear());
        daysInWeek = getDaysOfWeek(new Date(currentDate));
        hoursOfDay = getHoursOfDay();
        dataReady = true;
      }
    } catch (error) {
      triggerToast("error");
      console.log(error);
    }
  }

  const handleDeleteSchedule = async (event) => {
    console.log("Not yet implemented");
    console.log(selectedDate);
    let date = new Date(selectedDate);
    let token = sessionStorage.getItem("token");
    console.log(token);
    if (token !== null) {
      const response = await $call.call("client-dashboard", "delete-schedule", {
        URLQuery: {},
        URLParams: {},
        Params: {
          date: format(date, "yyyy-MM-dd"),
          start_time: event.timeStart + ":00",
          end_time: event.timeEnd + ":00",
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status === 200) {
        console.log(response);
        await getSchedule();
        modal = false;
      } else {
        console.log("Error", response);
      }
    }
  };

  const handleSearchInvitees = async (value) => {
    console.log("not yet implemented");
    console.log(value);
    return [];
    // try {
    //   const response = await axios.get(
    //     `${API_URL}/api/users?type=search&value=${value}&role=${inviteTypeSelected}&unit=${inviteUnitSelected}`,
    //     {
    //       headers: {
    //         Authorization: "Bearer " + userDataSession().jwt,
    //       },
    //     }
    //   );
    //   // console.log(response);

    //   if (response?.data?.status === 401) {
    //     alert("Session expired, please login again");
    //     window.location = ("/");
    //   }

    //   if (response?.data) {
    //     return response.data;
    //   } else {
    //     return [];
    //   }
    // } catch (error) {
    //   console.log(error);
    // }
  };

  const handleDeleteInvitee = (email) => {
    newSchedule.person = newSchedule.person.filter((person) => person.email !== email);
  };

  async function getPersons() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let page = 1;
      while (true) {
        let response = await $call.call("client-dashboard", "get-persons", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: 25,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        if (response.result.data.length > 0) {
          let promise = response.result.data.map((person) => {
            return {
              value: person.id,
              name:
                person.attributes.fullname !== null
                  ? person.attributes.fullname
                  : `<${person.attributes.email}>`,
            };
          });

          let data = await Promise.all(promise);
          personClient = personClient.concat(data);
          page++;
        } else {
          break;
        }
      }
    }
  }

  onMount(() => {
    getSchedule();
    getPersons();
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
  <div class="flex pb-4 justify-between items-center">
    <Heading
      tag="h3"
      class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
    >
      {$_("calendar")}
    </Heading>
    <Button
      class="text-nowrap gap-1"
      on:click={() => {
        modalMode = "form";
        modal = true;
      }}><PlusOutline class="size-5" /> {$_("addNewBooking")}</Button
    >
  </div>
  <div class="pb-8">
    <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />
  </div>
  {#if tabDatas.selected === "Calendar"}
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
              on:click={() => (currentCalendarMode = "day")}>{$_("day")}</Button
            >
            <Button
              color="light"
              class={currentCalendarMode === "week"
                ? "bg-brand-solid-hover-light dark:bg-brand-solid-hover-light text-primary-dark hover:bg-brand-primary-dark dark:hover:bg-brand-primary-dark"
                : ""}
              on:click={() => (currentCalendarMode = "week")}>{$_("week")}</Button
            >
            <Button
              color="light"
              class={currentCalendarMode === "month"
                ? "bg-brand-solid-hover-light dark:bg-brand-solid-hover-light text-primary-dark hover:bg-brand-primary-dark dark:hover:bg-brand-primary-dark"
                : ""}
              on:click={() => (currentCalendarMode = "month")}>{$_("month")}</Button
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
                  <P class="font-semibold">
                    {#if $locale === "id"}
                      {daysInMonth[15].dateObject.toLocaleDateString("id-ID", {
                        month: "long",
                        year: "numeric",
                      })}
                    {:else}
                      {daysInMonth[15].dateObject.toLocaleDateString("en-US", {
                        month: "long",
                        year: "numeric",
                      })}
                    {/if}
                  </P>
                {:else if currentCalendarMode === "week"}
                  <P class="font-semibold">
                    <!-- {format(startOfWeek(new Date(currentDate)), "dd MMMM")} - {format(
                      endOfWeek(new Date(currentDate)),
                      "dd MMMM yyyy",
                    )} -->
                    {#if $locale === "id"}
                      {format(startOfWeek(new Date(currentDate)), "dd MMMM", { locale: localeId })} -
                      {format(endOfWeek(new Date(currentDate)), "dd MMMM yyyy", {
                        locale: localeId,
                      })}
                    {:else}
                      {format(startOfWeek(new Date(currentDate)), "dd MMMM", { locale: localeEn })} -
                      {format(endOfWeek(new Date(currentDate)), "dd MMMM yyyy", {
                        locale: localeEn,
                      })}
                    {/if}
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
                      {#if $locale === "id"}
                        {currentDate.toLocaleDateString("id-ID", {
                          weekday: "long",
                        })}
                      {:else}
                        {currentDate.toLocaleDateString("en-US", {
                          weekday: "long",
                        })}
                      {/if}
                      <br />
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
  {:else if tabDatas.selected === "Availability Slots"}
    <div class="flex">
      <div class="w-1/3"><P class="font-semibold">Availability</P></div>
      <div class="w-2/3">
        {#if dataReady}
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
  <ToastContainer />
</div>

<Modal
  title={modalMode === "form" ? "Create new schedule" : "Event detail"}
  bind:open={modal}
  outsideclose
>
  {#if modalMode === "form"}
    <!-- create new schedule modal -->

    <div class="flex gap-4">
      <div class="flex-grow space-y-1">
        <Label class="ml-2">Date</Label>
        <DatePicker
          bind:isOpen
          bind:startDate
          bind:endDate
          isRange
          isMultipane
          enableFutureDates
          enablePastDates={false}
        >
          <button
            on:click={toggleDatePicker}
            class:open={isOpen}
            class="flex w-full p-2.5 border rounded-lg border-gray-300 dark:border-gray-500 gap-2 items-center"
          >
            <svelte:component
              this={CalendarPlus02}
              class="w-3 h-3 ms-2 text-gray-700 dark:text-gray-400"
            />
            <div class="date">
              {#if startDate != endDate}
                {formattedStartDate} - {formattedEndDate}
              {:else if startDate}
                {formattedStartDate}
              {:else}
                Pick a date, it can be a range
              {/if}
            </div>
          </button>
        </DatePicker>
        <Helper class="text-sm text-quaternary-light dark:text-quaternary-dark ms-1"
          >Pick a date, it can be a range</Helper
        >
      </div>
    </div>

    <div class="flex gap-4">
      <div class="flex-grow space-y-1">
        <Label class="ml-2">Time start</Label>
        <Input
          type="time"
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
        ></Input>
      </div>

      <div class="flex-grow space-y-1">
        <Label class="ml-2">Time end</Label>
        <Input
          type="time"
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
        ></Input>
      </div>
    </div>

    <div class="space-y-1">
      <Label class="ml-2">Title</Label>
      <Input
        required
        bind:value={newSchedule.title}
        maxlength="50"
        class="bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark focus:ring-0 border-gray-300 dark:border-gray-500 cursor-default"
      />
      {#if newSchedule.title.length === 0}
        <Helper class="text-sm !text-red-400 ms-1">Title is required</Helper>
      {/if}
    </div>

    <div class="space-y-1">
      <Label class="ml-2">Description</Label>
      <Textarea
        bind:value={newSchedule.description}
        class="bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark focus:ring-0 border-gray-300 focus:border-gray-300 dark:border-gray-500 focus:dark:border-gray-500 cursor-default resize-none"
      />
    </div>

    <div class="space-y-1">
      <Label class="ml-2">Invite Person</Label>
      <div
        class={`col-span-6 [&>*]:w-full [&>*]:!h-[2.75rem] [&>*]:m-0 [&>*>div]:rounded-lg
        `}
      >
        <MultiSelect
          bind:value={newSchedule.person}
          items={personClient}
          placeholder={personClient.length > 0 ? "" : "No person found"}
        />
      </div>
    </div>

    <!-- {#if newSchedule.person.length > 0}
      <div class="flex gap-2">
        {#each newSchedule.person as person}
          <P class="p-2 rounded-lg border"
            >{person.name} - {person.email}
            <span
              class="ml-2 cursor-pointer"
              on:click={() => {
                handleDeleteInvitee(person.email);
              }}>x</span
            ></P
          >
        {/each}
      </div>
    {/if} -->

    <!-- <div class="flex gap-4">
      <div class="flex-grow">
        <Label class="ml-2">Date start</Label>
        <DatePicker
          booking={true}
          on:datepicked={(evt) => {
            const selectedDate = new Date(evt.detail.datepicked);
            const currentDate = new Date();
            if (selectedDate < currentDate) {
              return;
            }
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
            const selectedDate = new Date(evt.detail.datepicked);
            const currentDate = new Date();
            if (selectedDate < currentDate) {
              return;
            }
            newSchedule.dateEvent[1].date = evt.detail.datepicked;
          }}
          customclass="block w-full disabled:cursor-not-allowed disabled:opacity-50 rtl:text-right p-2.5 bg-primary-light dark:bg-primary-dark text-gray-900 dark:text-white dark:placeholder-gray-400 border-gray-300 dark:border-gray-500 text-sm rounded-lg"
        />
      </div>
    </div> -->
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

          <Label class="mt-4 ml-2 text-lg">Description</Label>
          <Textarea
            bind:value={event.description}
            readonly
            class="bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark focus:ring-0 border-gray-300 focus:border-gray-300 dark:border-gray-500 focus:dark:border-gray-500 cursor-default resize-none"
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

          <Label class="mt-4 ml-2 text-lg" for="multiselect-disabled">Person</Label>
          <Input
            bind:value={event.person}
            readonly
            class="bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light hover:dark:bg-primary-hover-dark focus:ring-0 border-gray-300 focus:border-gray-300 dark:border-gray-500 focus:dark:border-gray-500 cursor-default"
          />

          <div class="flex justify-end mt-2">
            <Button
              class="text-nowrap"
              color="red"
              on:click={() => {
                handleDeleteSchedule(event);
              }}><TrashBinOutline class="size-6" /> Delete Schedule</Button
            >
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
      <Button disabled={buttonDisable} on:click={handleCreateNewSchedule}>Create</Button>
    {:else if modalMode === "event"}
      <!-- view detail event modal -->
      {#if !notAvailableBookings}
        <Button
          class="text-nowrap"
          on:click={() => {
            modalMode = "form";
            modal = true;
          }}><PlusOutline class="size-5" /> Add new booking</Button
        >
      {/if}
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
