export function toCamelCase(str) {
  return str
    .replace(/\s(\w)/g, function (match, p1) {
      return p1.toUpperCase();
    })
    .replace(/_./g, function (match) {
      return match.charAt(1).toUpperCase();
    })
    .replace(/\s/g, "")
    .replace(/^(.)/, function (match, p1) {
      return p1.toLowerCase();
    });
}