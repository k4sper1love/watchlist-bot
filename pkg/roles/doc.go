// Package roles provides utilities for managing user roles in the application.
//
// It defines a set of predefined roles and methods for working with them,
// including role comparison, access checks, and role transitions.
//
// Features:
// - Role Enumeration: Predefined roles (`User`, `Helper`, `Admin`, etc.) ordered by hierarchy.
// - String Representation: Each role has a human-readable name (e.g., "user", "admin").
// - Access Control: Methods to check if a role has sufficient access for a required role.
// - Role Comparison: Compare two roles to determine their relative hierarchy.
// - Role Transitions: Move to the next or previous role in the hierarchy.
//
// Usage:
//
// Roles are represented as integers (`Role` type), with higher values indicating greater privileges.
//
// Example:
//
//	role := roles.Admin
//	fmt.Println(role.String())       // Output: "admin"
//	fmt.Println(role.HasAccess(roles.Helper)) // Output: true
//	fmt.Println(role.NextRole().String())     // Output: "superAdmin"
//
// This package ensures consistent and flexible role management across the application.
package roles
