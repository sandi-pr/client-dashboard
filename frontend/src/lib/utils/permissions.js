import { isBakrie } from "./bakrieLibraryStore";

let isBakrieTmp = false;

isBakrie.subscribe((isBakrie) => {
  if (isBakrie) {
    isBakrieTmp = isBakrie;
  }
});

export function check(permission) {
  const user = JSON.parse(localStorage.getItem("user"));
  if (user.authority === "Admin") {
    return true;
  }

  if (isBakrieTmp) {
    if (user.authority === "User") {
      if (permission.includes("view")) {
        return true;
      } else {
        return false
      }
    }
  }
  const permissions = JSON.parse(localStorage.getItem("permissions"));

  for (const p of permissions) {
    if (p === permission) {
      return true;
    }
  }

  return false;
}

export function checkAnd(roles) {
  let ret = true;
  for (let p of roles) {
    ret = ret && check(p);
  }
  return ret;
}

export function checkOr(roles) {
  let ret = false;
  for (let p of roles) {
    ret = ret || check(p);
    // console.log(p, check(p));
  }

  return ret;
}