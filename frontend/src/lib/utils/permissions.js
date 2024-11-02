export function check(role) {
  const user = JSON.parse(localStorage.getItem("user"));
  return user.authority === role;
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