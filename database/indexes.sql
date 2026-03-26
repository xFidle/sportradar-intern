/* FILTERING QUERY INDEX */
CREATE INDEX idx_events_start_time ON events (start_time);

/* FOREIGN KEY INDEXES */
CREATE INDEX idx_events_competition_id ON events (_competition_id);
CREATE INDEX idx_events_venue_id ON events (_venue_id);
CREATE INDEX idx_competition_sport_id ON competitions (_sport_id);
CREATE INDEX idx_teams_sport_id ON teams (_sport_id);
CREATE INDEX idx_players_team_id ON players (_team_id);
CREATE INDEX idx_cities_country_id ON cities (_country_id);
CREATE INDEX idx_teams_city_id ON teams (_city_id);
CREATE INDEX idx_venues_city_id ON venues (_city_id);
CREATE INDEX idx_participants_team_event ON participants (_event_id, _team_id);
CREATE INDEX idx_scores_participant_id ON scores (_participant_id);
