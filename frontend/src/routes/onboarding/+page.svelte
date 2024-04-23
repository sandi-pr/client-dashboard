<script>
  import {
    P,
    Label,
    Input,
    Table,
    Button,
    Heading,
    Textarea,
    Checkbox,
    TableBody,
    TableHead,
    TableBodyRow,
    TableBodyCell,
    TableHeadCell,
  } from "flowbite-svelte";
  import Papa from "papaparse";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import upload from "../lib/icon/upload-cloud-02.svg";
  import csv from "../lib/icon/csv.png";
  import Timeline from "../lib/components/Timeline.svelte";
  import check from "../lib/icon/check-circle.svg";
  import { timelineStep } from "../lib/utils/timelineStep.js";
  import { userDataSession } from "../lib/utils/auth";
  import { emailBlast, checkOrganization } from "../lib/utils/setupcompany";

  onMount(async () => {
    const isOrganizationExist = await checkOrganization();

    if (isOrganizationExist) {
      console.log("udah ada");
      window.history.back();
    }

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

  const isValidEmail = (email) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

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

  const handleCsvUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = ".csv";

    inputElement.addEventListener("change", (event) => {
      const file = event.target.files[0];

      if (file && isValidFile(file)) {
        const reader = new FileReader();

        reader.onload = () => {
          Papa.parse(file, {
            complete: (result) => {
              uploadedFile = reader.result;
              uploadedFileName = file.name;
              uploadedFileSize = (file.size / 1024).toFixed(2);
              csvData = result.data;
              emailList = csvData
                .map((row) => row[Object.keys(row)[5]])
                .filter((cell) => isValidEmail(cell));
              rows = csvData
                .filter((row) => isValidEmail(row[Object.keys(row)[5]]))
                .map((row) => {
                  const newRow = {};
                  Object.keys(row).forEach((key) => {
                    const sanitizedKey = key.replace(/\s+/g, "").replace(/'/g, "").toLowerCase();
                    newRow[sanitizedKey] = row[key];
                  });
                  return newRow;
                });
            },
            header: true,
          });
        };

        reader.readAsText(file);
      } else {
        alert("Invalid file. Please select a valid CSV file (max. 200KB).");
      }
    });

    inputElement.click();
  };

  const isValidFile = (file) => {
    const allowedTypes = ["application/vnd.ms-excel", "text/csv"];
    const maxFileSize = 200 * 1024;
    return allowedTypes.includes(file.type) && file.size <= maxFileSize;
  };

  let currentStep = 1;
  const steps = ["1", "2", "3"];

  const unsubscribe = timelineStep.subscribe((value) => {
    currentStep = value;
  });

  let onboarding = true;
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
  <Timeline {currentStep} {steps} {onboarding} {endStep} {validateCurrentStep}></Timeline>

  <div class="flex flex-col items-center justify-center h-full w-full">
    {#if currentStep === 1}
      <div class="w-2/4">
        <form>
          <Heading tag="h2" class="mb-4">Company Data</Heading>
          <Label for="logo" class="mb-2">Company Logo</Label>
          <Button
            on:click={handleImageUploadWrapper}
            color="light"
            class="w-full rounded-xl border border-gray-200 flex-col justify-start items-center gap-1 flex py-8 mb-4"
          >
            {#if uploadedImage}
              <img src={uploadedImage} alt="Uploaded Logo" class="rounded-xl max-h-48 mb-4" />
              <P>{uploadedImageName}</P>
            {:else}
              <div
                class="rounded-xl shadow border border-hc-graylight-200 justify-center items-center inline-flex mb-4"
              >
                <img src={upload} alt="mail.svg" class="p-4 dark:invert" />
              </div>
              <div class="flex space-x-1">
                <P class="text-hc-brand-700">Click to upload</P><P>or drag and drop</P>
              </div>
              <P>SVG, PNG, JPG, or GIF (max. 800x400px)</P>
            {/if}
          </Button>
          <div class="mb-6">
            <Label for="company_name" class="mb-2">Company Name*</Label>
            <Input
              type="text"
              id="company_name"
              bind:value={companyName}
              placeholder="Enter your company name"
              required
            />
          </div>
          <div>
            <Label for="textarea-id" class="mb-2">Company Description</Label>
            <Textarea
              id="textarea-id"
              bind:value={companyDescription}
              placeholder="Enter your company description"
              rows="4"
              name="message"
            />
          </div>
        </form>
      </div>
    {:else if currentStep === 2}
      <div class="w-2/4">
        <form>
          <Heading tag="h2" class="mb-4">Upload your organization structure</Heading>
          <Label for="logo" class="mb-4"
            >Define your organization structure in .csv and we will automatically import it into our
            system</Label
          >
          <Button
            on:click={handleCsvUploadWrapper}
            color="light"
            class="w-full rounded-xl border border-hc-graylight-200 flex-col justify-start items-center gap-1 flex py-8 mb-4"
          >
            {#if uploadedFile}
              <div class="flex justify-between w-full">
                <div class="flex h-full items-center">
                  <img src={csv} alt="csv" class="mr-4" />
                  <div class="flex-col">
                    <P>{uploadedFileName}</P>
                    <P>{uploadedFileSize} KB - 100% uploaded</P>
                  </div>
                </div>
                <Checkbox checked />
              </div>
            {:else}
              <div
                class="rounded-xl shadow border border-gray-200 justify-center items-center inline-flex mb-4"
              >
                <img src={upload} alt="mail.svg" class="p-4 dark:invert" />
              </div>
              <div class="flex space-x-1">
                <P class="text-hc-brand-700">Click to upload</P><P>or drag and drop</P>
              </div>
              <P>CSV file (max. 200KB)</P>
            {/if}
          </Button>
        </form>
      </div>
    {:else if currentStep === 3}
      <div class="w-3/5">
        <Heading tag="h2" class="mb-4 text-center">Confirm your organization structure</Heading>
        {#if csvData.length > 0}
          <Table>
            <TableHead>
              {#each Object.keys(csvData[0]) as column}
                <TableHeadCell>{column}</TableHeadCell>
              {/each}
            </TableHead>
            <TableBody class="divide-y">
              {#each csvData as row}
                <TableBodyRow>
                  {#each Object.values(row) as cell, index}
                    {#if index == 2}
                      <TableBodyCell>
                        {#if cell == "Null"}
                          <Button
                            class="px-1.5 py-0.5 bg-hc-graylight-50 text-hc-graylight-700 border border-hc-graylight-200 rounded-md focus-within:ring-0 font-medium"
                          >
                            {cell}
                          </Button>
                        {:else}
                          <Button
                            class="px-1.5 py-0.5 bg-hc-brand-50 text-hc-brand-700 border border-hc-brand-200 rounded-md focus-within:ring-0 font-medium"
                          >
                            {cell}
                          </Button>
                        {/if}
                      </TableBodyCell>
                    {:else if index == 3}
                      <TableBodyCell>
                        {#if cell == "Null"}
                          <Button
                            class="px-1.5 py-0.5 bg-hc-graylight-50 text-hc-graylight-700 border border-hc-graylight-200 rounded-md focus-within:ring-0 font-medium"
                          >
                            {cell}
                          </Button>
                        {:else}
                          <Button
                            class="px-1.5 py-0.5 bg-blue-100 text-blue-800 border border-blue-300 rounded-md focus-within:ring-0 font-medium"
                          >
                            {cell}
                          </Button>
                        {/if}
                      </TableBodyCell>
                    {:else if index == 4}
                      <TableBodyCell>
                        {#if cell == "Null"}
                          <Button
                            class="px-1.5 py-0.5 bg-hc-graylight-50 text-hc-graylight-700 border border-hc-graylight-200 rounded-md focus-within:ring-0 font-medium"
                          >
                            {cell}
                          </Button>
                        {:else}
                          <Button
                            class="px-1.5 py-0.5 bg-blue-50 text-blue-700 border border-blue-200 rounded-md focus-within:ring-0 font-medium"
                          >
                            {cell}
                          </Button>
                        {/if}
                      </TableBodyCell>
                    {:else}
                      <TableBodyCell>{cell}</TableBodyCell>
                    {/if}
                  {/each}
                </TableBodyRow>
              {/each}
            </TableBody>
          </Table>
        {:else}
          <Heading tag="h6" class="text-center">No CSV data available</Heading>
        {/if}
      </div>
    {/if}
  </div>
{/if}
