package client

type TokenReviewReturn struct {
  ApiVersion  string  `json:"apiVersion"`
  Kind  string `json:"kind"`
  Status  TokenReviewStatus  `json:"status"`
}

// TokenReviewStatus is the result of the token authentication request.
type TokenReviewStatus struct {
	// Authenticated is true if the token is valid
	Authenticated bool `json:"authenticated"`
	// User contains information about the authenticated user.
	User UserInfo `json:"user"`
}

// UserInfo contains information about the user
type UserInfo struct {
	// The name that uniquely identifies this user among all active users.
	Username string `json:"username"`
	// A unique value that identifies this user across time. If this user is
	// deleted and another user by the same name is added, they will have
	// different UIDs.
	UID string `json:"uid,omitempty"`
	// The names of groups this user is a part of.
	Groups []string `json:"groups,omitempty"`
	// Any additional information provided by the authenticator.
	Extra map[string][]string `json:"extra,omitempty"`
}