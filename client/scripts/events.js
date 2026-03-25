import { apiFetch } from "./api.js"
import { dom } from "./dom.js"
import { state } from "./state.js"
import {
  runAsync,
  setMessage,
  titleCase,
  setSelectOptions,
  renderSportChips,
  renderTeamFilter,
  renderEvents
} from "./ui.js"

// INITIAL LOADERS 

async function loadSports() {
  const sports = await apiFetch("/api/sports/")
  state.sports = Array.isArray(sports) ? sports : []

  renderFilterSections()
  setSelectOptions(dom.sportCreate, state.sports, "sport_id", (sport) => titleCase(sport.name), "Select sport")
}

async function loadEvents() {
  const payload = {
    start_after: dom.startAfter.value,
    end_before: dom.endBefore.value
  }

  if (!payload.start_after || !payload.end_before) {
    throw new Error("Both filter dates are required.")
  }

  if (state.filter.sportID !== null) {
    payload.sport_id = state.filter.sportID
  }

  if (state.filter.status !== null) {
    payload.status = state.filter.status
  }

  if (state.filter.teamIDs.length > 0) {
    payload.team_ids = state.filter.teamIDs
  }

  const events = await apiFetch("/api/events/filter", {
    method: "POST",
    body: JSON.stringify(payload)
  })

  renderEvents(dom.eventsList, dom.eventsCount, Array.isArray(events) ? events : [])
}

// FILTERING

const statusOptions = [
  { value: null, label: "All" },
  { value: "scheduled", label: "Scheduled" },
  { value: "finished", label: "Finished" }
]

function setDefaultFilterDates() {
  const start = new Date("2026-01-01")
  const finish = new Date("2027-01-01")

  dom.startAfter.value = start.toISOString().slice(0, 10)
  dom.endBefore.value = finish.toISOString().slice(0, 10)
}

function renderStatusFilter() {
  dom.statusFilterChips.innerHTML = ""

  statusOptions.forEach((option) => {
    const btn = document.createElement("button")
    btn.type = "button"
    btn.className = `chip${state.filter.status === option.value ? " active" : ""}`
    btn.textContent = option.label
    btn.addEventListener("click", () => onFilterStatusSelect(option.value))
    dom.statusFilterChips.appendChild(btn)
  })
}

function onFilterStatusSelect(status) {
  state.filter.status = status
  renderStatusFilter()
}

function renderFilterSections() {
  renderSportChips({
    target: dom.sportFilterChips,
    sports: state.sports,
    activeID: state.filter.sportID,
    includeAll: true,
    onSelect: onFilterSportSelect
  })

  renderStatusFilter()

  renderTeamFilter({
    target: dom.teamFilterList,
    sportID: state.filter.sportID,
    teams: state.filterOptions.teams,
    selectedTeamIDs: state.filter.teamIDs,
    onToggle: (teamID, checked) => {
      if (checked) {
        state.filter.teamIDs.push(teamID)
        return
      }
      state.filter.teamIDs = state.filter.teamIDs.filter((id) => id !== teamID)
    }
  })
}

async function onFilterSportSelect(sportID) {
  state.filter.sportID = sportID
  state.filter.teamIDs = []
  state.filterOptions.teams = []

  if (sportID === null) {
    renderFilterSections()
    return
  }

  const teams = await apiFetch(`/api/sports/${sportID}/teams`)
  state.filterOptions.teams = Array.isArray(teams) ? teams : []
  renderFilterSections()
}

function clearFilters() {
  setDefaultFilterDates()
  state.filter.sportID = null
  state.filter.status = null
  state.filter.teamIDs = []
  state.filterOptions.teams = []
  renderFilterSections()
}

// CREATE FORM

const DEFAULT_STAGE_ID = 1
const teamNameLabel = (team) => (team.abbreviation ? `${team.name} (${team.abbreviation})` : team.name)

function clearCreateSelects({ competition, venue, teams }) {
  if (!competition) {
    setSelectOptions(dom.competitionCreate, [], "competition_id", (c) => c.name, "Select competition")
  }

  if (!venue) {
    setSelectOptions(dom.venueCreate, [], "venue_id", (v) => v.name, "Select venue")
  }

  if (!teams) {
    setSelectOptions(dom.teamOne, [], "team_id", (t) => t.name, "Select team")
    setSelectOptions(dom.teamTwo, [], "team_id", (t) => t.name, "Select team")
  }
}

function disableCreateInputs({ competition, venue, teams }) {
  dom.competitionCreate.disabled = !competition
  dom.venueCreate.disabled = !venue
  dom.teamOne.disabled = !teams
  dom.teamTwo.disabled = !teams
}

function populateCreateSelects({ sport, competition } = {}) {
  if (sport) {
    setSelectOptions(
      dom.competitionCreate,
      state.createOptions.competitions,
      "competition_id",
      (competition) => competition.name,
      "Select competition"
    )

    setSelectOptions(
      dom.venueCreate,
      state.createOptions.venues,
      "venue_id",
      (venue) => `${venue.name} (${venue.city?.name || "Unknown city"})`,
      "Select venue"
    )
  }

  if (competition) {
    setSelectOptions(dom.teamOne, state.createOptions.teams, "team_id", teamNameLabel, "Select team")
    setSelectOptions(dom.teamTwo, state.createOptions.teams, "team_id", teamNameLabel, "Select team")
  }
}

function resetCreateForm() {
  setMessage(dom.createMsg, "", "")
  dom.sportCreate.value = ""

  const disabledState = { competition: false, venue: false, teams: false }
  clearCreateSelects(disabledState)
  disableCreateInputs(disabledState)

  state.createOptions.competitions = []
  state.createOptions.venues = []
  state.createOptions.teams = []
}

function openCreateModal() {
  dom.createModal.classList.add("open")
  dom.createModal.setAttribute("aria-hidden", "false")
  document.body.style.overflow = "hidden"
}

function closeCreateModal() {
  dom.createModal.classList.remove("open")
  dom.createModal.setAttribute("aria-hidden", "true")
  document.body.style.overflow = ""
}

async function onCreateSportChange() {
  setMessage(dom.createMsg, "", "")
  const sportID = Number(dom.sportCreate.value)

  const toSelect = { competition: false, venue: false, teams: false }
  clearCreateSelects(toSelect)
  disableCreateInputs(toSelect)

  if (!sportID) {
    return
  }

  const options = await apiFetch(`/api/sports/${sportID}/event-options`)
  state.createOptions.competitions = Array.isArray(options.competitions) ? options.competitions : []
  state.createOptions.venues = Array.isArray(options.venues) ? options.venues : []

  toSelect.competition = true
  toSelect.venue = true

  const selected = { sport: true }
  populateCreateSelects(selected)
  disableCreateInputs(toSelect)
}

async function onCreateCompetitionChange() {
  setMessage(dom.createMsg, "", "")
  const competitionID = Number(dom.competitionCreate.value)

  const toSelect = { competition: true, venue: true, teams: false }
  clearCreateSelects(toSelect)
  disableCreateInputs(toSelect)

  if (!competitionID) {
    return
  }

  const teams = await apiFetch(`/api/competitions/${competitionID}/teams`)
  state.createOptions.teams = Array.isArray(teams) ? teams : []

  toSelect.teams = true

  const selected = { competition: true }
  populateCreateSelects(selected)
  disableCreateInputs(toSelect)
}

function localDateTimeToRFC3339(localValue) {
  if (!localValue) {
    throw new Error("Start time is required.")
  }

  const [datePart, timePart] = localValue.split("T")
  if (!datePart || !timePart) {
    throw new Error("Start time format is invalid.")
  }

  const [year, month, day] = datePart.split("-").map(Number)
  const [hour, minute] = timePart.split(":").map(Number)
  const d = new Date(year, month - 1, day, hour, minute, 0, 0)

  if (Number.isNaN(d.getTime())) {
    throw new Error("Start time value is invalid.")
  }

  return d.toISOString()
}

async function createEvent() {
  const competitionID = Number(dom.competitionCreate.value)
  const venueID = Number(dom.venueCreate.value)
  const team1 = Number(dom.teamOne.value)
  const team2 = Number(dom.teamTwo.value)

  if (!competitionID || !venueID) {
    throw new Error("Competition and venue are required.")
  }
  if (!team1 || !team2) {
    throw new Error("Please select Team 1 and Team 2.")
  }
  if (team1 === team2) {
    throw new Error("Team 1 and Team 2 must be different.")
  }

  const startTime = localDateTimeToRFC3339(dom.startTimeCreate.value)
  if (new Date(startTime) <= new Date()) {
    throw new Error("Start time must be in the future.")
  }

  dom.createEventBtn.disabled = true
  try {
    await apiFetch("/api/events/", {
      method: "POST",
      body: JSON.stringify({
        competition_id: competitionID,
        venue_id: venueID,
        stage_id: DEFAULT_STAGE_ID,
        start_time: startTime,
        team_ids: [team1, team2]
      })
    })
  } finally {
    dom.createEventBtn.disabled = false
  }

  setMessage(dom.createMsg, "ok", "Event created successfully.")
  await loadEvents()

  setTimeout(() => {
    closeCreateModal()
    resetCreateForm()
  }, 260)
}

function setupDefaults() {
  setDefaultFilterDates()
  const now = new Date()
  const in1h = new Date(now.getTime() + 60 * 60 * 1000)
  const yyyy = in1h.getFullYear()
  const mm = String(in1h.getMonth() + 1).padStart(2, "0")
  const dd = String(in1h.getDate()).padStart(2, "0")
  const hh = String(in1h.getHours()).padStart(2, "0")
  const min = String(in1h.getMinutes()).padStart(2, "0")
  dom.startTimeCreate.value = `${yyyy}-${mm}-${dd}T${hh}:${min}`
}

function bindModalEvents() {
  const closeAndReset = () => {
    closeCreateModal()
    resetCreateForm()
  }

  dom.openCreateModalBtn.addEventListener("click", openCreateModal)
  dom.closeCreateModalBtn.addEventListener("click", closeAndReset)
  dom.cancelCreateBtn.addEventListener("click", closeAndReset)
  dom.createModalBackdrop.addEventListener("click", closeAndReset)

  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && dom.createModal.classList.contains("open")) {
      closeAndReset()
    }
  })
}

function bindFilterEvents() {
  dom.loadEventsBtn.addEventListener("click", runAsync(loadEvents))
  dom.clearFiltersBtn.addEventListener("click", clearFilters)
}

function bindCreateEvents() {
  dom.sportCreate.addEventListener(
    "change",
    runAsync(onCreateSportChange, (err) => setMessage(dom.createMsg, "err", err.message))
  )

  dom.competitionCreate.addEventListener(
    "change",
    runAsync(onCreateCompetitionChange, (err) => setMessage(dom.createMsg, "err", err.message))
  )

  dom.createEventBtn.addEventListener(
    "click",
    runAsync(async () => {
      setMessage(dom.createMsg, "", "")
      await createEvent()
    }, (err) => setMessage(dom.createMsg, "err", err.message))
  )
}

async function bootstrap() {
  setupDefaults()
  resetCreateForm()
  renderFilterSections()

  bindFilterEvents()
  bindModalEvents()
  bindCreateEvents()

  try {
    await loadSports()
    await loadEvents()
  } catch (err) {
    console.error(err)
  }
}

bootstrap()
