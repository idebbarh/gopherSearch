const formElement = document.getElementById("search-form");
const inputElement = formElement.querySelector(".search-form__input");
const getDocsEndpoint = "/search?q=";
const getNextDocsEndpoint = "/nextSearch";
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
    })
    .catch((errorMessage) => {
      console.error("Server error:", errorMessage);
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

async function getNextDocs() {
  const res = await fetch(getNextDocsEndpoint);
  if (!res.ok) {
    const errorMessage = await res.text();
    console.error("Server error:", errorMessage);
    return;
  }

  const data = await res.json();

  renderDocs(data, false);
}

function getLinkElement(path, title) {
  const linkElement = document.createElement("a");
  linkElement.className = "search-result-container__doc-link";
  linkElement.innerText = title;

  linkElement.addEventListener("click", () => getDocContent(path));
  return linkElement;
}

function renderDocs(docs, newSearch = true) {
  const resultContainer = document.querySelector(".search-result-container");
  const getNextDocsBtn = document.getElementById("next-docs-btn");
  const isCompleteData = docs.IsCompleteData;

  if (newSearch) {
    resultContainer.innerHTML = "";
  }
  docs.Result.forEach((doc) => {
    if (!getNextDocsBtn) {
      resultContainer.appendChild(getLinkElement(doc.Path, doc.Title));
    } else {
      getNextDocsBtn.before(getLinkElement(doc.Path, doc.Title));
    }
  });

  if (!isCompleteData && !getNextDocsBtn) {
    const getNextDocsBtn = document.createElement("button");
    getNextDocsBtn.id = "next-docs-btn";
    getNextDocsBtn.innerText = "more results";
    getNextDocsBtn.addEventListener("click", getNextDocs);
    resultContainer.appendChild(getNextDocsBtn);
  }

  if (isCompleteData && getNextDocsBtn) {
    resultContainer.removeChild(getNextDocsBtn);
  }
}

formElement.addEventListener("submit", getDocs);
