const apiKey = process.env.GITHUB_TOKEN
const resp = await fetch("https://api.github.com/repos/agustin-del/proyectos/contents", {
  headers:{
    Authorization: `Bearer ${apiKey}`,
    Accept: "application/vnd.github+json",
    "User-Agent": "agustin-del"
  } 
})

console.log(await resp.json())
