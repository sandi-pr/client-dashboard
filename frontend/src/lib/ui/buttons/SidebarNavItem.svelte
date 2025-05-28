<!--
  CR3 Client Dashboard
  Package: gitlab.hexacode.org/public-repo/js/hexacomponents-svelte
  Maintainer:
  * Sandi Permana <sandy@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>

  Copyright (C) 2024 Hexacode Teknologi Indonesia
  All Rights Reserved.
-->

<script>
  import { goto } from "$app/navigation";
  import { Typography, getIcons, Badge } from "$lib";
  import { page } from "$app/stores";
  import { _ } from "svelte-i18n";
  import { sidebarLargeClosed } from "$lib/utils/sidebarLargeClosed";

  const isFrontend = process.env.NODE_ENV === "development";

  let {
    href = "/",
    commingSoon = false,
    classes = "flex w-full gap-3 items-center justify-between px-3 py-2 hover:no-underline rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700",
    classesActive = "flex w-full gap-3 items-center justify-between px-3 py-2 hover:no-underline bg-gray-100 dark:bg-gray-700 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700",
    icon = "",
    classesSvg = "w-5 h-5 flex text-hc-graylight-500 dark:text-hc-graydark-400",
    typographyData = {},
  } = $props();

  function extractPath(url) {
    if (!url) return;
    const match = url.match(/^(?:https?:\/\/)?[^\/]+(\/[^?#]*)/);
    return match ? "/" + match[1].split("/")[1] : url;
  }

  let activeUrl = $derived("/" + $page.url.pathname.split("/")[isFrontend ? 1 : 2]);

  // const newFunc = async () => {
  //   // !!! DANGER: NEVER EXPOSE THESE IN PRODUCTION CLIENT-SIDE CODE !!!
  //   const MOODLE_URL = "http://158.140.178.96:9191";
  //   const MOODLE_WEBSERVICE_TOKEN = "210c5c247ff3289766608f7b622aa4b0"; // YOU ARE EXPOSING THIS!
  //   const MOODLE_AUTH_USERKEY_SERVICE_FUNCTION = "auth_userkey_request_login_url";

  //   let authenticatedUser = {
  //     username: "john.doe",
  //     firstname: "John",
  //     lastname: "Doe",
  //     email: "john.doe@example.com",
  //   };

  //   let userIp = "127.0.0.1"; // Placeholder.

  //   if (!authenticatedUser || !MOODLE_WEBSERVICE_TOKEN) {
  //     console.error("User not authenticated or Moodle token missing.");
  //     alert("Authentication error. Please log in to your app first.");
  //     return;
  //   }

  //   const params = new URLSearchParams();
  //   params.append("wstoken", MOODLE_WEBSERVICE_TOKEN);
  //   params.append("wsfunction", MOODLE_AUTH_USERKEY_SERVICE_FUNCTION);
  //   params.append("moodlewsrestformat", "json");

  //   // --- THE CORRECT FIX FOR "Only arrays accepted" with 'user' parameter ---
  //   // The Web Service expects a parameter named 'user' (singular),
  //   // which then contains an array of user objects.
  //   // We achieve this using bracket notation.
  //   const userFields = {
  //     username: authenticatedUser.username,
  //     email: authenticatedUser.email,
  //     firstname: authenticatedUser.firstname,
  //     lastname: authenticatedUser.lastname,
  //     // Add other fields here if you need to update/create them in Moodle
  //   };

  //   // Iterate over the user object and append each field with the correct notation
  //   // to represent 'user' as an array of one object.
  //   Object.keys(userFields).forEach((key) => {
  //     // Here, we use 'user[0][key]' to represent the single user object
  //     // as the first element of the 'user' array parameter.
  //     params.append(`user[0][${key}]`, userFields[key]);
  //   });
  //   // --- END OF FIX ---

  //   if (userIp) {
  //     params.append("ip", userIp);
  //   }

  //   try {
  //     const response = await fetch(`${MOODLE_URL}/webservice/rest/server.php`, {
  //       method: "POST",
  //       headers: {
  //         "Content-Type": "application/x-www-form-urlencoded",
  //       },
  //       body: params.toString(),
  //     });

  //     const moodleResponse = await response.json();

  //     if (response.ok && moodleResponse && moodleResponse.loginurl) {
  //       console.log("Moodle login URL received:", moodleResponse.loginurl);
  //       window.location.href = moodleResponse.loginurl; // Using window.location.href for client-side redirect
  //     } else if (moodleResponse && moodleResponse.exception) {
  //       console.error("Moodle Web Service Error:", moodleResponse.message);
  //       alert(`Moodle SSO failed: ${moodleResponse.message}`);
  //     } else {
  //       console.error("Unexpected Moodle response:", moodleResponse);
  //       alert("Moodle SSO failed: Unexpected response from server.");
  //     }
  //   } catch (error) {
  //     console.error("Error during Moodle SSO:", error);
  //     alert("Moodle SSO failed due to a network or server error.");
  //   }
  // };

  const newFunc2 = async () => {
    let username = "test_user";
    let password = "^mGYk5dCCJcRSI";
    let errorMessage = "";
    let loading = false;
    let debugInfo = ""; // To show raw responses for debugging

    // --- WARNING: DO NOT USE THESE HARDCODED VALUES OR EXPOSE THEM IN PRODUCTION ---
    const moodleUrl = "http://158.140.178.96:9191"; // e.g., 'http://localhost/moodle' or 'https://your.moodle.site'
    const moodleServiceShortname = "moodle_mobile_app"; // Service for getting the token

    errorMessage = "";
    debugInfo = "";
    loading = true;

    if (!username || !password) {
      errorMessage = "Please enter username and password for the sanity check.";
      loading = false;
      return;
    }

    try {
      // --- Step 1: Get User Web Service Token ---
      const tokenParams = new URLSearchParams({
        username: username,
        password: password,
        service: moodleServiceShortname,
      });

      // This will trigger the CORS issue if you don't disable web security
      const tokenResponse = await fetch(`${moodleUrl}/login/token.php?${tokenParams.toString()}`);

      const tokenData = await tokenResponse.json();
      console.log("tokenData", tokenData);
      debugInfo += "Step 1 (Token) Response:\n" + JSON.stringify(tokenData, null, 2) + "\n\n";

      if (tokenData.error) {
        throw new Error(`Moodle Token Error: ${tokenData.error} (${tokenData.errorcode})`);
      }
      if (!tokenData.token) {
        throw new Error("Moodle did not return a valid token in Step 1.");
      }
      const userToken = tokenData.token;
      const privateToken = tokenData.privatetoken;

      const autologinParams = new URLSearchParams({
        wstoken: userToken,
        privatetoken: privateToken,
        wsfunction: "tool_mobile_get_autologin_key",
        moodlewsrestformat: "json",
      });

      const autologinResponse = await fetch(
        `${moodleUrl}/webservice/rest/server.php?${autologinParams.toString()}`,
      );

      const autologinData = await autologinResponse.json();
      debugInfo +=
        "Step 2 (Auto-Login URL) Response:\n" + JSON.stringify(autologinData, null, 2) + "\n\n";

      if (autologinData.exception) {
        throw new Error(
          `Moodle Autologin Error: ${autologinData.message} (${autologinData.errorcode})`,
        );
      }
      if (!autologinData.autologinurl) {
        throw new Error("Moodle did not return an autologinurl in Step 2.");
      }
      const autologinUrl = autologinData.autologinurl;
      console.log("Autologin URL:", autologinUrl);

      // --- Step 3: Redirect the Browser ---
      console.log("Auto-login successful. Redirecting to:", autologinUrl);
      window.location.href = autologinUrl;
    } catch (error) {
      console.error("Frontend auto-login test failed:", error);
      errorMessage = error.message || "An unexpected error occurred during auto-login test.";
    } finally {
      loading = false;
    }
  };
</script>

<button
  class={!$sidebarLargeClosed
    ? activeUrl === extractPath(href)
      ? classesActive
      : classes
    : "flex w-full py-2 m-0 !my-4 hover:no-underline rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 justify-center"}
  onclick={() => {
    if (href == "/e_learning") {
      newFunc2();
    } else {
      goto(href);
    }
  }}
>
  <div class="flex {$sidebarLargeClosed ? '' : 'gap-3'}">
    <div class="w-max">
      {#if icon}
        <svelte:component this={getIcons(icon)} svgClass={classesSvg} />
      {/if}
    </div>
    <Typography {...typographyData} />
  </div>
  {#if commingSoon}
    <Badge text={$_("soon")} type="brand" />
  {/if}
</button>

<!--
@component
## Config
@prop label: string; // Label of the item.
@prop href: string; // URL of the item.
@prop activeUrl: string; // Current active URL.
@prop classes: string; // Class of the item.
@prop classesActive: string; // Class of the item when active.
@prop icon: string; // Icon of the item.
@prop classesSvg: string; // Class of the icon.
@prop typographyData: object; // Data of Typography object.
-->
