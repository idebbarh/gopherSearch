const formElement = document.getElementById("search-form");
const inputElement = formElement.querySelector(".search-form__input");
const searchEndpoint = "/search?q=";
const isFirstSearch = true;

function animateFromMiddleToTop() {
  const searchContainer = document.querySelector(".search-container");
  searchContainer.classList.add("search-container--animate");
}

function getDocLink(path, title) {
  const linkElement = document.createElement("a");
  linkElement.className = "search-result-container__doc-link";
  linkElement.innerText = title;

  linkElement.addEventListener("click", (e) => {
    fetch("/file", {
      method: "POST",
      body: JSON.stringify({ filePath: path }),
    })
      .then((res) => {
        return res.blob();
      })
      .then((blob) => {
        const fileURL = URL.createObjectURL(blob);
        window.open(fileURL, "_blank");
      });
  });
  return linkElement;
}

function renderDocs(docs) {
  const resultContainer = document.querySelector(".search-result-container");
  resultContainer.innerHTML = "";

  docs.Result.forEach((doc) => {
    resultContainer.appendChild(getDocLink(doc.Path, doc.Title));
  });
}

async function submitHandler(e) {
  e.preventDefault();
  const searchQuery = inputElement.value.trim().split(" ").join("+");
  const res = await fetch(searchEndpoint + searchQuery);

  if (!res.ok) {
    const errorMessage = await res.text();
    console.error("Server error:", errorMessage);
    return;
  }

  if (isFirstSearch) {
    animateFromMiddleToTop();
  }

  const data = await res.json();

  setTimeout(() => {
    renderDocs(data);
  }, 2000);
}

formElement.addEventListener("submit", submitHandler);
