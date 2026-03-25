export function runAsync(fn, onError = console.error) {
  return (...args) => {
    Promise.resolve(fn(...args)).catch(onError)
  }
}

export function setMessage(el, type, text) {
  el.className = "message"
  el.textContent = text || ""
  if (!text) {
    return
  }
  el.classList.add(type === "ok" ? "ok" : "err")
}

export function titleCase(value) {
  return String(value || "")
    .trim()
    .toLowerCase()
    .replace(/\b\w/g, (c) => c.toUpperCase())
}

export function normalizeStatus(status) {
  return titleCase(status || "unknown")
}

export function statusClass(status) {
  const normalized = String(status || "")
    .trim()
    .toLowerCase()
    .replace(/\s+/g, "-")

  return normalized ? `status-${normalized}` : ""
}

export function formatDate(isoDate) {
  const d = new Date(isoDate)
  if (Number.isNaN(d.getTime())) {
    return isoDate
  }

  return new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "short"
  }).format(d)
}

export function setSelectOptions(selectEl, items, valueKey, labelBuilder, placeholder) {
  selectEl.innerHTML = ""

  const first = document.createElement("option")
  first.value = ""
  first.textContent = placeholder
  selectEl.appendChild(first)

  items.forEach((item) => {
    const opt = document.createElement("option")
    opt.value = String(item[valueKey])
    opt.textContent = labelBuilder(item)
    selectEl.appendChild(opt)
  })
}

function createCheckItem({ label, value, checked, inputName, onChange }) {
  const wrap = document.createElement("label")
  wrap.className = `check-item${checked ? " active" : ""}`

  const input = document.createElement("input")
  input.type = "checkbox"
  input.name = inputName
  input.value = String(value)
  input.checked = checked
  input.addEventListener("change", () => {
    wrap.classList.toggle("active", input.checked)
    onChange(input.checked)
  })

  const text = document.createElement("span")
  text.textContent = label

  wrap.appendChild(input)
  wrap.appendChild(text)
  return wrap
}

export function renderSportChips({ target, sports, activeID, includeAll, onSelect, onError = console.error }) {
  target.innerHTML = ""

  if (includeAll) {
    const allBtn = document.createElement("button")
    allBtn.type = "button"
    allBtn.className = `chip${activeID === null ? " active" : ""}`
    allBtn.textContent = "All"
    allBtn.addEventListener("click", runAsync(() => onSelect(null), onError))
    target.appendChild(allBtn)
  }

  sports.forEach((sport) => {
    const btn = document.createElement("button")
    btn.type = "button"
    btn.className = `chip${activeID === sport.sport_id ? " active" : ""}`
    btn.textContent = titleCase(sport.name)
    btn.addEventListener("click", runAsync(() => onSelect(sport.sport_id), onError))
    target.appendChild(btn)
  })
}

export function renderTeamFilter({ target, sportID, teams, selectedTeamIDs, onToggle }) {
  target.innerHTML = ""

  if (!sportID) {
    target.textContent = "Pick a sport to load teams."
    return
  }

  if (!teams.length) {
    target.textContent = "No teams available for selected sport."
    return
  }

  teams.forEach((team) => {
    const label = team.abbreviation ? `${team.name} (${team.abbreviation})` : team.name
    const item = createCheckItem({
      label,
      value: team.team_id,
      checked: selectedTeamIDs.includes(team.team_id),
      inputName: "team-filter",
      onChange: (checked) => onToggle(team.team_id, checked)
    })
    target.appendChild(item)
  })
}

function makeLogoNode(participant) {
  const fallback = participant.abbreviation || (participant.name || "?").slice(0, 2).toUpperCase()

  if (!participant.logo) {
    const span = document.createElement("span")
    span.className = "logo"
    span.textContent = fallback
    return span
  }

  const img = document.createElement("img")
  img.className = "logo"
  img.src = participant.logo
  img.alt = participant.name || "Team logo"
  img.addEventListener("error", () => {
    const replacement = document.createElement("span")
    replacement.className = "logo"
    replacement.textContent = fallback
    img.replaceWith(replacement)
  })
  return img
}

export function renderEvents(target, countTarget, events) {
  target.innerHTML = ""

  if (!events.length) {
    target.innerHTML = '<article class="event-card"><p>No events found for selected filter.</p></article>'
    countTarget.textContent = "0 events"
    return
  }

  events.forEach((event) => {
    const card = document.createElement("article")
    card.className = "event-card"

    const status = normalizeStatus(event.status)
    const statusCls = statusClass(event.status)
    card.innerHTML = `
      <div class="event-meta">
        <span class="badge">${titleCase(event.sport_name || "Unknown sport")}</span>
        <span class="badge ${statusCls}">${status}</span>
      </div>
      <h4>${event.competition_name || "Unknown competition"}</h4>
      <p>${formatDate(event.start_time)}</p>
    `

    const participants = document.createElement("div")
    participants.className = "participants"

    const list = Array.isArray(event.participants) ? event.participants : []
    if (!list.length) {
      participants.innerHTML = '<div class="participant"><span class="logo">--</span><span>TBD participants</span></div>'
    } else {
      list.forEach((participant) => {
        const row = document.createElement("div")
        row.className = "participant"
        row.appendChild(makeLogoNode(participant))

        const text = document.createElement("span")
        text.textContent = participant.abbreviation || participant.name || "Unknown"
        row.appendChild(text)

        participants.appendChild(row)
      })
    }

    card.appendChild(participants)
    target.appendChild(card)
  })

  countTarget.textContent = `${events.length} event${events.length === 1 ? "" : "s"}`
}
