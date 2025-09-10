import { triggerToast } from "$lib";

export default (_, error, action = "get") => {
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
