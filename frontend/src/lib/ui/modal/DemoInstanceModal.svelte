<script>
  import { _ } from "svelte-i18n";
  import { Modal, Button } from "flowbite-svelte";
  import { isDemoTrue, isLimitTrue, getIcons } from "$lib";
  import { Typography } from "$lib";

  const Send01Icon = getIcons("send01");
  const RequestDemoIcon = getIcons("requestDemo");
  const LimitReachedIcon = getIcons("limitReached");
</script>

<Modal
  bind:open={$isDemoTrue}
  autoclose={false}
  title={$_("requestADemo")}
  class="bg-primary-light dark:bg-primary-dark [&>*]:bg-primary-light dark:[&>*]:bg-primary-dark"
>
  <div class="flex flex-col items-center justify-center gap-4">
    <RequestDemoIcon svgClass="text-utility-brand-50-light dark:text-gray-800" />
    <Typography text={$_("requestADemoDesc")} propClass="mx-4 text-center" />
  </div>
  <svelte:fragment slot="footer" class="bg-primary-light dark:bg-primary-dark">
    <div class="w-full flex flex-col items-center gap-2">
      <Button
        class="w-full gap-2"
        onclick={() => {
          window.open("https://hexacode.ai/contact-us/", "_blank", "noopener noreferrer");
          isDemoTrue.set(false);
        }}
      >
        <Send01Icon svgClass="w-5 h-5 text-primary-onbrand-light dark:text-primary-onbrand-dark" />
        <Typography preset="Text md/Semibold" color="primary-onbrand" text={$_("requestADemo")} />
      </Button>
      <Button
        class="w-full gap-2 bg-secondary-light dark:bg-secondary-dark hover:bg-secondary-hover-light dark:hover:bg-secondary-hover-dark"
        color="none"
        onclick={() => {
          isDemoTrue.set(false);
        }}
      >
        <Typography preset="Text md/Semibold" color="fg-secondary" text={$_("cancel")} />
      </Button>
    </div>
  </svelte:fragment>
</Modal>

<Modal
  bind:open={$isLimitTrue}
  autoclose={false}
  title={$_("limitReached")}
  class="bg-primary-light dark:bg-primary-dark [&>*]:bg-primary-light dark:[&>*]:bg-primary-dark"
>
  <div class="flex flex-col items-center justify-center gap-4">
    <LimitReachedIcon svgClass="text-utility-brand-50-light dark:text-gray-800" />
    <Typography text={$_("limitReachedDesc")} propClass="mx-4 text-center" />
  </div>
  <svelte:fragment slot="footer" class="bg-primary-light dark:bg-primary-dark">
    <div class="w-full flex flex-col items-center gap-2">
      <Button
        class="w-full gap-2"
        onclick={() => {
          window.open("https://hexacode.ai/contact-us/", "_blank", "noopener noreferrer");
          isLimitTrue.set(false);
        }}
      >
        <Send01Icon svgClass="w-5 h-5 text-primary-onbrand-light dark:text-primary-onbrand-dark" />
        <Typography preset="Text md/Semibold" color="primary-onbrand" text={$_("requestADemo")} />
      </Button>
      <Button
        class="w-full gap-2 bg-secondary-light dark:bg-secondary-dark hover:bg-secondary-hover-light dark:hover:bg-secondary-hover-dark"
        color="none"
        onclick={() => {
          isLimitTrue.set(false);
        }}
      >
        <Typography preset="Text md/Semibold" color="fg-secondary" text={$_("cancel")} />
      </Button>
    </div>
  </svelte:fragment>
</Modal>
