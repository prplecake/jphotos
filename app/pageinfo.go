package app

// PageInfo holds pieces of information that slot into standard
// placeholders on pages.
// Any field nay be null if applicable
type PageInfo struct {
	// Username is populated if the user is logged in
	Username *string

	// Error is populated if an error message should be shown
	Error []string
}

// ErrorInfo is used to fill the error with redirect page
type ErrorInfo struct {
	Info, RedirectLink string
	RedirectTimer      int
}
