async function handleLogin(event) {
    event.preventDefault(); // Предотвращаем отправку формы
  
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;
    const errorMessage = document.getElementById("errorMessage");
  
    try {
      const response = await fetch("http://localhost:8000/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });
  
      if (response.ok) {
        const data = await response.json();
        console.log("Успешный вход:", data);
        // Редирект на главную страницу или другую
        window.location.href = "/dashboard";
      } else {
        errorMessage.style.display = "block";
        errorMessage.textContent = "Неверный email или пароль";
      }
    } catch (error) {
      console.error("Ошибка запроса:", error);
      errorMessage.style.display = "block";
      errorMessage.textContent = "Произошла ошибка. Попробуйте позже.";
    }
  }
  