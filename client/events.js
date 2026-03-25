import { apiFetch } from "./api.js"
import {
  runAsync,
  setMessage,
  titleCase,
  setSelectOptions,
  renderSportChips,
  renderTeamFilter,
  renderEvents
} from "./ui.js"

const dom = {
  startAfter: document.getElementById("startAfter"),
  endBefore: document.getElementById("endBefore"),
  sportFilterChips: document.getElementById("sportFilterChips"),
  competitionFilterList: document.getElementById("competitionFilterList"),
  teamFilterList: document.getElementById("teamFilterList"),
  loadEventsBtn: document.getElementById("loadEventsBtn"),
  clearFiltersBtn: document.getElementById("clearFiltersBtn"),
  eventsList: document.getElementById("eventsList"),
  eventsCount: document.getElementById("eventsCount"),
  openCreateModalBtn: document.getElementById("openCreateModalBtn"),
  createModal: document.getElementById("createModal"),
  createModalBackdrop: document.getElementById("createModalBackdrop"),
  closeCreateModalBtn: document.getElementById("closeCreateModalBtn"),
  cancelCreateBtn: document.getElementById("cancelCreateBtn"),
  sportCreate: document.getElementById("sportCreate"),
  competitionCreate: document.getElementById("competitionCreate"),
  venueCreate: document.getElementById("venueCreate"),
  teamOne: document.getElementById("teamOne"),
  teamTwo: document.getElementById("teamTwo"),
  stageId: document.getElementById("stageId"),
  startTimeCreate: document.getElementById("startTimeCreate"),
  createEventBtn: document.getElementById("createEventBtn"),
  createMsg: document.getElementById("createMsg")
}

const state = {
  sports: [],
  filter: {
    sportID: null,
    teamIDs: []
  },
  filterOptions: {
    teams: []
  },
  createOptions: {
    competitions: [],
    venues: [],
    teams: []
  }
}

function renderFilterSections() {
  renderSportChips({
    target: dom.sportFilterChips,
    sports: state.sports,
    activeID: state.filter.sportID,
    includeAll: true,
    onSelect: onFilterSportSelect
  })

  renderTeamFilter({
    target: dom.teamFilterList,
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

async function loadSports() {
  const sports = await apiFetch("/api/sports/")
  state.sports = Array.isArray(sports) ? sports : []

  renderFilterSections()
  setSelectOptions(dom.sportCreate, state.sports, "sport_id", (sport) => titleCase(sport.name), "Select sport")
}

async function onFilterSportSelect(sportID) {
  state.filter.sportID = sportID
  state.filter.teamIDs = []
  state.filterOptions.teams = []
  renderFilterSections()

  if (sportID === null) {
    return
  }

  const options = await apiFetch(`/api/sports/${sportID}/event-options`)
  state.filterOptions.competitions = Array.isArray(options.competitions) ? options.competitions : []
  renderFilterSections()
}

async function loadTeamsForFilter() {
  if (!state.filter.sportID) {
    state.filterOptions.teams = []
    renderFilterSections()
    return
  }

  const teams = await apiFetch(`/api/sports/${state.filter.sportID}/teams`)
  state.filterOptions.teams = Array.isArray(teams) ? teams : []
  renderFilterSections()
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
  if (state.filter.competitionID !== null) {
    payload.competition_id = state.filter.competitionID
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

function clearFilters() {
  state.filter.sportID = null
  state.filter.competitionID = null
  state.filter.teamIDs = []
  state.filterOptions.competitions = []
  state.filterOptions.teams = []
  renderFilterSections()
}

function resetCreateForm() {
  setMessage(dom.createMsg, "", "")
  dom.sportCreate.value = ""

  setSelectOptions(dom.competitionCreate, [], "competition_id", (c) => c.name, "Select competition")
  setSelectOptions(dom.venueCreate, [], "venue_id", (v) => v.name, "Select venue")
  setSelectOptions(dom.teamOne, [], "team_id", (t) => t.name, "Select team")
  setSelectOptions(dom.teamTwo, [], "team_id", (t) => t.name, "Select team")

  dom.competitionCreate.disabled = true
  dom.venueCreate.disabled = true
  dom.teamOne.disabled = true
  dom.teamTwo.disabled = true

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

  setSelectOptions(dom.competitionCreate, [], "competition_id", (c) => c.name, "Select competition")
  setSelectOptions(dom.venueCreate, [], "venue_id", (v) => v.name, "Select venue")
  setSelectOptions(dom.teamOne, [], "team_id", (t) => t.name, "Select team")
  setSelectOptions(dom.teamTwo, [], "team_id", (t) => t.name, "Select team")

  dom.competitionCreate.disabled = true
  dom.venueCreate.disabled = true
  dom.teamOne.disabled = true
  dom.teamTwo.disabled = true

  if (!sportID) {
    return
  }

  const options = await apiFetch(`/api/sports/${sportID}/event-options`)
  state.createOptions.competitions = Array.isArray(options.competitions) ? options.competitions : []
  state.createOptions.venues = Array.isArray(options.venues) ? options.venues : []

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

  dom.competitionCreate.disabled = false
  dom.venueCreate.disabled = false
}

async function onCreateCompetitionChange() {
  setMessage(dom.createMsg, "", "")
  const competitionID = Number(dom.competitionCreate.value)

  setSelectOptions(dom.teamOne, [], "team_id", (t) => t.name, "Select team")
  setSelectOptions(dom.teamTwo, [], "team_id", (t) => t.name, "Select team")
  dom.teamOne.disabled = true
  dom.teamTwo.disabled = true

  if (!competitionID) {
    return
  }

  const teams = await apiFetch(`/api/competitions/${competitionID}/teams`)
  state.createOptions.teams = Array.isArray(teams) ? teams : []
  const teamLabel = (team) => (team.abbreviation ? `${team.name} (${team.abbreviation})` : team.name)

  setSelectOptions(dom.teamOne, state.createOptions.teams, "team_id", teamLabel, "Select team")
  setSelectOptions(dom.teamTwo, state.createOptions.teams, "team_id", teamLabel, "Select team")
  dom.teamOne.disabled = false
  dom.teamTwo.disabled = false
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
  const stageID = Number(dom.stageId.value)
  const team1 = Number(dom.teamOne.value)
  const team2 = Number(dom.teamTwo.value)

  if (!competitionID || !venueID || !stageID) {
    throw new Error("Competition, venue and stage are required.")
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
        stage_id: stageID,
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
  const now = new Date()
  const in30Days = new Date(now)
  in30Days.setDate(now.getDate() + 30)

  dom.startAfter.value = now.toISOString().slice(0, 10)
  dom.endBefore.value = in30Days.toISOString().slice(0, 10)

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
