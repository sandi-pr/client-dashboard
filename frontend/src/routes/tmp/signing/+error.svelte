<script>
  import { page } from "$app/stores";
  import { Button } from "flowbite-svelte";
</script>

<div class="flex h-screen flex-col items-center justify-center">
  <div class="flex items-center justify-center">
    <p
      class="font-hexa absolute left-12 top-6 ml-4 text-center text-[18px] font-bold text-hexacode-grayone-800 dark:text-hexacode-grayone-100"
    >
      Bank Jasa Jakarta
    </p>
  </div>

  <div class="flex flex-col items-center justify-center">
    <h1
      class="text-center text-4xl font-light text-hexacode-grayone-800 dark:text-hexacode-grayone-100"
    >
      404
    </h1>
    <p
      class="text-center font-light text-hexacode-grayone-800 dark:text-hexacode-grayone-100"
    >
      Page not found
    </p>
    <div class="mt-4">
      <Button href="/" color="light">Go back home</Button>
    </div>
  </div>
</div>
