document.addEventListener('htmx:afterRequest', function (event) {
  if (event.detail.requestConfig.path === '/shorten') {
    const response = JSON.parse(event.detail.xhr.responseText);
    const resultDiv = document.getElementById('result');
    resultDiv.innerHTML = `
      <p>Shortened URL:</p>
      <a id="shortUrl" href="${response.short_url}" target="_blank">${response.short_url}</a>
    `;
  }
});