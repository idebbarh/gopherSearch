const formElement = document.getElementById("search-form");
const inputElement = formElement.querySelector(".search-form__input");
const getDocsEndpoint = "/search?q=";
const getDocContentEndpoint = "/file?path=";
const isFirstSearch = true;

function animateFromMiddleToTop() {
  const searchContainer = document.querySelector(".search-container");
  searchContainer.classList.add("search-container--animate");
}

function getDocContent(path) {
  fetch(getDocContentEndpoint + path)
    .then((res) => {
      return res.blob();
    })
    .then((blob) => {
      const fileURL = URL.createObjectURL(blob);
      window.open(fileURL, "_blank");
    });
}

function getLinkElement(path, title) {
  const linkElement = document.createElement("a");
  linkElement.className = "search-result-container__doc-link";
  linkElement.innerText = title;

  linkElement.addEventListener("click", () => getDocContent(path));
  return linkElement;
}

function renderDocs(docs) {
  const resultContainer = document.querySelector(".search-result-container");
  resultContainer.innerHTML = "";

  docs.Result.forEach((doc) => {
    resultContainer.appendChild(getLinkElement(doc.Path, doc.Title));
  });
}

async function getDocs(e) {
  e.preventDefault();
  const searchQuery = inputElement.value.trim().split(" ").join("+");
  const res = await fetch(getDocsEndpoint + searchQuery);

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

formElement.addEventListener("submit", getDocs);
