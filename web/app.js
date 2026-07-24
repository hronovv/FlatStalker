const tg = window.Telegram?.WebApp;

if (tg) {
  tg.ready();
  tg.expand();
  if (tg.setHeaderColor) tg.setHeaderColor("#eef3f8");
  if (tg.setBackgroundColor) tg.setBackgroundColor("#eef3f8");
}

const user = tg?.initDataUnsafe?.user;
const greeting = document.getElementById("greeting");
const lead = document.getElementById("lead");
const cta = document.getElementById("cta-main");

const botUsername = "flatstalker_bot"; // поменяй на username своего бота
cta.href = `https://t.me/${botUsername}`;

if (user?.first_name) {
  greeting.textContent = `${user.first_name}, добро пожаловать в FlatStalker`;
  lead.textContent =
    "Мы ловим свежие объявления об аренде в Беларуси и присылаем ссылку в Telegram раньше, чем квартира уйдёт.";
} else {
  greeting.textContent = "Ловите квартиру раньше остальных";
}
