<script>
  export let value = 50;
</script>

<input
  type="range"
  min="0"
  max="100"
  step="1"
  bind:value
  class="w-full h-2 rounded-lg appearance-none"
  id="slider"
  style="--slider-value: {value};"
/>

<style>
  /* thumb UI */
  #slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: #ffffff;
    cursor: pointer;
    border: 2px solid var(--primary-600); /* Purple border */
    margin-top: -6px; /* Center the thumb on the track */
    transition:
      background 0.3s,
      border-color 0.3s;
  }
  /* bar dynamic background */
  #slider::-webkit-slider-runnable-track {
    width: 100%;
    height: 8px;
    cursor: pointer;
    background: linear-gradient(
      to right,
      var(--primary-600) 0%,
      var(--primary-600) calc(var(--slider-value) * 1%),
      #eaecf0 calc(var(--slider-value) * 1%),
      #eaecf0 100%
    );
    border-radius: 4px;
    border: none;
  }
</style>
