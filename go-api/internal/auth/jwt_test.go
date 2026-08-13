package auth

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret"

func newTestApp(secret string) *fiber.App {
	app := fiber.New()
	app.Get("/protected", Middleware(secret), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	return app
}

func doRequest(t *testing.T, app *fiber.App, authHeader string) int {
	t.Helper()
	req := httptest.NewRequest("GET", "/protected", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	return resp.StatusCode
}

func TestIssueServiceToken_AcceptedByMiddleware(t *testing.T) {
	token, err := IssueServiceToken(testSecret)
	if err != nil {
		t.Fatalf("IssueServiceToken error: %v", err)
	}

	app := newTestApp(testSecret)
	status := doRequest(t, app, "Bearer "+token)
	if status != fiber.StatusOK {
		t.Errorf("status = %d, want 200", status)
	}
}

func TestMiddleware_MissingHeader(t *testing.T) {
	app := newTestApp(testSecret)
	status := doRequest(t, app, "")
	if status != fiber.StatusUnauthorized {
		t.Errorf("status = %d, want 401", status)
	}
}

func TestMiddleware_MalformedHeader(t *testing.T) {
	app := newTestApp(testSecret)
	status := doRequest(t, app, "NotBearer sometoken")
	if status != fiber.StatusUnauthorized {
		t.Errorf("status = %d, want 401", status)
	}
}

func TestMiddleware_WrongSecret(t *testing.T) {
	token, err := IssueServiceToken("other-secret")
	if err != nil {
		t.Fatalf("IssueServiceToken error: %v", err)
	}

	app := newTestApp(testSecret)
	status := doRequest(t, app, "Bearer "+token)
	if status != fiber.StatusUnauthorized {
		t.Errorf("status = %d, want 401", status)
	}
}

func TestMiddleware_ExpiredToken(t *testing.T) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   "go-api-service",
		IssuedAt:  jwt.NewNumericDate(now.Add(-10 * time.Minute)),
		ExpiresAt: jwt.NewNumericDate(now.Add(-5 * time.Minute)),
	}
	expired := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := expired.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("sign error: %v", err)
	}

	app := newTestApp(testSecret)
	status := doRequest(t, app, "Bearer "+token)
	if status != fiber.StatusUnauthorized {
		t.Errorf("status = %d, want 401", status)
	}
}
