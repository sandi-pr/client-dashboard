import { browser } from '$app/environment';
import { init, register, getLocaleFromNavigator } from 'svelte-i18n';

const defaultLocale = 'en';

register('en', () => import('./locales/en.json'));
register('id', () => import('./locales/id.json'));

init({
  fallbackLocale: defaultLocale,
  // initialLocale: browser ? getLocaleFromNavigator() : defaultLocale,
  initialLocale: defaultLocale,
});
