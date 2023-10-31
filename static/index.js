const formElement = document.getElementById("search-form");
const inputElement = formElement.querySelector(".search-form__input");
const getDocsEndpoint = "/search?q=";
const getNextDocsEndpoint = "/nextSearch";
const getDocContentEndpoint = "/file?path=";

function getDocContent(path) {
  fetch(getDocContentEndpoint + path)
    .then((res) => {
      return res.blob();
    })
    .then((blob) => {
      const fileURL = URL.createObjectURL(blob);
      return window.open(fileURL, "_blank");
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

  const data = await res.json();

  renderDocs(data, inputElement.value.trim(), true);
}

async function getNextDocs() {
  const res = await fetch(getNextDocsEndpoint);
  if (!res.ok) {
    const errorMessage = await res.text();
    console.error("Server error:", errorMessage);
    return;
  }

  const data = await res.json();

  renderDocs(data, inputElement.value, false);
}

function getLinkElement(path, title) {
  const linkElement = document.createElement("a");
  linkElement.className = "search-result-container__doc-link";
  linkElement.innerText = title;

  linkElement.addEventListener("click", () => getDocContent(path));
  return linkElement;
}

function renderDocs(docs, searchQuery, isNewSearch) {
  const resultContainer = document.querySelector(".search-result-container");

  if (isNewSearch) {
    while (resultContainer.firstChild) {
      resultContainer.removeChild(resultContainer.lastChild);
    }
    if (docs.Result.length === 0) {
      emptyResultElement = document.createElement("p");
      emptyResultElement.innerText = `Your search - ${searchQuery} - did not match any documents.`;
      emptyResultElement.className = "search-container__no-result";
      resultContainer.appendChild(emptyResultElement);
      return;
    }
  }

  const getNextDocsBtn = document.getElementById("next-docs-btn");
  const isCompleteData = docs.IsCompleteData;

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
