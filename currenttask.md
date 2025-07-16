# Next
- [ ] Delete sites
- [ ] Delete files
- [ ] Trace Id as separate middleware

# Bugs

- [ ] Non existing files return 200, should 404

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
