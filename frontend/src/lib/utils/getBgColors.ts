let getBgColor = (presetColor: String) => {
  switch (presetColor) {
    case "primary":
      return "bg-primary-light dark:bg-primary-dark";
    case "primary-alt":
      return "bg-primary-alt-light dark:bg-primary-alt-dark";
    case "primary-hover":
      return "bg-primary-hover-light dark:bg-primary-hover-dark";
    case "primary-solid":
      return "bg-primary-solid-light dark:bg-primary-solid-dark";
    case "secondary":
      return "bg-secondary-light dark:bg-secondary-dark";
    case "secondary-alt":
      return "bg-secondary-alt-light dark:bg-secondary-alt-dark";
    case "secondary-hover":
      return "bg-secondary-hover-light dark:bg-secondary-hover-dark";
    case "secondary-subtle":
      return "bg-secondary-subtle-light dark:bg-secondary-subtle-dark";
    case "secondary-solid":
      return "bg-secondary-solid-light dark:bg-secondary-solid-dark";
    case "tertiary":
      return "bg-tertiary-light dark:bg-tertiary-dark";
    case "quarterary":
      return "bg-quarterary-light dark:bg-quarterary-dark";
    case "active":
      return "bg-active-light dark:bg-active-dark";
    case "disabled":
      return "bg-disabled-light dark:bg-disabled-dark";
    case "disabled-subtle":
      return "bg-disabled-subtle-light dark:bg-disabled-subtle-dark";
    case "overlay":
      return "bg-overlay-light dark:bg-overlay-dark";
    case "brand-primary":
      return "bg-brand-primary-light dark:bg-brand-primary-dark";
    case "brand-secondary":
      return "bg-brand-secondary-light dark:bg-brand-secondary-dark";
    case "brand-solid":
      return "bg-brand-solid-light dark:bg-brand-solid-dark";
    case "brand-section":
      return "bg-brand-section-light dark:bg-brand-section-dark";
    case "success-primary":
      return "bg-success-primary-light dark:bg-success-primary-dark";
    case "success-secondary":
      return "bg-success-secondary-light dark:bg-success-secondary-dark";
    case "success-solid":
      return "bg-success-solid-light dark:bg-success-solid-dark";
    case "warning-primary":
      return "bg-warning-primary-light dark:bg-warning-primary-dark";
    case "error-primary":
      return "bg-error-primary-light dark:bg-error-primary-dark";
    default:
      return "bg-primary-light dark:bg-primary-dark";
  }
};

export default getBgColor;