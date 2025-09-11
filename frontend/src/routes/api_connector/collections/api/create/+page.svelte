<script>
  import { call } from "$lib";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { _ } from "svelte-i18n";

  import FormApiConnector from "$lib/pages/FormApiConnector.svelte";
  import FormApiParams from "$lib/pages/FormApiParams.svelte";

  const isFrontend = process.env.NODE_ENV === "development";

  let stateCreate = $state("formApi");

  let tabDatas = $state({
    selected: "API Information",
  });

  let paramsData = $state([
    // {
    //   id: 1,
    //   key: "user_id",
    //   type: "INT",
    //   description: "Unique identifier for the user",
    //   actions: {
    //     objects: [
    //       {
    //         onclick: deleteParam(1),
    //         name: "Delete",
    //       icon: "trash",
    //       },
    //     ],
    //   },
    // },
  ]);

  let bodyData = $state([
    // {
    //   id: 1,
    //   key: "user_id",
    //   type: "INT",
    //   description: "Unique identifier for the user",
    //   actions: {
    //     objects: [
    //       {
    //         onclick: deleteParam(1),
    //         name: "Delete",
    //       icon: "trash",
    //       },
    //     ],
    //   },
    // },
  ]);

  let data = $state({
    apiName: "",
    description: "",
    urlLink: "",
    method: "GET",
    authorization: "No Auth",
    username: "",
    password: "",
    token: "",
  });

  async function saveApi() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let newData = {
        name: data.apiName,
        description: data.description,
        url: data.urlLink,
        method: data.method,
        auth_type: data.authorization,
        username: data.username,
        password: data.password,
        token: data.token,
      };

      let collectionId = $page.url.searchParams.get("collectionId");
      newData.parent_id = parseInt(collectionId);
      if (paramsData.length > 0) {
        newData.params = JSON.stringify(paramsData);
      }
      if (bodyData.length > 0) {
        newData.body = JSON.stringify(bodyData);
      }
      // console.log("data", newData);

      let response = await $call.call("client-dashboard", "create-api", {
        URLQuery: {},
        URLParams: {},
        Params: newData,
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        if (response.status === 400) {
          console.error("Invalid parameters:", response);
          triggerToast($_("invalidParameters"), true, "error");
        } else if (response.status === 401) {
          console.error("Unauthorized access:", response);
          triggerToast($_("unauthorized"), true, "error");
        } else if (response.status === 500) {
          console.error("Server error:", response);
          triggerToast($_("serverError"), true, "error");
        } else {
          console.error("Unexpected status code:", response);
          triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
        }
        return;
      }

      // console.log(response);

      let link = isFrontend
        ? "/api_connector/collections/api"
        : await $call.getPage("client-dashboard", "api-connector-collections-api");

      const params = new URLSearchParams(window.location.search);

      const queryParams = {};
      for (const [key, value] of params.entries()) {
        queryParams[key] = value;
      }
      const queryString = new URLSearchParams(queryParams).toString();

      goto(link + "?" + queryString);
    }
  }
</script>

{#if stateCreate === "formApi"}
  <FormApiConnector bind:stateCreate {tabDatas} bind:paramsData bind:bodyData {saveApi} bind:data />
{:else if stateCreate === "formParams"}
  <FormApiParams bind:stateCreate bind:paramsData />
{:else if stateCreate === "formBody"}
  <FormApiParams bind:stateCreate bind:paramsData={bodyData} />
{/if}
