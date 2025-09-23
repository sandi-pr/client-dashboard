import { triggerToast } from "$lib";

const errorMessages = {
  get: {
    collection: "failedFetchCollections",
    case: "failedFetchCases",
  },
  create: {
    collection: "failedCreateCollection",
    case: "failedCreateCase",
  },
  update: {
    collection: "failedUpdateCollection",
    case: "failedUpdateCase",
  },
  delete: {
    collection: "failedDeleteCollection",
    case: "failedDeleteCase",
  },
};

export default (_, error, action = "get", target = "collection") => {
  if (action === "get") {
    console.error(`Error during fetching ${target}s:`, error);
  } else if (action === "create") {
    console.error(`Error during ${target} creation:`, error);
  } else if (action === "update") {
    console.error(`Error during updating ${target}:`, error);
  } else if (action === "delete") {
    console.error(`Error during deleting ${target}:`, error);
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
      triggerToast(_(errorMessages.get[target]), true, "error");
    } else if (action === "create") {
      triggerToast(_(errorMessages.create[target]), true, "error");
    } else if (action === "update") {
      triggerToast(_(errorMessages.update[target]), true, "error");
    } else if (action === "delete") {
      triggerToast(_(errorMessages.delete[target]), true, "error");
    } else {
      triggerToast(`${_("unexpectedError")}: ${error?.message}`, true, "error");
    }
  }
};
