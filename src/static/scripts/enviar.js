const btn = document.getElementById("enviar")
const inputs = Array.from(document.querySelectorAll("input, textarea"))
const feedback = document.getElementById("feedback");
const form = document.getElementById("form-contacto");

function checkInputs() {
  const llenos = inputs.every(el => el.value.trim() != "")
  btn.disabled = !llenos
}

inputs.forEach(el => {
  el.addEventListener("input", checkInputs)
})

checkInputs()

document.body.addEventListener("htmx:afterRequest", (evt) => {
  const elt = evt.detail.elt;
  const xhr = evt.detail.xhr;

  if (elt.id === "form-contacto" && xhr.status === 200) {
    elt.reset();
    checkInputs();
  }
});

document.body.addEventListener("htmx:afterSwap", (evt) => {
  if (evt.detail.target.id === "feedback") {
    const fb = evt.detail.target;

    fb.classList.remove("vh");

    setTimeout(() => {
      fb.classList.add("vh");
      fb.innerHTML = "";
    }, 3000);
  }
});
