async function handleRegister(event) {
  event.preventDefault();

  const data = {
    JobRoleId: parseInt(document.getElementById("jobRole").value),
    Name: document.getElementById("name").value,
    SecondName: document.getElementById("secondName").value,
    Surname: document.getElementById("surname").value,
    Email: document.getElementById("email").value,
    Password: document.getElementById("password").value,
    Birthday: new Date(document.getElementById("birthday").value).getTime() / 1000, // Unix timestamp
    BirthdayDate: new Date(document.getElementById("birthday").value).toISOString(), // ISO string
    Address: {
      SettlementTypeId: parseInt(document.getElementById("settlementType").value),
      Country: document.getElementById("country").value,
      Region: document.getElementById("region").value,
      District: document.getElementById("district").value,
      Settlement: document.getElementById("settlementName").value,
      Street: document.getElementById("street").value,
      HouseNumber: document.getElementById("house").value,
      FlatNumber: document.getElementById("apartment").value,
    },
  };

  const errorMessage = document.getElementById("errorMessage");

  try {
    const response = await fetch("http://localhost:8000/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });

    if (response.ok) {
      window.location.href = "http://localhost:8000/view/login.html";
    } else {
      errorMessage.style.display = "block";
      errorMessage.textContent = "Ошибка регистрации. Проверьте данные.";
    }
  } catch (error) {
    console.error("Ошибка запроса:", error);
    errorMessage.style.display = "block";
    errorMessage.textContent = "Произошла ошибка. Попробуйте позже.";
  }
}
