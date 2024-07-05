package user_handlers

import (
	"encoding/json"
	"log"
	"main/internal/database"
	"main/internal/models"
	"net/http"
	"strconv"
	"time"
)

// Workload represents the workload for a task
type Workload struct {
	TaskName string `json:"task_name"`
	Duration string `json:"duration"`
}

// GetUsers
// @Summary Get all users
// @Description Get a list of all users with filtering and pagination
// @Tags users
// @Produce json
// @Param first_name query string false "First name"
// @Param last_name query string false "Last name"
// @Param email query string false "Email"
// @Param passport_number query string false "Passport number"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of items per page"
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Router /users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Инициализация соединения с базой данных
	db := database.DB

	// Получение параметров запроса для фильтрации и пагинации
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")
	email := r.URL.Query().Get("email")
	passportNumber := r.URL.Query().Get("passport_number")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Создание SQL-запроса с фильтрацией
	query := "SELECT id, first_name, last_name, email, passport_number FROM users WHERE 1=1"
	var args []interface{}

	if firstName != "" {
		query += " AND first_name ILIKE ?"
		args = append(args, "%"+firstName+"%")
	}
	if lastName != "" {
		query += " AND last_name ILIKE ?"
		args = append(args, "%"+lastName+"%")
	}
	if email != "" {
		query += " AND email ILIKE ?"
		args = append(args, "%"+email+"%")
	}
	if passportNumber != "" {
		query += " AND passport_number ILIKE ?"
		args = append(args, "%"+passportNumber+"%")
	}

	// Добавление пагинации
	offset := (page - 1) * pageSize
	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	// Выполнение запроса к базе данных
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Println("Failed to execute query:", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Слайс для хранения пользователей
	var users []models.User

	// Обработка результатов запроса
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PassportNumber); err != nil {
			log.Println("Failed to scan row:", err)
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Обработка ошибок после завершения итерации по результатам
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Отправка ответа в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserWorkload
// @Summary Get user workload
// @Description Get user workload for a specific period
// @Tags users
// @Produce json
// @Param user_id query int true "User ID"
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {array} Workload
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/workload [get]
func GetUserWorkload(w http.ResponseWriter, r *http.Request) {
	// Получение параметров запроса
	userIDStr := r.URL.Query().Get("user_id")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// Проверка валидности параметров
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID < 1 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Invalid start date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Invalid end date format", http.StatusBadRequest)
		return
	}

	// Инициализация соединения с базой данных
	db := database.DB

	// Запрос к базе данных для получения задач
	query := `
		SELECT name, SUM(EXTRACT(EPOCH FROM (end_time - start_time)) / 3600) AS hours
		FROM tasks
		WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3
		GROUP BY name
		ORDER BY hours DESC`
	rows, err := db.Query(query, userID, startDate, endDate)
	if err != nil {
		log.Println("Failed to execute query:", err)
		http.Error(w, "Failed to fetch user workload", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Слайс для хранения трудозатрат
	var workloads []Workload

	// Обработка результатов запроса
	for rows.Next() {
		var workload Workload
		var hours float64
		if err := rows.Scan(&workload.TaskName, &hours); err != nil {
			log.Println("Failed to scan row:", err)
			http.Error(w, "Failed to fetch user workload", http.StatusInternalServerError)
			return
		}
		// Преобразование часов в формат "ЧЧ:ММ"
		h := int(hours)
		m := int((hours - float64(h)) * 60)
		workload.Duration = strconv.Itoa(h) + "h " + strconv.Itoa(m) + "m"
		workloads = append(workloads, workload)
	}

	// Обработка ошибок после завершения итерации по результатам
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, "Failed to fetch user workload", http.StatusInternalServerError)
		return
	}

	// Отправка ответа в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workloads)
}

// CreateUser
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User to create"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /usercreate [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUserRequest struct {
		PassportNumber string `json:"passportNumber"`
	}

	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if createUserRequest.PassportNumber == "" {
		http.Error(w, "Missing passport number", http.StatusBadRequest)
		return
	}

	// Создание нового пользователя с автоматически заполненными полями
	newUser := models.User{
		FirstName:      "Auto-generated first name", // Пример автозаполнения
		LastName:       "Auto-generated last name",  // Пример автозаполнения
		Email:          "autogenerated@example.com", // Пример автозаполнения
		PassportNumber: createUserRequest.PassportNumber,
	}

	// Вставка нового пользователя в базу данных
	db := database.DB
	_, err = db.Exec("INSERT INTO users (first_name, last_name, email, passport_number) VALUES ($1, $2, $3, $4)",
		newUser.FirstName, newUser.LastName, newUser.Email, newUser.PassportNumber)
	if err != nil {
		log.Println("Failed to insert user into database:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// DeleteUser удаляет пользователя по ID
// @Summary Delete a user by ID
// @Description Delete a user from the database based on the provided ID
// @Tags users
// @Produce json
// @Param id path int true "User ID to delete"
// @Success 200 {string} string "User deleted successfully"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /userdelete/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр ID из пути запроса
	idStr := r.URL.Path[len("/userdelete/"):]
	log.Println(idStr)
	id, err := strconv.Atoi(idStr)
	log.Println(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Удаление пользователя из базы данных
	db := database.DB
	result, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Println("Failed to delete user:", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Проверяем количество удаленных записей
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed to get rows affected:", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User deleted successfully")
}

// UpdateUser
// @Summary Update user data
// @Description Update the details of an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /userupdate/{userId} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/userupdate/"):]
	userId, err := strconv.Atoi(idStr)
	log.Println(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := database.DB
	_, err = db.Exec("UPDATE users SET first_name = $1, last_name = $2, email = $3, passport_number = $4 WHERE id = $5",
		user.FirstName, user.LastName, user.Email, user.PassportNumber, userId)
	if err != nil {
		log.Println("Failed to update user:", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
