package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bootdotdev/learn-cicd-starter/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerNotesGet(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := cfg.DB.GetNotesForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts for user", err)
		return
	}

	postsResp, err := databasePostsToPosts(posts)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert posts", err)
		return
	}

	respondWithJSON(w, http.StatusOK, postsResp)
}

func (cfg *apiConfig) handlerNotesCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Note string `json:"note"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// FIX 1: Convert UUID and Time to Strings if sqlc generated them as strings
	// If sqlc generated them as uuid.UUID/time.Time, remove the .String() and .Format()
	newNote, err := cfg.DB.CreateNote(r.Context(), database.CreateNoteParams{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Note:      params.Note,
		UserID:    user.ID,
	})

	// FIX 2: Added the missing 'err' argument to match the 'want' signature
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create note", err)
		return
	}

	// FIX 3: Use '=' instead of ':=' because 'err' was already declared above
	note, err := cfg.DB.GetNote(r.Context(), newNote.ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get note", err)
		return
	}

	noteResp, err := databaseNoteToNote(note)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert note", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, noteResp)
}
