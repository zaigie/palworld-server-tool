package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zaigie/palworld-server-tool/internal/auth"
	"github.com/zaigie/palworld-server-tool/internal/config"
)

func TestFirstVisitorInitializesAdministratorAndReadsSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store, err := config.Open(filepath.Join(t.TempDir(), "config.db"))
	if err != nil {
		t.Fatalf("open config store: %v", err)
	}
	defer store.Close()
	config.SetCurrent(store)

	router := gin.New()
	schedulerStarts := 0
	RegisterRouter(router, func() {
		schedulerStarts++
	})

	status := performJSONRequest(router, http.MethodGet, "/api/config/status", nil, "")
	if status.Code != http.StatusOK {
		t.Fatalf("configuration status code = %d, want 200: %s", status.Code, status.Body.String())
	}
	var statusBody struct {
		Initialized bool `json:"initialized"`
	}
	if err := json.Unmarshal(status.Body.Bytes(), &statusBody); err != nil {
		t.Fatalf("decode configuration status: %v", err)
	}
	if statusBody.Initialized {
		t.Fatal("fresh config.db must report uninitialized")
	}

	initialized := performJSONRequest(router, http.MethodPost, "/api/config/initialize", map[string]string{
		"password": "first-admin-password",
	}, "")
	if initialized.Code != http.StatusOK {
		t.Fatalf("initialize code = %d, want 200: %s", initialized.Code, initialized.Body.String())
	}
	var initBody struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(initialized.Body.Bytes(), &initBody); err != nil {
		t.Fatalf("decode initialize response: %v", err)
	}
	if initBody.Token == "" {
		t.Fatal("successful initialization must return an administrator token")
	}
	if schedulerStarts != 1 {
		t.Fatalf("scheduler start callbacks = %d, want 1", schedulerStarts)
	}

	settings := performJSONRequest(router, http.MethodGet, "/api/config", nil, initBody.Token)
	if settings.Code != http.StatusOK {
		t.Fatalf("settings code = %d, want 200: %s", settings.Code, settings.Body.String())
	}
	var settingsBody config.Config
	if err := json.Unmarshal(settings.Body.Bytes(), &settingsBody); err != nil {
		t.Fatalf("decode settings response: %v", err)
	}
	if settingsBody.Web.Port != 8080 || settingsBody.Save.SourceMode != "directory" {
		t.Fatalf("unexpected default settings: %#v", settingsBody)
	}
}

func TestConfigPortOverrideIsReportedAndCannotBeChanged(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store, err := config.Open(filepath.Join(t.TempDir(), "config.db"))
	if err != nil {
		t.Fatalf("open config store: %v", err)
	}
	defer store.Close()
	config.SetCurrent(store)
	if err := store.Initialize("admin-password"); err != nil {
		t.Fatalf("initialize administrator: %v", err)
	}
	token, err := auth.GenerateToken()
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	overridden := store.Config()
	overridden.Web.Port = 18080
	overridden.Web.PortSource = config.WebPortOverrideEnvironment
	if err := store.Update(overridden, ""); err != nil {
		t.Fatalf("persist overridden web port: %v", err)
	}

	router := gin.New()
	RegisterRouter(router, nil)
	response := performJSONRequest(router, http.MethodGet, "/api/config", nil, token)
	if response.Code != http.StatusOK {
		t.Fatalf("settings code = %d, want 200: %s", response.Code, response.Body.String())
	}
	var body struct {
		Web config.WebConfig `json:"web"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode settings response: %v", err)
	}
	if body.Web.Port != 18080 {
		t.Fatalf("persisted web port = %d, want 18080", body.Web.Port)
	}
	if body.Web.PortSource != config.WebPortOverrideEnvironment {
		t.Fatalf("web port source = %q, want %q", body.Web.PortSource, config.WebPortOverrideEnvironment)
	}

	attemptedUpdate := store.Config()
	attemptedUpdate.Web.Port = 28080
	response = performJSONRequest(router, http.MethodPut, "/api/config", map[string]any{
		"settings": attemptedUpdate,
	}, token)
	if response.Code != http.StatusConflict {
		t.Fatalf("settings update code = %d, want 409: %s", response.Code, response.Body.String())
	}
	if got := store.Config().Web.Port; got != 18080 {
		t.Fatalf("persisted web port after overridden update = %d, want 18080", got)
	}

	allowedUpdate := store.Config()
	allowedUpdate.Web.PortSource = config.WebPortOverrideNone
	allowedUpdate.Rcon.Address = "game-host:25575"
	response = performJSONRequest(router, http.MethodPut, "/api/config", map[string]any{
		"settings": allowedUpdate,
	}, token)
	if response.Code != http.StatusOK {
		t.Fatalf("non-port settings update code = %d, want 200: %s", response.Code, response.Body.String())
	}
	persisted := store.Config()
	if persisted.Web.PortSource != config.WebPortOverrideEnvironment {
		t.Fatalf("persisted web port source = %q, want %q", persisted.Web.PortSource, config.WebPortOverrideEnvironment)
	}
	if persisted.Rcon.Address != "game-host:25575" {
		t.Fatalf("persisted RCON address = %q, want game-host:25575", persisted.Rcon.Address)
	}
}

func TestChangingAdministratorPasswordInvalidatesExistingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store, err := config.Open(filepath.Join(t.TempDir(), "config.db"))
	if err != nil {
		t.Fatalf("open config store: %v", err)
	}
	defer store.Close()
	config.SetCurrent(store)
	if err := store.Initialize("old-password"); err != nil {
		t.Fatalf("initialize administrator: %v", err)
	}

	router := gin.New()
	RegisterRouter(router, nil)
	login := performJSONRequest(router, http.MethodPost, "/api/login", map[string]string{"password": "old-password"}, "")
	var loginBody struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(login.Body.Bytes(), &loginBody); err != nil {
		t.Fatalf("decode login response: %v", err)
	}

	update := performJSONRequest(router, http.MethodPut, "/api/config", map[string]any{
		"settings":     config.Default(),
		"new_password": "new-password",
	}, loginBody.Token)
	if update.Code != http.StatusOK {
		t.Fatalf("update code = %d, want 200: %s", update.Code, update.Body.String())
	}
	var updateBody struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(update.Body.Bytes(), &updateBody); err != nil {
		t.Fatalf("decode update response: %v", err)
	}
	if updateBody.Token == "" {
		t.Fatal("password change must return a replacement token")
	}

	oldToken := performJSONRequest(router, http.MethodGet, "/api/config", nil, loginBody.Token)
	if oldToken.Code != http.StatusUnauthorized {
		t.Fatalf("old token code = %d, want 401", oldToken.Code)
	}
	newToken := performJSONRequest(router, http.MethodGet, "/api/config", nil, updateBody.Token)
	if newToken.Code != http.StatusOK {
		t.Fatalf("replacement token code = %d, want 200: %s", newToken.Code, newToken.Body.String())
	}
}

func TestAdministratorCanBrowseHostDirectories(t *testing.T) {
	gin.SetMode(gin.TestMode)
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "Pal Saved"), 0755); err != nil {
		t.Fatalf("create test directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "not-a-directory.txt"), []byte("ignored"), 0644); err != nil {
		t.Fatalf("create test file: %v", err)
	}
	store, err := config.Open(filepath.Join(t.TempDir(), "config.db"))
	if err != nil {
		t.Fatalf("open config store: %v", err)
	}
	defer store.Close()
	config.SetCurrent(store)
	if err := store.Initialize("admin-password"); err != nil {
		t.Fatalf("initialize administrator: %v", err)
	}
	token, err := auth.GenerateToken()
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	router := gin.New()
	RegisterRouter(router, nil)
	response := performJSONRequest(router, http.MethodGet, "/api/config/directories?path="+url.QueryEscape(root), nil, token)
	if response.Code != http.StatusOK {
		t.Fatalf("directory listing code = %d, want 200: %s", response.Code, response.Body.String())
	}
	var body struct {
		Current string `json:"current"`
		Entries []struct {
			Name string `json:"name"`
			Path string `json:"path"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode directory listing: %v", err)
	}
	if body.Current != root {
		t.Fatalf("current directory = %q, want %q", body.Current, root)
	}
	if len(body.Entries) != 1 || body.Entries[0].Name != "Pal Saved" {
		t.Fatalf("directory entries = %#v, want only Pal Saved", body.Entries)
	}
}

func TestConfigurationConnectionStatuses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store, err := config.Open(filepath.Join(t.TempDir(), "config.db"))
	if err != nil {
		t.Fatalf("open config store: %v", err)
	}
	defer store.Close()
	config.SetCurrent(store)
	if err := store.Initialize("admin-password"); err != nil {
		t.Fatalf("initialize administrator: %v", err)
	}
	token, err := auth.GenerateToken()
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	saveDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(saveDir, "Level.sav"), []byte("fixture"), 0644); err != nil {
		t.Fatalf("create Level.sav: %v", err)
	}
	router := gin.New()
	RegisterRouter(router, nil)

	saveStatus := performJSONRequest(router, http.MethodPost, "/api/config/test/save", map[string]any{
		"save": map[string]any{"source_mode": "directory", "path": saveDir},
	}, token)
	if saveStatus.Code != http.StatusOK {
		t.Fatalf("save status code = %d, want 200: %s", saveStatus.Code, saveStatus.Body.String())
	}
	var saveBody struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(saveStatus.Body.Bytes(), &saveBody); err != nil {
		t.Fatalf("decode save status: %v", err)
	}
	if saveBody.Status != "normal" {
		t.Fatalf("save status = %q, want normal: %s", saveBody.Status, saveStatus.Body.String())
	}

	rconStatus := performJSONRequest(router, http.MethodPost, "/api/config/test/rcon", map[string]any{
		"rcon": map[string]any{},
	}, token)
	if rconStatus.Code != http.StatusOK {
		t.Fatalf("rcon status code = %d, want 200: %s", rconStatus.Code, rconStatus.Body.String())
	}
	var rconBody struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(rconStatus.Body.Bytes(), &rconBody); err != nil {
		t.Fatalf("decode rcon status: %v", err)
	}
	if rconBody.Status != "unconfigured" {
		t.Fatalf("rcon status = %q, want unconfigured", rconBody.Status)
	}
}

func TestConfigUpdateOnlyRequiresRestartForStartupAndScheduleFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store, err := config.Open(filepath.Join(t.TempDir(), "config.db"))
	if err != nil {
		t.Fatalf("open config store: %v", err)
	}
	defer store.Close()
	config.SetCurrent(store)
	if err := store.Initialize("admin-password"); err != nil {
		t.Fatalf("initialize administrator: %v", err)
	}
	token, err := auth.GenerateToken()
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	router := gin.New()
	RegisterRouter(router, nil)

	immediate := store.Config()
	immediate.Rcon.Address = "game-host:25575"
	response := performJSONRequest(router, http.MethodPut, "/api/config", map[string]any{"settings": immediate}, token)
	if response.Code != http.StatusOK {
		t.Fatalf("immediate update code = %d: %s", response.Code, response.Body.String())
	}
	var immediateBody struct {
		RestartRequired bool `json:"restart_required"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &immediateBody); err != nil {
		t.Fatalf("decode immediate update: %v", err)
	}
	if immediateBody.RestartRequired {
		t.Fatal("RCON-only update must apply immediately without restart")
	}

	restart := store.Config()
	restart.Web.Port = 18080
	restart.Save.SyncInterval = 300
	response = performJSONRequest(router, http.MethodPut, "/api/config", map[string]any{"settings": restart}, token)
	if response.Code != http.StatusOK {
		t.Fatalf("restart update code = %d: %s", response.Code, response.Body.String())
	}
	var restartBody struct {
		RestartRequired bool     `json:"restart_required"`
		RestartFields   []string `json:"restart_fields"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &restartBody); err != nil {
		t.Fatalf("decode restart update: %v", err)
	}
	if !restartBody.RestartRequired {
		t.Fatal("web port and save schedule changes must require restart")
	}
	wantFields := map[string]bool{"web.port": true, "save.sync_interval": true}
	for _, field := range restartBody.RestartFields {
		delete(wantFields, field)
	}
	if len(wantFields) != 0 {
		t.Fatalf("missing restart fields: %v; response=%s", wantFields, response.Body.String())
	}
}

func performJSONRequest(router http.Handler, method, path string, body any, token string) *httptest.ResponseRecorder {
	var payload bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&payload).Encode(body); err != nil {
			panic(err)
		}
	}
	req := httptest.NewRequest(method, path, &payload)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	return response
}
