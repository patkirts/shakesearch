import { h, hydrate } from 'https://unpkg.com/preact@latest?module'
import { useState } from 'https://unpkg.com/preact@latest/hooks/dist/hooks.module.js?module'
import { html } from 'https://unpkg.com/htm/preact/index.module.js?module'

// Initialize htm with Preact
const htm = html.bind(h)

hydrate(htm`<${App} name="World" />`, document.getElementById('app'))

function App() {
  const [results, setResults] = useState([])

  return htm`
  <h1>Welcome to ShakeSearch!</h1>
  <p>Where you can find Shakespeare's words.</p>
  <form id="form" onSubmit=${e => onSubmit(e, setResults)}>
    <input type="text" id="query" name="query"
      placeholder="Search Shakespeare's Complete Works"
    />
    <button type="submit">Search</button>
  </form>
  <${Results} results=${results} />
  `
}

async function onSubmit(e, setResults) {
  e.preventDefault()
  const form = document.getElementById("form")
  const data = Object.fromEntries(new FormData(form))

  try {
    const res = await (await fetch(`/search?q=${data.query}`)).json()
    setResults(res)

  } catch (err) {
    console.error("Something went wrong:", err)
  }
}

function Results({ results }) {
  if (!results.length) return null
  console.log(results)

  return htm`<ul>
    ${results.map(r => {
      return htm`<li>
        <div class="">${r.Lines.map(l => htm`${l.trim()}<br />`)}</div>
        <footer><em>${r.Title}</em></footer>
      </li>`
    })}
  </ul>`
}