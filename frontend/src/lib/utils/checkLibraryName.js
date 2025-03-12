export const checkLibraryName = async (name, parentId, type, call) => {
  let token = localStorage.getItem("token");
  if (!name || !token) {
    return null;
  }

  try {
    let resp = await call.call("client-dashboard", "is-library-name-available", {
      URLQuery: {},
      URLParams: {},
      Params: {
        name: name,
        parent_id: parseInt(parentId),
        type: type,
      },
      Headers: {
        Authorization: "Bearer " + token,
      },
    });

    console.log("cek resp", resp);

    if (resp === null) {
      throw new Error("Server sibuk");
    }

    if (resp.status === 200) {
      return resp.result.available; // Mengembalikan status (true: tersedia, false: tidak tersedia)
    } else {
      console.error("Error memeriksa name:", resp);
      return null;
    }
  } catch (error) {
    console.error("Terjadi kesalahan:", error);
    return null;
  }
};