<script>
    async function uploadFolder () {

    }
    let files = [];

    function handleFolderChange(event) {
        files = Array.from(event.target.files);
        
        let formData = new FormData()
        files.forEach((f) => {
            formData.append("files", f, f.webkitRelativePath)
            console.log(formData)
        })
    }


</script>


<input type="file" webkitdirectory multiple on:change={handleFolderChange} />
<button on:click={uploadFolder}> Upload Folder</button>