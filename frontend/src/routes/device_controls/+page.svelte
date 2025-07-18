<script>
  import { onMount } from "svelte";
  import { Typography } from "$lib";
  import { getIcons, call } from "$lib";

  const ByeIcon = getIcons("bye");
  const MeditationYogaIcon = getIcons("meditationYoga");

  const actionTitles = {
    restart: "Just a moment...",
    shutdown: "Until next time!",
  };

  const actionSubtitles = {
    restart: "Hexacube is restarting. We'll be back online for you in no time.",
    shutdown: "Hexacube is shutting down. We're looking forward to your return.",
  };

  let actionParam = "";

  const capitalizeFirstLetter = (val) => {
    return String(val).charAt(0).toUpperCase() + String(val).slice(1);
  };

  const restartHexacube = async () => {
    try {
      let token = localStorage.getItem("token");
      if (token) {
        // const response = await $call.call("client-dashboard", "restart", {
        //   URLQuery: {},
        //   URLParams: {},
        //   Params: {},
        //   Headers: {
        //     Authorization: `Bearer ${token}`,
        //   },
        // });

        const response = {
          status: 400,
        }; // dummy response for testing

        if (response.status !== 200) {
          console.error("Failed to restart Hexacube:", response);
        } else {
          console.log("Hexacube is restarting...");
        }
      }
    } catch (error) {
      console.log("error", error);
    }
  };

  const shutdownHexacube = async () => {
    try {
      let token = localStorage.getItem("token");
      if (token) {
        // const response = await $call.call("client-dashboard", "shutdown", {
        //   URLQuery: {},
        //   URLParams: {},
        //   Params: {},
        //   Headers: {
        //     Authorization: `Bearer ${token}`,
        //   },
        // });

        const response = {
          status: 400,
        }; // dummy response for testing

        if (response.status !== 200) {
          console.error("Failed to shutdown Hexacube:", response);
        } else {
          console.log("Hexacube is stopping...");
        }
      }
    } catch (error) {
      console.log("error", error);
    }
  };

  onMount(() => {
    const params = new URLSearchParams(window.location.search);
    actionParam = params.get("action");
    if (!actionParam || !["restart", "shutdown"].includes(actionParam)) {
      actionParam = "restart"; // Default to restart if no valid action is provided
    }

    if (actionParam === "restart") {
      restartHexacube();
    } else if (actionParam === "shutdown") {
      shutdownHexacube();
    }
  });
</script>

<svelte:head>
  <title>{capitalizeFirstLetter(actionParam)} Hexacube</title>
</svelte:head>

<main class="flex flex-col items-center justify-center bg-primary-light w-full h-screen">
  {#if actionParam === "restart"}
    <MeditationYogaIcon svgClass="w-[20rem] h-[20rem] md:w-[30rem] md:h-[30rem]" />
  {:else}
    <ByeIcon svgClass="w-[20rem] h-[20rem] md:w-[30rem] md:h-[30rem]" />
  {/if}
  <Typography
    text={actionTitles[actionParam]}
    preset="Display xl/semibold"
    propClass="m-8 text-primary-light dark:text-primary-light"
  />
  <Typography
    text={actionSubtitles[actionParam]}
    preset="Text xl/regular"
    propClass="text-center text-primary-light dark:text-primary-light"
  />
</main>
