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
    case "fg-primary":
      return "text-fg-primary-light dark:text-fg-primary-dark";
    case "fg-secondary":
      return "text-fg-secondary-light dark:text-fg-secondary-dark";
    case "fg-secondary-hover":
      return "text-fg-secondary-hover-light dark:text-fg-secondary-hover-dark";
    case "fg-tertiary":
      return "text-fg-tertiary-light dark:text-fg-tertiary-dark";
    case "fg-tertiary-hover":
      return "text-fg-tertiary-hover-light dark:text-fg-tertiary-hover-dark";
    case "fg-quarterary":
      return "text-fg-quarterary-light dark:text-fg-quarterary-dark";
    case "fg-quarterary-hover":
      return "text-fg-quarterary-hover-light dark:text-fg-quarterary-hover-dark";
    case "fg-quinary":
      return "text-fg-quinary-light dark:text-fg-quinary-dark";
    case "fg-quinary-hover":
      return "text-fg-quinary-hover-light dark:text-fg-quinary-hover-dark";
    case "fg-senary":
      return "text-fg-senary-light dark:text-fg-senary-dark";
    case "fg-white":
      return "text-fg-white";
    case "fg-disabled":
      return "text-fg-disabled-light dark:text-fg-disabled-dark";
    case "fg-disabled-subtle":
      return "text-fg-disabled-subtle-light dark:text-fg-disabled-subtle-dark";
    case "fg-brand-primary":
      return "text-fg-brand-primary-light dark:text-fg-brand-primary-dark";
    case "fg-brand-primary-alt":
      return "text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark";
    case "fg-brand-secondary":
      return "text-fg-brand-secondary-light dark:text-fg-brand-secondary-dark";
    case "fg-error-primary":
      return "text-fg-error-primary-light dark:text-fg-error-primary-dark";
    case "fg-error-secondary":
      return "text-fg-error-secondary-light dark:text-fg-error-secondary-dark";
    case "fg-warning-primary":
      return "text-fg-warning-primary-light dark:text-fg-warning-primary-dark";
    case "fg-warning-secondary":
      return "text-fg-warning-secondary-light dark:text-fg-warning-secondary-dark";
    case "fg-success-primary":
      return "text-fg-success-primary-light dark:text-fg-success-primary-dark";
    case "fg-success-secondary":
      return "text-fg-success-secondary-light dark:text-fg-success-secondary-dark";
    default:
      return "text-primary-light dark:text-primary-dark";
  }
};

export default getTextColor;
