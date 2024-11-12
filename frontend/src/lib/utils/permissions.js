export function check(permission) {
  const user = JSON.parse(localStorage.getItem("user"));
  if (user.authority === "Admin") {
    return true;
  }

  const permissions = JSON.parse(localStorage.getItem("permissions"));

  for (const p of permissions) {
    if (p === permission) {
      return true;
    }
  }

  return false;
  // const user = JSON.parse(localStorage.getItem("user"));
  // return user.authority === role;
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