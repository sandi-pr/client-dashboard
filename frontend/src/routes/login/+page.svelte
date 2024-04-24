<script>
  import { Label, Footer, Input, Button, Heading, FooterCopyright } from "flowbite-svelte";

  import { _ } from "svelte-i18n";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import bglogin from "../lib/image/bg_login.png";
  import LoadingSpinner from "../lib/components/LoadingSpinner.svelte";
  import HeaderRegister from "../lib/components/HeaderRegister.svelte";
  import Me from "../lib/utils/me";

  let email = "";
  let password = "";
  let isLoading = false;

  const onSubmit = async (event) => {
    isLoading = true;
    event.preventDefault();

    try {
      const response = await fetch(`/api/auth?email=${email}&password=${password}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        const data = await response.json();
        $Me = data;
        setTimeout(() => {
          isLoading = false;
          goto("/dashboard");
        }, 5000);
        // const result = await response.json();
        // const response2 = await fetch(`/api/auth?jwt=${result.jwt}`, {
        //   method: "GET",
        //   headers: {
        //     "Content-Type": "application/json",
        //   },
        // });
        // if (response2.ok) {
        //   // const result = await response2.json();
        //   // console.log(result);
        //   // if (result.organization !== null) {
        //   //   console.log(data);
        //   //   data.data.user = result;
        //   //   document.cookie = `crd_cookies=${data.data}; path=/; max-age=${60 * 60 * 24 * 7}`;
        //   //   goto("/dashboard");
        //   // } else {
        //   //   goto("/onboarding");
        //   // }
        // }
      } else {
        console.error("Login request failed:", response.status);
      }
    } catch (error) {
      console.error("An unexpected error occurred:", error);
    }
  };

  onMount(() => {
    // console.log(data);
  });
</script>

<LoadingSpinner {isLoading} />
<div class="flex">
  <div class="flex flex-col w-full lg:w-1/2 p-8 justify-between h-screen">
    <HeaderRegister />
    <div class="flex flex-col items-center justify-center h-full w-full">
      <form class="w-2/3" on:submit={onSubmit}>
        <Heading
          tag="h3"
          class="mb-4 font-semibold text-primary-light dark:text-primary-dark leading-[38px]"
          >{$_("signin")}</Heading
        >
        <div class="space-y-6">
          <div class="">
            <Label for="email" class="mb-2">{$_("companyemail")}</Label>
            <Input
              type="email"
              bind:value={email}
              id="email"
              placeholder={$_("placholderCompanyEmail")}
              required
            />
          </div>
          <div class="">
            <Label for="password" class="mb-2">{$_("password")}</Label>
            <Input
              type="password"
              bind:value={password}
              id="password"
              placeholder={$_("enterPassword")}
              required
            />
          </div>
        </div>
        <Button color="purple" type="submit" class="w-full mt-4">{$_("continue")}</Button>
        <!-- <Button color="light" class="w-full mb-4">
          <img src={windows} alt="windows" class="mr-2" />Continue with Active Directory
        </Button> -->
        <div class="flex w-full justify-center space-x-2">
          <!-- <P>Don’t have an account?</P><A color="text-purple-600" href="/signup">Sign up</A> -->
        </div>
      </form>
    </div>
    <Footer class="flex">
      <FooterCopyright by="Crossroad" year={2024} />
    </Footer>
  </div>

  <div class="w-1/2 h-screen hidden lg:block">
    <img
      class="w-full h-full object-cover rounded-tl-[80px] rounded-bl-[80px]"
      src={bglogin}
      alt="bg_login"
    />
  </div>
</div>
