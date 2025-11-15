package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:8080"

// Генерирует уникальный ID для теста
func uniqueID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

// Тестовые данные (шаблоны, имена будут уникальными)
func createTeamData(teamName string) map[string]interface{} {
	return map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": "u1", "username": "Alice", "is_active": true},
			{"user_id": "u2", "username": "Bob", "is_active": true},
			{"user_id": "u3", "username": "Charlie", "is_active": true},
			{"user_id": "u4", "username": "David", "is_active": false}, // неактивный
		},
	}
}

func TestMainConditions(t *testing.T) {
	t.Run("1. Создание команды и проверка назначения ревьюеров", func(t *testing.T) {
		// Создаём команду с уникальным именем
		teamName := uniqueID("team")
		teamData := createTeamData(teamName)
		createTeam(t, teamData)

		// Создаём PR от u1 с уникальным ID
		prID := uniqueID("pr")
		pr := createPR(t, prID, "Test PR", "u1")

		// Проверяем: назначены до 2 ревьюеров, исключая автора
		if len(pr["assigned_reviewers"].([]interface{})) > 2 {
			t.Errorf("Назначено больше 2 ревьюеров: %d", len(pr["assigned_reviewers"].([]interface{})))
		}

		reviewers := pr["assigned_reviewers"].([]interface{})
		for _, reviewer := range reviewers {
			if reviewer.(string) == "u1" {
				t.Error("Автор не должен быть назначен ревьюером")
			}
		}

		// Проверяем, что неактивный пользователь (u4) не назначен
		for _, reviewer := range reviewers {
			if reviewer.(string) == "u4" {
				t.Error("Неактивный пользователь не должен быть назначен ревьюером")
			}
		}
	})

	t.Run("2. Переназначение ревьюера из команды заменяемого", func(t *testing.T) {
		// Создаём команду с уникальным именем (нужно минимум 4 активных для переназначения)
		teamName := uniqueID("team")
		teamData := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": "u1", "username": "Alice", "is_active": true},
				{"user_id": "u2", "username": "Bob", "is_active": true},
				{"user_id": "u3", "username": "Charlie", "is_active": true},
				{"user_id": "u5", "username": "Eve", "is_active": true},    // дополнительный активный
				{"user_id": "u4", "username": "David", "is_active": false}, // неактивный
			},
		}
		createTeam(t, teamData)

		// Создаём PR с уникальным ID
		prID := uniqueID("pr")
		pr := createPR(t, prID, "Test PR 2", "u1")

		reviewers := pr["assigned_reviewers"].([]interface{})
		if len(reviewers) == 0 {
			t.Skip("Нет ревьюеров для переназначения")
		}

		oldReviewer := reviewers[0].(string)

		// Получаем команду старого ревьюера
		team := getTeam(t, teamName)
		members := team["members"].([]interface{})

		// Переназначаем
		reassignResp := reassignReviewer(t, prID, oldReviewer)

		newReviewer := reassignResp["replaced_by"].(string)

		// Проверяем, что новый ревьюер из той же команды
		found := false
		for _, member := range members {
			memberMap := member.(map[string]interface{})
			if memberMap["user_id"].(string) == newReviewer {
				found = true
				// Проверяем, что он активный
				if !memberMap["is_active"].(bool) {
					t.Error("Новый ревьюер должен быть активным")
				}
				break
			}
		}
		if !found {
			t.Error("Новый ревьюер должен быть из команды заменяемого")
		}
	})

	t.Run("3. Запрет изменения ревьюеров после MERGED", func(t *testing.T) {
		// Создаём команду с уникальным именем
		teamName := uniqueID("team")
		teamData := createTeamData(teamName)
		createTeam(t, teamData)

		prID := uniqueID("pr")
		createPR(t, prID, "Test PR 3", "u1")

		// Merge PR
		mergedPR := mergePR(t, prID)
		if mergedPR["status"].(string) != "MERGED" {
			t.Error("PR должен быть в статусе MERGED")
		}

		// Пытаемся переназначить после merge
		reviewers := mergedPR["assigned_reviewers"].([]interface{})
		if len(reviewers) > 0 {
			oldReviewer := reviewers[0].(string)
			resp, err := http.Post(
				fmt.Sprintf("%s/pullRequest/reassign", baseURL),
				"application/json",
				bytes.NewBuffer(mustJSON(map[string]string{
					"pull_request_id": prID,
					"old_user_id":     oldReviewer,
				})),
			)
			if err != nil {
				t.Fatalf("Ошибка запроса: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusConflict {
				t.Errorf("Ожидался статус 409 (Conflict), получен %d", resp.StatusCode)
			}

			var errorResp map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&errorResp)
			errorCode := errorResp["error"].(map[string]interface{})["code"].(string)
			if errorCode != "PR_MERGED" {
				t.Errorf("Ожидался код ошибки PR_MERGED, получен %s", errorCode)
			}
		}
	})

	t.Run("4. Назначение доступного количества ревьюеров (0/1)", func(t *testing.T) {
		// Создаём команду с одним активным пользователем
		teamName := uniqueID("team-solo")
		teamSmall := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": "s1", "username": "Solo", "is_active": true},
			},
		}
		createTeam(t, teamSmall)

		// Создаём PR от s1 (единственный активный, кроме автора никого нет)
		prID := uniqueID("pr-solo")
		pr := createPR(t, prID, "Solo PR", "s1")

		// Должно быть 0 ревьюеров
		reviewers := pr["assigned_reviewers"].([]interface{})
		if len(reviewers) != 0 {
			t.Errorf("Ожидалось 0 ревьюеров, получено %d", len(reviewers))
		}
	})

	t.Run("5. Неактивный пользователь не назначается на ревью", func(t *testing.T) {
		// Создаём команду с неактивным пользователем
		teamName := uniqueID("team-inactive")
		teamInactive := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": "i1", "username": "Inactive", "is_active": false},
				{"user_id": "i2", "username": "Active", "is_active": true},
			},
		}
		createTeam(t, teamInactive)

		// Создаём PR от i2
		prID := uniqueID("pr-inactive")
		pr := createPR(t, prID, "Inactive Test", "i2")

		reviewers := pr["assigned_reviewers"].([]interface{})
		for _, reviewer := range reviewers {
			if reviewer.(string) == "i1" {
				t.Error("Неактивный пользователь i1 не должен быть назначен ревьюером")
			}
		}
	})

	t.Run("6. Идемпотентность операции merge", func(t *testing.T) {
		// Создаём команду с уникальным именем
		teamName := uniqueID("team")
		teamData := createTeamData(teamName)
		createTeam(t, teamData)

		prID := uniqueID("pr-idempotent")
		createPR(t, prID, "Idempotent Test", "u1")

		// Первый merge
		pr1 := mergePR(t, prID)
		mergedAt1 := pr1["mergedAt"].(string)

		// Второй merge (должен быть идемпотентным)
		pr2 := mergePR(t, prID)

		if pr2["status"].(string) != "MERGED" {
			t.Error("PR должен оставаться в статусе MERGED")
		}

		mergedAt2 := pr2["mergedAt"].(string)
		if mergedAt1 != mergedAt2 {
			t.Error("Повторный merge не должен изменять mergedAt (или должен возвращать то же значение)")
		}
	})
}

// Вспомогательные функции

func createTeam(t *testing.T, team map[string]interface{}) {
	resp, err := http.Post(
		fmt.Sprintf("%s/team/add", baseURL),
		"application/json",
		bytes.NewBuffer(mustJSON(team)),
	)
	if err != nil {
		t.Fatalf("Ошибка создания команды: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		// Игнорируем ошибку, если команда уже существует (может быть при параллельных тестах)
		if resp.StatusCode == http.StatusBadRequest {
			errorCode := errorResp["error"].(map[string]interface{})["code"].(string)
			if errorCode == "TEAM_EXISTS" {
				t.Logf("Команда уже существует, пропускаем создание: %v", team["team_name"])
				return
			}
		}
		t.Fatalf("Ожидался статус 201, получен %d, ошибка: %v", resp.StatusCode, errorResp)
	}
}

func createPR(t *testing.T, prID, prName, authorID string) map[string]interface{} {
	req := map[string]string{
		"pull_request_id":   prID,
		"pull_request_name": prName,
		"author_id":         authorID,
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/pullRequest/create", baseURL),
		"application/json",
		bytes.NewBuffer(mustJSON(req)),
	)
	if err != nil {
		t.Fatalf("Ошибка создания PR: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		// Игнорируем ошибку, если PR уже существует (может быть при параллельных тестах)
		if resp.StatusCode == http.StatusBadRequest {
			errorCode := errorResp["error"].(map[string]interface{})["code"].(string)
			if errorCode == "PR_EXISTS" {
				// Генерируем новый уникальный ID и пробуем снова
				newPRID := uniqueID("pr")
				t.Logf("PR %s уже существует, пробуем с новым ID: %s", prID, newPRID)
				return createPR(t, newPRID, prName, authorID)
			}
		}
		t.Fatalf("Ошибка создания PR: статус %d, ответ %v", resp.StatusCode, errorResp)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result["pr"].(map[string]interface{})
}

func getTeam(t *testing.T, teamName string) map[string]interface{} {
	resp, err := http.Get(fmt.Sprintf("%s/team/get?team_name=%s", baseURL, teamName))
	if err != nil {
		t.Fatalf("Ошибка получения команды: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	var team map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&team)
	return team
}

func mergePR(t *testing.T, prID string) map[string]interface{} {
	req := map[string]string{
		"pull_request_id": prID,
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/pullRequest/merge", baseURL),
		"application/json",
		bytes.NewBuffer(mustJSON(req)),
	)
	if err != nil {
		t.Fatalf("Ошибка merge PR: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		t.Fatalf("Ошибка merge PR: статус %d, ответ %v", resp.StatusCode, errorResp)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result["pr"].(map[string]interface{})
}

func reassignReviewer(t *testing.T, prID, oldReviewerID string) map[string]interface{} {
	req := map[string]string{
		"pull_request_id": prID,
		"old_user_id":     oldReviewerID,
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/pullRequest/reassign", baseURL),
		"application/json",
		bytes.NewBuffer(mustJSON(req)),
	)
	if err != nil {
		t.Fatalf("Ошибка переназначения: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		t.Fatalf("Ошибка переназначения: статус %d, ответ %v", resp.StatusCode, errorResp)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func mustJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

// Проверка доступности сервера перед тестами
func TestServerHealth(t *testing.T) {
	// Ждём немного, чтобы сервер успел запуститься
	time.Sleep(1 * time.Second)

	resp, err := http.Get(fmt.Sprintf("%s/healthz", baseURL))
	if err != nil {
		t.Fatalf("Сервер недоступен: %v. Убедитесь, что сервер запущен на %s", err, baseURL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Health check не прошёл: статус %d", resp.StatusCode)
	}
}
