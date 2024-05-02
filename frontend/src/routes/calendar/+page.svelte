<script>
  // import { InlineCalendar, Datepicker } from "svelte-calendar";
  import InlineCalendar from "$lib/components/InlineCalendar/InlineCalendar.svelte";
  const theme = {
    calendar: {
      width: "600px",
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
          highlight: "#eb7400",
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
  };

  // import { DatePicker } from "date-picker-svelte";
  // let date = new Date();

  // import { DatePicker } from "@svelte-plugins/datepicker";
  // import { format } from "date-fns";

  // let startDate = new Date();
  // let dateFormat = "MM/dd/yy";
  // let isOpen = false;

  // const toggleDatePicker = () => (isOpen = !isOpen);

  // const formatDate = (dateString) => {
  //   return (dateString && format(new Date(dateString), dateFormat)) || "";
  // };

  // let formattedStartDate = formatDate(startDate);

  // const onChange = () => {
  //   startDate = new Date(formattedStartDate);
  // };

  // $: formattedStartDate = formatDate(startDate);
</script>

<!-- <Datepicker selected={date} /> -->
<InlineCalendar {theme} />

<!-- <DatePicker bind:value={date} /> -->

<!-- <DatePicker bind:isOpen bind:startDate>
  <input
    type="text"
    placeholder="Select date"
    bind:value={formattedStartDate}
    on:click={toggleDatePicker}
  />
</DatePicker>

<style>
  input[type="text"] {
    border: 1px solid #eee;
    border-radius: 4px;
    padding: 8px;
  }
</style> -->
