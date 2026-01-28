const btn = document.getElementById("enviar")
const inputs = Array.from(document.querySelectorAll("input, textarea"))

function checkInputs() {
  const llenos = inputs.every(el => el.value.trim() != "")
  btn.disabled = !llenos
}

inputs.forEach(el => {
  el.addEventListener("input", checkInputs)
})

checkInputs()

