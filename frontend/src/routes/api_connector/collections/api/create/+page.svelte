<script>
  import { call, errorHandling } from "$lib";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { _ } from "svelte-i18n";

  import FormApiConnector from "$lib/pages/FormApiConnector.svelte";
  import FormApiParams from "$lib/pages/FormApiParams.svelte";

  const isFrontend = process.env.NODE_ENV === "development";

  let isBody = $state(false);
  let stateCreate = $state("formApi");

  let tabDatas = $state({
    selected: "API Information",
  });

  let paramsData = $state([]);
  let bodyData = $state([]);

  let data = $state({
    apiName: "",
    description: "",
    urlLink: "",
    method: "GET",
    authorization: "No Auth",
    username: "",
    password: "",
    token: "",
    body_type: "raw",
  });

  const saveApi = async () => {
    try {
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
          body_type: data.body_type,
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

        if (response?.status !== 200) {
          const error = new Error("Failed to create API");
          error.response = response;
          error.status = response?.status;
          throw error;
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
    } catch (error) {
      errorHandling($_, error, "create");
    }
  };
</script>

{#if stateCreate === "formApi"}
  <FormApiConnector
    bind:stateCreate
    bind:data
    bind:paramsData
    bind:bodyData
    bind:isBody
    {tabDatas}
    {saveApi}
  />
{:else if stateCreate === "formParams"}
  <FormApiParams bind:stateCreate bind:paramsData />
{:else if stateCreate === "formBody"}
  <FormApiParams bind:stateCreate bind:data bind:paramsData={bodyData} bind:isBody />
{/if}
