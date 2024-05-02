<script>
  import dayjs from "dayjs";
  import datepickerStore from "./datepicker";
  import { keyControlsContextKey, storeContextKey } from "./context";
  import { setContext } from "svelte";
  import { derived, writable } from "svelte/store";
  import Theme from "./Theme.svelte";
  import Calendar from "./Calendar.svelte";
  import CrossfadeBoundary from "./CrossfadeBoundary.svelte";
  import { calendar as calendarDefaults } from "./defaults";

  export let selected = calendarDefaults.selected;
  export let start = calendarDefaults.start;
  export let end = calendarDefaults.end;
  export let format = calendarDefaults.format;
  export let formatted = "";
  export let theme = {};
  export let defaultTheme = undefined;
  export let startOfWeekIndex = 0;
  export let store = datepickerStore.get({ selected, start, end, startOfWeekIndex });
  export let handleClickDate = () => {};
  export let schedule = [];

  const focused = writable(false);

  setContext(storeContextKey, store);
  setContext(
    keyControlsContextKey,
    derived([store, focused], ([$s, $focused]) => ($focused ? $s.activeView : undefined))
  );

  const getFocusSetter = (bool) => () => focused.set(bool);

  $: selected = $store.selected;
  $: formatted = dayjs(selected).format(format);
</script>

<CrossfadeBoundary>
  <Theme {defaultTheme} {theme} let:style>
    <div
      {style}
      role="button"
      tabindex="0"
      on:focus={getFocusSetter(true)}
      on:blur={getFocusSetter(false)}
      on:mouseover={getFocusSetter(true)}
      on:mouseout={getFocusSetter(false)}
    >
      <Calendar {handleClickDate} {schedule} />
    </div>
  </Theme>
</CrossfadeBoundary>

<style>
  div {
    display: inline-block;
  }
</style>
