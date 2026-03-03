const shortBtn = document.querySelector(".btn")
const inputElement = document.querySelector(".url-input")
const outputElement = document.querySelector(".short-url")

const apiUrl = "http://localhost:8080"
var currentUrl = inputElement.value

async function shortUrl(url) {
    const response = await fetch(`${apiUrl}/short`, {
        "method": "POST",
        "body": JSON.stringify({ "url": url })
    })

    if (!response.ok) {
        outputElement.textContent = "Произошла ошибка при запросе"
        outputElement.href = "#"
    }

    body = await response.json()

    slugUrl = `${apiUrl}/s/${body["slug"]}`

    currentUrl = slugUrl
    outputElement.textContent = slugUrl
    outputElement.href = slugUrl
}

function onClickShortButton() {
    if (currentUrl != inputElement.value)
        shortUrl(inputElement.value)
}

shortUrl(currentUrl)

shortBtn.addEventListener("click", onClickShortButton, false)