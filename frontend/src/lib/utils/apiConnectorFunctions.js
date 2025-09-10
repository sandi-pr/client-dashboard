import { triggerToast } from "$lib";

const isFrontend = process.env.NODE_ENV === "development";

export const createCollectionAPIConnector = async (
  _,
  call,
  collectionName,
  collectionDescription,
) => {
  try {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await call.call("client-dashboard", "create-collection-api", {
        URLQuery: {},
        URLParams: {},
        Params: {
          name: collectionName,
          description: collectionDescription,
          parent_id: 0,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        const error = new Error("Failed to create collection");
        error.status = response.status;
        error.response = await response;
        throw error;
      }

      triggerToast(_("successCreateCollection"), true);

      collectionName = "";
      collectionDescription = "";

      return true;
    }
  } catch (error) {
    errorHandling(error, "create");
    return false;
  }
};

export const getCollectionAPIConnector = async (
  _,
  call,
  search = "",
  pagination = { page: 1, pageSize: 8 },
) => {
  try {
    let token = localStorage.getItem("token");
    let collections = [];

    if (token !== null) {
      let response = await call.call("client-dashboard", "get-collection-apis", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: pagination.page,
          page_size: pagination.pageSize,
          parent_id: 0,
          search: search,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        const error = new Error("Failed to fetch collections");
        error.status = response.status;
        error.response = await response;
        throw error;
      }

      if (response.result.data.length > 0) {
        collections = collections.concat(response.result.data);
      }

      let link = isFrontend
        ? "/api_connector/collections/api"
        : await call.getPage("client-dashboard", "api-connector-collections-api");

      let promise = collections.map(async (collection) => {
        let responseLibrary = await call.call("client-dashboard", "get-status-library", {
          URLQuery: {},
          URLParams: {
            id: collection.id.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (responseLibrary.status !== 200) {
          const error = new Error("Failed to fetch library status");
          error.status = responseLibrary.status;
          error.response = await responseLibrary;
          throw error;
        }

        return {
          id: collection.id,
          collectionName: collection.attributes.name,
          description: collection.attributes?.description || "",
          type: collection.attributes.type,
          href: link + "?collectionId=" + collection.id,
          parent_id: collection.attributes.parent_id,
        };
      });

      return await Promise.all(promise);
    }
  } catch (error) {
    errorHandling(error, "get");
    return [];
  }
};

export const updateCollectionAPIConnector = async (
  _,
  call,
  parentId,
  collectionId,
  collectionName,
  collectionDescription,
) => {
  try {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await call.call("client-dashboard", "update-collection-api", {
        URLQuery: {},
        URLParams: {
          id: collectionId.toString(),
        },
        Params: {
          name: collectionName,
          description: collectionDescription,
          parent_id: parentId ? parseInt(parentId) : 0,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        const error = new Error("Failed to update collection");
        error.status = response.status;
        error.response = await response;
        throw error;
      }

      parentId = "";
      collectionId = "";
      collectionName = "";
      collectionDescription = "";

      triggerToast(_("successUpdateCollection"), true);
      return true;
    }
  } catch (error) {
    errorHandling(error, "update");
    return false;
  }
};

export const deleteCollectionAPIConnector = async (_, call, collectionId) => {
  try {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await call.call("client-dashboard", "delete-collection-api", {
        URLQuery: {},
        URLParams: {
          id: collectionId.toString(),
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        const error = new Error("Failed to delete collection");
        error.status = response.status;
        error.response = await response;
        throw error;
      }

      collectionId = "";

      triggerToast(_("successfullyDeleted"), true);

      return true;
    }
  } catch (error) {
    errorHandling(error, "delete");
    return false;
  }
};

const errorHandling = (error, action = "get") => {
  if (action === "get") {
    console.error("Error during fetching collections:", error);
  } else if (action === "create") {
    console.error("Error during collection creation:", error);
  } else if (action === "update") {
    console.error("Error during updating collection:", error);
  } else if (action === "delete") {
    console.error("Error during deleting collection:", error);
  }

  if (error?.status) {
    if (error.status === 400) {
      triggerToast(_("invalidParameters"), true, "error");
    } else if (error.status === 401) {
      triggerToast(_("unauthorized"), true, "error");
    } else if (error.status === 500) {
      triggerToast(_("serverError"), true, "error");
    } else {
      triggerToast(`${_("unexpectedError")}: ${error?.status}`, true, "error");
    }
  } else {
    if (action === "get") {
      triggerToast(_("failedFetchCollections"), true, "error");
    } else if (action === "create") {
      triggerToast(_("failedCreateCollection"), true, "error");
    } else if (action === "update") {
      triggerToast(_("failedUpdateCollection"), true, "error");
    } else if (action === "delete") {
      triggerToast(_("failedDeleteCollection"), true, "error");
    } else {
      triggerToast(`${_("unexpectedError")}: ${error?.message}`, true, "error");
    }
  }
};
