// Package auth provides Open Policy Agent (OPA) Rego policies for authentication and authorization.
//
// Authentication Flow:
//
// The authentication.rego file implements a JWT validation policy using OPA's Rego language.
// Here's how it works:
//
// 1. Package Declaration:
//   - The policy belongs to the "sales.rego" package, defining its namespace.
//
// 2. Imports:
//   - It imports "rego.v1" to access functions like decode_verify for JWT validation.
//
// 3. Default Authentication:
//   - Sets a default "auth" variable to false, meaning authentication fails by default.
//
// 4. Authentication Rule:
//   - The "auth" rule returns true only when JWT verification succeeds.
//   - It extracts the verification result using the verify_jwt function.
//   - If the verification is valid (valid = true), then auth will be true.
//
// 5. JWT Verification:
//   - The verify_jwt function uses io.jwt.decode_verify to validate the JWT token.
//   - It checks:
//   - The token's signature against the provided key
//   - The issuer (iss) claim matches the expected value
//
// 6. Input Requirements:
//   - The policy expects an input object containing:
//   - Token: The JWT token to verify
//   - Key: The public key to verify the token's signature
//   - ISS: The expected issuer claim value
//
// This authentication mechanism verifies identity but doesn't determine authorization.
// After authentication succeeds, separate authorization policies would determine
// what actions the authenticated entity can perform.
//
// ------------------------------------------------------------------------------------------------
//
// Authorization Flow:
//
// The authorization.rego file implements role-based access control using OPA's Rego language.
// Here's how it works:
//
// 1. Package Declaration:
//   - Uses the same "sales.rego" package as authentication.
//
// 2. Default Rules:
//   - Sets all authorization rules to false by default:
//   - rule_any: Allows access to any authenticated user with valid roles
//   - rule_admin_only: Restricts access to administrators only
//   - rule_user_only: Restricts access to regular users only
//   - rule_admin_or_subject: Allows access to administrators or if the user is accessing their own data
//
// 3. Role Definitions:
//   - Defines constants for different roles:
//   - role_user: "USER" - Regular user role
//   - role_admin: "ADMIN" - Administrator role
//   - role_all: Set containing both USER and ADMIN roles
//
// 4. Authorization Rules:
//   - rule_any: Grants access if the user has any valid role defined in role_all
//   - rule_admin_only: Grants access only if the user has the ADMIN role
//
// 5. Input Requirements:
//   - The policy expects an input object containing:
//   - Roles: An array of roles assigned to the authenticated user
//
// The authorization mechanism works in conjunction with authentication to provide
// complete access control. After a user is authenticated, these rules determine
// what actions they are permitted to perform based on their assigned roles.
package auth
