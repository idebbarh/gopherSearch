const formElement = document.getElementById("search-form");
const inputElement = formElement.querySelector(".search-form__input");
const searchEndpoint = "/search?q=";
const isFirstSearch = true;

function animateFromMiddleToTop() {
  const searchContainer = document.querySelector(".search-container");
  searchContainer.classList.add("search-container--animate");
}

function renderDocs(docs) {
  const resultContainer = document.querySelector(".search-result-container");
  resultContainer.innerHTML = "";

  docs.Result.forEach((doc) => {
    const docElement = document.createElement("p");
    docElement.innerText = doc.Title;
    resultContainer.appendChild(docElement);
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
