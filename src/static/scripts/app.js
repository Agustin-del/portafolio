const btn = document.getElementById("volver")

console.log(btn)

btn.addEventListener("click", () => {
  history.back()
})
