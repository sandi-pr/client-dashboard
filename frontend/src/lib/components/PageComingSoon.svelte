<script>
  import { Typography } from "$lib";
  import ComingSoon from "$lib/image/coming_soon.svg";

  const location = globalThis.location.pathname;
  const lastSegment = location.split("/").filter(Boolean).pop(); // Get the last non-empty segment
  const capitalizedSegment = lastSegment
    ? lastSegment.charAt(0).toUpperCase() + lastSegment.slice(1)
    : "Home"; // Capitalize the first letter, fallback to "Home" if empty
</script>

<div
  class="flex flex-col items-center justify-center bg-primary-light dark:bg-primary-dark w-full h-screen"
>
  <img src={ComingSoon} alt="coming soon" class="w-[20rem] h-[20rem] md:w-[30rem] md:h-[30rem]" />

  <Typography
    text="Coming Soon {capitalizedSegment} Features"
    preset="Display xl/semibold"
    propClass="m-8 fluid-xl md:fluid-2xl"
  />
  <Typography
    color="tertiary"
    text="Get ready for smarter workflows, automation, and insights—all in one Hexa AI platform. Stay tuned for the launch!"
    preset="Text xl/regular"
    propClass="text-center"
  />
</div>
