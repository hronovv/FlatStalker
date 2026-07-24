const tg = window.Telegram?.WebApp;

if (tg) {
  tg.ready();
  tg.expand();
}

const user = tg?.initDataUnsafe?.user;
const greeting = document.getElementById("greeting");
const subtitle = document.getElementById("subtitle");

if (user?.first_name) {
  greeting.textContent = `Hello, ${user.first_name}!`;
  subtitle.textContent = "Telegram Mini App работает";
} else {
  greeting.textContent = "Hello World";
  subtitle.textContent =
    "Открой эту страницу из Telegram или через HTTPS-туннель";
}
