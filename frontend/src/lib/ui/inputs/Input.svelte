<script lang="ts">
  import { onMount } from 'svelte';
  import InputNumber from './InputNumber.svelte';
  import InputText from './InputText.svelte';
  import InputSelect from './InputSelect.svelte';
  import InputCheckbox from './InputCheckbox.svelte';
  import InputRadio from './InputRadio.svelte';
  import InputDate from './InputDate.svelte';
  import InputTime from './InputTime.svelte';
  import InputFile from './InputFile.svelte';
  import InputSearch from './InputSearch.svelte';
  import InputTextArea from './InputTextArea.svelte';
  import InputEmail from './InputEmail.svelte';
  import InputPassword from './InputPassword.svelte';
  import { Helper, Label } from 'flowbite-svelte';

  let {
    label,
    description,
    required,
    default: defaultValue,
    input_type,
    positionLabel = 'left',
    colorLabel = 'black',
    value,
    placeholder,
    options,
    readonly,
    maxSize,
  }: {
    label?: string;
    description?: string;
    required?: boolean;
    default?: any;
    input_type: string;
    positionLabel?: 'left' | 'top';
    colorLabel?: string;
    value?: any;
    placeholder?: string;
    options?: any[];
    readonly?: boolean;
    maxSize?: string;
  } = $props();

  let positionLabelState = positionLabel === 'top' ? '' : 'grid grid-cols-1 sm:grid-cols-5 gap-4';

  onMount(() => {
    value = value || defaultValue;
  });
</script>

<div class={`flex flex-col w-full h-full ${positionLabelState}`}>
  <div class={`${positionLabel === 'top' ? 'mb-2' : 'flex items-center'}`}>
    <Label class={`display-block text-${colorLabel}-600 text-wrap`}>{label}{required ? '*' : ''}</Label>
  </div>
  <div class={`${positionLabel === 'top' ? '' : 'col-span-4'}`}>
    {#if input_type === 'text'}
      <InputText bind:value={value} placeholder={placeholder} required={required} readonly={readonly} />
    {/if}
    {#if input_type === 'number'}
      <InputNumber bind:value={value} placeholder={placeholder} required={required} readonly={readonly} />
    {/if}
    {#if input_type === 'select'}
      <InputSelect bind:value={value} options={options} placeholder={placeholder} />
    {/if}
    {#if input_type === 'checkbox'}
      <InputCheckbox bind:checked={value} required={required} />
    {/if}
    {#if input_type === 'radio'}
      <InputRadio bind:group={value} options={options} />
    {/if}
    {#if input_type === 'date'}
      <InputDate bind:value={value} readonly={readonly} />
    {/if}
    {#if input_type === 'time'}
      <InputTime bind:value={value} readonly={readonly} />
    {/if}
    {#if input_type === 'file'}
      <InputFile bind:value={value} placeholder={placeholder} maxSize={maxSize} />
    {/if}
    {#if input_type === 'search'}
      <InputSearch bind:value={value} placeholder={placeholder} />
    {/if}
    {#if input_type === 'textarea'}
      <InputTextArea bind:value={value} placeholder={placeholder} required={required} readonly={readonly} />
    {/if}
    {#if input_type === 'email'}
      <InputEmail bind:value={value} placeholder={placeholder} required={required} readonly={readonly} />
    {/if}
    {#if input_type === 'password'}
      <InputPassword bind:value={value} placeholder={placeholder} required={required} readonly={readonly} />
    {/if}
  </div>
  {#if positionLabel !== 'top'}
    <div class="hidden sm:block col-span-1"></div>
  {/if}
  <Helper class="display-block w-full sm:w-max mt-1 sm:mt-0">{description}</Helper>
</div>

<!--
@component

## `<Input />`

- @prop {string} label
- @prop {string} description
- @prop {boolean} required
- @prop {any} default
- @prop {string} input_type
- @prop {string} positionLabel
- @prop {string} colorLabel
- @prop {any} value
- @prop {string} placeholder
- @prop {any[]} options if input_type is 'select' { value: string; name: string; }[] or radio { value: string; }[]
- @prop {boolean} readonly
-->