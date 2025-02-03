<script>
    import { writable } from "svelte/store";

    let iframeSrc = writable("http://192.168.18.13:54321");
    let ref;

    window.setUrlCsv = (newUrl) => {
        iframeSrc.set(newUrl);
    };


    var old = localStorage.getItem("color-theme")

    function ubah_tema(tema){
        console.log(tema)
    }

    setInterval(() => {
        let baru = localStorage.getItem("color-theme")
        if(old != baru){
            old = baru
            ubah_tema(baru)
        }
    }, 500)
    
</script>

<iframe bind:this={ref} class="w-full h-full" src={$iframeSrc}></iframe>
