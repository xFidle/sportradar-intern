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

function createLogoBadge({ logoUrl, fallback, className = "logo" }) {
  const fallbackText = fallback || "?"

  if (!logoUrl) {
    const span = document.createElement("span")
    span.className = className
    span.textContent = fallbackText
    return span
  }

  const img = document.createElement("img")
  img.className = className
  img.src = logoUrl
  img.alt = "Logo"
  img.addEventListener("error", () => {
    const replacement = document.createElement("span")
    replacement.className = className
    replacement.textContent = fallbackText
    img.replaceWith(replacement)
  })
  return img
}

function createCheckItem({ label, value, checked, inputName, onChange, logoUrl, fallback }) {
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
  text.className = "option-text"
  text.textContent = label

  wrap.appendChild(input)
  if (logoUrl || fallback) {
    wrap.appendChild(createLogoBadge({ logoUrl, fallback, className: "option-logo" }))
  }
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
      onChange: (checked) => onToggle(team.team_id, checked),
      logoUrl: team.logo,
      fallback: team.abbreviation || (team.name || "?").slice(0, 2).toUpperCase()
    })
    target.appendChild(item)
  })
}

function makeLogoNode(participant) {
  const fallback = participant.abbreviation || (participant.name || "?").slice(0, 2).toUpperCase()
  return createLogoBadge({ logoUrl: participant.logo, fallback, className: "logo" })
}

function makeCompetitionLogoNode(competition) {
  const name = competition?.name || "Unknown competition"
  const fallback = name.slice(0, 2).toUpperCase()
  return createLogoBadge({ logoUrl: competition?.logo, fallback, className: "competition-logo" })
}


export function renderEvents(target, countTarget, events) {
  target.innerHTML = ""

  if (!events.length) {
    target.innerHTML = '<article class="event-card"><p>No events found for selected filter.</p></article>'
    countTarget.textContent = "0 events"
    return
  }

  events.sort((a, b) => new Date(a.start_time) - new Date(b.start_time))

  events.forEach((event) => {
    const card = document.createElement("article")
    card.className = "event-card"

    const status = normalizeStatus(event.status)
    const statusCls = statusClass(event.status)
    const competition = event.competition || {}
    const competitionName = competition.name || "Unknown competition"
    card.innerHTML = `
      <div class="event-meta">
        <span class="badge">${titleCase(event.sport_name || "Unknown sport")}</span>
        <span class="badge ${statusCls}">${status}</span>
      </div>
      <div class="competition-row"></div>
      <p>${formatDate(event.start_time)}</p>
    `

    const competitionRow = card.querySelector(".competition-row")
    competitionRow.appendChild(makeCompetitionLogoNode(competition))
    const competitionText = document.createElement("h4")
    competitionText.textContent = competitionName
    competitionRow.appendChild(competitionText)

    const participants = document.createElement("div")
    participants.className = "participants"
    const isFinished = String(event.status || "").toLowerCase() === "finished"
    const scoreByTeamID = new Map(
      (Array.isArray(event.final_scores) ? event.final_scores : []).map((score) => [Number(score.team_id), score.agg_score])
    )

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

        if (isFinished) {
          const score = scoreByTeamID.get(Number(participant.team_id))
          const scoreText = document.createElement("span")
          scoreText.className = "participant-score"
          scoreText.textContent = score === undefined ? "-" : String(score)
          row.appendChild(scoreText)
        }

        participants.appendChild(row)
      })
    }

    card.appendChild(participants)
    target.appendChild(card)
  })

  countTarget.textContent = `${events.length} event${events.length === 1 ? "" : "s"}`
}
