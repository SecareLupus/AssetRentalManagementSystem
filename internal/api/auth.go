package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/desmond/rental-management-system/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("YOUR_SECRET_KEY") // TODO: Load from env

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

// Login authenticates a user and returns a JWT token.
// @Summary Login
// @Description Authenticates a user and returns a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Credentials"
// @Success 200 {object} AuthResponse
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt for username: %s", req.Username)
	user, err := h.repo.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		log.Printf("User not found: %s", req.Username)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if !user.IsEnabled {
		http.Error(w, "account disabled", http.StatusForbidden)
		return
	}

	// Generate Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	// Update Last Login
	now := time.Now()
	user.LastLoginAt = &now
	h.repo.UpdateUser(r.Context(), user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token: tokenString,
		User:  *user,
	})
}

// Register creates a new user account.
// @Summary Register
// @Description Creates a new user account.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Registration Info"
// @Success 201 {string} string "User created"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user.Username = req.Username
	user.Email = req.Email
	user.PasswordHash = string(hash)
	user.Role = domain.UserRoleViewer // Default role
	user.IsEnabled = true

	if err := h.repo.CreateUser(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Middleware

type contextKey string

const UserContextKey contextKey = "user"

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Admin Handlers

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/admin/users/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var u domain.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u.ID = id
	if err := h.repo.UpdateUser(r.Context(), &u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/v1/admin/users/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.repo.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.repo.GetSettings(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(settings)
}

func (h *Handler) UpdateSetting(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key   string          `json:"key"`
		Value json.RawMessage `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.repo.UpdateSetting(r.Context(), req.Key, req.Value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
