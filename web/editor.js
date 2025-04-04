function uploadFormSubmit(e) {
    e.preventDefault();
    const formData = new FormData(e.target);
    const fileInput = formData.get('file');
    const fileName = fileInput.name;
    const fileSize = fileInput.size;

    fetch('/upload', {
        method: 'POST',
        body: formData,
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    alert(`File "${fileName}" uploaded successfully!`);
}