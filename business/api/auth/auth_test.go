package auth_test

import (
	"bytes"
	"context"
	"fmt"
	"runtime/debug"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zucchini/services-golang/business/api/auth"
	"github.com/zucchini/services-golang/foundation/logger"
)

func TestAuth(t *testing.T) {
	testCases := []struct {
		name          string
		roles         []string
		tokenUserID   string
		requestUserID string // The userID to use in the Authorize call
		ruleTests     []struct {
			rule     string
			expected bool
			message  string
		}
	}{
		{
			name:          "Admin Role",
			roles:         []string{"ADMIN"},
			tokenUserID:   "5cf37266-3473-4006-984f-9325122678b7",
			requestUserID: "5cf37266-3473-4006-984f-9325122678b7", // Same as token
			ruleTests: []struct {
				rule     string
				expected bool
				message  string
			}{
				{rule: auth.RuleAdminOnly, expected: true, message: "Admin should be authorized for admin-only rule"},
				{rule: auth.RuleUserOnly, expected: false, message: "Admin should not be authorized for user-only rule"},
				{rule: auth.RuleAdminOrSubject, expected: true, message: "Admin should be authorized for admin-or-subject rule"},
			},
		},
		{
			name:          "User Role",
			roles:         []string{"USER"},
			tokenUserID:   "5cf37266-3473-4006-984f-9325122678b7",
			requestUserID: "5cf37266-3473-4006-984f-9325122678b7", // Same as token
			ruleTests: []struct {
				rule     string
				expected bool
				message  string
			}{
				{rule: auth.RuleAdminOnly, expected: false, message: "User should not be authorized for admin-only rule"},
				{rule: auth.RuleUserOnly, expected: true, message: "User should be authorized for user-only rule"},
				{rule: auth.RuleAdminOrSubject, expected: true, message: "User should be authorized for admin-or-subject rule"},
				{rule: auth.RuleAny, expected: true, message: "User should be authorized for any rule"},
			},
		},
		{
			name:          "User ID Mismatch",
			roles:         []string{"USER"},
			tokenUserID:   "5cf37266-3473-4006-984f-9325122678b7",
			requestUserID: "d5c65266-3473-4006-984f-9325122678c8", // Different from token
			ruleTests: []struct {
				rule     string
				expected bool
				message  string
			}{
				{rule: auth.RuleAdminOrSubject, expected: false, message: "Authorization should fail when userID doesn't match token subject"},
			},
		},
		{
			name:          "Different UserID with RuleAny - Allows access to any authenticated user with valid roles",
			roles:         []string{"USER", "ADMIN"},
			tokenUserID:   "5cf37266-3473-4006-984f-9325122678b7",
			requestUserID: "d5c65266-3473-4006-984f-9325122678c8", // Different from token
			ruleTests: []struct {
				rule     string
				expected bool
				message  string
			}{
				{rule: auth.RuleAny, expected: true, message: "RuleAny should pass regardless of userID mismatch"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			log, teardown := newUnit(t)
			defer func() {
				if r := recover(); r != nil {
					t.Log("recovered from test", r)
					t.Error(string(debug.Stack()))
				}
				teardown()
			}()

			a := auth.New(auth.Config{
				KeyLookup: &keyStore{},
				Issuer:    "service project",
				Log:       log,
			})

			// Generate the JWT with specified roles
			claims := auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "service project",
					Subject:   tc.tokenUserID,
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
				},
				Roles: tc.roles,
			}

			token, err := a.GenerateToken(kid, claims)
			if err != nil {
				t.Fatalf("Should be able to generate a JWT: %s", err)
			}

			parsedClaims, err := a.Authenticate(context.Background(), "Bearer "+token)
			if err != nil {
				t.Fatalf("Should be able to authenticate the claims: %s", err)
			}

			requestUserID := uuid.MustParse(tc.requestUserID)

			// Test each authorization rule
			for _, rt := range tc.ruleTests {
				err = a.Authorize(context.Background(), parsedClaims, requestUserID, rt.rule)

				if rt.expected && err != nil {
					t.Errorf("%s: %s", rt.message, err)
				}

				if !rt.expected && err == nil {
					t.Errorf("%s", rt.message)
				}
			}
		})
	}
}

func newUnit(t *testing.T) (*logger.Logger, func()) {
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo, "TEST", func(context.Context) string { return "00000000-0000-0000-0000-000000000000" })

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()

		fmt.Println("******************** LOGS ********************")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ********************")
	}

	return log, teardown
}

type keyStore struct{}

func (ks *keyStore) PrivateKey(kid string) (string, error) {
	return privateKeyPEM, nil
}

func (ks *keyStore) PublicKey(kid string) (string, error) {
	return publicKeyPEM, nil
}

const (
	kid = "s4sKIjD9kIRjxs2tulPqGLdxSfgPErRN1Mu3Hd9k9NQ"

	privateKeyPEM = `-----BEGIN PRIVATE KEY-----
MIIEpAIBAAKCAQEAuHIwYZo9v96wpJ0CwrSEmizSDO/41etP7Ovc0ML3EK2qumJR
JbPNxAu1OGswqHTgeBc2xlGy+eSLrvgKMxmykK8WSm0Myrw8J5JdUUrJef3ebL61
KEWT3zTRAxIvWwGXi0rf2T4fwgmZbv9DvW4RmJjreDoD0YRVdbKPanTLrxlxCeXX
9+6wR4PcDW2ymVfJWTQBP2YwB+lGoolimwatrahgZC64bgtTB9hGDvR5q44Fy26R
SV4Q8GvAP43h70V0oLk75ckU2ACtDQ3RioJeVHbMBf7Q/uzBfvMgC6ieJfYBDm7e
suoQqdfltLxvlbUXc4P82Mpp0CDUHySaKB9fRwIDAQABAoIBAA0fuY06/i7N/Rqm
vgtgRdkkhvOltY/Jv7T3brhUpGNwP/16wRRiL/9wPSjyXj12R4Heq6aT1JzO+W4Y
jGfovwRoI467JgBc7BjSuMdv1K3UoFlzfmILBVkJ9zSf5spebsvZ/DCHKOz+ZYej
gSHpmDXtQj/v53/7ybWSCtNZKSRgbPV/UA9LkkyNatiGYZwLUSq0MoNlB5jt33dq
AtGQgt6Mv1OpQ/nHn5yiPQ+OrypIBdQtnHadjeTZAHcOL4deQkjVllzErF6ooxRN
b1IJx70nDvek4aQzM/fQSLc+Pn6cm4VjopkdK2ZcwtD3iUbMlZDsEsaeBxWCJq1e
A9D7KO0CgYEA1XBoFyTNvK5ERDf9SqPXv8OeE8OcXYoZSblijfHRW02zsdRbem5g
R8t+rL5Lo3F9bJG++Up6SoLdZp7g8Xx59+/SanNX/7bfRqTRP0Dk7dSe8ZhjdZoD
llXe1L0sNt3ff5iSp5QeLDsHC8veeXg6qsMXs5BkX8Mzy81kTRKBDA0CgYEA3TnA
f06KxaRVy8CewqVcnc8B8kznIRSlr3Om+AnuH06QnkC9M7+fse30SqqIbOpZ78eg
+U8Mqikj4pYug1EW/cd/08dsH/1977YQN0uKQW2CVVfwgcXVS5q91jwDzmXPS5ie
3f0E8Ke33y+Hfgi7+wWXDzn9SpczrSs4Zj8Tv6MCgYEAp/ElhihoVfFNN9xSShu3
VGVoiaad76ANG8xp9sRydgQiw2cTf7c/vN4q02N8gqN3DCl1+hOXO+/So7+ZwYmv
Mt6aUzZk5ImRD7X6C2pVd6mYUwMUJ2HDNtRkXEJpetaD2JNFueQ7BZSAi9CjQjLO
/rQ5fwm7YPpoVBFNvbM0pTUCgYEA00jVEMFyyFCroG3XMw+75PxQX5oTJQUTOP67
+SmCw2PFu18ZVNMvMkCRkL5OjbdFLjz7ASD+d4XTQBUvVzubOcXPz/Qm0GbKYKcB
1c3Pva1ZTSkwCsFndU3VAUdQW0/hK9IX6Ow+S5njgsViIn01DAnKvEAmKZc/Q6SD
uSOFOQECgYA+o+q4SoPYXqGJmCrF1aVbSqUdTCrH1dF5YA6tB6T/gGr1Zrf4m+Ww
TCBG0rcl6+CwaFqfnvF10S9KiEwpscuCxy5M3BVsCKOK6+GbOwiBuqjH/0CTAtAq
isSdmgnbs9gCLby1Tjo6x8iallTG6dL7OnX4YZAVXPKoTRilt+gFCg==
-----END PRIVATE KEY-----
`

	publicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuHIwYZo9v96wpJ0CwrSE
mizSDO/41etP7Ovc0ML3EK2qumJRJbPNxAu1OGswqHTgeBc2xlGy+eSLrvgKMxmy
kK8WSm0Myrw8J5JdUUrJef3ebL61KEWT3zTRAxIvWwGXi0rf2T4fwgmZbv9DvW4R
mJjreDoD0YRVdbKPanTLrxlxCeXX9+6wR4PcDW2ymVfJWTQBP2YwB+lGoolimwat
rahgZC64bgtTB9hGDvR5q44Fy26RSV4Q8GvAP43h70V0oLk75ckU2ACtDQ3RioJe
VHbMBf7Q/uzBfvMgC6ieJfYBDm7esuoQqdfltLxvlbUXc4P82Mpp0CDUHySaKB9f
RwIDAQAB
-----END PUBLIC KEY-----
`
)
