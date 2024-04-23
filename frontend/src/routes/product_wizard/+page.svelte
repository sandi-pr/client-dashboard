<script>
  import { onMount } from "svelte";
  import {
    P,
    Card,
    Radio,
    Label,
    Input,
    // Modal,
    Button,
    Heading,
    Textarea,
    // Dropdown,
    // DropdownItem,
  } from "flowbite-svelte";
  import { goto } from "$app/navigation";
  // import ColorPicker, { ChromeVariant } from "svelte-awesome-color-picker";

  import upload from "../lib/icon/upload-cloud-02.svg";
  import Timeline from "../lib/components/Timeline.svelte";
  import { timelineStep } from "../lib/utils/timelineStep.js";
  import { userDataSession } from "../lib/utils/auth";
  import check from "../lib/icon/check-circle.svg";
  import MarketIcon from "../lib/icon/building-02.svg";
  import OfficeIcon from "../lib/icon/building-06.svg";
  import BankIcon from "../lib/icon/bank.svg";
  import UserIcon from "../lib/icon/user-03.svg";
  import PlusIcon from "../lib/icon/plus.svelte";
  import UsersRightIcon from "../lib/icon/users-right.svelte";
  import File04Icon from "../lib/icon/file-04.svelte";
  import BoxIcon from "../lib/icon/box.svelte";
  import BarChartSquareIcon from "../lib/icon/bar-chart-square-01.svelte";
  // import ArrowUpRightDownLeftOutline from "flowbite-svelte-icons";
  import { ChevronDownOutline } from "flowbite-svelte-icons";

  // import ProductCard from "../lib/components/ProductCard.svelte";

  onMount(async () => {
    // const isOrganizationExist = await checkOrganization();

    // if (isOrganizationExist) {
    //   console.log("udah ada");
    //   window.history.back();
    // }

    return () => {
      unsubscribe();
      uploadedImage = null;
      uploadedImageName = null;
      fileImage = null;
      companyName = "";
      companyDescription = "";

      uploadedFile = null;
      uploadedFileName = null;
      uploadedFileSize = null;

      csvData = [];
      emailList = [];
    };
  });

  let hex = "#7E56D8";

  let uploadedImage = null;
  let uploadedImageName = null;
  let fileImage = null;
  let companyName = "";
  let companyDescription = "";

  let uploadedFile = null;
  let uploadedFileName = null;
  let uploadedFileSize = null;

  let csvData = [];
  let emailList = [];
  let rows = [];

  const handleImageUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = "image/*";

    inputElement.addEventListener("change", (event) => {
      const file = event.target.files[0];

      if (file && isValidImage(file)) {
        const reader = new FileReader();

        reader.onload = () => {
          fileImage = file;
          uploadedImageName = file.name;
          uploadedImage = reader.result;
        };

        reader.readAsDataURL(file);
      } else {
        alert("Invalid file. Please select a valid image file (max. 800x400px).");
      }
    });
    inputElement.click();
  };

  const isValidImage = (image) => {
    const allowedTypes = ["image/svg+xml", "image/png", "image/jpeg", "image/gif"];
    const maxFileSize = 800 * 400;
    return allowedTypes.includes(image.type) && image.size <= maxFileSize;
  };

  const titleHeading = "Product Wizard";

  let currentStep = 1;
  const steps = ["1", "2", "3"];

  const unsubscribe = timelineStep.subscribe((value) => {
    currentStep = value;
  });

  let product_wizard = true;
  let done = false;

  const endStep = async () => {
    try {
      var userData = userDataSession();
      var jwt = userData.jwt;
      var data = {
        name: companyName,
        email: userData.user.email,
        description: companyDescription,
        users: [userData.user.id],
      };
      console.log(userData.user.id);
      try {
        var formData = new FormData();
        var data2 = JSON.stringify(data);

        formData.append("data", data2);
        formData.append("files.photo", fileImage);

        try {
          const response = await fetch(`/api/organization?jwt=${jwt}`, {
            method: "POST",
            body: formData,
          });

          if (response.ok) {
            const responseData = await response.json();
            console.log("Success create organization: ", responseData);
            var organization_id = responseData.data.id;

            const unitParentMap = {};
            const unitIds = {};

            for (const row of rows) {
              const unit = row.unit;
              const parentUnit = row.parentunit;

              unitParentMap[unit] = parentUnit;

              if (parentUnit === "null" || parentUnit === "Null" || parentUnit === "NULL") {
                await createAndSaveUnit(row, organization_id, jwt, unitIds);
              }
            }

            for (const unit in unitParentMap) {
              const parentUnit = unitParentMap[unit];

              if (parentUnit !== "null" || parentUnit !== "Null" || parentUnit !== "NULL") {
                const row = {
                  unit: unit,
                  parentunit: parentUnit,
                };

                await createAndSaveUnit(row, organization_id, jwt, unitIds);
              }
            }

            const roleIds = {};
            for (const row of rows) {
              await createAndSaveRole(row, jwt, roleIds);
            }

            for (const row of rows) {
              await createAndSaveUser(row, organization_id, jwt, unitIds, roleIds);
            }

            done = true;
          } else {
            const errorData = await response.json();
            console.error("Error Data:", errorData);
          }
        } catch (error) {
          console.error("An unexpected error occurred:", error);
        }
      } catch (error) {
        console.error("Error finishing step:", error);
      }
    } catch (error) {
      console.error("Error finishing step:", error);
      alert("An error occurred. Please try again later.");
    }
  };

  const createAndSaveUnit = async (row, organizationId, jwt, unitIds) => {
    const unitName = row.unit;
    const parentUnitId = unitIds[row.parentunit] || null;

    if (!unitIds.hasOwnProperty(unitName)) {
      const data = {
        name: unitName,
        organization_role: unitName,
        organization: organizationId,
        unit: parentUnitId,
      };

      const formData = new FormData();
      const dataJSON = JSON.stringify(data);

      formData.append("data", dataJSON);

      try {
        const response = await fetch(`/api/units?jwt=${jwt}`, {
          method: "POST",
          body: formData,
        });

        if (response.ok) {
          const responseData = await response.json();
          console.log(`Success create unit: ${unitName}, ID: ${responseData.data.id}`);
          unitIds[unitName] = responseData.data.id;
        } else {
          const errorData = await response.json();
          console.error("Error Data:", errorData);
        }
      } catch (e) {
        console.log(`Error create unit (${unitName}): `, e);
      }
    }
  };

  const createAndSaveRole = async (row, jwt, roleIds) => {
    const roleName = row.role;
    const unitshead = row.unitshead;

    if (!roleIds.hasOwnProperty(roleName)) {
      const data = {
        name: roleName,
        admin: roleName.toLowerCase() === "admin" ? true : false,
        unit_head: unitshead.toLowerCase() === "true" ? true : false,
      };

      const formData = new FormData();
      const dataJSON = JSON.stringify(data);

      formData.append("data", dataJSON);

      try {
        const response = await fetch(`/api/roles?jwt=${jwt}`, {
          method: "POST",
          body: formData,
        });

        if (response.ok) {
          const responseData = await response.json();
          console.log(`Success create role: ${roleName}, ID: ${responseData.data.id}`);
          roleIds[roleName] = responseData.data.id;
        } else {
          const errorData = await response.json();
          console.error("Error Data:", errorData);
        }
      } catch (e) {
        console.log(`Error create role (${roleName}): `, e);
      }
    }
  };

  function generateUniqueNumber() {
    return Math.floor(Math.random() * 10000);
  }

  function generateRandomPassword(length) {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    let password = "";
    for (let i = 0; i < length; i++) {
      const randomIndex = Math.floor(Math.random() * charset.length);
      password += charset[randomIndex];
    }
    return password;
  }

  const createAndSaveUser = async (row, organizationId, jwt, unitIds, roleIds) => {
    const uniqueNumber = generateUniqueNumber();
    const randomPassword = generateRandomPassword(10);
    const data = {
      username: `${row.email.split("@")[0].replace(/[^a-zA-Z0-9]/g, "")}${uniqueNumber}`,
      password: randomPassword,
      email: row.email,
      fullname: row.nama,
      organization: organizationId,
      unit: unitIds[row.unit] || null,
      user_roles: roleIds[row.role] || null,
      confirmed: true,
      role: 1,
    };

    const formData = new FormData();
    const dataJSON = JSON.stringify(data);

    formData.append("data", dataJSON);

    try {
      const response = await fetch(`/api/users?jwt=${jwt}`, {
        method: "POST",
        body: formData,
      });

      if (response.ok) {
        const responseData = await response.json();
        console.log(`Success create user: ${row.nama}, ID: ${responseData.id}`);
        // await emailBlast(row.email, randomPassword);
      } else {
        const errorData = await response.json();
        console.error("Error Data:", errorData);
      }
    } catch (e) {
      console.log(`Error create user (${row.username}): `, e);
    }
  };

  const validateCurrentStep = () => {
    if ((currentStep === 1 && !companyName.trim()) || uploadedImage === null) {
      return false;
    }
    if (currentStep === 2 && uploadedFile === null) {
      return false;
    }
    return true;
  };

  const handlesuccess = () => {
    goto("/dashboard");
  };

  let CompanyType = [
    {
      icon: MarketIcon,
      moduleName: "Start-ups and Small Businesses (1-50 employees)",
      desc: "Gain a competitive edge with tools that help you manage your growing business efficiently",
      option: false,
    },
    {
      icon: OfficeIcon,
      moduleName: "Medium Businesses (51-500 employees)",
      desc: "Enhance coordination and productivity across various departments and teams.",
      option: false,
    },
    {
      icon: BankIcon,
      moduleName: "Large Enterprises  (>500 employees)",
      desc: "Maximize your impact with streamlined operations and resource management.",
      option: false,
    },
    {
      icon: UserIcon,
      moduleName: "Freelancers and Consultants",
      desc: "Organize your projects, finances, and client interactions effectively in one place.",
      option: false,
    },
  ];

  // let moduleClient = [
  //   {
  //     icon: MarketIcon,
  //     moduleName: "Inventory System",
  //     desc: "Manage your inventory with our comprehensive module",
  //     option: false,
  //   },
  //   {
  //     icon: OfficeIcon,
  //     moduleName: "Forecasting System",
  //     desc: "Predict your sales",
  //     option: false,
  //   },
  //   {
  //     icon: BankIcon,
  //     moduleName: "Document Management",
  //     desc: "Organize your files",
  //     option: false,
  //   },
  //   {
  //     icon: UserIcon,
  //     moduleName: "CRM",
  //     desc: "Improve customer relationships with our CRM module",
  //     option: false,
  //   },
  // ];

  let popupModal = false;
  // let problemsFeature = [];
  // let title;
  // let desc;

  const companyCategories = [
    "Technology",
    "Healthcare",
    "Manufacturing",
    "Education",
    "Finance",
    "Retail",
    "Automotive",
    "Energy",
    "Telecommunications",
    "Entertainment",
    "Agriculture",
    "Construction",
    "Real Estate",
  ];
  let dropdownOpen = false;
  let selectedIndustry = "";
  // const dispatch = createEventDispatcher();

  function handleDropdownSelect(value) {
    selectedIndustry = value;
    // dispatch('industrySelected', { industry: selectedIndustry }); // Memancarkan event jika nilai berubah
    dropdownOpen = false;
  }

  // const addProblemFeature = () => {
  //   if (title.trim() !== "" && desc.trim() !== "") {
  //     problemsFeature = [...problemsFeature, { title: title.trim(), desc: desc.trim() }];
  //     title = "";
  //     desc = "";
  //     popupModal = false;
  //     console.log(problemsFeature);
  //   } else {
  //     alert("Please enter both title and description.");
  //   }
  // };

  // const deleteProblemFeature = (index) => {
  //   problemsFeature = problemsFeature.filter((_, i) => i !== index);
  // };

  let productData = [
    {
      title: "Document Management",
      description: "Organize your files",
      icon: File04Icon,
      isChecked: false,
      dropdownOpen: false,
      selectedOption: "",
    },
    {
      title: "CRM",
      description: "Improve customer relationships with our CRM module",
      icon: UsersRightIcon,
      isChecked: true,
      dropdownOpen: false,
      selectedOption: "",
    },
    {
      title: "Inventory",
      description: "Manage your inventory with our comprehensive module",
      icon: BoxIcon,
      isChecked: false,
      dropdownOpen: false,
      selectedOption: "",
    },
    {
      title: "Forecasting System",
      description: "Predict your sales",
      icon: BarChartSquareIcon,
      isChecked: false,
      dropdownOpen: false,
      selectedOption: "",
    },
  ];

  // let productSelected = [];

  // Fungsi untuk mengecek isChecked pada productData
  // function checkSelected() {
  //   productSelected = productData.filter((product) => product.isChecked);
  // }

  // Fungsi untuk menambahkan productData yang isChecked true ke productSelected
  // function addToSelected(index) {
  //   productData[index].isChecked = true;
  //   checkSelected();
  // }

  // function toggleChecked(index) {
  //   productData[index].isChecked = true;
  // }

  // function toggleChecked(index) {
  //   productData[index].isChecked = !productData[index].isChecked;
  //   if (productData[index].isChecked) {
  //     productSelected.push(productData[index]);
  //   } else {
  //     const selectedIndex = productSelected.findIndex(
  //       (item) => item.title === productData[index].title,
  //     );
  //     if (selectedIndex !== -1) {
  //       productSelected.splice(selectedIndex, 1);
  //     }
  //   }
  // }
</script>

{#if done}
  <div class="flex flex-col items-center justify-center h-full w-full">
    <div class="w-1/4 pt-28 flex flex-col items-center">
      <div
        class="rounded-xl shadow border border-gray-200 justify-center items-center inline-flex mb-4"
      >
        <img src={check} alt="mail.svg" class="p-4 dark:invert" />
      </div>
      <Heading tag="h2" class="mb-4 text-center">You successfully set up you company!</Heading>
      <P class="text-center mb-4"
        >Invitation email are being sent to members and you are ready to use our document management
        features!</P
      >
      <Button color="purple" class="w-full" on:click={handlesuccess}>Go to Dashboard</Button>
    </div>
  </div>
{:else}
  <Timeline {titleHeading} {currentStep} {steps} {product_wizard} {endStep} {validateCurrentStep}
  ></Timeline>
  <!-- <DarkMode /> -->

  <div class="flex flex-col items-center justify-center mt-16 w-full">
    {#if currentStep === 1}
      <div class="w-2/5">
        <form>
          <Heading tag="h2" class="mb-6">Company Types</Heading>
          <div class="grid grid-cols-1 gap-6 w-full items-center justify-center">
            {#each CompanyType as { icon, moduleName, desc, option }, i}
              <Card
                horizontal="true"
                class={`justify-between w-full !p-6 min-h-28 bg-secondary-light dark:bg-secondary-dark ${
                  option
                    ? "border-hc-brand-600 dark:border-hc-brand-500"
                    : "border-secondary-light dark:border-secondary-dark"
                } shadow-sm hover:no-underline`}
                size="none"
                on:click={() => {
                  CompanyType = CompanyType.map((module, index) => ({
                    ...module,
                    option: index === i,
                  }));
                }}
              >
                <div class="flex items-start gap-4">
                  <!-- <div
                    class="w-10 h-10 p-2.5 bg-brand-secondary-light dark:bg-hc-graydark-300 rounded-full justify-center items-center"
                  >
                    <img src={icon} alt="icon" width="20" height="20" />
                  </div> -->
                  <img src={icon} alt="icon" width="28" height="28" />
                  <div class="flex-col">
                    <div class="flex flex-col gap[2px]">
                      {#if option}
                        <P
                          class="text-base font-semibold leading-normal text-hc-brand-900 dark:text-hc-graydark-50"
                        >
                          {moduleName}
                        </P>
                        <P
                          class="text-sm font-normal leading-tight text-hc-brand-700 dark:text-hc-graydark-300"
                        >
                          {desc}
                        </P>
                      {:else}
                        <P
                          class="text-base font-semibold leading-normal text-primary-light dark:text-primary-dark"
                        >
                          {moduleName}
                        </P>
                        <P
                          class="text-sm font-normal leading-tight text-secondary-light dark:text-secondary-dark"
                        >
                          {desc}
                        </P>
                      {/if}
                    </div>
                  </div>
                </div>
                <Radio name="companyType" value={`companyType_${i}`} bind:checked={option}></Radio>
              </Card>
            {/each}
          </div>
        </form>
      </div>
    {:else if currentStep === 2}
      <div class="w-2/5">
        <form>
          <Heading tag="h2" class="mb-4">Company Data</Heading>
          <div class="mb-6">
            <Label for="company_name" class="mb-2">Company Name*</Label>
            <Input
              class="bg-primary-light dark:bg-primary-dark border-primary-light dark:border-primary-dark"
              type="text"
              id="company_name"
              bind:value={companyName}
              placeholder="Enter your company name"
              required
            />
          </div>
          <div class="grid gap-2 mb-6 md:grid-cols-2">
            <div>
              <Label for="company_industry" class="mb-2">Company Industry*</Label>
              <Button
                color="none"
                class="w-full justify-between flex-shrink-0 text-gray-900 bg-primary-light border border-primary-light dark:border-primary-dark dark:text-white focus:ring-gray-300 dark:bg-primary-dark  dark:focus:ring-gray-800"
              >
                {#if selectedIndustry}
                  {selectedIndustry}
                {:else}
                  Your company industry
                {/if}
                <ChevronDownOutline
                  class="w-3 h-3 ms-2 text-end text-hc-graylight-500 dark:hc-graydark-500"
                />
              </Button>

              <!-- <Dropdown containerClass="w-1/5" bind:open={dropdownOpen}>
                {#each companyCategories as item}
                  <DropdownItem on:click={() => handleDropdownSelect(item)}>{item}</DropdownItem>
                {/each}
              </Dropdown> -->
            </div>
            <div>
              <Label for="company_size" class="mb-2">Company Size*</Label>
              <Input
                class="bg-primary-light dark:bg-primary-dark border-primary-light dark:border-primary-dark"
                type="text"
                id="company_size"
                placeholder="How many employees?"
                required
              />
            </div>
          </div>
          <div class="mb-6">
            <div>
              <Label for="textarea-id" class="mb-2">Company Description</Label>
              <Textarea
                class="bg-primary-light dark:bg-primary-dark border-primary-light dark:border-primary-dark"
                id="textarea-id"
                bind:value={companyDescription}
                placeholder="Enter your company description"
                rows="4"
                name="message"
              />
            </div>
          </div>
          <div class="grid gap-2 mb-6 md:grid-cols-2">
            <div>
              <Label for="logo" class="mb-2">Company Logo</Label>
              <Button
                on:click={handleImageUploadWrapper}
                color="light"
                class="w-full rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark flex-col justify-start items-center gap-1 flex py-8 mb-4"
              >
                {#if uploadedImage}
                  <img src={uploadedImage} alt="Uploaded Logo" class="rounded-xl max-h-48 mb-4" />
                  <P>{uploadedImageName}</P>
                {:else}
                  <div
                    class="rounded-xl shadow border border-primary-light dark:border-primary-dark justify-center items-center inline-flex mb-4"
                  >
                    <img src={upload} alt="mail.svg" class="p-4 dark:invert" />
                  </div>
                  <div class="flex space-x-1">
                    <P class="text-hc-brand-700">Click to upload</P><P>or drag and drop</P>
                  </div>
                  <P>SVG, PNG, JPG, or GIF (max. 800x400px)</P>
                {/if}
              </Button>
            </div>
            <div>
              <Label for="brand_color" class="mb-2">Preferred Brand Color</Label>
              <div>No Can Do!</div>
              <!-- <div
                class="w-full px-3.5 py-2.5 bg-primary-light dark:bg-primary-dark rounded-lg border border-primary-light dark:border-primary-dark justify-start items-center gap-2 inline-flex [&>span>label]:text-primary-light dark:[&>span>label]:text-primary-dark"
              >
                <ColorPicker
                  bind:hex
                  components={ChromeVariant}
                  label={hex}
                  sliderDirection="horizontal"
                />
              </div> -->
              <!-- <Input
                class="bg-primary-light dark:bg-primary-dark border"
                type="text"
                id="brand_color"
                placeholder="#7F56D9"
                required
                pattern="^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$"
              >
              </Input> -->
            </div>
          </div>
        </form>
      </div>
    {:else if currentStep === 3}
      <div class="w-full px-20">
        <div class="flex justify-between items-end mb-6">
          <div>
            <Heading tag="h3" class="mb-1">Choose product</Heading>
            <P class="">Please select the product you’d like to include in your solution</P>
          </div>
          <Button
            on:click={() => (popupModal = true)}
            outline="true"
            size="sm"
            class="hover:bg-transparent dark:hover:bg-transparent hover:text-brand-primary-light dark:hover:text-hc-brand-500 h-max"
            ><svelte:component this={PlusIcon} svgClass="text-hc-brand-700 mr-1" />Add</Button
          >
        </div>
        <!-- {#if problemsFeature.length > 0}
          {#each problemsFeature as { title, desc }, i}
            <div class="bg-gray-50 rounded-md">
              <div class="flex justify-between px-8 py-4">
                <div>
                  <Heading tag="h5">{title}</Heading>
                  <P>{desc}</P>
                </div>
                <Button color="none" on:click={() => deleteProblemFeature(i)}>
                  <img src={CloseIcon} alt="close.svg" />
                </Button>
              </div>
            </div>
          {/each}
        {:else}
          <div class="bg-gray-50 rounded-md mb-6">
            <P class="px-8 py-4">You haven’t added product</P>
          </div>
        {/if} -->
        <div class="grid grid-cols-2 gap-4 w-full items-center justify-center">
          <div>No Card</div>
          <!-- {#each productData as product}
            {#if product.isChecked}
              <ProductCard dataCard={product} />
            {/if}
          {/each} -->
        </div>
      </div>
    {/if}
  </div>
{/if}

<!-- <Modal title="Choose product" bind:open={popupModal} size="xl">
  <div class="space-y-8">
    <div class="grid grid-cols-2 gap-4 w-full items-center justify-center">
      {#each productData as product, index}
        <ProductCard isChoice dataCard={product} />
      {/each}
    </div>
    <Button class="w-full"
      ><svelte:component this={PlusIcon} svgClass="text-white mr-1" />Add</Button
    >
  </div>
</Modal> -->
