<script>
	import FileUploader from "$lib/FileUploader.svelte";

</script>

<FileUploader 
  uploadUrl="/uploadFile"
  acceptedTypes=".html"
  maxSizeMB={5}
/>