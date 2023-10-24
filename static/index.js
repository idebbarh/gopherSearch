const formElement = document.getElementById("search-form");
const inputElement = formElement.querySelector(".search-form__input");
const searchEndpoint = "/search?q=";

async function submitHandler(e) {
  e.preventDefault();
  const searchQuery = inputElement.value.trim().split(" ").join("+");
  const res = await fetch(searchEndpoint + searchQuery);
  if (!res.ok) {
    const errorMessage = await res.text();
    console.error("Server error:", errorMessage);
    return;
  }
  const data = await res.json();
  console.log(data);
}

formElement.addEventListener("submit", submitHandler);
