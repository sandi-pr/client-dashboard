<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { format } from "date-fns";
  import { page } from "$app/stores";

  import { call, triggerToast } from "$lib";

  import FormApiConnector from "$lib/pages/FormApiConnector.svelte";
  import FormApiParams from "$lib/pages/FormApiParams.svelte";

  let apiID;

  let stateCreate = $state("formApi");

  let paramsData = $state([
    // {
    //   id: 1,
    //   key: "user_id",
    //   type: "INT",
    //   description: "Unique identifier for the user",
    // },
  ]);

  let bodyData = $state([
    // {
    //   id: 1,
    //   key: "user_id",
    //   type: "INT",
    //   description: "Unique identifier for the user",
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

  async function getApi() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "get-api", {
        URLQuery: {},
        URLParams: {
          id: apiID,
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      // console.log(response);

      if (response.status === 200) {
        data.apiName = response.result.attributes.name;
        data.description = response.result.attributes.description;
        data.urlLink = response.result.attributes.url;
        data.method = response.result.attributes.method;
        data.authorization = response.result.attributes.auth_type;
        data.username = response.result.attributes.username;
        data.password = response.result.attributes.password;
        data.token = response.result.attributes.token;
        data.createOn = format(new Date(response.result.attributes.createdAt), "dd-MM-yyyy");
        paramsData = response.result.attributes.params ? response.result.attributes.params : [];
      }
    }
  }

  async function deleteApi() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "delete-api", {
        URLQuery: {},
        URLParams: {
          id: apiID,
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

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

      if (response.status === 200) {
        let link = await $call.getPage("client-dashboard", "api-connector-collections-api");

        const params = new URLSearchParams(window.location.search);

        params.delete("apiId");

        const queryParams = {};
        for (const [key, value] of params.entries()) {
          queryParams[key] = value;
        }
        const queryString = new URLSearchParams(queryParams).toString();

        window.location.href = link + "?" + queryString;
      }
    }
  }

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
        newData.body_type = "raw";
      }
      // console.log("data", newData);

      let response = await $call.call("client-dashboard", "update-api", {
        URLQuery: {},
        URLParams: {
          id: apiID,
        },
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
    }
  }

  onMount(() => {
    apiID = $page.url.searchParams.get("apiId");
    getApi();
  });
</script>

{#if stateCreate === "formApi"}
  <FormApiConnector
    bind:data
    bind:stateCreate
    bind:paramsData
    bind:bodyData
    bind:apiID
    {saveApi}
    {deleteApi}
  />
{:else if stateCreate === "formParams"}
  <FormApiParams bind:stateCreate bind:paramsData />
{:else if stateCreate === "formBody"}
  <FormApiParams bind:stateCreate bind:paramsData={bodyData} />
{/if}
