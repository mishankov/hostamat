# APIs

- POST /uploadsession - creates a new upload session
- DELETE /uploadsession + H-Upload-Session - finish an upload session
- POST /files/api/upload  + H-Upload-Session - upload a file
# Platforma fixes

- [ ] Fix mandatory trailing slashes in URL

# Frontend tips
## Upload files

```javascript
// 1. Create upload session
const {uploadId} = await fetch('/uploads', {method: 'POST'});

// 2. Upload files
await Promise.all(files.map(file => 
  fetch(`/uploads/${uploadId}/files`, {
    method: 'POST',
    body: file
  })
));

// 3. Finalize
await fetch(`/uploads/${uploadId}/close`, {method: 'POST'});
```
