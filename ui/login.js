async function handleLogin(event) {
  event.preventDefault();

  const email = document.getElementById("email").value;
  const password = document.getElementById("password").value;
  const errorMessage = document.getElementById("errorMessage");

  const response = await fetch("http://localhost:8000/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  }).then((response) => response.json())
    .then((token) => {
      setTimeout(() => {
        location.href = "http://localhost:8088/";
      }, 500);
    })
    .catch((error) => {
      console.error("Ошибка запроса:", error);
      errorMessage.style.display = "block";
      errorMessage.textContent = "Произошла ошибка. Попробуйте позже.";
    });
}
  