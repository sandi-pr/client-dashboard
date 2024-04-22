<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: Azzis Arswendo <azzis@hexacode.org>

Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->

<script>
  import { Input, Label, Heading, Button } from "flowbite-svelte";
  import axios from "axios";
  import { createEventDispatcher } from "svelte";

  let a = 0;
  let b = 0;
  let result = 0;

  const dispatch = createEventDispatcher();

  function multiply() {
    var url = "/client-dashboard/multiply";
    axios
      .post(url, { a: a, b: b })
      .then((response) => {
        result = response.data.result.result;
      })
      .catch((error) => {
        dispatch("error", {
          error: error,
        });
      });
  }
</script>

<!-- <div class="w-1/3 mx-auto mt-10 flex flex-col">
    <div class="text-center">
        <Heading tag="h1">CR3 Client Dashboard</Heading>
        <Heading tag="h3">Multiply Page</Heading>
    </div>

    <div class="mt-10">
        <Label for="a" class="mb-2">Value A</Label>
        <Input bind:value={a} id="a" type="number" placeholder="Enter Value A" />
    </div>

    <div class="mt-10">
        <Label for="b" class="mb-2">Value B</Label>
        <Input bind:value={b} id="b" type="number" placeholder="Enter Value B" />
    </div>

    <div class="mt-10">
        <Button on:click={multiply}>Multiply</Button>
    </div>

    <div class="mt-10">
        <Label for="result" class="mb-2">Result</Label>
        <Input bind:value={result} id="result" type="text" placeholder="Result" readonly/>
    </div>
</div> -->
