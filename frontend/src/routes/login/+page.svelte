<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { Label, Footer, Input, Button, Heading, FooterCopyright } from "flowbite-svelte";
  import { goto } from "$app/navigation";
  import bglogin from "$lib/image/bg_login.png";
  import LoadingSpinner from "$lib/components/LoadingSpinner.svelte";
  import HeaderRegister from "$lib/components/HeaderRegister.svelte";

  let email = "";
  let password = "";
  let isLoading = false;

  const onSubmit = async (event) => {
    isLoading = true;
    event.preventDefault();

    try {
      // const response = await fetch(`/api/auth?email=${email}&password=${password}`, {
      //   method: "POST",
      //   headers: {
      //     "Content-Type": "application/json",
      //   },
      // });

      // if (response.status === 200) {
      if (isLoading) {
        sessionStorage.setItem(
          "user",
          JSON.stringify({
            jwt: "123abc",
            data: { email: email, username: "John Doe", brand: "Crossroads" },
          })
        );
        window.location.replace("/");
      } else {
        console.error("Login request failed:", response.status);
      }
    } catch (error) {
      console.error("An unexpected error occurred:", error);
    }
    isLoading = false;
  };
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
