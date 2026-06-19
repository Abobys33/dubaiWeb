const popup = document.getElementById('popup');
const openBtn = document.getElementById('openPopup');
const closeBtn = document.getElementById('popupClose');
const phoneInput = document.getElementById('phoneInput');
const nameInput = document.getElementById("nameInput");
const submitButton = document.getElementById("submitButton");

openBtn.addEventListener('click', () => {
  popup.classList.add('active');
  document.body.style.overflow = 'hidden';
});

closeBtn.addEventListener('click', closePopup);

popup.addEventListener('click', (e) => {
  if (e.target === popup) {
    closePopup();
  }
});

document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape') {
    closePopup();
  }
});

function closePopup() {
  popup.classList.remove('active');
  document.body.style.overflow = '';
}

phoneInput.addEventListener("input", maskPhone);

function maskPhone(e) {

  let input = e.target.value.replace(/\D/g, "");
  
  if (input.startsWith("8")) input = "7" + input.slice(1);
  if (input.startsWith("9")) input = "7" + input;

  let formatted = "+7 ";

  if (input.length > 1) {
    formatted += "(" + input.substring(1,4);
  }

  if (input.length >= 4) {
    formatted += ") " + input.substring(4,7);
  }

  if (input.length >= 7) {
    formatted += "-" + input.substring(7,9);
  }

  if (input.length >= 9) {
    formatted += "-" + input.substring(9,11);
  }

  e.target.value = formatted;
}

phoneInput.addEventListener("focus", () => {
  if (!phoneInput.value) {
    phoneInput.value = "+7 (";
  }
});

nameInput.addEventListener("input", validateForm);
phoneInput.addEventListener("input", validateForm);

function validateForm(){

  const nameValid = nameInput.value.trim().length >= 2;
  const phoneValid = phoneInput.value.length === 18;

  submitButton.disabled = !(nameValid && phoneValid);
}

nameInput.addEventListener("input", () => {
  nameInput.value = nameInput.value.replace(/[0-9]/g, "");
});
document.addEventListener("DOMContentLoaded", () => {

  const slider = document.getElementById("objectsSlider");

  slider.addEventListener("wheel", (e) => {
    e.preventDefault();
    slider.scrollLeft += e.deltaY;
  });

  const data = [
    {
      title: "Samana Boulevard Heights",
      subtitle: "Это комплекс в самом сердце растущего жилого<br>массива Дубая в 10 минутах от Дубай Аутлет Молл ",
      price: "от 804 000 AED",
      roi: "Q2 2027",
      image: "img/objects/obj1.png"
    },
    {
      title: "Binghatti Vintage",
      subtitle: "Это футуристический жилой комплекс в районе Маджан<br>с бассейном, спа и открытым кинотеатром",
      price: "от 670 000 AED",
      roi: "ACTIVE",
      image: "img/objects/obj2.png"
    },
     {
      title: "LuzOra Residences",
      subtitle: "Премиум-комплекс с бассейном и<br>залом в районе Дейра в Дубае",
      price: "от 1 590 000 AED",
      roi: "Q2 2027",
      image: "img/objects/obj3.png"
    },

    {
      title: "Celesto 2 Tower",
      subtitle: "Это 18-этажный жилой комплекс, расположенный<br >в районе Dubai Land Residence Complex",
      price: "от 1 032 350 AED",
      roi: "Q2 2028",
      image: "img/objects/obj4.png"
    },
  ];

  slider.innerHTML = data.map(item => `
    <div class="object-card">
      <div class="object-image">
        <img src="${item.image}">
        <div class="object-badges">
          <div class="badge">${item.price}</div>
          <div class="badge small">${item.roi}</div>
        </div>
      </div>

      <div class="object-info">
        <div>
          <div class="object-title">${item.title}</div>
          <div class="object-subtitle">${item.subtitle}</div>
        </div>

        <div class="object-link">
          Узнать подробнее
          <img src="img/objects/arrow.svg">
        </div>
      </div>
    </div>
  `).join('');

});
