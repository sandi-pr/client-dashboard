let getTextColor = (presetColor: String) => {
  switch (presetColor) {
    case "primary":
      return "text-primary-light dark:text-primary-dark";
    case "primary-onbrand":
      return "text-primary-onbrand-light dark:text-primary-onbrand-dark";
    case "secondary":
      return "text-secondary-light dark:text-secondary-dark";
    case "secondary-onbrand":
      return "text-secondary-onbrand-light dark:text-secondary-onbrand-dark";
    case "secondary-hover":
      return "text-secondary-hover-light dark:text-secondary-hover-dark";
    case "tertiary":
      return "text-tertiary-light dark:text-tertiary-dark";
    case "tertiary-onbrand":
      return "text-tertiary-onbrand-light dark:text-tertiary-onbrand-dark";
    case "tertiary-hover":
      return "text-tertiary-hover-light dark:text-tertiary-hover-dark";
    case "quarterary":
      return "text-quarterary-light dark:text-quarterary-dark";
    case "quarterary-onbrand":
      return "text-quarterary-onbrand-light dark:text-quarterary-onbrand-dark";
    case "disabled":
      return "text-disabled-light dark:text-disabled-dark";
    case "placeholder":
      return "text-placeholder-light dark:text-placeholder-dark";
    case "placeholder-subtle":
      return "text-placeholder-subtle-light dark:text-placeholder-subtle-dark";
    case "brand-primary":
      return "text-brand-primary-light dark:text-brand-primary-dark";
    case "brand-secondary":
      return "text-brand-secondary-light dark:text-brand-secondary-dark";
    case "brand-tertiary":
      return "text-brand-tertiary-light dark:text-brand-tertiary-dark";
    case "brand-tertiary-alt":
      return "text-brand-tertiary-alt-light dark:text-brand-tertiary-alt-dark";
    case "success-primary":
      return "text-success-primary-light dark:text-success-primary-dark";
    case "warning-primary":
      return "text-warning-primary-light dark:text-warning-primary-dark";
    case "error-primary":
      return "text-error-primary-light dark:text-error-primary-dark";
    default:
      return "text-primary-light dark:text-primary-dark";
  }
};

export default getTextColor;